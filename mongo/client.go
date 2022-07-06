package mongo

import (
	"context"
	"github.com/knadh/koanf"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func ConnectToMongo(config *koanf.Koanf, logger zerolog.Logger, ctx context.Context) (*mongo.Database, error) {
	connectionString, err := connstring.ParseAndValidate(config.String("mongodb.connection-string"))
	if err != nil {
		return nil, err
	}
	if connectionString.Database == "" {
		return nil, NoDatabaseConfigured
	}

	opts := options.Client().ApplyURI(connectionString.String())

	logger = logger.With().Str("component", "mongodb").Logger()
	logger.Debug().Str("connection-string", connectionString.String()).Msg("Connecting to MongoDB")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to MongoDB")
		return nil, err
	}

	logger.Info().Msg("Connected to MongoDB")

	return client.Database(connectionString.Database), nil
}
