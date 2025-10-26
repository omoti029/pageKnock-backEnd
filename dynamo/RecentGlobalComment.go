package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type RecentGlobalCommentRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRecentGlobalCommentRepository(client *dynamodb.Client, tableName string) *RecentGlobalCommentRepository {
	return &RecentGlobalCommentRepository{client: client, tableName: tableName}
}

func (r *RecentGlobalCommentRepository) PutRecentGlobalComment(item RecentGlobalCommentItem) error {
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

func (r *RecentGlobalCommentRepository) GetRecentGlobalComment() ([]RecentGlobalCommentItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("globalKey = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: "GLOBAL"},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(100),
	}

	out, err := r.client.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var comments []RecentGlobalCommentItem
	err = attributevalue.UnmarshalListOfMaps(out.Items, &comments)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return comments, nil
}
