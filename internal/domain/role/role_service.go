package role

type RoleService interface {
	Create(role *Role) error
	GetByID(id any) (*Role, error)
	GetByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id any) error
	ListAll() (Roles, error)
}
