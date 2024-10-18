// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"
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

func (g MyServiceServerImpl) GenerateVivoxToken(
	ctx context.Context, req *pb.GenerateVivoxTokenRequest,
) (*pb.GenerateVivoxTokenResponse, error) {
	var accessToken, uri string
	var err error

	switch req.Type.String() {
	case ActionLogin:
		accessToken, uri, err = GenerateVivocLoginToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			utils.RandomNumber(4),
			time.Now().Add(time.Duration(defaultExpiry)),
			g.claims,
		)

	case ActionJoin:
		accessToken, uri, err = GenerateVivoxJoinToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			req.ChannelType.String(),
			req.ChannelId,
			utils.RandomNumber(4),
			time.Now().Add(time.Duration(defaultExpiry)),
			g.claims,
		)

	case ActionJoinMuted:
		accessToken, uri, err = GenerateVivoxJoinMuteToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			req.ChannelType.String(),
			req.ChannelId,
			utils.RandomNumber(4),
			time.Now().Add(time.Duration(defaultExpiry)),
			g.claims,
		)

	case ActionMute:
		accessToken, uri, err = GenerateVivoxJoinMuteToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			req.ChannelType.String(),
			req.ChannelId,
			utils.RandomNumber(4),
			time.Now().Add(time.Duration(defaultExpiry)),
			g.claims,
		)

	case ActionKick:
		accessToken, uri, err = GenerateVivoxKickToken(
			signingKey,
			issuer,
			domain,
			req.Username,
			req.TargetUsername,
			req.ChannelType.String(),
			req.ChannelId,
			utils.RandomNumber(4),
			time.Now().Add(time.Duration(defaultExpiry)),
			g.claims,
		)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generate Vivox auth token: %v", err)
	}

	// Return the token
	return &pb.GenerateVivoxTokenResponse{AccessToken: accessToken, Uri: uri}, nil
}
