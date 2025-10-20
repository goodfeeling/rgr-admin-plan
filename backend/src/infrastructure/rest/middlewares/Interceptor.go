package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	operationRecordsDomain "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	operationRecordsRepository "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware(db *gorm.DB, logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody string
		var resp string
		// skip over upload api
		if checkURIIsUpload(c.Request.RequestURI) {
			reqBody = ""
			resp = ""
		} else {
			blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			buf := make([]byte, 4096)
			num, err := c.Request.Body.Read(buf)
			if err != nil && err.Error() != "EOF" {
				_ = fmt.Errorf("error reading buffer: %s", err.Error())
			}
			reqBody := string(buf[0:num])
			resp = blw.body.String()
			c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(reqBody)))
		}

		start := time.Now()

		c.Next()

		userId, _ := controllers.NewAppUtils(c).GetUserID()

		operationRecordsRepository := operationRecordsRepository.NewOperationRepository(db, logger)
		operationRecordsRepository.Create(&operationRecordsDomain.SysOperationRecord{
			IP:           c.ClientIP(),
			Method:       c.Request.Method,
			Path:         c.Request.RequestURI,
			Status:       int64(c.Writer.Status()),
			Agent:        c.Request.UserAgent(),
			Body:         reqBody,
			Resp:         resp,
			ErrorMessage: c.Errors.String(),
			UserID:       int64(userId),
			Latency:      time.Since(start).Milliseconds(),
			CreatedAt:    domain.CustomTime{Time: time.Now()},
		})

		_ = fmt.Sprintf("%v", map[string]any{
			"ruta":          c.FullPath(),
			"request_uri":   c.Request.RequestURI,
			"raw_request":   reqBody,
			"status_code":   c.Writer.Status(),
			"body_response": resp,
			"errors":        c.Errors.Errors(),
			"created_at":    time.Now().Format("2006-01-02T15:04:05"),
		})
	}
}

// check uri is upload file
func checkURIIsUpload(uri string) bool {
	if strings.HasPrefix(uri, "/v1/upload") {
		return true
	}
	if strings.HasSuffix(uri, "/excel/import") {
		return true
	}
	return false
}
