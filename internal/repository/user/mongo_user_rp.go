package repository

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoUserRepository implements the UserRepository interface using MongoDB.
type mongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository returns a new instance of mongoUserRepository.
func NewMongoUserRepository(db *mongo.Database, collectionName string) user.UserRepository {
	return &mongoUserRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new user into the database and sets the creation timestamp.
func (r *mongoUserRepository) Create(user *user.User) error {
	user.CreatedAt = time.Now()
	return mongo2.Create(user)
}

// FindByUsername returns a user document that matches the given username.
func (r *mongoUserRepository) FindByUsername(username string) (*user.User, error) {
	u := new(user.User)
	query := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"phonenumber": username},
			{"email": username},
		},
	}
	err := mongo2.FindOne(r.collection.Name(), query, u)
	return u, err
}

// FindByID retrieves a user by their ID (string or ObjectID).
func (r *mongoUserRepository) FindByID(id any) (*user.User, error) {
	u := new(user.User)
	err := mongo2.Get(r.collection.Name(), id, u)
	return u, err
}

// Update modifies an existing user and sets the update timestamp.
func (r *mongoUserRepository) Update(user *user.User) error {
	user.UpdatedAt = time.Now()
	return mongo2.Update(user)
}

// Delete removes a user by their ID after converting it to ObjectID.
func (r *mongoUserRepository) Delete(id any) error {
	objID, err := util.ToObjectID(id)
	if err != nil {
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

func (r *mongoUserRepository) List() (user.Users, error) {
	users := make(user.Users, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{}, &users)
	if err != nil {
		logger.Error("error while fetching users: ", err.Error())
		return nil, err
	}

	return users, nil
}
