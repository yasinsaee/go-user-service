package mongo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/structs"
	"github.com/yasinsaee/go-user-service/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI           string `json:"uri" yaml:"uri"`
	DB            string `json:"db" yaml:"db"`
	Username      string `json:"username" yaml:"username"`
	Password      string `json:"password" yaml:"password"`
	AuthMechanism string `json:"auth_mechanism" yaml:"auth_mechanism"`
	AuthSource    string `json:"auth_source" yaml:"auth_source"`
}

var (
	ErrorModelID  = errors.New("mogo: model id type error")
	DB            *MongoDB
	DefaultConfig = Config{
		URI:           "mongodb://localhost:27017",
		DB:            "ptm",
		Username:      "",
		Password:      "",
		AuthMechanism: "",
		AuthSource:    "",
	}
)

func Init(cfg Config) {
	var err error
	DB, err = connection(cfg)
	if err != nil {
		logger.Error("Error connecting to MongoDB, " + err.Error())
		for i := 0; i < 3; i++ {
			DB, err = connection(cfg)
			if err == nil {
				logger.Info("Connected to MongoDB")
				break
			} else {
				logger.Error("Error connecting to MongoDB, " + err.Error())
			}
		}
	} else {
		logger.Info("Connected to MongoDB")
	}
}

type MongoDB struct {
	*mongo.Database
}

// DBContext main connection struct
type DBContext struct {
	Session  *mongo.Session
	Database *mongo.Database
}

func connection(cfg Config) (*MongoDB, error) {
	opt := options.Client()

	opt.ApplyURI(cfg.URI)
	if cfg.Username != "" {
		opt.SetAuth(options.Credential{
			Username:      cfg.Username,
			Password:      cfg.Password,
			AuthMechanism: cfg.AuthMechanism,
			AuthSource:    cfg.AuthSource,
		})
	}
	client, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		logger.Error("NewClinet to MongoDB, " + err.Error())
		return nil, err
	}

	db := client.Database(cfg.DB)

	_, err = db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		logger.Error("MongoDB, " + err.Error())
		return nil, err
	}

	return &MongoDB{db}, nil

}

// Create a document in DBContext
func Create(model interface{}) error {
	col := DB.Collection(colName(model))
	setID(model)
	_, err := col.InsertOne(context.Background(), &model)
	return err
}

// Collection return mgo collection from model
func Collection(model interface{}) *mongo.Collection {
	return DB.Database.Collection(colName(model))
}

// Collection return mgo collection from model
func CollectionString(model string) *mongo.Collection {
	return DB.Database.Collection(model)
}

// Update a Document
func Update(model interface{}) error {
	collection := DB.Collection(colName(model))
	ctx := context.Background()
	id, err := getID(model)
	if err != nil {
		return err
	}
	query := bson.M{"_id": id}
	fieldsUpdate := parseBson(model)
	if _, err := collection.UpdateOne(ctx, query, bson.M{"$set": fieldsUpdate}); err != nil {
		return err
	}
	return Get(collection.Name(), id, model)
}

func UpdateMany(collection string, filter bson.M, update bson.M) error {
	ctx := context.Background()
	if _, err := DB.Collection(collection).UpdateMany(ctx, filter, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func UpdateOne(collection string, filter bson.M, update bson.M) error {
	ctx := context.Background()
	if _, err := DB.Collection(collection).UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// Find start generating query
// order of options are:  limit, page, sort
func Find(collection string, q bson.M, result interface{}, params ...interface{}) error {
	col := DB.Collection(collection)
	sort := ""
	var page int64
	var limit int64
	opt := new(options.FindOptions)
	if len(params) > 0 {
		limit = int64(params[0].(int))
		if limit > 0 {
			opt.Limit = &limit
		}
	}
	if len(params) > 1 {
		page = int64(params[1].(int))
		if page > 0 {
			skip := (page - 1) * limit
			opt.Skip = &skip
		}
	}
	if len(params) > 2 {
		sort = params[2].(string)
		// @TODO some ducument dont have created_at
		// if sort == "" {
		// 	sort = "-created_at"
		// }
		if sort != "" {
			if strings.HasPrefix(sort, "-") {
				sort = strings.ReplaceAll(sort, "-", "")
				opt.SetSort(bson.D{{Key: sort, Value: -1}})
			} else {
				opt.SetSort(bson.D{{Key: sort, Value: 1}})
			}
		}

	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := col.Find(ctx, q, opt)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(context.Background(), result)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func FindOne(collection string, q bson.M, result interface{}) error {
	col := DB.Collection(collection)
	return col.FindOne(context.Background(), q).Decode(result)
}

func Count(collection string, filter bson.M) (int, error) {
	ctx := context.Background()
	c, err := DB.Collection(collection).CountDocuments(ctx, filter)
	return int(c), err
}

func RemoveOne(collection string, filter bson.M) error {
	ctx := context.Background()
	_, err := DB.Collection(collection).DeleteOne(ctx, filter)
	return err
}

func RemoveMany(collection string, filter bson.M) error {
	ctx := context.Background()
	_, err := DB.Collection(collection).DeleteMany(ctx, filter)
	return err
}

func Get(collection string, id interface{}, resultObj interface{}) error {
	if id == nil {
		return errors.New("invalid ID")
	}
	ctx := context.Background()
	var err error
	var objID primitive.ObjectID
	if val, ok := id.(string); ok {
		objID, err = primitive.ObjectIDFromHex(val)
		if err != nil {
			return err
		}
	} else {
		objID = id.(primitive.ObjectID)
	}
	if objID.IsZero() {
		return errors.New("invalid ID")
	}
	singleResult := DB.Collection(collection).FindOne(ctx, bson.M{"_id": objID})
	if singleResult.Err() != nil {
		return singleResult.Err()
	}
	return singleResult.Decode(resultObj)
}

func setID(model interface{}) {
	m := structs.Map(model)
	var keyID string
	if _, ok := m["Id"]; ok {
		m["Id"] = primitive.NewObjectID()
		keyID = "Id"
	}
	if _, ok := m["ID"]; ok {
		m["ID"] = primitive.NewObjectID()
		keyID = "ID"
	}
	s := structs.New(model)
	field := s.Field(keyID)
	field.Set(m[keyID])

}

func getID(model interface{}) (primitive.ObjectID, error) {
	m := structs.Map(model)
	var (
		idInterface interface{}
		id          primitive.ObjectID
		ok          bool
	)
	if val, ok := m["Id"]; ok {
		idInterface = val
	}
	if val, ok := m["ID"]; ok {
		idInterface = val
	}
	id, ok = idInterface.(primitive.ObjectID)
	if !ok {
		return id, ErrorModelID
	}
	if id.IsZero() {
		return id, ErrorModelID
	}
	return id, nil
}

func parseBson(model interface{}) bson.M {
	b, _ := bson.Marshal(model)
	var body bson.M
	bson.Unmarshal(b, &body)
	return body
}

type coller interface {
	CollectionName() string
}

func colName(model interface{}) string {
	if c, ok := model.(coller); ok {
		return c.CollectionName()
	}
	tmp := fmt.Sprintf("%T", model)
	tmp = strings.Replace(tmp, "*", "", -1)
	tmp = strings.Replace(tmp, "]", "", -1)
	tmp = strings.Replace(tmp, "[", "", -1)
	ts := strings.Split(tmp, ".")
	if len(ts) < 2 {
		return toSnake(tmp)
	}
	return toSnake(ts[1])
}

func toSnake(s string) string {
	var (
		res  string
		last int
	)
	ls := []rune(s)
	for i, char := range ls {
		if (i == 0 || !unicode.IsUpper(char)) && i+1 != len(s) {
			continue
		}
		if i+1 != len(s) {
			res += strings.ToLower(s[last:i]) + "_"
		} else {
			res += strings.ToLower(s[last : i+1])
		}
		last = i
	}
	return res
}

func FindOneWithOptions(collection string, q bson.M, result interface{}, options options.FindOneOptions) error {
	col := DB.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return col.FindOne(ctx, q, &options).Decode(result)
}

func FindWithOptions(collection string, q bson.M, result interface{}, opt options.FindOptions, params ...interface{}) error {
	col := DB.Collection(collection)
	sort := ""
	var page int64
	var limit int64
	if len(params) > 0 {
		if l, ok := params[0].(int); ok {
			limit = int64(l)
		} else if lstr, ok := params[0].(string); ok {
			l, _ = strconv.Atoi(lstr)
			limit = int64(l)
		}
		if limit > 0 {
			opt.Limit = &limit
		}
	}
	if len(params) > 1 {
		if l, ok := params[1].(int); ok {
			page = int64(l)
		} else if lstr, ok := params[1].(string); ok {
			l, _ = strconv.Atoi(lstr)
			page = int64(l)
		}
		if page > 0 {
			skip := (page - 1) * limit
			opt.Skip = &skip
		}
	}
	if len(params) > 2 {
		sort = params[2].(string)
		// @TODO some ducument dont have created_at
		// if sort == "" {
		// 	sort = "-created_at"
		// }
		if sort != "" {
			if strings.HasPrefix(sort, "-") {
				sort = strings.ReplaceAll(sort, "-", "")
				opt.SetSort(bson.D{{Key: sort, Value: -1}})
			} else {
				opt.SetSort(bson.D{bson.E{Key: sort, Value: 1}})
			}
		}

	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := col.Find(ctx, q, &opt)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	err = cursor.All(context.Background(), result)
	if err != nil {
		return err
	}
	return nil
}
