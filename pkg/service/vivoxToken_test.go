// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestGenerateTokenLogin(t *testing.T) {
	userID := "baldeagle.1973"
	serialNumber := 10047
	expiredAt, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	if err != nil {
		t.Errorf("error parse time: %v", err)

		return
	}
	loginToken, _, err := GenerateVivocLoginToken(signingKey, issuer, domain, userID, int64(serialNumber), expiredAt)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJpc3MiOiJkZW1vIiwiZXhwIjoxNDUxNjA2NDAwLCJ2eGEiOiJsb2dpbiIsInZ4aSI6MTAwNDcsImYiOiJzaXA6LmRlbW8uYmFsZGVhZ2xlLjE5NzMuQHRsYS52aXZveC5jb20ifQ.LYDtnS20oRhk6mue5gZK0PkKjLWhPoB7w17fLr0cnCw", loginToken)
}
func TestGenerateTokenJoin(t *testing.T) {
	userID := "baldeagle.1973"
	serialNumber := 446905
	expiredAt, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	if err != nil {
		t.Errorf("error parse time: %v", err)

		return
	}
	channelID := "Qe3MHlbSq"
	loginToken, _, err := GenerateVivoxJoinToken(signingKey, issuer, domain, userID, ChannelNonPositional, channelID, int64(serialNumber), expiredAt)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJpc3MiOiJkZW1vIiwiZXhwIjoxNDUxNjA2NDAwLCJ2eGEiOiJqb2luIiwidnhpIjo0NDY5MDUsImYiOiJzaXA6LmRlbW8uYmFsZGVhZ2xlLjE5NzMuQHRsYS52aXZveC5jb20iLCJ0Ijoic2lwOmNvbmZjdGwtZy1kZW1vLlFlM01IbGJTcUB0bGEudml2b3guY29tIn0.U_g8ZuxpJAy66myU-DdhCdjOcsdtT_Rce7cbawWBkxU",
		loginToken)
}
func TestGenerateTokenKick(t *testing.T) {
	fromUserID := "Demo-Admin"
	toUserID := "kingfisher.1364"
	serialNumber := 303167
	expiredAt, err := time.Parse(time.RFC3339, "2016-01-01T00:00:00Z")
	if err != nil {
		t.Errorf("error parse time: %v", err)
		return
	}
	channelID := "Qe3MHlbSq"
	loginToken, _, err := GenerateVivoxKickToken(signingKey, issuer, domain, fromUserID, toUserID, ChannelNonPositional, channelID, int64(serialNumber), expiredAt)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJpc3MiOiJkZW1vIiwiZXhwIjoxNDUxNjA2NDAwLCJ2eGEiOiJraWNrIiwidnhpIjozMDMxNjcsInN1YiI6InNpcDouZGVtby5raW5nZmlzaGVyLjEzNjQuQHRsYS52aXZveC5jb20iLCJmIjoic2lwOi5kZW1vLkRlbW8tQWRtaW4uQHRsYS52aXZveC5jb20iLCJ0Ijoic2lwOmNvbmZjdGwtZy1kZW1vLlFlM01IbGJTcUB0bGEudml2b3guY29tIn0.rdH3AMkfhl0zB-dOKtMlHjkyQit5Lsc-WJJYhUK08Yc",
		loginToken)
}
