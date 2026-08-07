package permission

import (
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
)

// permissionServiceImpl is the concrete implementation of PermissionService.
type permissionServiceImpl struct {
	repo permission.PermissionRepository
}

// NewPermissionService creates a new instance of PermissionService.
func NewPermissionService(repo permission.PermissionRepository) permission.PermissionService {
	return &permissionServiceImpl{
		repo: repo,
	}
}

func (s *permissionServiceImpl) Create(permission *permission.Permission) error {
	return s.repo.Create(permission)
}

func (s *permissionServiceImpl) GetByID(id any) (*permission.Permission, error) {
	return s.repo.FindByID(id)
}

func (s *permissionServiceImpl) GetByName(name string) (*permission.Permission, error) {
	return s.repo.FindByName(name)
}

func (s *permissionServiceImpl) Update(permission *permission.Permission) error {
	return s.repo.Update(permission)
}

func (s *permissionServiceImpl) Delete(id any) error {
	return s.repo.SoftDelete(id)
}

func (s *permissionServiceImpl) ListAll() (permission.Permissions, error) {
	return s.repo.List()
}

func (s *permissionServiceImpl) GetByIDs(ids []string) (permission.Permissions, error) {
	if len(ids) == 0 {
		return make(permission.Permissions, 0), nil
	}
	return s.repo.GetByIDs(ids)
}

func (s *permissionServiceImpl) GetByKey(key string) (*permission.Permission, error) {
	return s.repo.FindOneByFilter(permission.PermissionFilter{Key: key, IsDelete: "false"})
}
