package group

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/group"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"github.com/yasinsaee/go-user-service/pkg/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGroupRepository struct {
	collection *mongo.Collection
}

func NewMongoGroupRepository(db *mongo.Database, collectionName string) group.GroupRepository {
	return &mongoGroupRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *mongoGroupRepository) Create(group *group.Group) error {
	group.CreatedAt = time.Now().UTC()
	group.UpdatedAt = time.Now().UTC()
	group.UniqueID = uuid.Must(uuid.NewV7()).String()
	return mongo2.Create(group)
}

func (r *mongoGroupRepository) FindByID(id any) (*group.Group, error) {
	rData := new(group.Group)
	query := bson.M{
		"unique_id": id,
		"$or": []bson.M{
			{"is_deleted": bson.M{"$exists": false}},
			{"is_deleted": false},
		},
	}
	err := mongo2.FindOne(r.collection.Name(), query, &rData)
	return rData, err
}

func (r *mongoGroupRepository) Update(group *group.Group) error {
	group.UpdatedAt = time.Now().UTC()
	return mongo2.Update(group)
}

func (r *mongoGroupRepository) Delete(id any) error {
	objID, err := util.ToObjectID(id)
	if err != nil {
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

func (r *mongoGroupRepository) List() (group.Groups, error) {
	groups := make(group.Groups, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{}, &groups)
	if err != nil {
		logger.Error("error while fetching groups: ", err.Error())
		return nil, err
	}

	return groups, nil
}

// Count returns number of OTPs matching a query.
func (r *mongoGroupRepository) Count(q group.GroupFilter) (int, error) {
	query := q.GetFilters()
	return mongo2.Count(r.collection.Name(), query)
}

// PaginationList returns paginated categories.
func (r *mongoGroupRepository) PaginationList(metaData context.MetaData, q group.GroupFilter) (group.Groups, error) {
	groups := make(group.Groups, 0)
	query := q.GetFilters()
	err := mongo2.Find(r.collection.Name(), query, &groups, metaData.Limit, metaData.CurrentPage, metaData.Sort)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *mongoGroupRepository) SoftDelete(id any) error {
	return mongo2.UpdateOne(r.collection.Name(), bson.M{"unique_id": id}, bson.M{"$set": bson.M{"is_deleted": true, "deleted_at": time.Now().UTC()}})
}

func (r *mongoGroupRepository) FindOneByFilter(filter group.GroupFilter) (*group.Group, error) {
	g := new(group.Group)
	err := mongo2.FindOne(r.collection.Name(), filter.GetFilters(), &g)
	if err != nil {
		logger.Error("error while fetching group: ", err.Error())
		return nil, err
	}
	return g, nil
}
