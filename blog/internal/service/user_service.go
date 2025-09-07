package service

import (
	"blog/internal/model"
	"blog/internal/pkg/jwt"
	"blog/internal/repository"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo       *repository.UserRepository
	jwtIssuer      string
	jwtExpireHours int
}

func NewUserService(repo *repository.UserRepository, issuer string, expire int) *UserService {
	return &UserService{
		userRepo:       repo,
		jwtIssuer:      issuer,
		jwtExpireHours: expire,
	}
}

func (s *UserService) Register(username, password, email string) error {
	// 检查用户是否已存在
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建用户
	user := &model.User{Username: username, Email: email}
	if err := user.SetPassword(password); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found or invalid credentials")
		}
		return "", err
	}

	if !user.CheckPassword(password) {
		return "", errors.New("user not found or invalid credentials")
	}

	return jwt.GenerateToken(user.ID, user.Username, s.jwtIssuer, s.jwtExpireHours)
}
