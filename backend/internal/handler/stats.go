package handler

import (
	"database/sql"
	"net/http"

	"personal-blog/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

type statsResponse struct {
	TotalArticles  int64 `json:"total_articles"`
	PublishedCount int64 `json:"published_count"`
	DraftCount     int64 `json:"draft_count"`
	TotalViews     int64 `json:"total_views"`
}

// Stats 返回后台仪表盘所需聚合统计。
func (h *StatsHandler) Stats(c *gin.Context) {
	var resp statsResponse
	h.db.Model(&model.Article{}).Count(&resp.TotalArticles)
	h.db.Model(&model.Article{}).Where("status = ?", model.StatusPublish).Count(&resp.PublishedCount)
	h.db.Model(&model.Article{}).Where("status = ?", model.StatusDraft).Count(&resp.DraftCount)
	// 总浏览量：用 Row().Scan 直接读取聚合值，避免 GORM 结构体字段名映射差异
	var total sql.NullInt64
	row := h.db.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Row()
	_ = row.Scan(&total)
	resp.TotalViews = total.Int64

	// 近期 5 篇文章
	var recent []model.ArticleListItem
	h.db.Model(&model.Article{}).Order("created_at DESC").Limit(5).Find(&recent)

	c.JSON(http.StatusOK, gin.H{
		"stats":  resp,
		"recent": recent,
	})
}
