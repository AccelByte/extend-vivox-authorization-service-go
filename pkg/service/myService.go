// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"
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

func (g MyServiceServerImpl) GenerateVivoxToken(
	ctx context.Context, req *pb.GenerateVivoxTokenRequest,
) (*pb.GenerateVivoxTokenResponse, error) {
	uri := ""
	accessToken, err := g.generateToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error generate Vivox auth token: %v", err)
	}

	// Return the token
	return &pb.GenerateVivoxTokenResponse{AccessToken: accessToken, Uri: uri}, nil
}

func (g MyServiceServerImpl) generateToken() (string, error) {

	return "", nil
}
