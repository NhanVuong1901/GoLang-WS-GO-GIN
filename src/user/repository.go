package user

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection("users")
	if err := CreateUserIndexes(collection); err != nil {
		panic(err)
	}
	return &Repository{
		Collection: collection,
	}
}

func (r *Repository) Create(user *User) error {
	user.CreatedAt = time.Now()

	_, err := r.Collection.InsertOne(context.TODO(), user)
	return err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

func CreateUserIndexes(collection *mongo.Collection) error {
	names, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}
	fmt.Println("Created Indexes:", names)
	return nil
}

func (r *Repository) UpdateByID(userID bson.ObjectID, update bson.M) error {
	_, err := r.Collection.UpdateByID(context.TODO(), userID, bson.M{
		"$set": update,
	})
	return err
}
