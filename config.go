package cors

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Config represents all available options for the middleware.
type Config struct {
	AllowAllOrigins bool

	// AllowOrigins是一个跨域请求可以执行的源列表。
	// 如果列表中存在特殊的“*”值，则允许所有起源。默认值为[]
	AllowOrigins []string

	// AllowOriginFunc是一个用于验证origin的自定义函数。
	// 它将origin作为参数，如果允许则返回true，否则返回false。
	// 如果设置了该选项，AllowOrigins的内容将被忽略。
	AllowOriginFunc func(origin string) bool

	// AllowMethods是允许客户端跨域请求使用的方法列表。
	// 默认值为简单方法(GET、POST、PUT、PATCH、DELETE、HEAD和OPTIONS)
	AllowMethods []string

	// AllowHeaders是允许客户端跨域请求使用的非简单头的列表。
	AllowHeaders []string

	// AllowCredentials指示请求是否可以包括用户凭据，
	// 如cookie、HTTP身份验证或客户端SSL证书。
	AllowCredentials bool

	// ExposeHeaders指示向CORS API规范的API公开哪些头是安全的
	ExposeHeaders []string

	// MaxAge表示preflight请求的结果可以缓存多长时间(秒)
	MaxAge time.Duration

	// 允许添加像
	// http://some-domain/*，
	// https://api.*或http://some.*.subdomain.com的origin
	AllowWildcard bool

	// 允许使用流行浏览器的扩展模式
	AllowBrowserExtensions bool

	// 允许使用websocket协议
	AllowWebSockets bool

	// // Allows usage of file:// schema (dangerous!) use it only when you 100% sure it's needed
	AllowFiles bool
}

// AddAllowMethods 允许添加自定义方法
func (c *Config) AddAllowMethods(methods ...string) {
	c.AllowMethods = append(c.AllowMethods, methods...)
}

// AddAllowHeaders 允许添加自定义 header 列表
func (c *Config) AddAllowHeaders(headers ...string) {
	c.AllowHeaders = append(c.AllowHeaders, headers...)
}

// AddExposeHeaders 允许添加自定义公开头
func (c *Config) AddExposeHeaders(headers ...string) {
	c.ExposeHeaders = append(c.ExposeHeaders, headers...)
}

func (c Config) getAllowedSchemas() []string {
	allowedSchemas := DefaultSchemas
	if c.AllowBrowserExtensions {
		allowedSchemas = append(allowedSchemas, ExtensionSchemas...)
	}
	if c.AllowWebSockets {
		allowedSchemas = append(allowedSchemas, WebSocketSchemas...)
	}
	if c.AllowFiles {
		allowedSchemas = append(allowedSchemas, FileSchemas...)
	}
	return allowedSchemas
}

func (c Config) validateAllowedSchemas(origin string) bool {
	allowedSchemas := c.getAllowedSchemas()
	for _, schema := range allowedSchemas {
		if strings.HasPrefix(origin, schema) {
			return true
		}
	}
	return false
}

// Validate 检查用户定义的配置
func (c Config) Validate() error {
	if c.AllowAllOrigins && (c.AllowOriginFunc != nil || len(c.AllowOrigins) > 0) {
		return errors.New("conflict settings: all origins are allowed. AllowOriginFunc or AllowOrigins is not needed")
	}
	if !c.AllowAllOrigins && c.AllowOriginFunc == nil && len(c.AllowOrigins) == 0 {
		return errors.New("conflict settings: all origins disabled")
	}
	for _, origin := range c.AllowOrigins {
		if !strings.Contains(origin, "*") && !c.validateAllowedSchemas(origin) {
			return errors.New("bad origin: origins must contain '*' or include " + strings.Join(c.getAllowedSchemas(), ","))
		}
	}
	return nil
}

// parseWildcardRules 解析包含通配符的url
func (c Config) parseWildcardRules() [][]string {
	var wRules [][]string

	if !c.AllowWildcard {
		return wRules
	}

	for _, o := range c.AllowOrigins {
		if !strings.Contains(o, "*") {
			continue
		}

		if c := strings.Count(o, "*"); c > 1 {
			panic(errors.New("only one * is allowed").Error())
		}

		i := strings.Index(o, "*")
		if i == 0 {
			wRules = append(wRules, []string{"*", o[1:]})
			continue
		}
		if i == (len(o) - 1) {
			wRules = append(wRules, []string{o[:i-1], "*"})
			continue
		}
		wRules = append(wRules, []string{o[:i], o[i+1:]})
	}

	return wRules
}

func DefaultConfig() Config {
	return Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
}

func Default() gin.HandlerFunc {
	config := DefaultConfig()
	config.AllowAllOrigins = true
	return New(config)
}

func New(config Config) gin.HandlerFunc {
	cors := newCors(config)
	return func(c *gin.Context) {
		cors.applyCors(c)
	}
}
