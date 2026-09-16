package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"personal-blog/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	db *gorm.DB
}

func NewArticleHandler(db *gorm.DB) *ArticleHandler {
	return &ArticleHandler{db: db}
}

type listResponse struct {
	Total int64                  `json:"total"`
	Items []model.ArticleListItem `json:"items"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}

func parsePaging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	return page, size
}

// ListPublic 返回已发布文章分页列表，按创建时间倒序，不含正文。
func (h *ArticleHandler) ListPublic(c *gin.Context) {
	page, size := parsePaging(c)
	var items []model.ArticleListItem
	var total int64
	q := h.db.Model(&model.Article{}).Where("status = ?", model.StatusPublish)
	q.Count(&total)
	if err := q.Order("created_at DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, listResponse{Total: total, Items: items, Page: page, Size: size})
}

// GetPublic 返回单篇文章详情，同时累加浏览量。草稿不对外暴露。
func (h *ArticleHandler) GetPublic(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var a model.Article
	if err := h.db.Where("id = ? AND status = ?", id, model.StatusPublish).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	if err := h.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update view count failed"})
		return
	}
	a.ViewCount++
	c.JSON(http.StatusOK, a)
}

// --- 管理 API ---

type createArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Summary    string `json:"summary"`
	Content    string `json:"content" binding:"required"`
	CoverImage string `json:"cover_image"`
	Tags       string `json:"tags"`
	Status     int8   `json:"status"`
}

func (r *createArticleRequest) normalize() {
	r.Title = strings.TrimSpace(r.Title)
	r.Summary = strings.TrimSpace(r.Summary)
	r.Content = strings.TrimSpace(r.Content)
	if r.Status != model.StatusDraft && r.Status != model.StatusPublish {
		r.Status = model.StatusPublish
	}
}

// ListAdmin 管理端列表，可选按 status 过滤，含草稿，不含正文。
func (h *ArticleHandler) ListAdmin(c *gin.Context) {
	page, size := parsePaging(c)
	q := h.db.Model(&model.Article{})
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			q = q.Where("status = ?", status)
		}
	}
	var total int64
	q.Count(&total)
	var items []model.ArticleListItem
	if err := q.Order("created_at DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, listResponse{Total: total, Items: items, Page: page, Size: size})
}

// GetAdmin 管理端详情，返回完整正文，不计浏览量。
func (h *ArticleHandler) GetAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var a model.Article
	if err := h.db.First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// Create 创建文章。
func (h *ArticleHandler) Create(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.normalize()
	a := model.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		Tags:       req.Tags,
		Status:     req.Status,
	}
	if err := h.db.Create(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, a)
}

// Update 更新文章。
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	req.normalize()
	var a model.Article
	if err := h.db.First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	a.Title = req.Title
	a.Summary = req.Summary
	a.Content = req.Content
	a.CoverImage = req.CoverImage
	a.Tags = req.Tags
	a.Status = req.Status
	if err := h.db.Save(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}
	c.JSON(http.StatusOK, a)
}

// Delete 删除文章。
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.db.Delete(&model.Article{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
