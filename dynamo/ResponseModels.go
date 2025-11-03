package dynamo

type PageGlobalStructureResponse struct {
	SiteDomain string `json:"siteDomain"`
	UrlCount   int    `json:"urlCount"`
}