package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"oss.navercorp.com/metis/metis-server/server/database"
)

type Client struct {
	client *mongo.Client
}

const (
	uri         = "mongodb://localhost:27017"
	dbName      = "metis"
	dialTimeout = 10
)

func NewClient() *Client {
	return &Client{}
}

func Dial(ctx context.Context) (*Client, error) {
	cli := NewClient()
	return cli, cli.Dial(ctx)
}

func (d *Client) Dial(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, dialTimeout*time.Second)
	defer cancel()

	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("Could not connect to MongoDB: %s\n", err.Error())
		return err
	}
	fmt.Println("Connected to MongoDB")

	d.client = client
	return nil
}

func (d *Client) CreateDiagram(ctx context.Context, name string) (*database.Diagram, error) {
	result, err := d.client.Database(dbName).Collection("diagrams").InsertOne(ctx, bson.M{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return &database.Diagram{
		ID:   result.InsertedID.(primitive.ObjectID),
		Name: name,
	}, nil
}

// Close closes the client, releasing any open resources.
func (d *Client) Close(ctx context.Context) error {
	if err := d.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}