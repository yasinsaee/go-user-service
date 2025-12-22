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

func (s *userService) Register(username string, user *user.User) error {
	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashed := util.HashPassword(user.Password)

	user.Password = hashed
	user.CreatedAt = time.Now().UTC()
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

	user.LastLogin = time.Now().UTC()
	if err = s.Update(user); err != nil {
		return nil, errors.New("update failed")
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
	user.UpdatedAt = time.Now().UTC()
	return s.repo.Update(user)
}

func (s *userService) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *userService) ListAll() (user.Users, error) {
	return s.repo.List()
}

func (s *userService) ResetPassword(user *user.User, currentPassword, password, rePassword string) error {
	if !util.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("password_is_not_ok")
	}
	if password != rePassword {
		return errors.New("password_is_not_matched")
	}
	user.Password = util.HashPassword(password)
	return s.Update(user)
}
