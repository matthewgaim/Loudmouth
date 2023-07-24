package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type GetCommentsResponse struct {
	Time    int    `json:"time"`
	VideoID string `json:"videoid"`
}

func GetDBItem(tableName string, svc *dynamodb.DynamoDB, time int) []Item {
	filt := expression.Name("time").GreaterThanEqual(expression.Value(time))
	proj := expression.NamesList(expression.Name("uuid"), expression.Name("time"), expression.Name("comment"), expression.Name("videoid"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %s", err)
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		log.Printf("Query API call failed: %s", err)
	}
	var items []Item
	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}
		items = append(items, item)
	}
	return items
}
