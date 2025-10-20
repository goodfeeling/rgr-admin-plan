// src/infrastructure/cache/sts_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain/sys/files"
	"github.com/redis/go-redis/v9"
)

type STSCacheService struct {
	redisClient *redis.Client
}

func NewSTSCacheService(redisClient *redis.Client) *STSCacheService {
	return &STSCacheService{
		redisClient: redisClient,
	}
}

// 缓存STS Token
func (s *STSCacheService) SetSTSToken(
	ctx context.Context,
	key string, token *files.STSTokenCache,
	expiration time.Duration,
) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, data, expiration).Err()
}

// 获取缓存的STS Token
func (s *STSCacheService) GetSTSToken(ctx context.Context, key string) (*files.STSTokenCache, error) {
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var token files.STSTokenCache
	err = json.Unmarshal([]byte(data), &token)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if token.Expiration.Before(time.Now()) {
		s.DeleteSTSToken(ctx, key)
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}

// 删除STS Token
func (s *STSCacheService) DeleteSTSToken(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, key).Err()
}

// 生成并缓存Refresh Token
func (s *STSCacheService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	refreshToken := generateRandomToken() // 实现随机token生成函数

	rt := &files.RefreshToken{
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时有效期
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(rt)
	if err != nil {
		return "", err
	}

	// 缓存Refresh Token
	err = s.redisClient.Set(ctx, fmt.Sprintf("refresh_token:%s", refreshToken), data, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// 验证Refresh Token
func (s *STSCacheService) ValidateRefreshToken(ctx context.Context, token string) (*files.RefreshToken, error) {
	key := fmt.Sprintf("refresh_token:%s", token)
	data, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var rt files.RefreshToken
	err = json.Unmarshal([]byte(data), &rt)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if rt.ExpiresAt.Before(time.Now()) {
		s.redisClient.Del(ctx, key)
		return nil, fmt.Errorf("refresh token expired")
	}

	return &rt, nil
}

// 删除Refresh Token
func (s *STSCacheService) DeleteRefreshToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("refresh_token:%s", token)
	return s.redisClient.Del(ctx, key).Err()
}

// 生成随机token的辅助函数
func generateRandomToken() string {
	// 实现随机token生成逻辑
	// 可以使用crypto/rand或者uuid
	return fmt.Sprintf("rt_%d", time.Now().UnixNano())
}
