package models

import (
	"context"
	"pet-search-backend-go/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type member struct {
	UserID User   `bson:"user_id" json:"user_id"`
	Role   string `bson:"role" json:"role"`
}

type Group struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	GroupName   string             `bson:"group_name" json:"group_name"`
	Description string             `bson:"description" json:"description"`
	Members     []member           `bson:"members" json:"members"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

var groupsCollection = db.GetClient().Database("petsearch").Collection("groups")

func updateGroupReturnResult(groupId primitive.ObjectID, filter, update primitive.D) (Group, error) {
	_, err := groupsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return Group{}, err
	}
	group, err := FindGroup(groupId)
	if err != nil {
		return Group{}, err
	}
	return group, nil
}

func FindAllGroups() ([]Group, error) {
	cursor, err := groupsCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return []Group{}, err
	}
	var groups []Group
	if err = cursor.All(context.Background(), &groups); err != nil {
		return []Group{}, err
	}
	return groups, nil
}

func FindGroup(groupId primitive.ObjectID) (Group, error) {
	filter := bson.D{{Key: "_id", Value: groupId}}
	var result Group
	err := groupsCollection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Group{}, err
		}
	}
	return result, nil
}

func (g *Group) Create() (Group, error) {
	newGroup := Group{ID: primitive.NewObjectID(), GroupName: g.GroupName, Description: g.Description, Members: g.Members, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_, err := groupsCollection.InsertOne(context.Background(), newGroup)
	if err != nil {
		return Group{}, err
	}
	return newGroup, nil
}

func (g *Group) Update(updatedGroup Group) (Group, error) {
	filter := bson.D{{Key: "_id", Value: updatedGroup.ID}}
	update := bson.D{
		{Key: "$set", Value: Group{ID: g.ID, GroupName: g.GroupName, Description: g.Description, Members: g.Members, CreatedAt: g.CreatedAt, UpdatedAt: time.Now()}},
	}
	return updateGroupReturnResult(g.ID, filter, update)
}

func (g *Group) Delete() (*mongo.DeleteResult, error) {
	filter := bson.D{{Key: "_id", Value: g.ID}}
	result, err := groupsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return result, err
}
