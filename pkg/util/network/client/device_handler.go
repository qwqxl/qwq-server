package client

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// 设备类型常量
const (
	DevicePC          = "PC"
	DevicePhone       = "Phone"
	DeviceTablet      = "Tablet"
	DeviceTV          = "TV"
	DeviceConsole     = "Game Console"
	DeviceWearable    = "Wearable"
	DeviceBot         = "Bot"
	DeviceUnknown     = "Unknown"
	DeviceCar         = "Car"
	DeviceSmartScreen = "Smart Screen"
)

// 操作系统常量
const (
	OSWindows     = "Windows"
	OSMacOS       = "macOS"
	OSLinux       = "Linux"
	OSChromeOS    = "Chrome OS"
	OSAndroid     = "Android"
	OSiOS         = "iOS"
	OSiPadOS      = "iPadOS"
	OSTvOS        = "tvOS"
	OSWatchOS     = "watchOS"
	OSPlayStation = "PlayStation"
	OSXbox        = "Xbox"
	OSNintendo    = "Nintendo"
	OSUnknown     = "Unknown"
)

// 浏览器常量
const (
	BrowserChrome      = "Chrome"
	BrowserSafari      = "Safari"
	BrowserFirefox     = "Firefox"
	BrowserEdge        = "Edge"
	BrowserIE          = "Internet Explorer"
	BrowserOpera       = "Opera"
	BrowserSamsung     = "Samsung Internet"
	BrowserUC          = "UC Browser"
	BrowserBot         = "Bot"
	BrowserUnknown     = "Unknown"
	BrowserPlayStation = "PlayStation Browser"
	BrowserNintendo    = "Nintendo Browser"
)

// DeviceInfo 包含详细的设备信息
type DeviceInfo struct {
	DeviceType     string `json:"device_type"`
	OS             string `json:"os"`
	OSVersion      string `json:"os_version,omitempty"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version,omitempty"`
	IsMobile       bool   `json:"is_mobile"`
	IsTablet       bool   `json:"is_tablet"`
	IsBot          bool   `json:"is_bot"`
	Brand          string `json:"brand,omitempty"`
	Model          string `json:"model,omitempty"`
}

func DeviceHandler(w http.ResponseWriter, r *http.Request) string {
	userAgent := r.UserAgent()
	if userAgent == "" {
		http.Error(w, "无法获取 User-Agent", http.StatusBadRequest)
		return "Not User-Agent"
	}

	info := detectDeviceInfo(userAgent)

	return fmt.Sprintf(`
		<div class="card">
			<h2>设备检测详情</h2>
			<div class="info-grid">
				<div class="info-label">设备类型:</div>
				<div><span class="device-badge %s">%s</span></div>
				
				<div class="info-label">操作系统:</div>
				<div>%s %s</div>
				
				<div class="info-label">浏览器:</div>
				<div>%s %s</div>
				
				<div class="info-label">品牌:</div>
				<div>%s</div>
				
				<div class="info-label">型号:</div>
				<div>%s</div>
				
				<div class="info-label">移动设备:</div>
				<div>%v</div>
				
				<div class="info-label">平板设备:</div>
				<div>%v</div>
				
				<div class="info-label">机器人:</div>
				<div>%v</div>
			</div>
		</div>
		
		<div class="card">
			<h3>原始 User-Agent</h3>
			<div class="raw-ua">%s</div>
		</div>
	`,
		getBadgeClass(info.DeviceType),
		info.DeviceType,
		info.OS,
		info.OSVersion,
		info.Browser,
		info.BrowserVersion,
		info.Brand,
		info.Model,
		info.IsMobile,
		info.IsTablet,
		info.IsBot,
		userAgent,
	)
}

func JsonDeviceHandler(w http.ResponseWriter, r *http.Request) *DeviceInfo {
	userAgent := r.UserAgent()
	if userAgent == "" {
		http.Error(w, `{"error": "无法获取 User-Agent"}`, http.StatusBadRequest)
		return nil
	}

	info := detectDeviceInfo(userAgent)
	return &info
}

func getBadgeClass(deviceType string) string {
	switch deviceType {
	case DevicePC:
		return "pc-badge"
	case DevicePhone:
		return "mobile-badge"
	case DeviceTablet:
		return "tablet-badge"
	case DeviceTV:
		return "tv-badge"
	case DeviceConsole:
		return "console-badge"
	case DeviceBot:
		return "bot-badge"
	default:
		return ""
	}
}

func detectDeviceInfo(userAgent string) DeviceInfo {
	ua := strings.ToLower(userAgent)
	info := DeviceInfo{
		DeviceType: DeviceUnknown,
		OS:         OSUnknown,
		Browser:    BrowserUnknown,
	}

	// 1. 检测机器人
	if isBot(ua) {
		info.DeviceType = DeviceBot
		info.IsBot = true
		info.Browser = BrowserBot
		return info
	}

	// 2. 检测操作系统
	info.OS, info.OSVersion = detectOS(ua)

	// 3. 检测浏览器
	info.Browser, info.BrowserVersion = detectBrowser(ua)

	// 4. 检测设备类型
	info.DeviceType, info.IsMobile, info.IsTablet = detectDeviceType(ua, info.OS)

	// 5. 检测品牌和型号
	info.Brand, info.Model = detectBrandAndModel(ua)

	return info
}

func isBot(ua string) bool {
	botPatterns := []string{
		`googlebot`, `bingbot`, `slurp`, `duckduckbot`, `baiduspider`,
		`yandexbot`, `facebot`, `ia_archiver`, `bot`, `spider`, `crawler`,
		`archive\.org_bot`, `ahrefsbot`, `semrushbot`, `mj12bot`,
	}

	for _, pattern := range botPatterns {
		if matched, _ := regexp.MatchString(pattern, ua); matched {
			return true
		}
	}
	return false
}

func detectOS(ua string) (string, string) {
	// Windows
	if strings.Contains(ua, "windows") {
		version := ""
		if strings.Contains(ua, "windows nt 10.0") {
			version = "10"
		} else if strings.Contains(ua, "windows nt 6.3") {
			version = "8.1"
		} else if strings.Contains(ua, "windows nt 6.2") {
			version = "8"
		} else if strings.Contains(ua, "windows nt 6.1") {
			version = "7"
		} else if strings.Contains(ua, "windows nt 6.0") {
			version = "Vista"
		} else if strings.Contains(ua, "windows nt 5.1") {
			version = "XP"
		}
		return OSWindows, version
	}

	// macOS
	if strings.Contains(ua, "mac os x") {
		version := ""
		if strings.Contains(ua, "mac os x 10_15") || strings.Contains(ua, "mac os x 10.15") {
			version = "Catalina"
		} else if strings.Contains(ua, "mac os x 10_14") || strings.Contains(ua, "mac os x 10.14") {
			version = "Mojave"
		} else if strings.Contains(ua, "mac os x 10_13") || strings.Contains(ua, "mac os x 10.13") {
			version = "High Sierra"
		}
		return OSMacOS, version
	}

	// Linux
	if strings.Contains(ua, "linux") {
		version := ""
		if strings.Contains(ua, "ubuntu") {
			version = "Ubuntu"
		} else if strings.Contains(ua, "fedora") {
			version = "Fedora"
		} else if strings.Contains(ua, "debian") {
			version = "Debian"
		}
		return OSLinux, version
	}

	// Chrome OS
	if strings.Contains(ua, "cros") {
		return OSChromeOS, ""
	}

	// Android
	if strings.Contains(ua, "android") {
		versionRe := regexp.MustCompile(`android (\d+(?:\.\d+)*)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return OSAndroid, matches[1]
		}
		return OSAndroid, ""
	}

	// iOS
	if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipod") {
		versionRe := regexp.MustCompile(`os (\d+(_\d+)*) like`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return OSiOS, strings.Replace(matches[1], "_", ".", -1)
		}
		return OSiOS, ""
	}

	// iPadOS
	if strings.Contains(ua, "ipad") {
		versionRe := regexp.MustCompile(`os (\d+(_\d+)*) like`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return OSiPadOS, strings.Replace(matches[1], "_", ".", -1)
		}
		return OSiPadOS, ""
	}

	// tvOS
	if strings.Contains(ua, "apple tv") {
		return OSTvOS, ""
	}

	// watchOS
	if strings.Contains(ua, "watch") {
		return OSWatchOS, ""
	}

	// PlayStation
	if strings.Contains(ua, "playstation") {
		version := ""
		if strings.Contains(ua, "playstation 5") {
			version = "5"
		} else if strings.Contains(ua, "playstation 4") {
			version = "4"
		}
		return OSPlayStation, version
	}

	// Xbox
	if strings.Contains(ua, "xbox") {
		return OSXbox, ""
	}

	// Nintendo
	if strings.Contains(ua, "nintendo") {
		return OSNintendo, ""
	}

	return OSUnknown, ""
}

func detectBrowser(ua string) (string, string) {
	// Chrome
	if strings.Contains(ua, "chrome") && !strings.Contains(ua, "chromium") && !strings.Contains(ua, "edg") {
		versionRe := regexp.MustCompile(`chrome/(\d+\.\d+\.\d+\.\d+)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserChrome, matches[1]
		}
		return BrowserChrome, ""
	}

	// Safari
	if strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome") && !strings.Contains(ua, "crios") {
		versionRe := regexp.MustCompile(`version/(\d+\.\d+\.\d+)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserSafari, matches[1]
		}
		return BrowserSafari, ""
	}

	// Firefox
	if strings.Contains(ua, "firefox") {
		versionRe := regexp.MustCompile(`firefox/(\d+\.\d+\.?\d*)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserFirefox, matches[1]
		}
		return BrowserFirefox, ""
	}

	// Edge
	if strings.Contains(ua, "edg") {
		versionRe := regexp.MustCompile(`edg/(\d+\.\d+\.\d+\.\d+)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserEdge, matches[1]
		}
		return BrowserEdge, ""
	}

	// Internet Explorer
	if strings.Contains(ua, "msie") || strings.Contains(ua, "trident") {
		version := ""
		if strings.Contains(ua, "msie 10.0") {
			version = "10"
		} else if strings.Contains(ua, "msie 9.0") {
			version = "9"
		}
		return BrowserIE, version
	}

	// Opera
	if strings.Contains(ua, "opera") || strings.Contains(ua, "opr/") {
		versionRe := regexp.MustCompile(`(?:opera|opr)/(\d+\.\d+\.?\d*)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserOpera, matches[1]
		}
		return BrowserOpera, ""
	}

	// Samsung Internet
	if strings.Contains(ua, "samsungbrowser") {
		versionRe := regexp.MustCompile(`samsungbrowser/(\d+\.\d+)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserSamsung, matches[1]
		}
		return BrowserSamsung, ""
	}

	// UC Browser
	if strings.Contains(ua, "ucbrowser") {
		versionRe := regexp.MustCompile(`ucbrowser/(\d+\.\d+\.\d+\.\d+)`)
		matches := versionRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return BrowserUC, matches[1]
		}
		return BrowserUC, ""
	}

	// PlayStation Browser
	if strings.Contains(ua, "playstation") {
		return BrowserPlayStation, ""
	}

	// Nintendo Browser
	if strings.Contains(ua, "nintendo") {
		return BrowserNintendo, ""
	}

	return BrowserUnknown, ""
}

func detectDeviceType(ua, os string) (string, bool, bool) {
	// 平板设备检测
	isTablet := false
	tabletPatterns := []string{
		`ipad`, `tablet`, `playbook`, `kindle fire`, `nexus 7`, `nexus 10`, `tab`, `galaxy tab`,
	}

	for _, pattern := range tabletPatterns {
		if matched, _ := regexp.MatchString(pattern, ua); matched {
			isTablet = true
		}
	}

	// Android平板特殊处理
	if strings.Contains(ua, "android") && !strings.Contains(ua, "mobile") {
		isTablet = true
	}

	// 手机设备检测
	isMobile := false
	mobilePatterns := []string{
		`iphone`, `ipod`, `android.*mobile`, `blackberry`, `windows phone`, `iemobile`, `mobile`,
	}

	for _, pattern := range mobilePatterns {
		if matched, _ := regexp.MatchString(pattern, ua); matched {
			isMobile = true
		}
	}

	// 排除平板设备
	if isTablet {
		isMobile = false
	}

	// 电视设备
	if strings.Contains(ua, "tv") || strings.Contains(ua, "smart-tv") ||
		strings.Contains(ua, "appletv") || strings.Contains(ua, "googletv") ||
		strings.Contains(ua, "crkey") || strings.Contains(ua, "hbbtv") ||
		strings.Contains(ua, "netcast") || strings.Contains(ua, "roku") {
		return DeviceTV, false, false
	}

	// 游戏机
	if strings.Contains(ua, "playstation") || strings.Contains(ua, "xbox") ||
		strings.Contains(ua, "nintendo") || strings.Contains(ua, "wii") {
		return DeviceConsole, false, false
	}

	// 可穿戴设备
	if strings.Contains(ua, "watch") || strings.Contains(ua, "galaxy fit") {
		return DeviceWearable, true, false
	}

	// 汽车系统
	if strings.Contains(ua, "car") || strings.Contains(ua, "automotive") {
		return DeviceCar, false, false
	}

	// 智能屏幕
	if strings.Contains(ua, "smart-screen") || strings.Contains(ua, "smart display") {
		return DeviceSmartScreen, false, false
	}

	// 最终分类
	if isTablet {
		return DeviceTablet, false, true
	}

	if isMobile {
		return DevicePhone, true, false
	}

	return DevicePC, false, false
}

func detectBrandAndModel(ua string) (string, string) {
	// Apple 设备
	if strings.Contains(ua, "iphone") {
		modelRe := regexp.MustCompile(`iphone(\d+,\d+)`)
		matches := modelRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return "Apple", "iPhone " + matches[1]
		}
		return "Apple", "iPhone"
	}

	if strings.Contains(ua, "ipad") {
		modelRe := regexp.MustCompile(`ipad(\d+,\d+)`)
		matches := modelRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return "Apple", "iPad " + matches[1]
		}
		return "Apple", "iPad"
	}

	if strings.Contains(ua, "mac") {
		return "Apple", "Mac"
	}

	// Samsung 设备
	if strings.Contains(ua, "samsung") {
		modelRe := regexp.MustCompile(`sm-(\w+)`)
		matches := modelRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return "Samsung", "Galaxy " + strings.ToUpper(matches[1])
		}

		if strings.Contains(ua, "galaxy") {
			modelRe := regexp.MustCompile(`galaxy ([\w\s]+)`)
			matches := modelRe.FindStringSubmatch(ua)
			if len(matches) > 1 {
				return "Samsung", "Galaxy " + matches[1]
			}
			return "Samsung", "Galaxy"
		}

		return "Samsung", ""
	}

	// Google Pixel 设备
	if strings.Contains(ua, "pixel") {
		modelRe := regexp.MustCompile(`pixel (\d)`)
		matches := modelRe.FindStringSubmatch(ua)
		if len(matches) > 1 {
			return "Google", "Pixel " + matches[1]
		}
		return "Google", "Pixel"
	}

	// Microsoft Surface
	if strings.Contains(ua, "surface") {
		return "Microsoft", "Surface"
	}

	// Sony PlayStation
	if strings.Contains(ua, "playstation") {
		if strings.Contains(ua, "playstation 5") {
			return "Sony", "PlayStation 5"
		}
		if strings.Contains(ua, "playstation 4") {
			return "Sony", "PlayStation 4"
		}
		return "Sony", "PlayStation"
	}

	// Microsoft Xbox
	if strings.Contains(ua, "xbox") {
		return "Microsoft", "Xbox"
	}

	// Nintendo
	if strings.Contains(ua, "nintendo") {
		if strings.Contains(ua, "switch") {
			return "Nintendo", "Switch"
		}
		return "Nintendo", ""
	}

	return "", ""
}
