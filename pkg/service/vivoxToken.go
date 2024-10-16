// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Payload struct {
	Iss string `json:"iss"`
	Exp int64  `json:"exp"`
	Vxa string `json:"vxa"`
	Vxi int64  `json:"vxi"`
	Sub string `json:"sub,omitempty"`
	F   string `json:"f,omitempty"`
	T   string `json:"t,omitempty"`
}
type Token struct {
	AccessToken string
	Uri         string
}

const (
	ActionJoin           = "join"
	ActionJoinMuted      = "join_muted"
	ActionKick           = "kick"
	ActionLogin          = "login"
	ActionMute           = "mute"
	ChannelEcho          = "-e-"
	ChannelNonPositional = "-g-"
	ChannelPositional    = "-d-"
)

func GenerateVivocLoginToken(
	signingKey, issuer, domain, username string,
	serialNumber int64,
	expiredAt time.Time) (token string, uri string, err error) {
	header := make(map[string]any)
	expireAtFloat := float64(expiredAt.Unix())
	payload := Payload{
		Iss: issuer,
		Exp: int64(expireAtFloat),
		Vxa: ActionLogin,
		Vxi: serialNumber,
		F:   protocol + ":" + userName(issuer, username) + "@" + domain,
	}

	t, e := makeVivoxToken(signingKey, header, payload)
	if e != nil {
		logrus.Error(e)

		return "", "", e
	}

	return t, payload.T, nil
}
func GenerateVivoxJoinToken(
	signingKey, issuer, domain, username, channelType, channelID string,
	uniqueNumber int64,
	expiredAt time.Time) (token string, uri string, err error) {
	header := make(map[string]any)
	expireAtFloat := float64(expiredAt.Unix())
	payload := Payload{
		Iss: issuer,
		Exp: int64(expireAtFloat),
		Vxa: ActionJoin,
		Vxi: uniqueNumber,
		F:   protocol + ":" + userName(issuer, username) + "@" + domain,
		T:   protocol + ":" + channelName(channelType, issuer, channelID) + "@" + domain,
	}

	t, e := makeVivoxToken(signingKey, header, payload)
	if e != nil {
		logrus.Error(e)

		return "", "", e
	}

	return t, payload.T, nil
}
func GenerateVivoxJoinMuteToken(
	signingKey, issuer, domain, username, channelType, channelID string,
	serialNumber int64,
	expiredAt time.Time) (token string, uri string, err error) {
	header := make(map[string]any)
	expireAtFloat := float64(expiredAt.Unix())
	payload := Payload{
		Iss: issuer,
		Exp: int64(expireAtFloat),
		Vxa: ActionJoinMuted,
		Vxi: serialNumber,
		F:   protocol + ":" + userName(issuer, username) + "@" + domain,
		T:   protocol + ":" + channelName(channelType, issuer, channelID) + "@" + domain,
	}

	t, e := makeVivoxToken(signingKey, header, payload)
	if e != nil {
		logrus.Error(e)

		return "", "", e
	}

	return t, payload.T, nil
}
func GenerateVivoxKickToken(
	signingKey, issuer, domain, fromUserID, toUserID, channelType, channelID string,
	serialNumber int64,
	expiredAt time.Time) (token string, uri string, err error) {
	header := make(map[string]any)
	expireAtFloat := float64(expiredAt.Unix())
	payload := Payload{
		Iss: issuer,
		Exp: int64(expireAtFloat),
		Vxa: ActionKick,
		Vxi: serialNumber,
		Sub: protocol + ":" + userName(issuer, toUserID) + "@" + domain,
		F:   protocol + ":" + userName(issuer, fromUserID) + "@" + domain,
		T:   protocol + ":" + channelName(channelType, issuer, channelID) + "@" + domain,
	}

	t, e := makeVivoxToken(signingKey, header, payload)
	if e != nil {
		logrus.Error(e)

		return "", "", e
	}

	return t, payload.T, nil
}

func channelName(channelType, issuer, channelID string) string {
	return channelPrefix + "-" + channelType + "-" + issuer + "." + channelID
}
func userName(issuer, userID string) string {
	return "." + issuer + "." + userID + "."
}
func makeVivoxToken(signingKey string,
	header map[string]any,
	payload Payload) (string, error) {
	headerMarshal, err := json.Marshal(header)
	if err != nil {
		text := fmt.Sprintf("error make token: %v", err)
		logrus.Error(text)

		return "", errors.New(text)
	}
	encodedHeader := Base64URLEncode(string(headerMarshal))
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		text := fmt.Sprintf("error make token: %v", err)
		logrus.Error(text)

		return "", errors.New(text)
	}
	payloadString := string(payloadMarshal)
	encodedPayload := Base64URLEncode(payloadString)
	signature, err := Sign(header, payload, signingKey)
	if err != nil {
		text := fmt.Sprintf("error make token: %v", err)
		logrus.Error(text)

		return "", errors.New(text)
	}

	return strings.Join([]string{encodedHeader, encodedPayload, signature}, "."), nil
}

func Sign(header map[string]any, payload Payload, key string) (string, error) {
	headerMarshal, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to encode header when signing token with error %w", err)
	}
	headerString := string(headerMarshal)
	payloadMarshal, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode payload when signing token with error %w", err)
	}
	payloadString := string(payloadMarshal)
	base64Header := Base64URLEncode(headerString)
	base64Payload := Base64URLEncode(payloadString)

	return HmacBase64Encode(base64Header+"."+base64Payload, key), nil
}
func Base64URLEncode(str string) string {
	return base64EncodeAndReplaceChar([]byte(str))
}
func HmacBase64Encode(seed, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(seed))

	return base64EncodeAndReplaceChar(h.Sum(nil))
}
func base64EncodeAndReplaceChar(byteArray []byte) string {
	encoded := base64.StdEncoding.EncodeToString(byteArray)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.Trim(encoded, "=")

	return encoded
}
