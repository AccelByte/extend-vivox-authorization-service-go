// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package common

import (
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Trim(value, "\"")
	}

	return fallback
}

func GetEnvInt(key string, fallback int) int {
	str := GetEnv(key, strconv.Itoa(fallback))
	val, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}

	return val
}

func getBasePath() string {
	basePath := GetEnv("BASE_PATH", "/vivoxauth")
	if !strings.HasPrefix(basePath, "/") {
		slog.Default().Error("BASE_PATH envar is invalid, no leading '/' found. Valid example: /basePath")
		os.Exit(1)
	}

	return basePath
}

func SetEnv(key, value string) error {
	err := os.Setenv(key, value)
	if err == nil {
		return nil
	}

	return err
}

func RandomNumber(n int) int64 {
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
