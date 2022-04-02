package mongodb

import (
	"context"
	"time"

	"errors"

	paginate "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	pb "choirulanwar/user-svc/internal/core/domain"
	"choirulanwar/user-svc/internal/core/ports"
	"choirulanwar/user-svc/pkg/constants"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// Create mongo client
func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, err
}

// Create repository
func NewMongoRepository(
	mongoURL, mongoDB string,
	mongoTimeout int,
) (ports.UserRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}

	collection := client.Database(repo.database).Collection("users")

	_, err = collection.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		}, {
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		}},
	)
	if err != nil {
		return nil, err
	}

	repo.client = client

	return repo, nil
}

// Find user by id
func (r *mongoRepository) Find(id string) (*pb.FindRes, error) {
	user := &pb.FindRes{}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{"$or": []bson.M{{"id": id}, {"email": id}, {"username": id}}}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{"password": 0})

	err := collection.FindOne(ctx, filter, findOptions).Decode(&user.Data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, constants.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// Store new user
func (r *mongoRepository) Store(data *pb.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("users")

	_, err := collection.InsertOne(
		ctx,
		&data,
	)
	if err != nil {
		exists := mongo.IsDuplicateKeyError(err)
		if exists {
			return constants.ErrExists
		}
		return err
	}

	return nil
}

// Update user
func (r *mongoRepository) Update(id string, data *pb.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{"id": id}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := collection.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return constants.ErrNotFound
		}
		return err
	}

	return nil
}

// Find all user
func (r *mongoRepository) FindAll(page int64, limit int64, orderBy string, orderType string) (*pb.FindAllRes, error) {
	var users []*pb.User

	var pageVal int64 = page
	if page < 1 {
		pageVal = 1
	}
	var limitVal int64 = limit
	if limit < 1 {
		limitVal = 1
	}
	var orderByVal string = "created_at"
	var orderTypeVal int64 = 1

	switch orderBy {
	case "CREATED_AT":
		orderByVal = "created_at"

	case "UPDATED_AT":
		orderByVal = "updated_at"

	default:
		orderByVal = "created_at"
	}

	switch orderType {
	case "ASC":
		orderTypeVal = 1

	case "DESC":
		orderTypeVal = -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{}
	projection := bson.D{
		{Key: "password", Value: 0},
	}

	paginatedData, err := paginate.New(collection).Context(ctx).Limit(limitVal).Page(pageVal).Sort(orderByVal, orderTypeVal).Select(projection).Filter(filter).Decode(&users).Find()
	if err != nil {
		return nil, err
	}

	results := pb.FindAllRes{
		TotalDatas: paginatedData.Pagination.Total,
		Limit:      paginatedData.Pagination.PerPage,
		Page:       paginatedData.Pagination.Page,
		TotalPages: paginatedData.Pagination.TotalPage,
		Datas:      users,
		NextPage:   paginatedData.Pagination.Next,
		PrevPage:   paginatedData.Pagination.Prev,
	}

	return &results, nil
}

// Delete user
func (r *mongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("users")
	filter := bson.M{"id": id}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return constants.ErrNotFound
		}
		return err
	}

	return nil
}
