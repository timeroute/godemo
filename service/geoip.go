package service

import (
	"godemo/config"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var (
	ip2regionSearcher *xdb.Searcher
	ip2regionOnce     sync.Once
)

// InitGeoIP 初始化 IP2Region
func InitGeoIP() error {
	var err error
	ip2regionOnce.Do(func() {
		if config.AppConfig.Logging.GeoIPDBPath == "" {
			log.Println("⚠️  IP2Region 数据库路径未配置，地理位置功能将不可用")
			return
		}

		// 使用 VectorIndex 缓存策略创建查询器
		// 先加载 VectorIndex 缓存
		vIndex, vErr := xdb.LoadVectorIndexFromFile(config.AppConfig.Logging.GeoIPDBPath)
		if vErr != nil {
			log.Printf("⚠️  无法加载 IP2Region VectorIndex: %v (地理位置功能将不可用)", vErr)
			err = vErr
			return
		}

		// 创建查询器（IPv4）
		ip2regionSearcher, err = xdb.NewWithVectorIndex(xdb.IPv4, config.AppConfig.Logging.GeoIPDBPath, vIndex)
		if err != nil {
			log.Printf("⚠️  无法创建 IP2Region 查询器: %v (地理位置功能将不可用)", err)
			return
		}

		log.Println("✅ IP2Region 数据库加载成功")
	})

	return err
}

// CloseGeoIP 关闭 IP2Region
func CloseGeoIP() {
	if ip2regionSearcher != nil {
		ip2regionSearcher.Close()
	}
}

// GetGeoLocation 获取 IP 的地理位置信息
// 返回: 国家, 城市, 纬度, 经度
// 注意: ip2region 不提供经纬度信息，返回 0, 0
func GetGeoLocation(ipStr string) (country, city string, latitude, longitude float64) {
	// 如果 IP2Region 未初始化，返回空值
	if ip2regionSearcher == nil {
		return "", "", 0, 0
	}

	// 解析 IP
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", 0, 0
	}

	// 跳过本地 IP
	if isLocalIP(ip) {
		return "本地", "本地", 0, 0
	}

	// 查询 IP 信息
	// 返回格式: 国家|区域|省份|城市|ISP
	// 例如: 中国|0|北京|北京市|鹏博士
	region, err := ip2regionSearcher.SearchByStr(ipStr)
	if err != nil {
		return "", "", 0, 0
	}

	// 解析返回的地理信息
	parts := strings.Split(region, "|")
	if len(parts) >= 5 {
		country = parts[0]
		province := parts[2]
		city = parts[3]

		// 如果城市为 0，使用省份
		if city == "0" || city == "" {
			city = province
		}

		// 如果省份也为 0，使用国家
		if city == "0" || city == "" {
			city = country
		}
	}

	// ip2region 不提供经纬度信息
	return country, city, 0, 0
}

// isLocalIP 判断是否为本地 IP
func isLocalIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}

	// 检查是否为本地 IPv4 地址
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}

	return false
}
