package models

import (
	"context"
	"pet-search-backend-go/db"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	ID        primitive.ObjectID   `bson:"_id" json:"_id"`
	Creator   primitive.ObjectID   `bson:"user" json:"user"`
	Content   string               `bson:"content" json:"content"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

type Post struct {
	ID        primitive.ObjectID   `bson:"_id" json:"_id"`
	Title     string               `bson:"title" json:"title"`
	ImageUrl  string               `bson:"imageUrl" json:"imageUrl"`
	Content   string               `bson:"content" json:"content"`
	Creator   primitive.ObjectID   `bson:"creator" json:"creator"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	Comments  []Comment            `bson:"comments" json:"comments"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

var postsCollection = db.GetClient().Database("petsearch").Collection("posts")

func FindAllPosts() ([]Post, error) {
	cursor, err := postsCollection.Find(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	var posts []Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func FindPost(postId primitive.ObjectID) (Post, error) {
	filter := bson.D{{Key: "_id", Value: postId}}
	var result Post
	err := postsCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			panic(err)
		}
	}
	return result, nil
}

func (p *Post) Create() (Post, error) {
	newPost := Post{ID: primitive.NewObjectID(), Title: p.Title, ImageUrl: p.ImageUrl, Content: p.Content, Creator: p.Creator, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_, err := postsCollection.InsertOne(context.Background(), newPost)
	if err != nil {
		return Post{}, err
	}
	return newPost, nil
}

func (p *Post) Update(updatedPost Post) (Post, error) {
	filter := bson.D{{Key: "_id", Value: updatedPost.ID}}
	update := bson.D{
		{Key: "$set", Value: Post{ID: p.ID, Title: p.Title, ImageUrl: p.ImageUrl, Content: p.Content, Creator: p.Creator, CreatedAt: p.CreatedAt, UpdatedAt: time.Now()}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func Delete(postId primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.D{{Key: "_id", Value: postId}}
	result, err := postsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *Post) Like(userId primitive.ObjectID) (Post, error) {
	filter := bson.D{{Key: "_id", Value: p.ID}}
	var userLikes []primitive.ObjectID
	if !(slices.Contains(p.Likes, userId)) {
		userLikes = append(p.Likes, userId)
	} else {
		for _, id := range p.Likes {
			if id != userId {
				userLikes = append(userLikes, id)
			}
		}
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "likes", Value: userLikes}}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p *Post) AddComment(comment Comment) (Post, error) {
	newComment := Comment{ID: primitive.NewObjectID(), Creator: comment.Creator, Content: comment.Content, Likes: comment.Likes, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	commentsList := append(p.Comments, newComment)
	filter := bson.D{{Key: "_id", Value: p.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "comments", Value: commentsList}}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p *Post) LikeComment(commentId, userId primitive.ObjectID) (Post, error) {
	filter := bson.D{{Key: "_id", Value: p.ID}}
	for index, c := range p.Comments {
		if c.ID == commentId {
			var userCommentLikes []primitive.ObjectID
			if !(slices.Contains(c.Likes, userId)) {
				userCommentLikes = append(c.Likes, userId)
			} else {
				for _, id := range c.Likes {
					if id != userId {
						userCommentLikes = append(userCommentLikes, id)
					}
				}
			}
			c = Comment{ID: c.ID, Creator: c.Creator, Content: c.Content, Likes: userCommentLikes, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
			p.Comments[index] = c
		}
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "comments", Value: p.Comments}}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p *Post) UpdateComment(comment Comment) (Post, error) {
	filter := bson.D{{Key: "_id", Value: p.ID}}
	for index, c := range p.Comments {
		if c.ID == comment.ID {
			c = Comment{ID: comment.ID, Creator: c.Creator, Content: comment.Content, Likes: c.Likes, CreatedAt: c.CreatedAt, UpdatedAt: time.Now()}
			p.Comments[index] = c
		}
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "comments", Value: p.Comments}}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p *Post) DeleteComment(commentId primitive.ObjectID) (Post, error) {
	filter := bson.D{{Key: "_id", Value: p.ID}}
	var newCommentsList []Comment
	for _, c := range p.Comments {
		if c.ID != commentId {
			newCommentsList = append(newCommentsList, c)
		}
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "comments", Value: newCommentsList}}},
	}
	_, err := postsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Post{}, err
	}
	post, err := FindPost(p.ID)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}
