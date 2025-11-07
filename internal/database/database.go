package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db *mongo.Client
}

var (
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	db *mongo.Client
	err error
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

func (m *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func InitDB() {
	dsn := "mongodb://localhost:27017"
	if os.Getenv("DB_CONNECTION") != "" {
		dsn = os.Getenv("DB_CONNECTION")
	}

	db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("Failed to connect to the Database")
	}
	log.Printf("🚀 Connected Successfully to the Database: %v\n", dsn)
}

func GetDatabase() *mongo.Database {
	dbName := "invoice_db"
	if os.Getenv("DB_NAME") != "" {
		dbName = os.Getenv("DB_NAME")
	}
	return db.Database(dbName)
}
