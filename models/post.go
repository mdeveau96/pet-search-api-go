package models

import (
	"context"
	"pet-search-backend-go/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type comment struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`
	
}

type Post struct {
	ID        primitive.ObjectID   `bson:"_id" json:"_id"`
	Title     string               `bson:"title" json:"title"`
	ImageUrl  string               `bson:"imageUrl" json:"imageUrl"`
	Content   string               `bson:"content" json:"content"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	Comments comment `bson:"comments" json:"comments"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

var collection = db.GetClient().Database("petsearch").Collection("posts")

func FindAllPosts() ([]Post, error) {
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	var posts []Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		panic(err)
	}
	return posts, nil
}

func FindPost(postId primitive.ObjectID) (Post, error) {
	filter := bson.D{{Key: "_id", Value: postId}}
	var result Post
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			panic(err)
		}
	}
	return result, nil
}

func (p *Post) Create() (Post, error) {
	newPost := Post{ID: primitive.NewObjectID(), Title: p.Title, ImageUrl: p.ImageUrl, Content: p.Content, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_, err := collection.InsertOne(context.Background(), newPost)
	if err != nil {
		return Post{}, err
	}
	return newPost, nil
}

func (p *Post) Update(updatedPost Post) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: updatedPost.ID}}
	update := bson.D{
		{Key: "$set", Value: Post{ID: p.ID, Title: p.Title, ImageUrl: p.ImageUrl, Content: p.Content, CreatedAt: p.CreatedAt, UpdatedAt: time.Now()}},
	}
	post, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *Post) Delete(postId primitive.ObjectID) {

}
