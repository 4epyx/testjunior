package mongorepo

import (
	"context"

	"github.com/4epyx/testtask/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRefreshTokenRepository implements the RefreshTokenRepository and may be used for work with MongoDB
type MongoRefreshTokenRepository struct {
	collection *mongo.Collection
}

// NewRefreshTokenRepository returns a new MongoRefreshTokenRepository
// collection is a collection of refresh tokens in database
func NewMongoRTRepo(collection *mongo.Collection) MongoRefreshTokenRepository {
	return MongoRefreshTokenRepository{
		collection: collection,
	}
}

// GetTokenById returns refresh token's data from the database and error if something goes wrong
// tokenId - UUID of the token
func (r MongoRefreshTokenRepository) GetTokenById(ctx context.Context, tokenId string) (model.RefreshToken, error) {
	filter := bson.D{{Key: "_id", Value: tokenId}}
	res := model.RefreshToken{}
	err := r.collection.FindOne(ctx, filter).Decode(&res)
	return res, err
}

func (r MongoRefreshTokenRepository) CreateToken(ctx context.Context, token model.RefreshToken) error {
	_, err := r.collection.InsertOne(ctx, &token)

	return err
}

func (r MongoRefreshTokenRepository) DeleteToken(ctx context.Context, tokenId string) error {
	filter := bson.D{{Key: "_id", Value: tokenId}}
	_, err := r.collection.DeleteOne(ctx, filter)

	return err
}
