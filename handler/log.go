package handler

import (
	"godemo/database"
	"godemo/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListRequestLogs godoc
// @Summary 获取请求日志列表
// @Description 获取请求日志列表，支持分页和筛选
// @Tags 日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param method query string false "HTTP方法"
// @Param status query int false "状态码"
// @Param ip query string false "客户端IP"
// @Param username query string false "用户名"
// @Param country query string false "国家"
// @Param city query string false "城市"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/logs [get]
func ListRequestLogs(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 构建查询
	query := database.DB.Model(&models.RequestLog{})

	// 筛选条件
	if method := c.Query("method"); method != "" {
		query = query.Where("method = ?", method)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if ip := c.Query("ip"); ip != "" {
		query = query.Where("client_ip = ?", ip)
	}
	if username := c.Query("username"); username != "" {
		query = query.Where("username = ?", username)
	}
	if country := c.Query("country"); country != "" {
		query = query.Where("country = ?", country)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var logs []models.RequestLog
	query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetRequestLog godoc
// @Summary 获取请求日志详情
// @Description 根据ID获取请求日志详细信息
// @Tags 日志
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Security Bearer
// @Success 200 {object} models.RequestLog
// @Router /api/logs/{id} [get]
func GetRequestLog(c *gin.Context) {
	id := c.Param("id")
	var log models.RequestLog

	if err := database.DB.First(&log, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "日志不存在"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetLogStatistics godoc
// @Summary 获取日志统计信息
// @Description 获取请求日志的统计数据
// @Tags 日志
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/logs/statistics [get]
func GetLogStatistics(c *gin.Context) {
	var stats struct {
		TotalRequests  int64
		TotalUsers     int64
		TotalCountries int64
		AvgLatency     float64
	}

	// 总请求数
	database.DB.Model(&models.RequestLog{}).Count(&stats.TotalRequests)

	// 独立用户数
	database.DB.Model(&models.RequestLog{}).
		Where("user_id IS NOT NULL").
		Distinct("user_id").
		Count(&stats.TotalUsers)

	// 独立国家数
	database.DB.Model(&models.RequestLog{}).
		Where("country != ''").
		Distinct("country").
		Count(&stats.TotalCountries)

	// 平均延迟
	database.DB.Model(&models.RequestLog{}).
		Select("AVG(latency)").
		Row().
		Scan(&stats.AvgLatency)

	// 按状态码统计
	var statusStats []struct {
		Status int
		Count  int64
	}
	database.DB.Model(&models.RequestLog{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Order("count DESC").
		Limit(10).
		Scan(&statusStats)

	// 按国家统计
	var countryStats []struct {
		Country string
		Count   int64
	}
	database.DB.Model(&models.RequestLog{}).
		Select("country, COUNT(*) as count").
		Where("country != ''").
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&countryStats)

	// 按路径统计
	var pathStats []struct {
		Path  string
		Count int64
	}
	database.DB.Model(&models.RequestLog{}).
		Select("path, COUNT(*) as count").
		Group("path").
		Order("count DESC").
		Limit(10).
		Scan(&pathStats)

	c.JSON(http.StatusOK, gin.H{
		"overview":             stats,
		"status_distribution":  statusStats,
		"country_distribution": countryStats,
		"top_paths":            pathStats,
	})
}
