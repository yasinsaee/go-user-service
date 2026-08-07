package group

import "github.com/yasinsaee/go-user-service/internal/context"

type GroupService interface {
	Create(group *Group) error
	GetByID(id any) (*Group, error)
	GetByName(name string) (*Group, error)
	Update(group *Group) error
	Delete(id any) error
	ListAll() (Groups, error)
	//totalCounts, totalPages,model,error
	PaginationList(metaData context.MetaData, q GroupFilter) (context.MetaData, Groups, error)
	Count(q GroupFilter) (int, error)
	GetByKey(key string) (*Group, error)
	GetByFilter(filter GroupFilter) (*Group, error)
}
