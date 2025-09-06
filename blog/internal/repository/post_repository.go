package repository

import (
	"blog/internal/model"
	"gorm.io/gorm"
)

// PostRepository 定义了文章数据仓库的接口
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建一个新的 PostRepository 实例
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create 创建一篇文章
func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// FindByID 根据ID查找文章，并预加载用户信息和评论
func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	// Preload User and Comments to get related data
	err := r.db.Preload("User").Preload("Comments").First(&post, id).Error
	return &post, err
}

// FindAll 获取所有文章列表
func (r *PostRepository) FindAll() ([]model.Post, error) {
	var posts []model.Post
	// Find all posts and preload the user for each post
	err := r.db.Preload("User").Order("created_at desc").Find(&posts).Error
	return posts, err
}

// Update 更新一篇文章
func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete 根据ID删除一篇文章
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}
