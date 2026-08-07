package repository

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"github.com/yasinsaee/go-user-service/pkg/uuid"
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
	permission.CreatedAt = time.Now().UTC()
	permission.UpdatedAt = time.Now().UTC()
	permission.UniqueID = uuid.Must(uuid.NewV7()).String()
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
	query := bson.M{
		"unique_id": id,
		"$or": []bson.M{
			{"is_deleted": bson.M{"$exists": false}},
			{"is_deleted": false},
		},
	}
	err := mongo2.FindOne(r.collection.Name(), query, p)
	return p, err
}

// Update modifies an existing permission and sets the update timestamp.
func (r *mongoPermissionRepository) Update(permission *permission.Permission) error {
	permission.UpdatedAt = time.Now().UTC()
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
	query := bson.M{"is_deleted": bson.M{"$ne": true}}
	err := mongo2.Find(r.collection.Name(), query, &permissions)
	if err != nil {
		logger.Error("error while fetching permissions: ", err.Error())
		return nil, err
	}

	return permissions, nil
}

func (r *mongoPermissionRepository) SoftDelete(id any) error {
	return mongo2.UpdateOne(r.collection.Name(), bson.M{"unique_id": id}, bson.M{"$set": bson.M{
		"is_deleted": true, "deleted_at": time.Now(),
	}})
}

func (r *mongoPermissionRepository) GetByIDs(ids []string) (permission.Permissions, error) {
	permissions := make(permission.Permissions, 0)

	filter := bson.M{
		"unique_id": bson.M{"$in": ids},
		"$or": []bson.M{
			{"is_deleted": bson.M{"$exists": false}},
			{"is_deleted": false},
		},
	}
	err := mongo2.Find(r.collection.Name(), filter, &permissions)
	if err != nil {
		logger.Error("error while fetching permissions: ", err.Error())
		return nil, err
	}

	return permissions, nil
}

func (r *mongoPermissionRepository) FindOneByFilter(filter permission.PermissionFilter) (*permission.Permission, error) {
	p := new(permission.Permission)
	err := mongo2.FindOne(r.collection.Name(), filter.GetFilters(), &p)
	if err != nil {
		logger.Error("error while fetching permission: ", err.Error())
		return nil, err
	}
	return p, nil
}
