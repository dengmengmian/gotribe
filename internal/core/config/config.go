// Package config defines application configuration structures, defaults, and loading logic.
package config

// 本文件定义配置结构、默认值和配置加载逻辑。

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 汇总应用运行所需的全部配置。
type Config struct {
	App           AppConfig           `mapstructure:"app"`
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Upload        UploadConfig        `mapstructure:"upload"`
	Auth          AuthConfig          `mapstructure:"auth"`
	RateLimit     RateLimitConfig     `mapstructure:"rate_limit"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	AI            AIConfig            `mapstructure:"ai"`
	Admin         AdminConfig         `mapstructure:"admin"`
}

// AppConfig 定义应用标识、运行环境和默认项目等基础配置。
type AppConfig struct {
	Name             string `mapstructure:"name"`
	Env              string `mapstructure:"env"`
	DefaultProjectID string `mapstructure:"default_project_id"`
}

// ServerConfig 定义 HTTP 服务监听与超时配置。
type ServerConfig struct {
	Port                  int        `mapstructure:"port"`
	ReadTimeoutSeconds    int        `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds   int        `mapstructure:"write_timeout_seconds"`
	IdleTimeoutSeconds    int        `mapstructure:"idle_timeout_seconds"`
	ShutdownTimeoutSecond int        `mapstructure:"shutdown_timeout_seconds"`
	MaxHeaderBytes        int        `mapstructure:"max_header_bytes"`
	MaxRequestBodyBytes   int64      `mapstructure:"max_request_body_bytes"`
	CORS                  CORSConfig `mapstructure:"cors"`
}

// DatabaseConfig 定义数据库连接和连接池配置。
type DatabaseConfig struct {
	Type                string `mapstructure:"type"`
	Host                string `mapstructure:"host"`
	Port                int    `mapstructure:"port"`
	Username            string `mapstructure:"username"`
	Password            string `mapstructure:"password"`
	Database            string `mapstructure:"database"`
	Query               string `mapstructure:"query"`
	LogMode             bool   `mapstructure:"log_mode"`
	SSLMode             string `mapstructure:"sslmode"`
	Charset             string `mapstructure:"charset"`
	Collation           string `mapstructure:"collation"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns"`
	ConnMaxIdleTimeMins int    `mapstructure:"conn_max_idle_time_mins"`
	ConnMaxLifetimeMins int    `mapstructure:"conn_max_lifetime_mins"`
}

// RedisConfig 定义 Redis 连接配置。
type RedisConfig struct {
	Addr                string `mapstructure:"addr"`
	Password            string `mapstructure:"password"`
	DB                  int    `mapstructure:"db"`
	PoolSize            int    `mapstructure:"pool_size"`
	DefaultCacheTTLMins int    `mapstructure:"default_cache_ttl_mins"`
}

// AuthConfig 定义统一认证配置：共享 issuer / secret，按 audience 分配独立 TTL。
type AuthConfig struct {
	Issuer string             `mapstructure:"issuer"`
	Secret string             `mapstructure:"secret"`
	User   AuthAudienceConfig `mapstructure:"user"`
	Admin  AuthAudienceConfig `mapstructure:"admin"`
}

// AuthAudienceConfig 定义单个 audience（user / admin）的 token 时长与 audience 字符串。
type AuthAudienceConfig struct {
	Audience              string `mapstructure:"audience"`
	AccessTokenTTLMinutes int    `mapstructure:"access_token_ttl_minutes"`
	RefreshTokenTTLHours  int    `mapstructure:"refresh_token_ttl_hours"`
}

// AccessTokenTTL 返回访问令牌有效期。
func (c AuthAudienceConfig) AccessTokenTTL() time.Duration {
	return time.Duration(c.AccessTokenTTLMinutes) * time.Minute
}

// RefreshTokenTTL 返回刷新令牌有效期。
func (c AuthAudienceConfig) RefreshTokenTTL() time.Duration {
	return time.Duration(c.RefreshTokenTTLHours) * time.Hour
}

// RateLimitConfig 定义不同接口场景的限流阈值配置。
type RateLimitConfig struct {
	Enabled       bool  `mapstructure:"enabled"`
	AuthPerMinute int64 `mapstructure:"auth_per_minute"`
	APIForMinute  int64 `mapstructure:"api_per_minute"`
	EventPerMin   int64 `mapstructure:"event_per_minute"`
}

// ObservabilityConfig 定义 metrics 与 tracing 的基础配置。
type ObservabilityConfig struct {
	MetricsEnabled bool `mapstructure:"metrics_enabled"`
	TracingEnabled bool `mapstructure:"tracing_enabled"`
}

// AIConfig 定义后台 AI 生成能力配置，兼容 DeepSeek/OpenAI 风格 chat completions。
type AIConfig struct {
	Provider       string `mapstructure:"provider"`
	BaseURL        string `mapstructure:"base_url"`
	APIKey         string `mapstructure:"api_key"`
	Model          string `mapstructure:"model"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

// AdminConfig 汇总 Admin 端专属配置。
type AdminConfig struct {
	Mode          string               `mapstructure:"mode"`
	Host          string               `mapstructure:"host"`
	Port          int                  `mapstructure:"port"`
	UrlPathPrefix string               `mapstructure:"url_path_prefix"`
	InitData      bool                 `mapstructure:"init_data"`
	SyncSeedData  bool                 `mapstructure:"sync_seed_data"`
	CDNDomain     string               `mapstructure:"cdn_domain"`
	EnableMigrate bool                 `mapstructure:"enable_migrate"`
	EnableOss     bool                 `mapstructure:"enable_oss"`
	Logs          AdminLogsConfig      `mapstructure:"logs"`
	Casbin        AdminCasbinConfig    `mapstructure:"casbin"`
	JWT           AdminJWTConfig       `mapstructure:"jwt"`
	RateLimit     AdminRateLimitConfig `mapstructure:"rate_limit"`
	Upload        UploadConfig         `mapstructure:"upload"`
	BaiduPush     BaiduPushConfig      `mapstructure:"baidu_push"`
	Jobs          AdminJobsConfig      `mapstructure:"jobs"`
	Lockout       AdminLockoutConfig   `mapstructure:"lockout"`
	TOTP          AdminTOTPConfig      `mapstructure:"totp"`
}

// AdminLockoutConfig 定义登录失败锁定策略。
type AdminLockoutConfig struct {
	Enabled            bool `mapstructure:"enabled"`
	AccountMaxFails    int  `mapstructure:"account_max_fails"`
	AccountLockMinutes int  `mapstructure:"account_lock_minutes"`
	IPMaxFails         int  `mapstructure:"ip_max_fails"`
	IPLockMinutes      int  `mapstructure:"ip_lock_minutes"`
}

// AccountLockDuration 返回账户级锁定时长。
func (c AdminLockoutConfig) AccountLockDuration() time.Duration {
	return time.Duration(c.AccountLockMinutes) * time.Minute
}

// IPLockDuration 返回 IP 级锁定时长。
func (c AdminLockoutConfig) IPLockDuration() time.Duration {
	return time.Duration(c.IPLockMinutes) * time.Minute
}

// AdminTOTPConfig 定义 TOTP 二次校验相关配置。
// 注意：SecretEncryptionKey 必须为 32 字节字符串（AES-256-GCM），生产环境通过环境变量注入。
type AdminTOTPConfig struct {
	Issuer              string `mapstructure:"issuer"`
	SecretEncryptionKey string `mapstructure:"secret_encryption_key"`
	RecoveryCodesCount  int    `mapstructure:"recovery_codes_count"`
	StepTokenTTLSeconds int    `mapstructure:"step_token_ttl_seconds"`
	Period              int    `mapstructure:"period"`
	Digits              int    `mapstructure:"digits"`
	Skew                int    `mapstructure:"skew"`
	// Required 为 true 时，未绑 TOTP 的 admin 登录后必须强制完成绑定才能拿到 access_token；
	// 为 false（默认）时只在登录成功响应里带 mfa_reminder=true，由用户自行决定是否绑定。
	Required bool `mapstructure:"required"`
}

// StepTokenTTL 返回 step_token 时长。
func (c AdminTOTPConfig) StepTokenTTL() time.Duration {
	return time.Duration(c.StepTokenTTLSeconds) * time.Second
}

// AdminLogsConfig 定义 Admin 端日志配置。
type AdminLogsConfig struct {
	Level      int    `mapstructure:"level"`
	Path       string `mapstructure:"path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// AdminCasbinConfig 定义 Casbin RBAC 配置。
type AdminCasbinConfig struct {
	ModelPath string `mapstructure:"model_path"`
}

// AdminJWTConfig 定义 Admin 端 JWT 配置。
type AdminJWTConfig struct {
	Realm       string `mapstructure:"realm"`
	Key         string `mapstructure:"key"`
	Timeout     int    `mapstructure:"timeout"`
	MaxRefresh  int    `mapstructure:"max_refresh"`
	TokenLookup string `mapstructure:"token_lookup"`
}

// AdminRateLimitConfig 定义 Admin 端令牌桶限流配置。
type AdminRateLimitConfig struct {
	FillInterval int64 `mapstructure:"fill_interval"`
	Capacity     int64 `mapstructure:"capacity"`
}

// UploadConfig 定义公用文件上传配置，ToC API 和 Admin 都可以复用。
type UploadConfig struct {
	Provider  string `mapstructure:"provider"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	CDNDomain string `mapstructure:"cdn_domain"`
}

// BaiduPushConfig 定义百度 SEO 推送配置。
type BaiduPushConfig struct {
	Token string `mapstructure:"token"`
}

// AdminJobsConfig 定义 Admin 端定时任务配置。
type AdminJobsConfig struct {
	Enabled bool                 `mapstructure:"enabled"`
	List    map[string]JobConfig `mapstructure:"list"`
}

// JobConfig 定义单个定时任务配置。
type JobConfig struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Schedule    string `mapstructure:"schedule"`
	Enabled     bool   `mapstructure:"enabled"`
	Timeout     string `mapstructure:"timeout"`
	RetryCount  int    `mapstructure:"retry_count"`
}

// CORSConfig 定义跨域资源共享的允许来源和头信息。
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
}

// Load 读取并校验应用配置。
func Load() (Config, error) {
	v := viper.New()
	setDefaults(v)

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.SetEnvPrefix("GOTRIBE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// setDefaults 为配置项设置默认值。
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "gotribe-api")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.default_project_id", "")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout_seconds", 10)
	v.SetDefault("server.write_timeout_seconds", 15)
	v.SetDefault("server.idle_timeout_seconds", 60)
	v.SetDefault("server.shutdown_timeout_seconds", 10)
	v.SetDefault("server.max_header_bytes", 1048576)
	v.SetDefault("server.max_request_body_bytes", 10485760)
	v.SetDefault("server.cors.allowed_origins", []string{"*"})
	v.SetDefault("server.cors.allowed_headers", []string{"Authorization", "Content-Type", "X-Project-ID", "X-Request-ID"})
	v.SetDefault("server.cors.allowed_methods", []string{"GET", "POST", "PATCH", "OPTIONS"})

	v.SetDefault("database.type", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "develop")
	v.SetDefault("database.password", "")
	v.SetDefault("database.database", "develop")
	v.SetDefault("database.query", "connect_timeout=10 TimeZone=Asia/Shanghai")
	v.SetDefault("database.log_mode", true)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 30)
	v.SetDefault("database.conn_max_idle_time_mins", 10)
	v.SetDefault("database.conn_max_lifetime_mins", 30)

	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.default_cache_ttl_mins", 5)

	v.SetDefault("upload.provider", "qiniu")
	v.SetDefault("upload.cdn_domain", "")

	v.SetDefault("auth.issuer", "gotribe")
	v.SetDefault("auth.secret", "")
	v.SetDefault("auth.user.audience", "gotribe.user")
	v.SetDefault("auth.user.access_token_ttl_minutes", 120)
	v.SetDefault("auth.user.refresh_token_ttl_hours", 168)
	v.SetDefault("auth.admin.audience", "gotribe.admin")
	v.SetDefault("auth.admin.access_token_ttl_minutes", 60)
	v.SetDefault("auth.admin.refresh_token_ttl_hours", 24)

	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.auth_per_minute", 10)
	v.SetDefault("rate_limit.api_per_minute", 120)
	v.SetDefault("rate_limit.event_per_minute", 60)

	v.SetDefault("observability.metrics_enabled", true)
	v.SetDefault("observability.tracing_enabled", true)

	v.SetDefault("ai.provider", "deepseek")
	v.SetDefault("ai.base_url", "https://api.deepseek.com")
	v.SetDefault("ai.model", "deepseek-v4-flash")
	v.SetDefault("ai.timeout_seconds", 30)

	v.SetDefault("admin.sync_seed_data", false)

	v.SetDefault("admin.lockout.enabled", true)
	v.SetDefault("admin.lockout.account_max_fails", 5)
	v.SetDefault("admin.lockout.account_lock_minutes", 15)
	v.SetDefault("admin.lockout.ip_max_fails", 20)
	v.SetDefault("admin.lockout.ip_lock_minutes", 60)

	v.SetDefault("admin.totp.issuer", "GoTribe 管理后台")
	v.SetDefault("admin.totp.recovery_codes_count", 10)
	v.SetDefault("admin.totp.step_token_ttl_seconds", 300)
	v.SetDefault("admin.totp.period", 30)
	v.SetDefault("admin.totp.digits", 6)
	v.SetDefault("admin.totp.skew", 1)
	v.SetDefault("admin.totp.required", false)
}

// Address 返回 HTTP 服务监听地址。
func (c ServerConfig) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

// IsDevelopment 判断当前环境是否属于开发环境。
func (c AppConfig) IsDevelopment() bool {
	switch normalizeEnv(c.Env) {
	case "", "development", "dev", "local":
		return true
	default:
		return false
	}
}

// IsProduction 判断当前环境是否属于生产环境。
func (c AppConfig) IsProduction() bool {
	switch normalizeEnv(c.Env) {
	case "production", "prod":
		return true
	default:
		return false
	}
}

// ReadTimeout 返回服务读取请求的超时时间。
func (c ServerConfig) ReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeoutSeconds) * time.Second
}

// WriteTimeout 返回服务写入响应的超时时间。
func (c ServerConfig) WriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeoutSeconds) * time.Second
}

// IdleTimeout 返回长连接空闲超时时间。
func (c ServerConfig) IdleTimeout() time.Duration {
	return time.Duration(c.IdleTimeoutSeconds) * time.Second
}

// ShutdownTimeout 返回优雅停机的等待时间。
func (c ServerConfig) ShutdownTimeout() time.Duration {
	return time.Duration(c.ShutdownTimeoutSecond) * time.Second
}

// DSN 组装数据库连接所需的 DSN 字符串。
func (c DatabaseConfig) DSN() string {
	parts := []string{
		fmt.Sprintf("host=%s", c.Host),
		fmt.Sprintf("port=%d", c.Port),
		fmt.Sprintf("user=%s", c.Username),
		fmt.Sprintf("password=%s", c.Password),
		fmt.Sprintf("dbname=%s", c.Database),
		fmt.Sprintf("sslmode=%s", c.SSLMode),
	}
	if strings.TrimSpace(c.Query) != "" {
		parts = append(parts, strings.TrimSpace(c.Query))
	}
	return strings.Join(parts, " ")
}

// ConnMaxLifetime 返回数据库连接的最大复用时长。
func (c DatabaseConfig) ConnMaxLifetime() time.Duration {
	return time.Duration(c.ConnMaxLifetimeMins) * time.Minute
}

// ConnMaxIdleTime 返回数据库空闲连接最大复用时长。
func (c DatabaseConfig) ConnMaxIdleTime() time.Duration {
	return time.Duration(c.ConnMaxIdleTimeMins) * time.Minute
}

// validate 对配置内容进行启动前校验。
func validate(cfg Config) error {
	secret := strings.TrimSpace(cfg.Auth.Secret)
	if secret == "" {
		return fmt.Errorf("auth.secret is required")
	}
	if len(secret) < 32 {
		return fmt.Errorf("auth.secret must be at least 32 characters")
	}
	weakSecrets := map[string]struct{}{
		"change-me-in-production":                  {},
		"local-dev-secret":                         {},
		"replace-with-your-own-long-random-secret": {},
	}
	if _, found := weakSecrets[secret]; found {
		return fmt.Errorf("auth.secret uses a placeholder value")
	}
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return fmt.Errorf("database.port must be between 1 and 65535")
	}
	if cfg.Database.Username == "" {
		return fmt.Errorf("database.username is required")
	}
	if cfg.Database.Database == "" {
		return fmt.Errorf("database.database is required")
	}
	if cfg.Database.MaxIdleConns <= 0 || cfg.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("database connection pool values must be positive")
	}
	if cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return fmt.Errorf("database.max_idle_conns must not exceed max_open_conns")
	}
	if cfg.Database.ConnMaxIdleTimeMins <= 0 {
		return fmt.Errorf("database.conn_max_idle_time_mins must be positive")
	}
	if cfg.Database.ConnMaxLifetimeMins <= 0 {
		return fmt.Errorf("database.conn_max_lifetime_mins must be positive")
	}
	if strings.TrimSpace(cfg.Redis.Addr) == "" {
		return fmt.Errorf("redis.addr is required")
	}
	if _, _, err := net.SplitHostPort(strings.TrimSpace(cfg.Redis.Addr)); err != nil {
		return fmt.Errorf("redis.addr must use host:port format")
	}
	if cfg.Redis.DB < 0 {
		return fmt.Errorf("redis.db must be non-negative")
	}
	if cfg.Redis.PoolSize <= 0 {
		return fmt.Errorf("redis.pool_size must be positive")
	}
	if cfg.Redis.DefaultCacheTTLMins <= 0 {
		return fmt.Errorf("redis.default_cache_ttl_mins must be positive")
	}
	if len(cfg.Server.CORS.AllowedOrigins) == 0 {
		return fmt.Errorf("server.cors.allowed_origins must not be empty")
	}
	for _, origin := range cfg.Server.CORS.AllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "*" {
			continue
		}
		u, err := neturl.Parse(origin)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("server.cors.allowed_origins contains invalid origin: %s", origin)
		}
	}
	if len(cfg.Server.CORS.AllowedHeaders) == 0 {
		return fmt.Errorf("server.cors.allowed_headers must not be empty")
	}
	if len(cfg.Server.CORS.AllowedMethods) == 0 {
		return fmt.Errorf("server.cors.allowed_methods must not be empty")
	}
	for _, method := range cfg.Server.CORS.AllowedMethods {
		method = strings.ToUpper(strings.TrimSpace(method))
		if _, ok := allowedHTTPMethods[method]; !ok {
			return fmt.Errorf("server.cors.allowed_methods contains invalid method: %s", method)
		}
	}
	if cfg.Server.MaxHeaderBytes <= 0 {
		cfg.Server.MaxHeaderBytes = 1 << 20 // 1 MB default
	}
	if err := validateAuthAudience("auth.user", cfg.Auth.User); err != nil {
		return err
	}
	if err := validateAuthAudience("auth.admin", cfg.Auth.Admin); err != nil {
		return err
	}
	if cfg.Admin.Lockout.Enabled {
		if cfg.Admin.Lockout.AccountMaxFails <= 0 {
			return fmt.Errorf("admin.lockout.account_max_fails must be positive")
		}
		if cfg.Admin.Lockout.AccountLockMinutes <= 0 {
			return fmt.Errorf("admin.lockout.account_lock_minutes must be positive")
		}
		if cfg.Admin.Lockout.IPMaxFails <= 0 {
			return fmt.Errorf("admin.lockout.ip_max_fails must be positive")
		}
		if cfg.Admin.Lockout.IPLockMinutes <= 0 {
			return fmt.Errorf("admin.lockout.ip_lock_minutes must be positive")
		}
	}
	if err := validateTOTP(cfg.Admin.TOTP); err != nil {
		return err
	}
	return nil
}

// validateTOTP 校验 TOTP 配置：encryption key 必须是 32 字节（AES-256-GCM 要求），其余正数即可。
func validateTOTP(cfg AdminTOTPConfig) error {
	key := strings.TrimSpace(cfg.SecretEncryptionKey)
	if key == "" {
		return fmt.Errorf("admin.totp.secret_encryption_key is required (32 bytes for AES-256-GCM)")
	}
	if len(key) != 32 {
		return fmt.Errorf("admin.totp.secret_encryption_key must be exactly 32 bytes, got %d", len(key))
	}
	if strings.TrimSpace(cfg.Issuer) == "" {
		return fmt.Errorf("admin.totp.issuer is required")
	}
	if cfg.RecoveryCodesCount <= 0 {
		return fmt.Errorf("admin.totp.recovery_codes_count must be positive")
	}
	if cfg.StepTokenTTLSeconds <= 0 {
		return fmt.Errorf("admin.totp.step_token_ttl_seconds must be positive")
	}
	if cfg.Period <= 0 {
		return fmt.Errorf("admin.totp.period must be positive")
	}
	if cfg.Digits != 6 && cfg.Digits != 8 {
		return fmt.Errorf("admin.totp.digits must be 6 or 8")
	}
	if cfg.Skew < 0 {
		return fmt.Errorf("admin.totp.skew must be >= 0")
	}
	return nil
}

// validateAuthAudience 校验单个 audience 配置必填字段。
func validateAuthAudience(prefix string, cfg AuthAudienceConfig) error {
	if strings.TrimSpace(cfg.Audience) == "" {
		return fmt.Errorf("%s.audience is required", prefix)
	}
	if cfg.AccessTokenTTLMinutes <= 0 {
		return fmt.Errorf("%s.access_token_ttl_minutes must be positive", prefix)
	}
	if cfg.RefreshTokenTTLHours <= 0 {
		return fmt.Errorf("%s.refresh_token_ttl_hours must be positive", prefix)
	}
	return nil
}

var allowedHTTPMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
}

// normalizeEnv 统一整理环境名称，避免不同写法造成分支差异。
func normalizeEnv(env string) string {
	return strings.ToLower(strings.TrimSpace(env))
}
