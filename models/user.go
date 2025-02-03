package models

import (
	"context"
	"pet-search-backend-go/db"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Login struct {
	Email    string
	Password string
}

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email" json:"email"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Password    string             `bson:"password" json:"password"`
	Posts       []Post             `bson:"posts" json:"posts"`
	MemberOf    []Group            `bson:"member_of" json:"member_of"`
}

var usersCollection = db.GetClient().Database("petsearch").Collection("users")

func FindUserByEmail(email string) (User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	var result User
	err := usersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return User{}, nil
	}
	return result, err
}

func (u *User) AddUser() (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	u.Password = string(hashedPassword)
	newUser := User{ID: primitive.NewObjectID(), Username: u.Username, Email: u.Email, PhoneNumber: u.PhoneNumber, Password: u.Password, Posts: u.Posts, MemberOf: u.MemberOf}
	_, err = usersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}
