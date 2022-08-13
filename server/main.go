package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// (Testing variables)
var title = "Trailer Park Boys"
var episode = "3"
var username = "CoolioSchmoolio"
var message = "Did ricky just say that!?"

func main() {
	// Get Secrets
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	// Connecting to MongoDB
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.lldfuec.mongodb.net/?retryWrites=true&w=majority", user, password)
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("netflix").Collection("movies")

	// (Testing) Show & Add Comments
	commentData := showComments(coll, title, episode)
	if commentData != nil {
		comments, err := json.MarshalIndent(commentData, "", "    ")
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s \n", comments)
		}
	}

	err = addComment(coll, username, message, title, episode)
	if err != nil {
		fmt.Print(err)
	}
}

func addComment(coll *mongo.Collection, author string, comment string, title string, episode string) error {
	if coll == nil {
		return nil
	}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.M{"title": title, "episode": episode}
	newComment := bson.D{{Key: "author", Value: author}, {Key: "message", Value: comment}, {Key: "likes", Value: 1}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "comments", Value: newComment}}}}
	var updatedDocument bson.M
	err := coll.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	} else {
		return nil
	}
}

func showComments(coll *mongo.Collection, title string, episode string) bson.M {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "title", Value: title}, {Key: "episode", Value: episode}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title: %s\n", title)
		return nil
	} else if err != nil {
		panic(err)
	}
	return result
}
