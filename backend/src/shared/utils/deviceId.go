package utils

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateDefaultDeviceID(ctx *gin.Context) string {
	// 优先使用客户端提供的设备ID
	if deviceID := getProvidedDeviceID(ctx); deviceID != "" {
		return deviceID
	}

	// 基于客户端特征生成设备ID
	return generateDeviceIDFromClientInfo(ctx)
}

func getProvidedDeviceID(ctx *gin.Context) string {
	// 检查各种可能的设备ID头部
	headers := []string{
		"X-Device-ID",
		"Device-ID",
		"X-Client-ID",
		"Client-ID",
	}

	for _, header := range headers {
		if deviceID := ctx.GetHeader(header); deviceID != "" {
			return deviceID
		}
	}

	// 检查查询参数
	if deviceID := ctx.Query("device_id"); deviceID != "" {
		return deviceID
	}

	return ""
}

func generateDeviceIDFromClientInfo(ctx *gin.Context) string {
	var components []string

	// 收集客户端信息
	if ip := ctx.ClientIP(); ip != "" {
		components = append(components, ip)
	}

	if userAgent := ctx.GetHeader("User-Agent"); userAgent != "" {
		components = append(components, userAgent)
	}

	if accept := ctx.GetHeader("Accept"); accept != "" {
		components = append(components, accept)
	}

	if acceptLanguage := ctx.GetHeader("Accept-Language"); acceptLanguage != "" {
		components = append(components, acceptLanguage)
	}

	// 如果没有足够的信息，使用UUID
	if len(components) == 0 {
		return uuid.New().String()
	}

	// 组合所有组件并生成哈希
	combined := strings.Join(components, "|")
	hash := md5.Sum([]byte(combined))
	return fmt.Sprintf("%x", hash)
}
