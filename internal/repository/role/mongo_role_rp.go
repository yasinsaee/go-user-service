package role

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"github.com/yasinsaee/go-user-service/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRoleRepository struct {
	collection *mongo.Collection
}

func NewMongoRoleRepository(db *mongo.Database, collectionName string) role.RoleRepository {
	return &mongoRoleRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoRoleRepository) Create(role *role.Role) error {
	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()
	role.UniqueID = uuid.Must(uuid.NewV7()).String()
	return mongo2.Create(role)
}

func (r *mongoRoleRepository) FindByID(id any) (*role.Role, error) {
	rData := new(role.Role)
	query := bson.M{
		"unique_id": id,
		"$or": []bson.M{
			{"is_deleted": bson.M{"$exists": false}},
			{"is_deleted": false},
		},
	}
	err := mongo2.FindOne(r.collection.Name(), query, rData)
	return rData, err
}

func (r *mongoRoleRepository) FindByName(name string) (*role.Role, error) {
	rData := new(role.Role)
	query := bson.M{"name": name}
	err := mongo2.FindOne(r.collection.Name(), query, rData)
	return rData, err
}

func (r *mongoRoleRepository) Update(role *role.Role) error {
	role.UpdatedAt = time.Now().UTC()
	return mongo2.Update(role)
}

func (r *mongoRoleRepository) Delete(id any) error {
	objID, err := util.ToObjectID(id)
	if err != nil {
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

func (r *mongoRoleRepository) List() (role.Roles, error) {
	roles := make(role.Roles, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{}, &roles)
	if err != nil {
		logger.Error("error while fetching roles: ", err.Error())
		return nil, err
	}

	return roles, nil
}

// Count returns number of OTPs matching a query.
func (r *mongoRoleRepository) Count(q role.RoleFilter) (int, error) {
	query := q.GetFilters()
	return mongo2.Count(r.collection.Name(), query)
}

// PaginationList returns paginated categories.
func (r *mongoRoleRepository) PaginationList(metaData context.MetaData, q role.RoleFilter) (role.Roles, error) {
	roles := make(role.Roles, 0)
	query := q.GetFilters()

	err := mongo2.Find(r.collection.Name(), query, &roles, metaData.Limit, metaData.CurrentPage, metaData.Sort)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *mongoRoleRepository) SoftDelete(id any) error {
	return mongo2.UpdateOne(r.collection.Name(), bson.M{"unique_id": id}, bson.M{"$set": bson.M{
		"is_deleted": true, "deleted_at": time.Now(),
	}})
}

func (r *mongoRoleRepository) GetByIDs(ids []string) (role.Roles, error) {
	roles := make(role.Roles, 0)

	filter := bson.M{
		"unique_id":  bson.M{"$in": ids},
		"is_deleted": false,
	}
	err := mongo2.Find(r.collection.Name(), filter, &roles)
	if err != nil {
		logger.Error("error while fetching roles: ", err.Error())
		return nil, err
	}

	return roles, nil
}

func (r *mongoRoleRepository) FindOneByFilter(filter role.RoleFilter) (*role.Role, error) {
	ro := new(role.Role)
	err := mongo2.FindOne(r.collection.Name(), filter.GetFilters(), &ro)
	if err != nil {
		logger.Error("error while fetching role: ", err.Error())
		return nil, err
	}
	return ro, nil
}
