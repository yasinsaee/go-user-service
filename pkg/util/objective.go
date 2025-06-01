package util

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectID(id interface{}) (primitive.ObjectID, error) {
	var (
		objID primitive.ObjectID
		err   error
	)
	if val, ok := id.(string); ok {
		objID, err = primitive.ObjectIDFromHex(val)
		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	return objID, nil
}
