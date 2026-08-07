package role

import "github.com/yasinsaee/go-user-service/internal/context"

type RoleRepository interface {
	Create(role *Role) error
	FindByID(id any) (*Role, error)
	FindByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id any) error
	SoftDelete(id any) error
	List() (Roles, error)
	Count(q RoleFilter) (int, error)
	//totalCounts, totalPages,model,error
	PaginationList(metaData context.MetaData, q RoleFilter) (Roles, error)
	GetByIDs(ids []string) (Roles, error)
	FindOneByFilter(filter RoleFilter) (*Role, error)
}
