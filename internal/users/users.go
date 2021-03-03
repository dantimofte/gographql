package users

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	database "gqlexample/v2/internal/pkg/db/mongodb"
	"reflect"
)

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() {
	collection, err := database.GetDBCollection("gql_db", "users")
	if err != nil {
		logrus.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	user.Password = hashedPassword
	resp, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		logrus.Fatal(err)
	}
	newID := resp.InsertedID.(primitive.ObjectID)
	logrus.Info("InsertOne() newID:", newID.Hex())
	logrus.Info("InsertOne() newID type:", reflect.TypeOf(newID))
}

//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (string, error) {

	// 2. get users db collection
	collection, err := database.GetDBCollection("gql_db","users")
	if err != nil {
		return "", err
	}

	// 3. check if user already exists
	var result User
	err = collection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&result)

	if err == nil {
		return "", err
	}

	return result.ID, nil
}


func (user *User) Authenticate() bool {
	collection, err := database.GetDBCollection("gql_db","users")
	if err != nil {
		return false
	}

	var result User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		return false
	}

	return CheckPasswordHash(user.Password, result.Password)
}


//HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}