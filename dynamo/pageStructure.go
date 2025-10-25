package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutStructure(client *dynamodb.Client, tableName string, item PageStructureItem) error {
	av, err := attributevalue.MarshalMap(item)

	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	return err
}

func GetStructureBydomain(client *dynamodb.Client, tableName string, domain string) ([]PageStructureItem, error) {

	out, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("domain = :d"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":d": &types.AttributeValueMemberS{Value: domain},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var records []PageStructureItem
	err = attributevalue.UnmarshalListOfMaps(out.Items, &records)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return records, nil
}

func IncrementStructureCountByURL(client *dynamodb.Client, tableName string, domain string, url string) error {

	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"domain": &types.AttributeValueMemberS{Value: domain},
			"url":    &types.AttributeValueMemberS{Value: url},
		},
		UpdateExpression: aws.String("ADD #c :inc"),
		ExpressionAttributeNames: map[string]string{
			"#c": "count",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return fmt.Errorf("failed to increment count for %s: %w", url, err)
	}

	return nil
}

func ExistsStructureByDomainAndURL(client *dynamodb.Client, tableName string, domain string, url string) (bool, error) {

	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"domain": &types.AttributeValueMemberS{Value: domain},
			"url":    &types.AttributeValueMemberS{Value: url},
		},
		ProjectionExpression: aws.String("domain"),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}

	if out.Item == nil {
		return false, nil
	}

	return true, nil
}
