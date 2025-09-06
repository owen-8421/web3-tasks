package service

import (
	"blog/internal/model"
	"blog/internal/repository"
	"errors"
)

// CommentService 定义了评论服务的接口
type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository // 依赖 PostRepository 检查文章是否存在
}

// NewCommentService 创建一个新的 CommentService 实例
func NewCommentService(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment 创建一条新评论
func (s *CommentService) CreateComment(content string, postID, userID uint) (*model.Comment, error) {
	// 检查文章是否存在
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comment := &model.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	err = s.commentRepo.Create(comment)
	return comment, err
}

// GetCommentsByPostID 根据文章ID获取所有评论
func (s *CommentService) GetCommentsByPostID(postID uint) ([]model.Comment, error) {
	return s.commentRepo.FindByPostID(postID)
}

// DeleteComment 删除评论，并校验操作权限
// 在这个简单的博客系统中，我们允许文章作者或评论者本人删除评论
func (s *CommentService) DeleteComment(commentID, userID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	post, err := s.postRepo.FindByID(comment.PostID)
	if err != nil {
		return err
	}

	// 权限校验：评论者本人 或 文章作者 可以删除评论
	if comment.UserID != userID && post.UserID != userID {
		return errors.New("permission denied: you are not authorized to delete this comment")
	}

	return s.commentRepo.Delete(commentID)
}
