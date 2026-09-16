package router

import (
	"net/http"
	"time"

	"personal-blog/internal/config"
	"personal-blog/internal/handler"
	"personal-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, auth *handler.AuthHandler, article *handler.ArticleHandler, stats *handler.StatsHandler) *gin.Engine {
	r := gin.Default()

	// CORS：开发环境允许 Vite dev server 访问
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.POST("/auth/login", auth.Login)

		api.GET("/articles", article.ListPublic)
		api.GET("/articles/:id", article.GetPublic)

		admin := api.Group("/admin", middleware.AuthMiddleware(cfg))
		{
			admin.GET("/articles", article.ListAdmin)
			admin.GET("/articles/:id", article.GetAdmin)
			admin.POST("/articles", article.Create)
			admin.PUT("/articles/:id", article.Update)
			admin.DELETE("/articles/:id", article.Delete)
			admin.GET("/stats", stats.Stats)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", time.Hour.String())
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
