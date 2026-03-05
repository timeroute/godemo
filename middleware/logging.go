package middleware

import (
	"bytes"
	"encoding/json"
	"godemo/config"
	"godemo/database"
	"godemo/models"
	"godemo/service"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter 包装 gin.ResponseWriter 以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware 日志记录中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体（如果需要记录）
		var requestBody string
		if config.AppConfig.Logging.LogRequestBody && shouldLogBody(c) {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// 重新设置请求体，以便后续处理器可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// 限制请求体长度
			if len(requestBody) > config.AppConfig.Logging.MaxBodySize {
				requestBody = requestBody[:config.AppConfig.Logging.MaxBodySize] + "...(truncated)"
			}
		}

		// 包装 ResponseWriter 以捕获响应体
		var responseBody string
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		if config.AppConfig.Logging.LogResponseBody {
			c.Writer = blw
		}

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(startTime)

		// 获取客户端 IP
		clientIP := getClientIP(c)

		// 获取用户信息（如果已认证）
		var userID *uint
		var username string
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}
		if uname, exists := c.Get("username"); exists {
			if name, ok := uname.(string); ok {
				username = name
			}
		}

		// 获取响应体
		if config.AppConfig.Logging.LogResponseBody && blw.body.Len() > 0 {
			responseBody = blw.body.String()
			// 限制响应体长度
			if len(responseBody) > config.AppConfig.Logging.MaxBodySize {
				responseBody = responseBody[:config.AppConfig.Logging.MaxBodySize] + "...(truncated)"
			}
		}

		// 获取错误信息
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// 获取地理位置信息
		country, city, latitude, longitude := service.GetGeoLocation(clientIP)

		// 创建日志记录
		requestLog := models.RequestLog{
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			Status:       c.Writer.Status(),
			Latency:      latency.Milliseconds(),
			ClientIP:     clientIP,
			UserAgent:    c.Request.UserAgent(),
			ErrorMessage: errorMessage,
			RequestBody:  requestBody,
			ResponseBody: responseBody,
			UserID:       userID,
			Username:     username,
			Country:      country,
			City:         city,
			Latitude:     latitude,
			Longitude:    longitude,
		}

		// 异步保存到数据库
		if config.AppConfig.Logging.SaveToDB {
			go func() {
				if err := database.DB.Create(&requestLog).Error; err != nil {
					log.Printf("Failed to save request log: %v", err)
				}
			}()
		}

		// 控制台日志输出
		if config.AppConfig.Logging.LogToConsole {
			logToConsole(requestLog, latency)
		}
	}
}

// getClientIP 获取真实客户端 IP
func getClientIP(c *gin.Context) string {
	// 优先从 X-Forwarded-For 获取
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 从 X-Real-IP 获取
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// 从 RemoteAddr 获取
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}

// shouldLogBody 判断是否应该记录请求体
func shouldLogBody(c *gin.Context) bool {
	// 只记录 POST, PUT, PATCH 请求的请求体
	method := c.Request.Method
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}

	// 不记录文件上传
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return false
	}

	return true
}

// logToConsole 输出到控制台
func logToConsole(log models.RequestLog, latency time.Duration) {
	// 根据状态码选择颜色
	var statusColor string
	switch {
	case log.Status >= 200 && log.Status < 300:
		statusColor = "\033[32m" // 绿色
	case log.Status >= 300 && log.Status < 400:
		statusColor = "\033[36m" // 青色
	case log.Status >= 400 && log.Status < 500:
		statusColor = "\033[33m" // 黄色
	default:
		statusColor = "\033[31m" // 红色
	}
	resetColor := "\033[0m"

	// 格式化输出
	logMsg := map[string]interface{}{
		"time":     time.Now().Format("2006-01-02 15:04:05"),
		"method":   log.Method,
		"path":     log.Path,
		"status":   log.Status,
		"latency":  latency.String(),
		"ip":       log.ClientIP,
		"location": log.Country + "/" + log.City,
	}

	if log.Username != "" {
		logMsg["user"] = log.Username
	}

	if log.ErrorMessage != "" {
		logMsg["error"] = log.ErrorMessage
	}

	jsonLog, _ := json.Marshal(logMsg)
	println(statusColor + string(jsonLog) + resetColor)
}
