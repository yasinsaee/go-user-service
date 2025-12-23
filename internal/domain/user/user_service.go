package user

// UserService defines business logic operations related to users.
type UserService interface {
	Register(username string, user *User) error
	Login(username, password string) (*User, error)
	GetByID(id any) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id any) error
	ListAll() (Users, error)
	ResetPassword(user *User, currentPassword, password, rePassword string) error
	UpdatePassword(user *User, password, rePassword string) error
}
