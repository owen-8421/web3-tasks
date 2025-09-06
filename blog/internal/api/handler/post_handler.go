package handler

import (
	"blog/internal/api/response"
	"blog/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{postService: svc}
}

// CreatePostRequest 定义了创建文章的请求体结构
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdatePostRequest 定义了更新文章的请求体结构
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// Create 处理创建文章的请求
func (h *PostHandler) Create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	// 从中间件获取当前登录用户的ID
	userID := c.MustGet("userID").(uint)

	post, err := h.postService.CreatePost(req.Title, req.Content, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, post)
}

// GetByID 处理根据ID获取文章的请求
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Post not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, post)
}

// GetAll 处理获取所有文章列表的请求
func (h *PostHandler) GetAll(c *gin.Context) {
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, posts)
}

// Update 处理更新文章的请求
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	userID := c.MustGet("userID").(uint)

	post, err := h.postService.UpdatePost(uint(id), userID, req.Title, req.Content)
	if err != nil {
		if err.Error() == "permission denied: you are not the author of this post" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Post not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, post)
}

// Delete 处理删除文章的请求
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.postService.DeletePost(uint(id), userID)
	if err != nil {
		if err.Error() == "permission denied: you are not the author of this post" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Post not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Post deleted successfully"})
}
