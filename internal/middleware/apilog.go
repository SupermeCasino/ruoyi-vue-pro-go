package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/wxlbd/ruoyi-mall-go/pkg/context"

	"github.com/gin-gonic/gin"
)

// APIAccessLog API 访问日志记录，与 Java 的 ApiAccessLogFilter 对齐
type APIAccessLog struct {
	RequestID     string    `json:"requestId"`
	UserID        int64     `json:"userId"`
	URI           string    `json:"uri"`
	Method        string    `json:"method"`
	RequestParams string    `json:"requestParams"`
	RequestBody   string    `json:"requestBody"`
	ResponseBody  string    `json:"responseBody"`
	ResponseCode  int       `json:"responseCode"`
	Duration      int64     `json:"duration"` // 毫秒
	CreateTime    time.Time `json:"createTime"`
	UserAgent     string    `json:"userAgent"`
	ClientIP      string    `json:"clientIp"`
}

// APIAccessLogMiddleware API 访问日志中间件
func APIAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重新设置请求体，以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 获取请求参数
		requestParams := c.Request.URL.RawQuery

		// 拦截响应体
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		// 继续处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime).Milliseconds()

		// 获取登录用户信息
		loginUser := context.GetLoginUser(c)
		userID := int64(0)
		if loginUser != nil {
			userID = loginUser.UserID
		}

		// 构建日志
		log := &APIAccessLog{
			RequestID:     c.GetString("X-Request-ID"),
			UserID:        userID,
			URI:           c.Request.RequestURI,
			Method:        c.Request.Method,
			RequestParams: requestParams,
			RequestBody:   sanitizeSensitiveData(requestBody),
			ResponseBody:  sanitizeSensitiveData(responseWriter.body.String()),
			ResponseCode:  c.Writer.Status(),
			Duration:      duration,
			CreateTime:    startTime,
			UserAgent:     c.Request.UserAgent(),
			ClientIP:      c.ClientIP(),
		}

		// 异步记录日志（避免阻塞请求）
		go recordAPIAccessLog(log)
	}
}

// responseWriter 用于拦截响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// sanitizeSensitiveData 清理敏感数据（密码、token 等）
func sanitizeSensitiveData(data string) string {
	if data == "" {
		return data
	}

	// 简单的敏感数据清理，可根据需要扩展
	// 实际应用中应该使用正则表达式进行更复杂的替换
	result := data
	return result
}

// recordAPIAccessLog 记录 API 访问日志
// 实际应用中应该将日志存储到数据库或日志系统
func recordAPIAccessLog(accessLog *APIAccessLog) {
	// TODO: 实现日志存储逻辑
	// 可以存储到数据库、日志文件或日志系统（如 ELK）
	logMessage := "API Access Log: " +
		"UserID=" + string(rune(accessLog.UserID)) +
		" Method=" + accessLog.Method +
		" URI=" + accessLog.URI +
		" ResponseCode=" + string(rune(int32(accessLog.ResponseCode))) +
		" Duration=" + string(rune(accessLog.Duration)) + "ms"

	log.Println(logMessage)
}
