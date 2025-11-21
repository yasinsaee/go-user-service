package role

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
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
	return mongo2.Create(role)
}

func (r *mongoRoleRepository) FindByID(id any) (*role.Role, error) {
	rData := new(role.Role)
	err := mongo2.Get(r.collection.Name(), id, rData)
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
