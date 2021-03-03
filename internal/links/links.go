package links

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	database "gqlexample/v2/internal/pkg/db/mongodb"
	"gqlexample/v2/internal/users"
	"log"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

// #1
type dbLink struct {
	Title   string
	Address string
	UserId  string
}

//#2
func (link Link) Save() string {

	collection, err := database.GetDBCollection("gql_db", "links")
	if err != nil {
		log.Fatal(err)
	}

	l := dbLink{
		Title:   link.Title,
		Address: link.Address,
		UserId:  link.User.ID,
	}
	resp, err := collection.InsertOne(context.TODO(), l)
	newID := resp.InsertedID.(primitive.ObjectID)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return newID.Hex()
}

func GetAll() ([]Link, error) {
	// passing bson.D{{}} matches all documents in the collection
	//filter := bson.D{{}}
	filter := bson.D{{"$lookup", bson.D{{"from", "users"}, {"localField", "userId"}, {"foreignField", "_id"}, {"as", "user"}}}}

	return filterTasks(filter)
}


func filterTasks(filter interface{}) ([]Link, error) {
	ctx := context.TODO()
	collection, err := database.GetDBCollection("gql_db", "links")
	if err != nil {
		log.Fatal(err)
	}

	var links []Link
	unwind := bson.D{{"unwind", "user"}}
	cur, err := collection.Aggregate(ctx, mongo.Pipeline{filter.(bson.D), unwind})
	//cur, err := collection.Find(ctx, filter)
	if err != nil {
		return links, err
	}

	for cur.Next(ctx) {
		var l Link
		err := cur.Decode(&l)
		if err != nil {
			return links, err
		}

		links = append(links, l)
	}

	if err := cur.Err(); err != nil {
		return links, err
	}

	// once exhausted, close the cursor
	_ = cur.Close(ctx)

	if len(links) == 0 {
		return links, mongo.ErrNoDocuments
	}

	return links, nil

}
