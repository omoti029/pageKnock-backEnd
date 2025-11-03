package dynamo

type PageGlobalStructureResponse struct {
	SiteDomain string `json:"siteDomain"`
	UrlCount   int    `json:"urlCount"`
}

type PageStructureResponse struct {
	Url            string `json:"urlCount"`
	CommentCount   int    `json:"commentCount"`
	LatestUnixTime int64  `json:"latestUnixTime"`
}

type RecentGlobalCommentResponse struct {
	UnixTime  int64  `dynamodbav:"unixTime"` //Sort
	Comment   string `dynamodbav:"comment"`
	CommentId string `dynamodbav:"commentId"`
	Url       string `dynamodbav:"url"`
	UserID    string `dynamodbav:"userId"`
}
