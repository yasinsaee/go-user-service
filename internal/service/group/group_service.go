package group

import (
	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/group"
)

type GroupServiceImpl struct {
	repo group.GroupRepository
}

func NewGroupService(repo group.GroupRepository) group.GroupService {
	return &GroupServiceImpl{
		repo: repo,
	}
}

func (s *GroupServiceImpl) Create(group *group.Group) error {
	return s.repo.Create(group)
}

func (s *GroupServiceImpl) GetByID(id any) (*group.Group, error) {
	return s.repo.FindByID(id)
}

func (s *GroupServiceImpl) GetByName(name string) (*group.Group, error) {
	return s.repo.FindOneByFilter(group.GroupFilter{IsDelete: "false", Name: name})
}

func (s *GroupServiceImpl) Update(group *group.Group) error {
	return s.repo.Update(group)
}

func (s *GroupServiceImpl) Delete(id any) error {
	return s.repo.SoftDelete(id)
}

func (s *GroupServiceImpl) ListAll() (group.Groups, error) {
	return s.repo.List()
}

func (s *GroupServiceImpl) Count(q group.GroupFilter) (int, error) {
	return s.repo.Count(q)
}

func (s *GroupServiceImpl) PaginationList(metaData context.MetaData, q group.GroupFilter) (context.MetaData, group.Groups, error) {
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

func (s *GroupServiceImpl) GetByKey(key string) (*group.Group, error) {
	return s.repo.FindOneByFilter(group.GroupFilter{Key: key, IsDelete: "false"})
}

func (s *GroupServiceImpl) GetByFilter(filter group.GroupFilter) (*group.Group, error) {
	return s.repo.FindOneByFilter(filter)
}
