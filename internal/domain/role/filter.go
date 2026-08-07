package role

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

type RoleFilter struct {
	filter.SearchFilter
	ID       string `query:"id"`
	IsDelete string `query:"is_delete"`
	Key      string `query:"key"`
	Scope    string `query:"scope"`
}

func (filter RoleFilter) GetFilters() map[string]interface{} {
	q := bson.M{}

	filter.SearchKeys = searchKeys
	filter.ApplySearch(q)

	if id, err := primitive.ObjectIDFromHex(filter.ID); err == nil {
		q["_id"] = id
	}

	if delete, err := strconv.ParseBool(filter.IsDelete); err == nil {
		q["is_deleted"] = delete
	}

	if filter.Key != "" {
		q["key"] = filter.Key
	}

	if filter.Scope != "" {
		q["scope"] = filter.Scope
	}

	return q
}
