package user

import (
	"strconv"

	"github.com/yasinsaee/go-user-service/internal/context/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	searchKeys []string = []string{
		"first_name",
		"last_name",
		"phone_number",
		"email",
		"username",
	}
)

type UserFilter struct {
	filter.SearchFilter
	ID          string `query:"id"`
	Username    string `query:"username"`
	Email       string `query:"email"`
	Phonenumber string `query:"phonenumber"`
	IsDelete    string `query:"is_delete"`
	TenantID    string `query:"tenant_id"`
	Group       string `query:"group"`
	NeGroup     string `query:"ne_group"`
}

func (filter UserFilter) GetFilters() map[string]interface{} {
	q := bson.M{}

	filter.SearchKeys = searchKeys
	filter.ApplySearch(q)

	if id, err := primitive.ObjectIDFromHex(filter.ID); err == nil {
		q["_id"] = id
	}

	if delete, err := strconv.ParseBool(filter.IsDelete); err == nil {
		q["is_deleted"] = delete
	}

	if filter.Username != "" {
		q["username"] = filter.Username
	}

	if filter.Email != "" {
		q["email"] = filter.Email
	}

	if filter.Phonenumber != "" {
		q["email"] = filter.Email
	}

	if filter.TenantID != "" {
		q["tenant_id"] = filter.TenantID
	}

	if filter.Group != "" {
		q["group"] = filter.Group
	}

	if filter.NeGroup != "" {
		q["group"] = bson.M{
			"$ne": filter.NeGroup,
		}
	}

	return q
}
