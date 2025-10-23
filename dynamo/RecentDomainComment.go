package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutRecentDomainComment(client *dynamodb.Client, tableName string, item RecentDomainCommentItem) error {
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

func GetRecentDomainComment(client *dynamodb.Client, tableName string, siteDomain string) ([]RecentDomainCommentItem, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("siteDomain = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: siteDomain},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(100),
	}

	out, err := client.Query(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var comments []RecentDomainCommentItem
	err = attributevalue.UnmarshalListOfMaps(out.Items, &comments)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return comments, nil
}
