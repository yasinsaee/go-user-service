package group

import "github.com/yasinsaee/go-user-service/internal/context"

// GroupRepository defines the interface for Group data access operations.
type GroupRepository interface {
	Create(group *Group) error       // Creates a new Group in the database
	FindByID(id any) (*Group, error) // Retrieves a Group by their ID
	Update(group *Group) error       // Updates an existing Group
	Delete(id any) error             // Deletes a Group by ID
	List() (Groups, error)
	Count(q GroupFilter) (int, error)
	//totalCounts, totalPages,model,error
	PaginationList(metaData context.MetaData, q GroupFilter) (Groups, error)
	SoftDelete(id any) error
	FindOneByFilter(filter GroupFilter) (*Group, error)
}
