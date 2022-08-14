package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// (Testing variables)
var coll *mongo.Collection

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	http.HandleFunc("/ws", serveWs)
}

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
	coll = client.Database("netflix").Collection("movies")

	// (Testing) Show Comments
	// commentData := showComments(coll, title, episode)
	// if commentData != nil {
	// 	comments, err := json.MarshalIndent(commentData, "", "    ")
	// 	if err != nil {
	// 		panic(err)
	// 	} else {
	// 		fmt.Printf("%s \n", comments)
	// 	}
	// }

	setupRoutes()
	http.ListenAndServe(":8080", nil)
}

func addComment(coll *mongo.Collection, author string, comment string, title string, episode string) error {
	if coll == nil {
		fmt.Println("Coll is nil")
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

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	author  string `json:"author"`
	comment string `json:"comment"`
	title   string `json:"title"`
	episode string `json:"episode"`
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		// fmt.Println(string(p))

		var jsonMap map[string]interface{}
		json.Unmarshal(p, &jsonMap)
		author := fmt.Sprintf("%v", jsonMap["author"])
		comment := fmt.Sprintf("%v", jsonMap["comment"])
		title := fmt.Sprintf("%v", jsonMap["title"])
		episode := fmt.Sprintf("%v", jsonMap["episode"])
		fmt.Println(author, comment, title, episode)
		err = addComment(coll, author, comment, title, episode)
		if err != nil {
			fmt.Print(err)
		} else {
			log.Println("Created comment")
		}
	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}
