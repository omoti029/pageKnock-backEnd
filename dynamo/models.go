package dynamo

type CommentItem struct {
	URL       string `dynamodbav:"url"`      //PartitionKey
	UnixTime  int64  `dynamodbav:"unixTime"` //Sort
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
	UserID    string `dynamodbav:"userId"`
}

type CommentLogItem struct {
	Global    string `dynamodbav:"global"`   //PartitionKey
	UnixTime  int64  `dynamodbav:"unixTime"` //Sort
	CommentId string `dynamodbav:"commentId"`
	Ip        string `dynamodbav:"ip"`
	UserAgent string `dynamodbav:"userAgent"`
}

type PageGlobalStructureItem struct {
	Global string `dynamodbav:"global"` //PartitionKey
	Domain string `dynamodbav:"domain"` //Sort
	Count  int    `dynamodbav:"count"`
}
type PageStructureItem struct {
	Domain string `dynamodbav:"domain"` //PartitionKey
	URL    string `dynamodbav:"url"`    //Sort
	Count  int    `dynamodbav:"count"`
}

type RecentDomainCommentItem struct {
	Domain    string `dynamodbav:"domain"`   //PartitionKey
	UnixTime  int64  `dynamodbav:"unixTime"` //Sort
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
	URL       string `dynamodbav:"url"`
	UserID    string `dynamodbav:"userId"`
}

type RecentGlobalCommentItem struct {
	Global    string `dynamodbav:"global"` //PartitionKey
	URL       string `dynamodbav:"url"`    //Sort
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
	UnixTime  int64  `dynamodbav:"unixTime"`
	UserID    string `dynamodbav:"userId"`
}
