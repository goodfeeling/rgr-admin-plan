package upload

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	domain "github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainFiles "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	shareUtils "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	cacheService "github.com/gbrayhan/microservices-go/src/infrastructure/lib/cache"
	"github.com/redis/go-redis/v9"
)

type IUploadController interface {
	Single(ctx *gin.Context)
	Multiple(ctx *gin.Context)
	GetSTSToken(ctx *gin.Context)
	RefreshSTSToken(ctx *gin.Context)
}

type STSTokenResponse struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SecurityToken   string `json:"security_token"`
	Expiration      string `json:"expiration"`
	BucketName      string `json:"bucket_name"`
	Region          string `json:"region"`
	RefreshToken    string `json:"refresh_token,omitempty"`
}

type UploadController struct {
	sysFilesUseCase domainFiles.ISysFilesService
	Logger          *logger.Logger
	RedisClient     *redis.Client
	CacheService    *cacheService.STSCacheService
}

// MultipleUpload
// @Summary multiple files upload
// @Description upload multiple files get files info
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "fileResources" collectionFormat(multi)
// @Success 200 {object} domain.CommonResponse[[]domainFiles.SysFiles]
// @Router /v1/upload/multiple [post]
func (u *UploadController) Multiple(ctx *gin.Context) {
	// 获取多文件表单
	form, err := ctx.MultipartForm()
	if err != nil {
		u.Logger.Error("Failed to get multipart form", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	files := form.File["file"]
	var uploadedFiles []domainFiles.SysFiles

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		ext := filepath.Ext(filename)
		newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(os.Getenv("NATIVE_STORAGE_UPLOAD_DIR"), newFilename)

		if err := os.MkdirAll(os.Getenv("NATIVE_STORAGE_UPLOAD_DIR"), os.ModePerm); err != nil {
			u.Logger.Error("Error creating dir", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			u.Logger.Error("Error save file", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		md5Value, err := shareUtils.CalculateFileMD5(savePath)
		if err != nil {
			u.Logger.Error("calculate file to md5", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		fileInfo := domainFiles.SysFiles{
			FileName:       newFilename,
			FilePath:       savePath,
			FileMD5:        md5Value,
			FileOriginName: filename,
			StorageEngine:  "local",
		}

		res, err := u.sysFilesUseCase.Create(&fileInfo)
		if err != nil {
			u.Logger.Error("insert file info to database", zap.Error(err))
			appError := domainErrors.NewAppError(err, domainErrors.UploadError)
			_ = ctx.Error(appError)
			return
		}

		uploadedFiles = append(uploadedFiles, *res)
	}

	response := &domain.CommonResponse[[]domainFiles.SysFiles]{
		Data:    uploadedFiles,
		Message: "Upload success",
		Status:  200,
	}

	u.Logger.Info("multiple upload successful", zap.Int("fileCount", len(files)))

	ctx.JSON(http.StatusOK, response)
}

// SingleUpload
// @Summary single file upload
// @Description upload single file get file info
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "fileResource"
// @Success 200 {object} domain.CommonResponse[domainFiles.SysFiles]
// @Router /v1/upload/single [post]
func (u *UploadController) Single(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		u.Logger.Error("Failed to get file", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	filename := filepath.Base(file.Filename)
	// only name
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// join name
	savePath := filepath.Join(os.Getenv("NATIVE_STORAGE_UPLOAD_DIR"), newFilename)

	// create save file dir
	if err := os.MkdirAll(os.Getenv("NATIVE_STORAGE_UPLOAD_DIR"), os.ModePerm); err != nil {
		u.Logger.Error("Error creating dir", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		u.Logger.Error("Error save file", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// calculate md5 file
	md5Value, err := shareUtils.CalculateFileMD5(savePath)
	if err != nil {
		u.Logger.Error("calculate  file to md5", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// insert to database
	files := domainFiles.SysFiles{
		FileName:       newFilename,
		FilePath:       savePath,
		FileMD5:        md5Value,
		FileOriginName: filename,
		StorageEngine:  "local",
	}
	res, err := u.sysFilesUseCase.Create(&files)
	if err != nil {
		u.Logger.Error("insert file info to database", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}
	response := &domain.CommonResponse[domainFiles.SysFiles]{
		Data:    *res,
		Message: "Upload success",
		Status:  200,
	}

	u.Logger.Info("upload successful", zap.String("filename", newFilename))

	ctx.JSON(http.StatusOK, response)
}

// GetSTSToken
// @Summary get sts token with aliyun
// @Description get sts token
// @Tags sts token
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/upload/sts-token [get]
func (u *UploadController) GetSTSToken(ctx *gin.Context) {
	userID := ctx.GetString("user_id") // 从JWT中获取用户ID
	if userID == "" {
		userID = "anonymous" // 匿名用户
	}

	cacheKey := fmt.Sprintf("sts_token:%s", userID)
	ctxBg := context.Background()

	// 尝试从缓存获取
	if cachedToken, err := u.CacheService.GetSTSToken(ctxBg, cacheKey); err == nil {
		// 检查token是否即将过期（提前5分钟刷新）
		if time.Until(cachedToken.Expiration) > 5*time.Minute {
			result := controllers.NewCommonResponseBuilder[*STSTokenResponse]().
				Data(&STSTokenResponse{
					AccessKeyId:     cachedToken.AccessKeyId,
					AccessKeySecret: cachedToken.AccessKeySecret,
					SecurityToken:   cachedToken.SecurityToken,
					Expiration:      cachedToken.Expiration.Format(time.RFC3339),
					BucketName:      cachedToken.BucketName,
					Region:          cachedToken.Region,
				}).
				Message("success").
				Status(0).
				Build()
			ctx.JSON(http.StatusOK, result)
			return
		}
	}

	// 如果缓存中没有有效token，则生成新的token
	u.processSTSToken(ctx, userID, cacheKey)
}

// 添加刷新Token的接口
// @Summary refresh token
// @Description refresh sts token
// @Tags sts token
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/upload/refresh-sts [get]
func (u *UploadController) RefreshSTSToken(ctx *gin.Context) {
	refreshToken := ctx.Query("refresh_token")
	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	ctxBg := context.Background()

	// 验证Refresh Token
	rt, err := u.CacheService.ValidateRefreshToken(ctxBg, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// 删除已使用的refresh token（一次性使用）
	u.CacheService.DeleteRefreshToken(ctxBg, refreshToken)

	// 生成新的STS Token
	userID := rt.UserID
	cacheKey := fmt.Sprintf("sts_token:%s", userID)

	u.processSTSToken(ctx, userID, cacheKey)
}
func NewAuthController(sysFilesUseCase domainFiles.ISysFilesService, loggerInstance *logger.Logger, redisClient *redis.Client) IUploadController {
	return &UploadController{
		sysFilesUseCase: sysFilesUseCase,
		Logger:          loggerInstance,
		RedisClient:     redisClient,
		CacheService:    cacheService.NewSTSCacheService(redisClient),
	}
}

// generateSTSToken 生成新的STS Token
func (u *UploadController) generateSTSToken(sessionName string, roleArn string) (*sts20150401.AssumeRoleResponseBodyCredentials, error) {
	// 从环境变量中获取步骤1.1生成的RAM用户的访问密钥（AccessKey ID和AccessKey Secret）。
	accessKeyId := os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")
	serviceAddress := os.Getenv("ALIYUN_OSS_SECURITY_SERVICE_ADDRESS")

	// 创建权限策略客户端。
	config := &openapi.Config{
		// 必填，步骤1.1获取到的 AccessKey ID。
		AccessKeyId: tea.String(accessKeyId),
		// 必填，步骤1.1获取到的 AccessKey Secret。
		AccessKeySecret: tea.String(accessKeySecret),
	}

	// Endpoint 请参考 https://api.aliyun.com/product/Sts
	config.Endpoint = tea.String(serviceAddress)
	client, err := sts20150401.NewClient(config)
	if err != nil {
		return nil, err
	}

	// 使用RAM用户的AccessKey ID和AccessKey Secret向STS申请临时访问凭证。
	request := &sts20150401.AssumeRoleRequest{
		// 指定STS临时访问凭证过期时间为3600秒。
		DurationSeconds: tea.Int64(3600),
		// 从环境变量中获取步骤1.3生成的RAM角色的RamRoleArn。
		RoleArn: tea.String(roleArn),
		// 指定自定义角色会话名称
		RoleSessionName: tea.String(sessionName),
	}

	response, err := client.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		return nil, err
	}

	return response.Body.Credentials, nil
}

// processSTSToken 处理STS Token的完整流程
func (u *UploadController) processSTSToken(ctx *gin.Context, userID string, cacheKey string) {
	roleArn := os.Getenv("ALIYUN_OSS_RAM_ROLE_ARN")
	// 生成唯一的会话名称
	sessionName := fmt.Sprintf("upload-session-%d", time.Now().Unix())
	//  *gin.Context, userID生成STS Token
	credentials, err := u.generateSTSToken(sessionName, roleArn)
	ctxBg := context.Background()

	// 从环境变量中获取步骤1.3生成的RAM角色的RamRole.generateSTSToken(sessionName, roleArn)
	if err != nil {
		u.Logger.Error("Failed to generate STS token:", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.UploadError)
		_ = ctx.Error(appError)
		return
	}

	// 生成Refresh Token
	refreshToken, err := u.CacheService.GenerateRefreshToken(ctxBg, userID)
	if err != nil {
		u.Logger.Error("Failed to generate refresh token:", zap.Error(err))
		// 即使生成refresh token失败，也继续返回STS token
	}

	// 缓存STS Token
	expirationTime, _ := time.Parse(time.RFC3339, *credentials.Expiration)
	tokenCache := &domainFiles.STSTokenCache{
		AccessKeyId:     *credentials.AccessKeyId,
		AccessKeySecret: *credentials.AccessKeySecret,
		SecurityToken:   *credentials.SecurityToken,
		Expiration:      expirationTime,
		BucketName:      os.Getenv("ALIYUN_OSS_BUCKET_NAME"),
		Region:          os.Getenv("ALIYUN_OSS_SECURITY_REGION_ID"),
		CreatedAt:       time.Now(),
	}

	// 缓存token，设置为过期时间减去1分钟，确保在过期前刷新
	cacheDuration := time.Until(expirationTime) - time.Minute
	if cacheDuration > 0 {
		u.CacheService.SetSTSToken(ctxBg, cacheKey, tokenCache, cacheDuration)
	}

	result := controllers.NewCommonResponseBuilder[*STSTokenResponse]().
		Data(&STSTokenResponse{
			AccessKeyId:     *credentials.AccessKeyId,
			AccessKeySecret: *credentials.AccessKeySecret,
			SecurityToken:   *credentials.SecurityToken,
			Expiration:      *credentials.Expiration,
			BucketName:      os.Getenv("ALIYUN_OSS_BUCKET_NAME"),
			Region:          os.Getenv("ALIYUN_OSS_SECURITY_REGION_ID"),
			RefreshToken:    refreshToken,
		}).
		Message("success").
		Status(0).
		Build()

	ctx.JSON(http.StatusOK, result)
}
