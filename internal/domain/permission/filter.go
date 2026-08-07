package permission

import (
	"strconv"

	"github.com/yasinsaee/go-user-service/internal/context/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	searchKeys []string = []string{
		"name",
	}
)

type PermissionFilter struct {
	filter.SearchFilter
	ID       string `query:"id"`
	Key      string `query:"key"`
	IsDelete string `query:"is_delete"`
}

func (filter PermissionFilter) GetFilters() map[string]interface{} {
	q := bson.M{}

	filter.SearchKeys = searchKeys
	filter.ApplySearch(q)

	if id, err := primitive.ObjectIDFromHex(filter.ID); err == nil {
		q["_id"] = id
	}

	if filter.Key != "" {
		q["key"] = filter.Key
	}

	if delete, err := strconv.ParseBool(filter.IsDelete); err == nil {
		q["is_deleted"] = delete
	}

	return q
}
