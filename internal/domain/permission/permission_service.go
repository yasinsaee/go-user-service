package permission

type PermissionService interface {
	Create(permission *Permission) error
	GetByID(id any) (*Permission, error)
	GetByName(name string) (*Permission, error)
	Update(permission *Permission) error
	Delete(id any) error
	ListAll() (Permissions, error)
	GetByIDs(ids []string) (Permissions, error)
	GetByKey(key string) (*Permission, error)
}
