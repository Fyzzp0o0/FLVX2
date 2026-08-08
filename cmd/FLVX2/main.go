package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fyzzp0o0/FLVX2/internal/api"
	"github.com/Fyzzp0o0/FLVX2/internal/config"
	"github.com/Fyzzp0o0/FLVX2/internal/db"
	"github.com/Fyzzp0o0/FLVX2/internal/service"
	"github.com/Fyzzp0o0/FLVX2/internal/ws"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FLVX2] 配置错误: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Init(ctx, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword)
	if err != nil {
		log.Fatalf("[FLVX2] 数据库初始化失败: %v", err)
	}
	defer pool.Close()

	users := service.NewUserService(pool)
	nodes := service.NewNodeService(pool)
	hub := ws.NewHub(pool, nodes, cfg.JWTSecret)
	flows := service.NewFlowService(pool, hub)
	handler := api.NewHandler(users, flows, hub, cfg.JWTSecret)
	router := api.NewRouter(handler)

	// 双端口监听:6636 后端(API/WS/flow,agent 对接) 与 6635 前端(静态+同源 API)
	servers := []*http.Server{
		{Addr: ":" + cfg.BackendPort, Handler: router},
		{Addr: ":" + cfg.FrontendPort, Handler: router},
	}
	for _, s := range servers {
		ln, err := net.Listen("tcp", s.Addr)
		if err != nil {
			log.Fatalf("[FLVX2] 监听 %s 失败: %v", s.Addr, err)
		}
		go func(srv *http.Server, l net.Listener) {
			log.Printf("[FLVX2] 已监听 %s", l.Addr())
			if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("[FLVX2] 服务 %s 异常退出: %v", l.Addr(), err)
			}
		}(s, ln)
	}
	log.Printf("[FLVX2] 面板已启动: 前端 http://127.0.0.1:%s · 后端(agent) http://127.0.0.1:%s", cfg.FrontendPort, cfg.BackendPort)

	// 优雅退出(对应原 server.shutdown=graceful)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("[FLVX2] 正在关闭...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for _, s := range servers {
		_ = s.Shutdown(shutdownCtx)
	}
	fmt.Println("[FLVX2] 已退出")
}
