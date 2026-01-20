// Copyright (c) 2024-2026 AccelByte Inc. All Rights Reserved.
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
	loginToken, _, err := GenerateVivocLoginToken(signingKey, issuer, domain, userID, int64(serialNumber), expiredAt, nil)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJ2eGkiOjEwMDQ3LCJmIjoic2lwOi5kZW1vLmJhbGRlYWdsZS4xOTczLkB0bGEudml2b3guY29tIiwiaXNzIjoiZGVtbyIsInZ4YSI6ImxvZ2luIiwiZXhwIjoxNDUxNjA2NDAwfQ.yJIgDg_l4hvkofzDXQEzuCELuLhurn_DVgF2mmUZls8", loginToken)
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
	loginToken, _, err := GenerateVivoxJoinToken(signingKey, issuer, domain, userID, ChannelNonPositional, channelID, int64(serialNumber), expiredAt, nil)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJ2eGkiOjQ0NjkwNSwiZiI6InNpcDouZGVtby5iYWxkZWFnbGUuMTk3My5AdGxhLnZpdm94LmNvbSIsImlzcyI6ImRlbW8iLCJ2eGEiOiJqb2luIiwidCI6InNpcDpjb25mY3RsZGVtby5RZTNNSGxiU3FAdGxhLnZpdm94LmNvbSIsImV4cCI6MTQ1MTYwNjQwMH0.cz1dH_FDUprLmrOS86R3VIh9h16qAgnbCRkl2Pxp-eI",
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
	loginToken, _, err := GenerateVivoxKickToken(signingKey, issuer, domain, fromUserID, toUserID, ChannelNonPositional, channelID, int64(serialNumber), expiredAt, nil)

	assert.Nil(t, err)
	assert.Equal(t, "e30.eyJ2eGkiOjMwMzE2Nywic3ViIjoic2lwOi5kZW1vLmtpbmdmaXNoZXIuMTM2NC5AdGxhLnZpdm94LmNvbSIsImYiOiJzaXA6LmRlbW8uRGVtby1BZG1pbi5AdGxhLnZpdm94LmNvbSIsImlzcyI6ImRlbW8iLCJ2eGEiOiJraWNrIiwidCI6InNpcDpjb25mY3RsZGVtby5RZTNNSGxiU3FAdGxhLnZpdm94LmNvbSIsImV4cCI6MTQ1MTYwNjQwMH0.AetRLye3w7pYpfhZWudGci8W3bgCET5y0ShZ7hkCHs8",
		loginToken)
}
