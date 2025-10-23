package dynamo

type CommentItem struct {
	URL       string `dynamodbav:"url"`
	UnixTime  int64  `dynamodbav:"unixTime"`
	UserID    string `dynamodbav:"user_id"`
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
}

type PageStructureItem struct {
	SiteDomain string `dynamodbav:"siteDomain"`
	URL        string `dynamodbav:"url"`
	Count      int    `dynamodbav:"count"`
}
