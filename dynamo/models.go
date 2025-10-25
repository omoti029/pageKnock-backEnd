package dynamo

type CommentItem struct {
	URL       string `dynamodbav:"url"`
	UnixTime  int64  `dynamodbav:"unixTime"`
	UserID    string `dynamodbav:"user_id"`
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
}

type RecentGlobalCommentItem struct {
	GlobalRecent string `dynamodbav:"globalRecent"`
	URL          string `dynamodbav:"url"`
	UnixTime     int64  `dynamodbav:"unixTime"`
	UserID       string `dynamodbav:"user_id"`
	Comment      string `dynamodbav:"comment"`
	CommentId    string `dynamodbav:"commentId"`
}

type RecentDomainCommentItem struct {
	SiteDomain string `dynamodbav:"siteDomain"`
	URL        string `dynamodbav:"url"`
	UnixTime   int64  `dynamodbav:"unixTime"`
	UserID     string `dynamodbav:"user_id"`
	Comment    string `dynamodbav:"comment"`
	CommentId  string `dynamodbav:"commentId"`
}

type CommentLogItem struct {
	Global    string `dynamodbav:"global"`
	CommentId string `dynamodbav:"commentId"`
	Ip        string `dynamodbav:"ip"`
	UserAgent string `dynamodbav:"userAgent"`
	UnixTime  int64  `dynamodbav:"unixTime"`
}

type PageStructureItem struct {
	SiteDomain string `dynamodbav:"siteDomain"`
	URL        string `dynamodbav:"url"`
	Count      int    `dynamodbav:"count"`
}

type PageGlobalStructureItem struct {
	GlobalSiteDomain string `dynamodbav:"globalSiteDomain"`
	SiteDomain       string `dynamodbav:"siteDomain"`
	Count            int    `dynamodbav:"count"`
}
