package repository

import (
	"time"

	"github.com/yasinsaee/go-user-service/internal/domain/otp"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	mongo2 "github.com/yasinsaee/go-user-service/pkg/mongo"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoOTPRepository implements the OTPRepository interface using MongoDB.
type mongoOTPRepository struct {
	collection *mongo.Collection
}

// NewMongoOTPRepository creates a new MongoDB-based OTP repository.
func NewMongoOTPRepository(db *mongo.Database, collectionName string) otp.OTPRepository {
	return &mongoOTPRepository{
		collection: db.Collection(collectionName),
	}
}

// Create inserts a new OTP into the database and sets timestamps.
func (r *mongoOTPRepository) Create(o *otp.Otp) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return mongo2.Create(o)
}

// FindByID retrieves an OTP by its ID.
func (r *mongoOTPRepository) FindByID(id any) (*otp.Otp, error) {
	o := new(otp.Otp)
	err := mongo2.Get(r.collection.Name(), id, o)
	return o, err
}

// FindByName retrieves an OTP by a name field (if your model supports it).
func (r *mongoOTPRepository) FindByName(name string) (*otp.Otp, error) {
	o := new(otp.Otp)
	query := bson.M{"receiver": name}
	err := mongo2.FindOne(r.collection.Name(), query, o)
	return o, err
}

// Update updates an existing OTP record.
func (r *mongoOTPRepository) Update(o *otp.Otp) error {
	o.UpdatedAt = time.Now()
	return mongo2.Update(o)
}

// Delete removes an OTP entry by ID.
func (r *mongoOTPRepository) Delete(id any) error {
	objID, err := util.ToObjectID(id)
	if err != nil {
		logger.Error("error while deleting otp:", err.Error())
		return err
	}
	return mongo2.RemoveOne(r.collection.Name(), bson.M{"_id": objID})
}

// List returns all stored OTP entries.
func (r *mongoOTPRepository) List() (otp.Otps, error) {
	list := make(otp.Otps, 0)
	err := mongo2.Find(r.collection.Name(), bson.M{}, &list)
	if err != nil {
		logger.Error("error while fetching otps:", err.Error())
		return nil, err
	}
	return list, nil
}
