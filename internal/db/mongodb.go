package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"time"
)

type Database struct {
	Client     *mongo.Client
	collection *mongo.Collection
}

type RefreshToken struct {
	Guid    string    `json:"guid" bson:"_id"`
	Refresh string    `json:"refresh" bson:"refresh"`
	Time    time.Time `json:"time" bson:"time"`
}

func NewClient(ctx context.Context, url string) (*Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	collection := client.Database("auth").Collection("refreshTokens")

	fmt.Println("Successfully connected to DB")

	return &Database{
		Client:     client,
		collection: collection,
	}, nil
}

func (d *Database) PostRefreshToken(ctx context.Context, guid, refresh string) error {
	expTime := time.Now().Add(1 * time.Hour)

	refreshToken := RefreshToken{
		Guid:    guid,
		Refresh: refresh,
		Time:    expTime,
	}

	_, err := d.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) ReadRefreshToken(ctx context.Context, guid string) (*RefreshToken, error) {

	filter := bson.D{{Key: "_id", Value: guid}}
	var refreshToken *RefreshToken

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return refreshToken, fmt.Errorf("failed to find one refresh token by id : %s due to error: %w", guid, result.Err())
	}

	if err := result.Decode(&refreshToken); err != nil {
		return refreshToken, fmt.Errorf("failed to decode refresh token by id : %s due to error: %w", guid, err)
	}

	return refreshToken, nil
}

func (d *Database) UpdateRefreshToken(ctx context.Context, guid, refresh string) error {

	filter := bson.D{{Key: "_id", Value: guid}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "refresh", Value: refresh},
		{Key: "time", Value: time.Now().Add(1 * time.Hour)},
	}},
	}

	if _, err := d.collection.UpdateOne(ctx, filter, update); err != nil {
		slog.Info("msg", err)
		return err
	}
	return nil
}
