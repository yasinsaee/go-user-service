package repository

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoPermissionRepository implements the PermissionRepository interface using MongoDB.
type mongoPermissionRepository struct {
	collection *mongo.Collection
}

// NewMongoPermissionRepository returns a new instance of mongoPermissionRepository.
func NewMongoPermissionRepository(db *mongo.Database, collectionName string) permission.PermissionRepository {
	return &mongoPermissionRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new permission into the database and sets the creation timestamp.
func (r *mongoPermissionRepository) Create(permission *permission.Permission) error {
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()
	return mongo2.Create(permission)
}

// FindByName returns a permission document that matches the given name.
func (r *mongoPermissionRepository) FindByName(name string) (*permission.Permission, error) {
	p := new(permission.Permission)
	query := bson.M{"name": name}
	err := mongo2.FindOne(r.collection.Name(), query, p)
	return p, err
}

// FindByID retrieves a permission by their ID (string or ObjectID).
func (r *mongoPermissionRepository) FindByID(id any) (*permission.Permission, error) {
	p := new(permission.Permission)
	err := mongo2.Get(r.collection.Name(), id, p)
	return p, err
}

// Update modifies an existing permission and sets the update timestamp.
func (r *mongoPermissionRepository) Update(permission *permission.Permission) error {
	permission.UpdatedAt = time.Now()
	return mongo2.Update(permission)
}

// Delete removes a permission by their ID after converting it to ObjectID.
func (r *mongoPermissionRepository) Delete(id any) error {
	objID, err := util.ToObjectID(id)
	if err != nil {
		logger.Error("error while delete permission: ", err.Error())
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

// List returns all permissions from the collection.
func (r *mongoPermissionRepository) List() (permission.Permissions, error) {
	permissions := make(permission.Permissions, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{}, &permissions)
	if err != nil {
		logger.Error("error while fetching permissions: ", err.Error())
		return nil, err
	}

	return permissions, nil
}
