package user

import (
	"github.com/yasinsaee/go-user-service/internal/context"
)

// UserRepository defines the interface for user data access operations.
type UserRepository interface {
	Create(user *User) error                       // Creates a new user in the database
	FindByUsername(username string) (*User, error) // Retrieves a user by their username
	FindByID(id any) (*User, error)                // Retrieves a user by their ID
	Update(user *User) error                       // Updates an existing user
	Delete(id any) error                           // Deletes a user by ID
	List() (Users, error)
	Count(q UserFilter) (int, error)
	//totalCounts, totalPages,model,error
	PaginationList(metaData context.MetaData, q UserFilter) (Users, error)
	SoftDelete(id any) error
	FindOneByFilter(filter UserFilter) (*User, error)
}
