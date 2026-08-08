package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var indexHTML []byte

//go:embed dist
var distFS embed.FS

// MIME 映射(embed.FS 无 Content-Type,需自维护;完整映射后续按 dist 产物补充)
var mimeTypes = map[string]string{
	".html":  "text/html; charset=utf-8",
	".js":    "application/javascript",
	".css":   "text/css",
	".json":  "application/json",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".jpeg":  "image/jpeg",
	".gif":   "image/gif",
	".svg":   "image/svg+xml",
	".ico":   "image/x-icon",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
}

// ServeIndex 返回 index.html(SPA 入口;优先 Vue 构建产物,缺失时回退占位页)
func ServeIndex(c *gin.Context) {
	if data, err := readDistFile("index.html"); err == nil {
		c.Data(http.StatusOK, mimeTypes[".html"], data)
		return
	}
	c.Data(http.StatusOK, mimeTypes[".html"], indexHTML)
}

// ServeAsset 从 embed.FS 读取 dist/assets 静态资源
func ServeAsset(c *gin.Context) {
	rel := strings.TrimPrefix(c.Param("filepath"), "/")
	if rel == "" || strings.Contains(rel, "..") {
		c.Status(http.StatusNotFound)
		return
	}
	data, err := readDistFile("assets/" + rel)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	ext := ""
	if i := strings.LastIndex(rel, "."); i >= 0 {
		ext = rel[i:]
	}
	ct := mimeTypes[ext]
	if ct == "" {
		ct = "application/octet-stream"
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable") // 对应 nginx expires 1y
	c.Data(http.StatusOK, ct, data)
}

// readDistFile 从 embed dist 读取文件
func readDistFile(rel string) ([]byte, error) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(sub, rel)
}
