package service

import (
	"blog/internal/model"
	"blog/internal/repository"
	"errors"
)

// PostService 定义了文章服务的接口
type PostService struct {
	postRepo *repository.PostRepository
}

// NewPostService 创建一个新的 PostService 实例
func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

// CreatePost 创建一篇新文章
func (s *PostService) CreatePost(title, content string, userID uint) (*model.Post, error) {
	post := &model.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}
	err := s.postRepo.Create(post)
	return post, err
}

// GetPostByID 根据ID获取文章详情
func (s *PostService) GetPostByID(id uint) (*model.Post, error) {
	return s.postRepo.FindByID(id)
}

// GetAllPosts 获取所有文章列表
func (s *PostService) GetAllPosts() ([]model.Post, error) {
	return s.postRepo.FindAll()
}

// UpdatePost 更新文章，并校验操作权限
func (s *PostService) UpdatePost(id, userID uint, title, content string) (*model.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 权限校验：确保是文章作者本人在操作
	if post.UserID != userID {
		return nil, errors.New("permission denied: you are not the author of this post")
	}

	post.Title = title
	post.Content = content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePost 删除文章，并校验操作权限
func (s *PostService) DeletePost(id, userID uint) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 权限校验：确保是文章作者本人在操作
	if post.UserID != userID {
		return errors.New("permission denied: you are not the author of this post")
	}

	return s.postRepo.Delete(id)
}
