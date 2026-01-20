// Copyright (c) 2024-2026 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"context"
	"fmt"
	"testing"

	"extend-rtu-vivox-authorization-service/pkg/common"
	pb "extend-rtu-vivox-authorization-service/pkg/pb"
	"extend-rtu-vivox-authorization-service/pkg/service/mocks"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination ./mocks/server_mock.go -package mocks extend-custom-guild-service/pkg/pb myServiceServer
//go:generate mockgen -destination ./mocks/repo_mock.go -package mocks github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository TokenRepository,ConfigRepository,RefreshTokenRepository

func TestMyServiceServerImpl_GenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		req           *pb.GenerateVivoxTokenRequest
		claims        *Claims
		wantErr       bool
		expectedErr   error
		expectedToken string
		expectedUri   string
	}{
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-login-token.htm
			name: "login test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_login,
				Username:       "jerky",
				ChannelId:      "933000",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionLogin,
				Vxi: 933000,
				F:   "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjkzMzAwMCwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5qZXJreS5AdGxhLnZpdm94LmNvbSIsImlzcyI6ImJsaW5kbWVsb24tQXBwTmFtZS1kZXYiLCJ2eGEiOiJsb2dpbiIsImV4cCI6MTYwMDM0OTQwMH0.YJwjX0P2Pjk1dzFpIo1fjJM21pphfBwHm8vShJib8ds",
			expectedUri:   fmt.Sprintf("sip:%s@tla.vivox.com", channelName(ChannelEcho, "blindmelon-AppName-dev", "933000")),
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-join-token.htm
			name: "join test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join,
				Username:       "jerky",
				ChannelId:      "sip:confctl-g-blindmelon.testchannel@tla.vivox.com",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_nonpositional,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionJoin,
				Vxi: 444000,
				F:   "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjQ0NDAwMCwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5qZXJreS5AdGxhLnZpdm94LmNvbSIsImlzcyI6ImJsaW5kbWVsb24tQXBwTmFtZS1kZXYiLCJ2eGEiOiJqb2luIiwidCI6InNpcDpjb25mY3RsLWctYmxpbmRtZWxvbi1BcHBOYW1lLWRldi50ZXN0Y2hhbm5lbEB0bGEudml2b3guY29tIiwiZXhwIjoxNjAwMzQ5NDAwfQ.u7us5eCxOBtuEZuDg1HapEEgxLedLaliIy7gOMfbeko",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-join-muted-token.htm
			name: "join muted test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join_muted,
				Username:       "jerky",
				ChannelId:      "sip:confctl-g-blindmelon.testchannel@tla.vivox.com",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_nonpositional,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionJoinMuted,
				Vxi: 542680,
				F:   "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjU0MjY4MCwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5qZXJreS5AdGxhLnZpdm94LmNvbSIsImlzcyI6ImJsaW5kbWVsb24tQXBwTmFtZS1kZXYiLCJ2eGEiOiJqb2luX211dGVkIiwidCI6InNpcDpjb25mY3RsLWctYmxpbmRtZWxvbi1BcHBOYW1lLWRldi50ZXN0Y2hhbm5lbEB0bGEudml2b3guY29tIiwiZXhwIjoxNjAwMzQ5NDAwfQ.N6sZL3F3e-p2KLQlMweXnbGNzE7Qc91rn_uqCEtRjsc",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-kick-token.htm
			name: "user kick user test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_kick,
				Username:       "beef",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Sub: "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				Exp: 1600349400,
				Vxa: ActionKick,
				Vxi: 665000,
				F:   "sip:.blindmelon-AppName-dev.beef.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjY2NTAwMCwic3ViIjoic2lwOi5ibGluZG1lbG9uLUFwcE5hbWUtZGV2Lmplcmt5LkB0bGEudml2b3guY29tIiwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5iZWVmLkB0bGEudml2b3guY29tIiwiaXNzIjoiYmxpbmRtZWxvbi1BcHBOYW1lLWRldiIsInZ4YSI6ImtpY2siLCJ0Ijoic2lwOmNvbmZjdGwtZy1ibGluZG1lbG9uLUFwcE5hbWUtZGV2LnRlc3RjaGFubmVsQHRsYS52aXZveC5jb20iLCJleHAiOjE2MDAzNDk0MDB9.kKnWD3smth6KUuRaY11O-yqAbXy2L2wDZeIoDK_098c",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-kick-token.htm
			name: "admin kick user from channel test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_kick,
				Username:       "blindmelon-AppName-dev-Admin",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Sub: "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				Exp: 1600349400,
				Vxa: ActionKick,
				Vxi: 8000,
				F:   "sip:blindmelon-AppName-dev-Admin@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjgwMDAsInN1YiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5qZXJreS5AdGxhLnZpdm94LmNvbSIsImYiOiJzaXA6YmxpbmRtZWxvbi1BcHBOYW1lLWRldi1BZG1pbkB0bGEudml2b3guY29tIiwiaXNzIjoiYmxpbmRtZWxvbi1BcHBOYW1lLWRldiIsInZ4YSI6ImtpY2siLCJ0Ijoic2lwOmNvbmZjdGwtZy1ibGluZG1lbG9uLUFwcE5hbWUtZGV2LnRlc3RjaGFubmVsQHRsYS52aXZveC5jb20iLCJleHAiOjE2MDAzNDk0MDB9.7Fn08cctqltxNxPAAeOhPQd4KCsmT1ue1EDIxUNQ3gg",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-kick-token.htm
			name: "admin kick user from server test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_kick,
				Username:       "blindmelon-AppName-dev-Admin",
				ChannelId:      "", // entire server
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Sub: "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				Exp: 1600349400,
				Vxa: ActionKick,
				Vxi: 613642,
				F:   "sip:blindmelon-AppName-dev-Admin@tla.vivox.com",
				T:   "sip:blindmelon-AppName-dev-service@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjYxMzY0Miwic3ViIjoic2lwOi5ibGluZG1lbG9uLUFwcE5hbWUtZGV2Lmplcmt5LkB0bGEudml2b3guY29tIiwiZiI6InNpcDpibGluZG1lbG9uLUFwcE5hbWUtZGV2LUFkbWluQHRsYS52aXZveC5jb20iLCJpc3MiOiJibGluZG1lbG9uLUFwcE5hbWUtZGV2IiwidnhhIjoia2ljayIsInQiOiJzaXA6YmxpbmRtZWxvbi1BcHBOYW1lLWRldi1zZXJ2aWNlQHRsYS52aXZveC5jb20iLCJleHAiOjE2MDAzNDk0MDB9.jinc73lQ_ZSN4Mb8WLFK7Clu-Se9LG-QifXKfpaa3g4",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-drop-all-token.htm
			name: "admin drop all test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_kick,
				Username:       "blindmelon-AppName-dev-Admin",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionKick,
				Vxi: 729614,
				F:   "sip:blindmelon-AppName-dev-Admin@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjcyOTYxNCwiZiI6InNpcDpibGluZG1lbG9uLUFwcE5hbWUtZGV2LUFkbWluQHRsYS52aXZveC5jb20iLCJpc3MiOiJibGluZG1lbG9uLUFwcE5hbWUtZGV2IiwidnhhIjoia2ljayIsInQiOiJzaXA6Y29uZmN0bC1nLWJsaW5kbWVsb24tQXBwTmFtZS1kZXYudGVzdGNoYW5uZWxAdGxhLnZpdm94LmNvbSIsImV4cCI6MTYwMDM0OTQwMH0.bvNDTcXIHRpckAiFzy3wPiwgYGks4E9_WuEA5HMGpGE",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-mute-token.htm
			name: "user mute test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join_muted,
				Username:       "beef",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Sub: "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				Exp: 1600349400,
				Vxa: ActionMute,
				Vxi: 123456,
				F:   "sip:.blindmelon-AppName-dev.beef.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjEyMzQ1Niwic3ViIjoic2lwOi5ibGluZG1lbG9uLUFwcE5hbWUtZGV2Lmplcmt5LkB0bGEudml2b3guY29tIiwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5iZWVmLkB0bGEudml2b3guY29tIiwiaXNzIjoiYmxpbmRtZWxvbi1BcHBOYW1lLWRldiIsInZ4YSI6Im11dGUiLCJ0Ijoic2lwOmNvbmZjdGwtZy1ibGluZG1lbG9uLUFwcE5hbWUtZGV2LnRlc3RjaGFubmVsQHRsYS52aXZveC5jb20iLCJleHAiOjE2MDAzNDk0MDB9.vM9zkCXTORjgv8w7eiMHHHkc4DumTwR_-I06y4SnpHA",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-mute-token.htm
			name: "admin mute test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join_muted,
				Username:       "beef",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Sub: "sip:.blindmelon-AppName-dev.jerky.@tla.vivox.com",
				Exp: 1600349400,
				Vxa: ActionMute,
				Vxi: 654321,
				F:   "sip:blindmelon-AppName-dev-Admin@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjY1NDMyMSwic3ViIjoic2lwOi5ibGluZG1lbG9uLUFwcE5hbWUtZGV2Lmplcmt5LkB0bGEudml2b3guY29tIiwiZiI6InNpcDpibGluZG1lbG9uLUFwcE5hbWUtZGV2LUFkbWluQHRsYS52aXZveC5jb20iLCJpc3MiOiJibGluZG1lbG9uLUFwcE5hbWUtZGV2IiwidnhhIjoibXV0ZSIsInQiOiJzaXA6Y29uZmN0bC1nLWJsaW5kbWVsb24tQXBwTmFtZS1kZXYudGVzdGNoYW5uZWxAdGxhLnZpdm94LmNvbSIsImV4cCI6MTYwMDM0OTQwMH0.ix0mFGS1XDXCBXH044f6B2JxutExbH2hZjGqZAwoHH8",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-mute-all-token.htm
			name: "user mute all test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join_muted,
				Username:       "beef",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionMute,
				Vxi: 19283,
				F:   "sip:.blindmelon-AppName-dev.beef.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjE5MjgzLCJmIjoic2lwOi5ibGluZG1lbG9uLUFwcE5hbWUtZGV2LmJlZWYuQHRsYS52aXZveC5jb20iLCJpc3MiOiJibGluZG1lbG9uLUFwcE5hbWUtZGV2IiwidnhhIjoibXV0ZSIsInQiOiJzaXA6Y29uZmN0bC1nLWJsaW5kbWVsb24tQXBwTmFtZS1kZXYudGVzdGNoYW5uZWxAdGxhLnZpdm94LmNvbSIsImV4cCI6MTYwMDM0OTQwMH0.fARLW2eX10ZbiIl_5WIg4bhPbYIhn2xfCcUNySfwBMs",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-mute-all-token.htm
			name: "admin mute all test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_join_muted,
				Username:       "blindmelon-AppName-dev-Admin",
				ChannelId:      "testchannel",
				ChannelType:    pb.GenerateVivoxTokenRequestChannelType_echo,
				TargetUsername: "", // all users
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: ActionMute,
				Vxi: 825647,
				F:   "sip:blindmelon-AppName-dev-Admin@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjgyNTY0NywiZiI6InNpcDpibGluZG1lbG9uLUFwcE5hbWUtZGV2LUFkbWluQHRsYS52aXZveC5jb20iLCJpc3MiOiJibGluZG1lbG9uLUFwcE5hbWUtZGV2IiwidnhhIjoibXV0ZSIsInQiOiJzaXA6Y29uZmN0bC1nLWJsaW5kbWVsb24tQXBwTmFtZS1kZXYudGVzdGNoYW5uZWxAdGxhLnZpdm94LmNvbSIsImV4cCI6MTYwMDM0OTQwMH0.zpyvlBbVAKatuCeELb0Q1PsCb4tg0yacneL_sYHIVaw",
			expectedUri:   "",
		},
		{
			// Test values taken from:
			// https://docs.vivox.com/v5/general/unity/15_1_160000/en-us/access-token-guide/access-token-examples/example-transcription-token.htm
			name: "transcription test",
			req: &pb.GenerateVivoxTokenRequest{
				Type:           pb.GenerateVivoxTokenRequestType_login,
				Username:       "beef",
				ChannelId:      "testChannel",
				ChannelType:    0,
				TargetUsername: "jerky",
			},
			claims: &Claims{
				Iss: "blindmelon-AppName-dev",
				Exp: 1600349400,
				Vxa: "trxn",
				Vxi: 542680,
				F:   "sip:.blindmelon-AppName-dev.beef.@tla.vivox.com",
				T:   "sip:confctl-g-blindmelon-AppName-dev.testchannel@tla.vivox.com",
			},
			wantErr:       false,
			expectedToken: "e30.eyJ2eGkiOjU0MjY4MCwiZiI6InNpcDouYmxpbmRtZWxvbi1BcHBOYW1lLWRldi5iZWVmLkB0bGEudml2b3guY29tIiwiaXNzIjoiYmxpbmRtZWxvbi1BcHBOYW1lLWRldiIsInZ4YSI6InRyeG4iLCJ0Ijoic2lwOmNvbmZjdGwtZy1ibGluZG1lbG9uLUFwcE5hbWUtZGV2LnRlc3RjaGFubmVsQHRsYS52aXZveC5jb20iLCJleHAiOjE2MDAzNDk0MDB9.-A0w_fcPCZaG5NMksnbSrGSVXNNt25YqlRjcKcLkGnA",
			expectedUri:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := common.SetEnv("VIVOX_ISSUER", "blindmelon-AppName-dev")
			err = common.SetEnv("VIVOX_DOMAIN", "tla.vivox.com")
			err = common.SetEnv("VIVOX_SIGNING_KEY", "secret!")

			if err != nil {
				t.Fatalf("Could not set env variable: %v", err)
			}

			tokenRepo := mocks.NewMockTokenRepository(ctrl)
			refreshRepo := mocks.NewMockRefreshTokenRepository(ctrl)
			configRepo := mocks.NewMockConfigRepository(ctrl)
			service := NewMyServiceServer(tokenRepo, configRepo, refreshRepo, tt.claims)

			// when
			res, err := service.GenerateVivoxToken(context.Background(), tt.req)

			// then
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedToken, res.AccessToken)
			}
		})
	}
}
