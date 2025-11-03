package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PageGlobalStructureRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewPageGlobalStructureRepository(client *dynamodb.Client, tableName string) *PageGlobalStructureRepository {
	return &PageGlobalStructureRepository{client: client, tableName: tableName}
}

func (r *PageGlobalStructureRepository) PutGlobalStructure(item PageGlobalStructureItem) error {
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

func (r *PageGlobalStructureRepository) GetGlobalStructure() ([]PageGlobalStructureItem, error) {

	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("globalKey = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: "GLOBAL"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var records []PageGlobalStructureItem
	err = attributevalue.UnmarshalListOfMaps(out.Items, &records)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return records, nil
}

func (r *PageGlobalStructureRepository) IncrementGlobalStructureUrlCountByURL(siteDomain string) error {

	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"globalKey":  &types.AttributeValueMemberS{Value: "GLOBAL"},
			"siteDomain": &types.AttributeValueMemberS{Value: siteDomain},
		},
		UpdateExpression: aws.String("ADD #c :inc"),
		ExpressionAttributeNames: map[string]string{
			"#c": "urlCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return fmt.Errorf("failed to increment urlCount for %s", err)
	}

	return nil
}

func (r *PageGlobalStructureRepository) ExistsGlobalStructureBySiteDomainAndURL(siteDomain string) (bool, error) {

	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"globalKey":  &types.AttributeValueMemberS{Value: "GLOBAL"},
			"siteDomain": &types.AttributeValueMemberS{Value: siteDomain},
		},
		ProjectionExpression: aws.String("globalKey"),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get item: %w", err)
	}

	if out.Item == nil {
		return false, nil
	}

	return true, nil
}
