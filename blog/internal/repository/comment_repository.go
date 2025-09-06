package repository

import (
	"blog/internal/model"
	"gorm.io/gorm"
)

// CommentRepository 定义了评论数据仓库
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建一个新的 CommentRepository 实例
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create 创建一条评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// FindByPostID 根据文章ID查找所有评论
func (r *CommentRepository) FindByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	// Find all comments for a given post and preload the user for each comment
	err := r.db.Where("post_id = ?", postID).Preload("User").Order("created_at asc").Find(&comments).Error
	return comments, err
}

// Delete 根据ID删除一条评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// FindByID 根据ID查找评论
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, id).Error
	return &comment, err
}
