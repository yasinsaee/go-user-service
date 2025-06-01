package permission

type PermissionRepository interface {
	Create(permission *Permission) error
	FindByID(id any) (*Permission, error)
	FindByName(name string) (*Permission, error)
	Update(permission *Permission) error
	Delete(id any) error
	List() (Permissions, error)
}
