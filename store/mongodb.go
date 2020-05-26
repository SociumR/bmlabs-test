package store

import (
	"context"

	"github.com/SociumR/bmlabs-test/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionUser = "user"
const CollectionGame = "user_games"

// Store ...
type Store struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Store {
	return &Store{
		db: db,
	}
}

// Insert ...
func (s *Store) Insert(m interface{}, c string) (string, error) {

	res, err := s.db.Collection(c).InsertOne(context.Background(), m)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// InsertMany ...
func (s *Store) InsertMany(c string, m ...interface{}) error {

	if _, err := s.db.Collection(c).InsertMany(context.Background(), m); err != nil {
		return err
	}

	return nil
}

func (s *Store) Group(c string, group bson.M) ([]primitive.M, error) {
	l := []bson.M{
		group,
	}

	var models []bson.M

	cur, err := s.db.Collection(c).Aggregate(context.Background(), l, &options.AggregateOptions{})

	if err != nil {
		return models, err
	}

	for cur.Next(context.Background()) {
		var result primitive.M
		err := cur.Decode(&result)

		if err != nil {
			return models, err
		}

		models = append(models, result)

	}
	if err := cur.Err(); err != nil {
		return models, err
	}

	return models, nil
}

// Count ...
func (s *Store) Count(limit, pos int64, args map[string]interface{}, c string) (int64, error) {

	opt := &options.CountOptions{}

	if limit > 0 {
		opt.Limit = &limit
	}

	if pos > 0 {
		opt.Skip = &pos
	}

	return s.db.Collection(c).CountDocuments(context.Background(), helpers.MapStringToBSonD(args), opt)

}

// Find ...
func (s *Store) Find(c string, args primitive.D, opt *options.FindOptions) ([]primitive.M, error) {

	var models []primitive.M

	cur, err := s.db.Collection(c).Find(context.Background(), args, opt)

	if err != nil {
		return models, err
	}

	for cur.Next(context.Background()) {
		var result primitive.M
		err := cur.Decode(&result)

		if err != nil {
			return models, err
		}

		models = append(models, result)

	}
	if err := cur.Err(); err != nil {
		return models, err
	}

	return models, nil
}
