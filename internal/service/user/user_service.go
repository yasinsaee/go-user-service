package user

import (
	"errors"
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/util"
)

type userService struct {
	repo user.UserRepository
}

// NewUserService returns a new instance of UserService.
func NewUserService(repo user.UserRepository) user.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(user *user.User) error {
	_, err := s.repo.FindByUsername(user.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashed := util.HashPassword(user.Password)

	user.Password = hashed
	user.CreatedAt = time.Now()
	return s.repo.Create(user)
}

func (s *userService) Login(username, password string) (*user.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return nil, errors.New("invalid username or password")
	}

	if !util.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid username or password")
	}
	user.LastLogin = time.Now()
	if err = s.Update(user); err != nil {
		return nil, errors.New("update missing")
	}

	return user, nil
}

func (s *userService) GetByID(id any) (*user.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetByUsername(username string) (*user.User, error) {
	return s.repo.FindByUsername(username)
}

func (s *userService) Update(user *user.User) error {
	user.UpdatedAt = time.Now()
	return s.repo.Update(user)
}

func (s *userService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *userService) ListAll() (user.Users, error) {
	return s.repo.List()
}
