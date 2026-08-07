package group

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

type GroupFilter struct {
	filter.SearchFilter
	ID       string `query:"id"`
	IsDelete string `query:"is_delete"`
	Name     string `query:"name"`
	TenantID string `query:"tenant_id"`
	Key      string `query:"key"`
}

func (filter GroupFilter) GetFilters() map[string]interface{} {
	q := bson.M{}

	filter.SearchKeys = searchKeys
	filter.ApplySearch(q)

	if id, err := primitive.ObjectIDFromHex(filter.ID); err == nil {
		q["_id"] = id
	}

	if delete, err := strconv.ParseBool(filter.IsDelete); err == nil {
		q["is_deleted"] = delete
	}

	if filter.Name != "" {
		q["name"] = filter.Name
	}

	if filter.TenantID != "" {
		q["tenant_id"] = filter.TenantID
	}

	if filter.Key != "" {
		q["key"] = filter.Key
	}

	return q
}
