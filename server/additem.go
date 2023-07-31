package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	VideoID string `json:"videoid"`
	UUID    string `json:"uuid"`
	Time    int    `json:"time"`
	Comment string `json:"comment"`
}

func AddDBItem(tableName string, item Item, svc *dynamodb.DynamoDB) {
	av, err := dynamodbattribute.MarshalMap(item)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("81153184"),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatal("Got error calling PutItem", err.Error())
	}
	fmt.Println("Successfully added item to table")
}
