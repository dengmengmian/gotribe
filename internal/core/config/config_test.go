package config

// 本文件验证环境识别等配置辅助逻辑。

import "testing"

// TestAppConfigIsDevelopment 验证配置相关逻辑是否符合预期。
func TestAppConfigIsDevelopment(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		env  string
		want bool
	}{
		{env: "", want: true},
		{env: "development", want: true},
		{env: "dev", want: true},
		{env: "local", want: true},
		{env: "production", want: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.env, func(t *testing.T) {
			t.Parallel()
			if got := (AppConfig{Env: tc.env}).IsDevelopment(); got != tc.want {
				t.Fatalf("IsDevelopment(%q) = %v, want %v", tc.env, got, tc.want)
			}
		})
	}
}

// TestAppConfigIsProduction 验证配置相关逻辑是否符合预期。
func TestAppConfigIsProduction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		env  string
		want bool
	}{
		{env: "production", want: true},
		{env: "prod", want: true},
		{env: "development", want: false},
		{env: "local", want: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.env, func(t *testing.T) {
			t.Parallel()
			if got := (AppConfig{Env: tc.env}).IsProduction(); got != tc.want {
				t.Fatalf("IsProduction(%q) = %v, want %v", tc.env, got, tc.want)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	valid := Config{
		App: AppConfig{Name: "gotribe", Env: "development"},
		Server: ServerConfig{
			Port:                8080,
			MaxHeaderBytes:      1024,
			MaxRequestBodyBytes: 1024,
			CORS: CORSConfig{
				AllowedOrigins: []string{"http://localhost:3000"},
				AllowedHeaders: []string{"Authorization"},
				AllowedMethods: []string{"GET", "POST"},
			},
		},
		Database: DatabaseConfig{
			Type:                "postgres",
			Host:                "localhost",
			Port:                5432,
			Username:            "user",
			Database:            "db",
			MaxIdleConns:        5,
			MaxOpenConns:        10,
			ConnMaxIdleTimeMins: 5,
			ConnMaxLifetimeMins: 30,
		},
		Redis: RedisConfig{
			Addr:                "127.0.0.1:6379",
			DB:                  0,
			PoolSize:            10,
			DefaultCacheTTLMins: 5,
		},
		Auth: AuthConfig{
			Secret: "this-is-a-valid-test-secret-with-enough-length",
			User:   AuthAudienceConfig{Audience: "gotribe.user", AccessTokenTTLMinutes: 120, RefreshTokenTTLHours: 168},
			Admin:  AuthAudienceConfig{Audience: "gotribe.admin", AccessTokenTTLMinutes: 60, RefreshTokenTTLHours: 24},
		},
	}

	if err := validate(valid); err != nil {
		t.Fatalf("validate(valid) error = %v, want nil", err)
	}

	testCases := []struct {
		name string
		mut  func(*Config)
	}{
		{
			name: "invalid redis addr",
			mut: func(cfg *Config) {
				cfg.Redis.Addr = "localhost"
			},
		},
		{
			name: "invalid cors origin",
			mut: func(cfg *Config) {
				cfg.Server.CORS.AllowedOrigins = []string{"not-a-url"}
			},
		},
		{
			name: "invalid cors method",
			mut: func(cfg *Config) {
				cfg.Server.CORS.AllowedMethods = []string{"FETCH"}
			},
		},
		{
			name: "idle gt open",
			mut: func(cfg *Config) {
				cfg.Database.MaxIdleConns = 20
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg := valid
			tc.mut(&cfg)
			if err := validate(cfg); err == nil {
				t.Fatalf("validate(%s) = nil, want error", tc.name)
			}
		})
	}
}
