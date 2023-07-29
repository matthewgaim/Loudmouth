package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	envCreds := credentials.NewEnvCredentials()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: envCreds,
	})
	if err != nil {
		log.Fatal("Error connecting to AWS")
	}
	svc := dynamodb.New(sess)

	HandleHTTPRequests(sess, svc)
}

func HandleHTTPRequests(sess *session.Session, svc *dynamodb.DynamoDB) {

	http.HandleFunc("/addcomment", func(w http.ResponseWriter, r *http.Request) {
		addcomment(sess, svc, w, r)
	})

	http.HandleFunc("/getcomments", func(w http.ResponseWriter, r *http.Request) {
		getComments(svc, w, r)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func addcomment(sess *session.Session, svc *dynamodb.DynamoDB, w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println("/addcomment called")
	switch r.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "POST":
		if svc != nil {
			var item Item
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			createmydbtable(item.VideoID, sess, svc)
			err = svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: aws.String(item.VideoID)})
			if err != nil {
				fmt.Printf("Got error waiting for table %s to exist:\n", item.VideoID)
				fmt.Println(err.Error())
			} else {
				AddDBItem(item.VideoID, item, svc)
			}
		}
	}
}

func getComments(svc *dynamodb.DynamoDB, w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Println("/getcomments called")
	switch r.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	case "POST":
		if svc == nil {
			http.Error(w, "Problem connecting to AWS", http.StatusBadRequest)
		} else {
			var response GetCommentsResponse
			err := json.NewDecoder(r.Body).Decode(&response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				items := GetDBItem(response.VideoID, svc, response.Time)
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(items)
			}
		}
	default:
		w.WriteHeader(http.StatusForbidden)
	}
}

func createmydbtable(tableName string, sess *session.Session, svc *dynamodb.DynamoDB) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("uuid"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("time"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("uuid"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("time"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}
	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Printf("Got error calling CreateTable for table '%s': %s\n", tableName, err.Error())
	}
}

func enableCors(w *http.ResponseWriter) {
	header := (*w).Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}
