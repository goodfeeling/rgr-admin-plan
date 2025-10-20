package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedisClient(loggerInstance *logger.Logger) (*redis.Client, error) {
	// 从环境变量获取配置
	host := getEnv("REDIS_HOST", "127.0.0.1")
	port := getEnv("REDIS_PORT", "6379")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := getEnv("REDIS_DB", "0")
	poolSizeStr := getEnv("REDIS_POOL_SIZE", "10")

	// 转换配置值
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		poolSize = 10
	}

	// 构建Redis地址
	addr := fmt.Sprintf("%s:%s", host, port)

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		loggerInstance.Error("Failed to connect to Redis", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	loggerInstance.Info("Successfully connected to Redis",
		zap.String("host", host),
		zap.String("port", port),
		zap.Int("db", db))

	return client, nil
}

// 辅助函数：获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
