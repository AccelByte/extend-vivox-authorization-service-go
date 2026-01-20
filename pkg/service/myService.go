// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"
	"strings"
	"time"

	utils "extend-rtu-vivox-authorization-service/pkg/common"
	pb "extend-rtu-vivox-authorization-service/pkg/pb"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MyServiceServerImpl struct {
	pb.UnimplementedServiceServer
	tokenRepo   repository.TokenRepository
	configRepo  repository.ConfigRepository
	refreshRepo repository.RefreshTokenRepository
	claims      *Claims
}

func NewMyServiceServer(
	tokenRepo repository.TokenRepository,
	configRepo repository.ConfigRepository,
	refreshRepo repository.RefreshTokenRepository,
	claims *Claims,
) *MyServiceServerImpl {
	return &MyServiceServerImpl{
		tokenRepo:   tokenRepo,
		configRepo:  configRepo,
		refreshRepo: refreshRepo,
		claims:      claims,
	}
}

var (
	issuer     = utils.GetEnv("VIVOX_ISSUER", "")
	domain     = utils.GetEnv("VIVOX_DOMAIN", "")
	signingKey = utils.GetEnv("VIVOX_SIGNING_KEY", "")

	// optional
	expiry   = utils.GetEnvInt("VIVOX_DEFAULT_EXPIRY", 90)
	protocol = utils.GetEnv("VIVOX_PROTOCOL", "sip")
	cPrefix  = utils.GetEnv("VIVOX_CHANNEL_PREFIX", "confctl")
)

func (g MyServiceServerImpl) GenerateVivoxToken(
	ctx context.Context, req *pb.GenerateVivoxTokenRequest,
) (*pb.GenerateVivoxTokenResponse, error) {
	var accessToken, uri string
	var err error

	if errValidate := g.validateRequest(req); errValidate != nil {
		return nil, errValidate
	}

	expiry := time.Now().Add(time.Duration(expiry) * time.Second)
	uniqueNum := utils.RandomNumber(4)
	cTypeStr := req.ChannelType.String()

	// Route based on Enum
	switch req.Type {
	case pb.GenerateVivoxTokenRequestType_login:
		accessToken, uri, err = GenerateVivocLoginToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			uniqueNum,
			expiry,
			g.claims,
		)

	case pb.GenerateVivoxTokenRequestType_join:
		accessToken, uri, err = GenerateVivoxJoinToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			cTypeStr,
			req.ChannelId,
			uniqueNum,
			expiry,
			g.claims,
		)

	case pb.GenerateVivoxTokenRequestType_join_muted:
		accessToken, uri, err = GenerateVivoxJoinMuteToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			cTypeStr,
			req.ChannelId,
			uniqueNum,
			expiry,
			g.claims,
		)

	case pb.GenerateVivoxTokenRequestType_kick:
		accessToken, uri, err = GenerateVivoxKickToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			req.TargetUsername,
			cTypeStr,
			req.ChannelId,
			uniqueNum,
			expiry,
			g.claims,
		)

	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported action type: %s", req.Type.String())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generate Vivox auth token: %v", err)
	}

	// Return the token
	return &pb.GenerateVivoxTokenResponse{AccessToken: accessToken, Uri: uri}, nil
}

func (g *MyServiceServerImpl) validateRequest(req *pb.GenerateVivoxTokenRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request body cannot be nil")
	}

	isInvalid := func(s string) bool {
		return s == "" || strings.ToLower(s) == "string"
	}

	if isInvalid(signingKey) || isInvalid(issuer) || isInvalid(domain) {
		return status.Error(codes.Internal, "vivox configuration (key/issuer/domain) is missing")
	}

	if req.Type == pb.GenerateVivoxTokenRequestType_generatevivoxtokenrequest_type_unknown {
		return status.Error(codes.InvalidArgument, "a valid action type must be provided")
	}

	if isInvalid(req.Username) {
		return status.Error(codes.InvalidArgument, "username is required")
	}

	switch req.Type {
	case pb.GenerateVivoxTokenRequestType_join,
		pb.GenerateVivoxTokenRequestType_join_muted:
		if isInvalid(req.ChannelId) {
			return status.Error(codes.InvalidArgument, "channel_id is required for join actions")
		}
		if isInvalid(req.Username) {
			return status.Error(codes.InvalidArgument, "username is required for kick action")
		}
		cType := req.ChannelType.String()
		if isInvalid(cType) || strings.Contains(strings.ToLower(cType), "unknown") {
			return status.Error(codes.InvalidArgument, "valid channel_type is required. Please use one of these values: echo, positional, or nonpositional.")
		}

	case pb.GenerateVivoxTokenRequestType_kick:
		if isInvalid(req.ChannelId) {
			return status.Error(codes.InvalidArgument, "channel_id are required for kick")
		}
		if isInvalid(req.TargetUsername) {
			return status.Error(codes.InvalidArgument, "target_username are required for kick")
		}
	}

	return nil
}
