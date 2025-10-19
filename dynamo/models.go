package dynamo

type CommentItem struct {
	URL      string `dynamodbav:"url"`
	UnixTime int64  `dynamodbav:"unixTime"`
	UserID   string `dynamodbav:"user_id"`
	Comment  string `dynamodbav:"comment"`
}

type PageStructureItem struct {
	Domain string `dynamodbav:"domain"`
	URL    string `dynamodbav:"url"`
}
