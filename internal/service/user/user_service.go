package user

import (
	"errors"
	"time"

	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/util"
)

type userService struct {
	repo       user.UserRepository
	tokenStore user.RefreshTokenStore // Redis-based limiter
}

// NewUserService returns a new instance of UserService.
func NewUserService(repo user.UserRepository, tokenStore user.RefreshTokenStore) user.UserService {
	return &userService{
		repo:       repo,
		tokenStore: tokenStore,
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
	return s.repo.SoftDelete(id)
}

func (s *userService) ListAll() (user.Users, error) {
	return s.repo.List()
}

func (s *userService) ResetPassword(user *user.User, currentPassword, password, rePassword string) error {
	if !util.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("password_is_not_ok")
	}
	if err := s.UpdatePassword(user, password, rePassword); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdatePassword(user *user.User, password, rePassword string) error {
	if password != rePassword {
		return errors.New("password_is_not_matched")
	}
	user.Password = util.HashPassword(password)
	return s.Update(user)
}

func (s *userService) StoreRefreshToken(userID string, refreshToken string) error {
	return s.tokenStore.Set(userID, refreshToken)
}

func (s *userService) ValidateRefreshToken(userID string, refreshToken string) (bool, error) {
	return s.tokenStore.Exists(userID, refreshToken)
}

func (s *userService) RevokeRefreshToken(userID string, refreshToken string) error {
	return s.tokenStore.Delete(userID, refreshToken)
}

func (s *userService) Count(q user.UserFilter) (int, error) {
	return s.repo.Count(q)
}

func (s *userService) PaginationList(metaData context.MetaData, q user.UserFilter) (context.MetaData, user.Users, error) {
	totalCount, err := s.Count(q)
	if err != nil {
		return metaData, nil, err
	}

	totalPages := 0
	if metaData.Limit > 0 {
		totalPages = totalCount / metaData.Limit
		if totalCount%metaData.Limit != 0 {
			totalPages++
		}
	}

	if metaData.CurrentPage < 1 {
		metaData.CurrentPage = 1
	} else if metaData.CurrentPage > totalPages && totalPages > 0 {
		metaData.CurrentPage = totalPages
	}

	nextPage := 0
	if metaData.CurrentPage < totalPages {
		nextPage = metaData.CurrentPage + 1
	}

	metaData.TotalCounts = totalCount
	metaData.TotalPages = totalPages
	metaData.NextPage = nextPage

	cats, err := s.repo.PaginationList(metaData, q)
	if err != nil {
		return metaData, nil, err
	}

	return metaData, cats, nil
}

func (s *userService) BanUser(id string) (*user.User, error) {
	var (
		err error
	)
	usr := new(user.User)
	if usr, err = s.GetByID(id); err != nil {
		return nil, err
	}

	usr.IsBanned = !usr.IsBanned
	usr.BannedAt = time.Now().UTC()
	if err = s.Update(usr); err != nil {
		return nil, err
	}
	return usr, nil
}
