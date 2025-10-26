package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CommentLogRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewCommentLogRepository(client *dynamodb.Client, tableName string) *CommentLogRepository {
	return &CommentLogRepository{client: client, tableName: tableName}
}

func (r *CommentLogRepository) PutCommentLog(item CommentLogItem) error {
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
