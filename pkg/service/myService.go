// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"
	pb "extend-rtu-vivox-authorization-service/pkg/pb"
	"math/rand"
	"os"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MyServiceServerImpl struct {
	pb.UnimplementedServiceServer
	tokenRepo   repository.TokenRepository
	configRepo  repository.ConfigRepository
	refreshRepo repository.RefreshTokenRepository
}

func NewMyServiceServer(
	tokenRepo repository.TokenRepository,
	configRepo repository.ConfigRepository,
	refreshRepo repository.RefreshTokenRepository,
) *MyServiceServerImpl {
	return &MyServiceServerImpl{
		tokenRepo:   tokenRepo,
		configRepo:  configRepo,
		refreshRepo: refreshRepo,
	}
}

// TODO move this to config
var (
	signingKey    = ""
	issuer        = "demo"
	domain        = "tla.vivox.com"
	channelPrefix = "confctl"
	protocol      = "sip"
	defaultExpiry = 90
	vIssuer       = os.Getenv("VIVOX_ISSUER")
	vDomain       = os.Getenv("VIVOX_DOMAIN")
	vSigningKey   = os.Getenv("VIVOX_SIGNING_KEY")
)

func (g MyServiceServerImpl) GenerateVivoxToken(
	ctx context.Context, req *pb.GenerateVivoxTokenRequest,
) (*pb.GenerateVivoxTokenResponse, error) {

	if vIssuer != "" {
		issuer = vIssuer
	}

	if vDomain != "" {
		domain = vDomain
	}

	if vSigningKey != "" {
		signingKey = vSigningKey
	}

	// TODO check the serial number
	accessToken, uri, err := GenerateVivoxJoinToken(
		signingKey,
		issuer,
		domain,
		req.Username,
		req.ChannelType.String(),
		req.ChannelId,
		randomNumber(4),
		time.Now().Add(time.Duration(defaultExpiry)),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generate Vivox auth token: %v", err)
	}

	// Return the token
	return &pb.GenerateVivoxTokenResponse{AccessToken: accessToken, Uri: uri}, nil
}

func randomNumber(n int) int64 {
	if n <= 0 {
		return 0
	}

	lowerBound := int64(1)
	upperBound := int64(9)

	for i := 1; i < n; i++ {
		lowerBound *= 10
		upperBound = upperBound*10 + 9
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNum := lowerBound + r.Int63n(upperBound-lowerBound+1)

	return randomNum
}
