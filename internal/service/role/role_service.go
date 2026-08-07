package role

import (
	"github.com/yasinsaee/go-user-service/internal/context"
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
	return s.repo.SoftDelete(id)
}

func (s *roleServiceImpl) ListAll() (role.Roles, error) {
	return s.repo.List()
}

func (s *roleServiceImpl) Count(q role.RoleFilter) (int, error) {
	return s.repo.Count(q)
}

func (s *roleServiceImpl) PaginationList(metaData context.MetaData, q role.RoleFilter) (context.MetaData, role.Roles, error) {
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

	//default
	q.IsDelete = "false"
	cats, err := s.repo.PaginationList(metaData, q)
	if err != nil {
		return metaData, nil, err
	}

	return metaData, cats, nil
}

func (r *roleServiceImpl) GetByIDs(ids []string) (role.Roles, error) {
	if len(ids) == 0 {
		return make(role.Roles, 0), nil
	}

	return r.repo.GetByIDs(ids)
}

func (s *roleServiceImpl) GetByKey(key string) (*role.Role, error) {
	return s.repo.FindOneByFilter(role.RoleFilter{Key: key, IsDelete: "false"})
}
