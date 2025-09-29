package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/omatheuscaetano/planus-api/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds MongoDB configuration
type Config struct {
	URI             string
	Database        string
	MaxPoolSize     uint64
	MinPoolSize     uint64
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration
	SocketTimeout   time.Duration
}

// MongoDB represents the MongoDB database connection
type MongoDB struct {
	C        *mongo.Client
	DB       *mongo.Database
	config   *Config
}

func NewProductionMongoDb() *MongoDB {
	return NewMongoDB(&Config{
		URI:             "mongodb://" + env.DBUser() + ":" + env.DBPassword() + "@" + env.DBHost() + ":" + env.DBPort() + "/" + env.DBName() + "?authSource=admin",
		Database:        env.DBName(),
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 30 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		SocketTimeout:   30 * time.Second,
	})
}

func NewTestMongoDb() *MongoDB {
	return NewMongoDB(&Config{
		URI:             "mongodb://" + env.DBTestUser() + ":" + env.DBTestPassword() + "@" + env.DBTestHost() + ":" + env.DBTestPort() + "/" + env.DBTestName() + "?authSource=admin",
		Database:        env.DBTestName(),
		MaxPoolSize:     100,
		MinPoolSize:     10,
		MaxConnIdleTime: 30 * time.Minute,
		ConnectTimeout:  10 * time.Second,
		SocketTimeout:   30 * time.Second,
	})
}

func NewMongoDB(config *Config) *MongoDB {
	return &MongoDB{
		config: config,
	}
}

// Connect establishes connection to MongoDB with connection pooling
func (m *MongoDB) Connect(ctx context.Context) error {
	// Set default values if not provided
	if m.config.MaxPoolSize == 0 {
		m.config.MaxPoolSize = 100
	}
	if m.config.MinPoolSize == 0 {
		m.config.MinPoolSize = 5
	}
	if m.config.MaxConnIdleTime == 0 {
		m.config.MaxConnIdleTime = 30 * time.Minute
	}
	if m.config.ConnectTimeout == 0 {
		m.config.ConnectTimeout = 10 * time.Second
	}
	if m.config.SocketTimeout == 0 {
		m.config.SocketTimeout = 30 * time.Second
	}

	// Configure client options with connection pooling
	clientOptions := options.Client().
		ApplyURI(m.config.URI).
		SetMaxPoolSize(m.config.MaxPoolSize).
		SetMinPoolSize(m.config.MinPoolSize).
		SetMaxConnIdleTime(m.config.MaxConnIdleTime).
		SetConnectTimeout(m.config.ConnectTimeout).
		SetSocketTimeout(m.config.SocketTimeout).
		SetRetryWrites(true).   // Enable retryable writes
		SetRetryReads(true).    // Enable retryable reads
		SetHeartbeatInterval(10 * time.Second)

	// Create context with timeout for connection
	connectCtx, cancel := context.WithTimeout(ctx, m.config.ConnectTimeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	m.C = client
	m.DB = client.Database(m.config.Database)

	return nil
}

// Disconnect closes the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.C == nil {
		return nil
	}

	disconnectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := m.C.Disconnect(disconnectCtx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

// Health checks the database connection health
func (m *MongoDB) Health(ctx context.Context) error {
	if m.C == nil {
		return fmt.Errorf("mongodb client is not initialized")
	}

	healthCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := m.C.Ping(healthCtx, readpref.Primary()); err != nil {
		return fmt.Errorf("mongodb health check failed: %w", err)
	}

	return nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}

// GetDatabase returns the MongoDB database instance
func (m *MongoDB) Database() *mongo.Database {
	return m.DB
}

// GetClient returns the MongoDB client instance
func (m *MongoDB) Client() *mongo.Client {
	return m.C
}

// Drop drops the entire database (use with caution!)
func (m *MongoDB) Drop(ctx context.Context) error {
	if m.DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	dropCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := m.DB.Drop(dropCtx); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	log.Printf("Database %s dropped successfully", m.config.Database)
	return nil
}

// CreateIndexes creates indexes for better query performance
// This is a placeholder - you'll need to implement specific indexes for your collections
func (m *MongoDB) CreateIndexes(ctx context.Context) error {
	// Example: Create index for users collection
	// usersCollection := m.GetCollection("users")
	// indexModel := mongo.IndexModel{
	// 	Keys: bson.D{{"email", 1}}, // 1 for ascending, -1 for descending
	// 	Options: options.Index().SetUnique(true),
	// }
	//
	// _, err := usersCollection.Indexes().CreateOne(ctx, indexModel)
	// if err != nil {
	// 	return fmt.Errorf("failed to create user email index: %w", err)
	// }

	log.Println("Database indexes created successfully")
	return nil
}

// Stats returns database statistics
// func (m *MongoDB) Stats(ctx context.Context) (*mongo.DatabaseStats, error) {
// 	if m.DB == nil {
// 		return nil, fmt.Errorf("database is not initialized")
// 	}

// 	statsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	var stats mongo.DatabaseStats
// 	err := m.DB.RunCommand(statsCtx, map[string]interface{}{"dbStats": 1}).Decode(&stats)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get database stats: %w", err)
// 	}

// 	return &stats, nil
// }
