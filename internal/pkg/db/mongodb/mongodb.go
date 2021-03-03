package mongodb

/*
https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.4.6#readme-usage
https://www.digitalocean.com/community/tutorials/how-to-use-go-with-mongodb-using-the-mongodb-go-driver
https://levelup.gitconnected.com/working-with-mongodb-using-golang-754ead0c10c

*/


import (
"context"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"sync"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var MongoDB *mongo.Client

//Used to execute client creation procedure only once.
var mongoOnce sync.Once
//I have used below constants just to hold required database config's.

//GetMongoClient - Return mongodb connection to work with
func InitMongoDBClient(source,username,password string) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// set credentials
		mdbCredentials := options.Credential{
			AuthSource: source,
			Username:   username,
			Password:   password,
		}
		// Set client options
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(mdbCredentials)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		check(err)
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		check(err)
		MongoDB = client
	})
}

func GetDBCollection(dbName, collectionName string) (*mongo.Collection, error) {
	collection := MongoDB.Database(dbName).Collection(collectionName)
	return collection, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
