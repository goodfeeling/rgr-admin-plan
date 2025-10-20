package email

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	UserIdTokenKeyPrefix = "reset_passwd__token:%s"
)

var (
	UserTokenExpireDuration = time.Hour * 1
)

func GetUserIdTokenKey(userId int64) string {
	return fmt.Sprintf(UserIdTokenKeyPrefix, userId)
}
func GetEnvAsInt64OrDefault(key string, defaultValue int64) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(intValue)
		}
	}
	return time.Duration(defaultValue)
}
