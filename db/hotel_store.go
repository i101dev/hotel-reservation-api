package db

import (
	"context"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	Insert(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(ctx context.Context, id string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	DBNAME := os.Getenv(MongoDBEnvName)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Update(ctx context.Context, filter Map, update Map) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {

	resp, err := s.coll.InsertOne(ctx, hotel)

	if err != nil {
		return nil, err
	}

	hotel.ID = resp.InsertedID.(primitive.ObjectID)

	return hotel, nil
}
func (s *MongoHotelStore) GetHotels(ctx context.Context, filter Map, p *Pagination) ([]*types.Hotel, error) {

	opts := options.FindOptions{}
	opts.SetSkip((p.Page - 1) * p.Limit)
	opts.SetLimit(p.Limit)

	resp, err := s.coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel

	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}
func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {

	var hotel types.Hotel

	if err := s.coll.FindOne(ctx, Map{"_id": id}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}
