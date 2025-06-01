package role

import (
	"github.com/yasinsaee/go-user-service/internal/domain/role"
)

type roleServiceImpl struct {
	repo role.RoleRepository
}

func NewRoleService(repo role.RoleRepository) role.RoleService {
	return &roleServiceImpl{
		repo: repo,
	}
}

func (s *roleServiceImpl) Create(role *role.Role) error {
	return s.repo.Create(role)
}

func (s *roleServiceImpl) GetByID(id any) (*role.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleServiceImpl) GetByName(name string) (*role.Role, error) {
	return s.repo.FindByName(name)
}

func (s *roleServiceImpl) Update(role *role.Role) error {
	return s.repo.Update(role)
}

func (s *roleServiceImpl) Delete(id any) error {
	return s.repo.Delete(id)
}

func (s *roleServiceImpl) ListAll() (role.Roles, error) {
	return s.repo.List()
}
