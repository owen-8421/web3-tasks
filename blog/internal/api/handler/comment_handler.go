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

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: svc}
}

// CreateCommentRequest 定义了创建评论的请求体结构
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// Create 处理在某篇文章下创建评论的请求
func (h *CommentHandler) Create(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	userID := c.MustGet("userID").(uint)

	comment, err := h.commentService.CreateComment(req.Content, uint(postID), userID)
	if err != nil {
		if err.Error() == "post not found" {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, comment)
}

// GetByPostID 处理获取某篇文章下所有评论的请求
func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post ID")
		return
	}

	comments, err := h.commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, comments)
}

// Delete 处理删除评论的请求
// 注意：此处的 :id 是 commentID
func (h *CommentHandler) Delete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.commentService.DeleteComment(uint(commentID), userID)
	if err != nil {
		if err.Error() == "permission denied: you are not authorized to delete this comment" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "Comment not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Comment deleted successfully"})
}
