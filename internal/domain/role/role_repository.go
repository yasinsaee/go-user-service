package role

type RoleRepository interface {
	Create(role *Role) error
	FindByID(id any) (*Role, error)
	FindByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id any) error
	List() ([]*Role, error)
}
