package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	UserTokenKeyPrefix        = "user_token:%d"
	UserRefreshTokenKeyPrefix = "user_refresh_token:%d"
)

var (
	UserTokenExpireDuration    = time.Minute * getEnvAsInt64OrDefault("JWT_ACCESS_TIME_MINUTE", 60)
	RefreshTokenExpireDuration = getEnvAsInt64OrDefault("JWT_REFRESH_TIME_HOUR", 24) * time.Hour
)

func GetUserTokenKey(userID int64) string {
	return fmt.Sprintf(UserTokenKeyPrefix, userID)
}

func GetUserRefreshTokenKey(userID int64) string {
	return fmt.Sprintf(UserRefreshTokenKeyPrefix, userID)
}

func getEnvAsInt64OrDefault(key string, defaultValue int64) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(intValue)
		}
	}
	return time.Duration(defaultValue)
}
