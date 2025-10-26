package dynamo

type CommentItem struct {
	Url       string `dynamodbav:"url"`      //PartitionKey
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
	Global     string `dynamodbav:"global"`     //PartitionKey
	SiteDomain string `dynamodbav:"siteDomain"` //Sort
	Count      int    `dynamodbav:"count"`
}
type PageStructureItem struct {
	SiteDomain string `dynamodbav:"siteDomain"` //PartitionKey
	Url        string `dynamodbav:"url"`        //Sort
	Count      int    `dynamodbav:"count"`
}

type RecentDomainCommentItem struct {
	SiteDomain string `dynamodbav:"siteDomain"` //PartitionKey
	UnixTime   int64  `dynamodbav:"unixTime"`   //Sort
	Comment    string `dynamodbav:"comment"`
	CommentId  string `dynamodbav:"commentId"`
	Url        string `dynamodbav:"url"`
	UserID     string `dynamodbav:"userId"`
}

type RecentGlobalCommentItem struct {
	Global    string `dynamodbav:"global"`   //PartitionKey
	UnixTime  int64  `dynamodbav:"unixTime"` //Sort
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
	Url       string `dynamodbav:"url"`
	UserID    string `dynamodbav:"userId"`
}
