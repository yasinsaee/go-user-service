package role

import "github.com/yasinsaee/go-user-service/internal/context"

type RoleService interface {
	Create(role *Role) error
	GetByID(id any) (*Role, error)
	GetByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id any) error
	ListAll() (Roles, error)
	//totalCounts, totalPages,model,error
	PaginationList(metaData context.MetaData, q RoleFilter) (context.MetaData, Roles, error)
	Count(q RoleFilter) (int, error)
	GetByIDs(ids []string) (Roles, error)
	GetByKey(key string) (*Role, error)
}
