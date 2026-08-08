package config

import (
	"fmt"
	"os"
)

// Config 面板配置(环境变量驱动,对照原 application.yml 与 00-overview.md 命名规范)
type Config struct {
	FrontendPort string // FRONTEND_PORT 前端页面入口(默认 6635)
	BackendPort  string // BACKEND_PORT  API/WS/flow(agent 对接,默认 6636)

	DBHost     string // DB_HOST
	DBPort     string // DB_PORT
	DBName     string // DB_NAME
	DBUser     string // DB_USER
	DBPassword string // DB_PASSWORD

	JWTSecret string // JWT_SECRET(必须设置;生产沿用原值保证存量 token 有效)
	LogDir    string // LOG_DIR
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	c := &Config{
		FrontendPort: env("FRONTEND_PORT", "6635"),
		BackendPort:  env("BACKEND_PORT", "6636"),
		DBHost:       env("DB_HOST", "127.0.0.1"),
		DBPort:       env("DB_PORT", "5432"),
		DBName:       env("DB_NAME", "flvx2"),
		DBUser:       env("DB_USER", "flvx2"),
		DBPassword:   env("DB_PASSWORD", ""),
		JWTSecret:    env("JWT_SECRET", ""),
		LogDir:       env("LOG_DIR", "logs"),
	}
	if c.DBPassword == "" {
		return nil, fmt.Errorf("缺少环境变量 DB_PASSWORD")
	}
	if c.JWTSecret == "" {
		return nil, fmt.Errorf("缺少环境变量 JWT_SECRET(生产必须沿用原值,否则存量 token 失效)")
	}
	return c, nil
}
