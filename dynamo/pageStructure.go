package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PageStructureRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewPageStructureRepository(client *dynamodb.Client, tableName string) *PageStructureRepository {
	return &PageStructureRepository{client: client, tableName: tableName}
}

func (r *PageStructureRepository) PutStructure(item PageStructureItem) error {
	av, err := attributevalue.MarshalMap(item)

	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	return err
}

func (r *PageStructureRepository) GetStructureBySiteDomain(siteDomain string) ([]PageStructureItem, error) {

	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("siteDomain = :d"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":d": &types.AttributeValueMemberS{Value: siteDomain},
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

func (r *PageStructureRepository) IncrementStructureCommentCountByURL(siteDomain string, url string) error {

	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"siteDomain": &types.AttributeValueMemberS{Value: siteDomain},
			"url":        &types.AttributeValueMemberS{Value: url},
		},
		UpdateExpression: aws.String("ADD #c :inc"),
		ExpressionAttributeNames: map[string]string{
			"#c": "commentCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return fmt.Errorf("failed to increment commentCount for %s: %w", url, err)
	}

	return nil
}

func (r *PageStructureRepository) ExistsStructureBySiteDomainAndURL(siteDomain string, url string) (bool, error) {

	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"siteDomain": &types.AttributeValueMemberS{Value: siteDomain},
			"url":        &types.AttributeValueMemberS{Value: url},
		},
		ProjectionExpression: aws.String("siteDomain"),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}

	if out.Item == nil {
		return false, nil
	}

	return true, nil
}
