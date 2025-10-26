package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"pageknock-backend/dynamo"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

var (
	client                   *dynamodb.Client
	commentTable             string
	commentLogTable          string
	pageGlobalStructureTable string
	pageStructureTable       string
	recentDomainCommentTable string
	recentGlobalCommentTable string
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	region := os.Getenv("AWS_REGION")
	commentTable = os.Getenv("DYNAMO_TABLE_NAME_COMMENT")
	commentLogTable = os.Getenv("DYNAMO_TABLE_NAME_COMMENTLOG")
	pageGlobalStructureTable = os.Getenv("DYNAMO_TABLE_NAME_PAGEGLOBALSTRUCTURE")
	pageStructureTable = os.Getenv("DYNAMO_TABLE_NAME_PAGESTRUCTURE")
	recentDomainCommentTable = os.Getenv("DYNAMO_TABLE_NAME_RECENTDOMAINCOMMENT")
	recentGlobalCommentTable = os.Getenv("DYNAMO_TABLE_NAME_RECENTGLOBALCOMMENT")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

func main() {
	http.HandleFunc("/comment", handlePostComment)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePostComment(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL     string `json:"url"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.URL == "" || req.Comment == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	commentId := dynamo.GenerateCommentId()
	nowUnix := dynamo.GetUnixMillsecound()

	comment := dynamo.CommentItem{
		URL:       req.URL,
		UnixTime:  nowUnix,
		Comment:   req.Comment,
		CommentId: commentId,
		UserID:    "0",
	}

	err := dynamo.PutComment(client, commentTable, comment)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	recentGlobalComment := dynamo.RecentGlobalCommentItem{
		Global:    "GLOBAL",
		UnixTime:  nowUnix,
		Comment:   req.Comment,
		CommentId: commentId,
		URL:       req.URL,
		UserID:    "0",
	}

	err = dynamo.PutRecentGlobalComment(client, recentGlobalCommentTable, recentGlobalComment)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	domain, err := dynamo.GetDomainWithScheme(req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("URL変換処理失敗: %v", err), http.StatusInternalServerError)
		return
	}

	recentDomainComment := dynamo.RecentDomainCommentItem{
		SiteDomain: domain,
		UnixTime:   nowUnix,
		Comment:    req.Comment,
		CommentId:  commentId,
		URL:        req.URL,
		UserID:     "0",
	}

	err = dynamo.PutRecentDomainComment(client, recentDomainCommentTable, recentDomainComment)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	CommentLog := dynamo.CommentLogItem{
		Global:    "GOLBAL",
		UnixTime:  nowUnix,
		CommentId: commentId,
		Ip:        dynamo.GetIpAddress(w, r),
		UserAgent: dynamo.GetUserAgent(w, r),
	}

	err = dynamo.PutCommentLog(client, commentLogTable, CommentLog)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	err = handleStructureProcess(req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"message": "Insert succeeded!",
		"url":     req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleStructureProcess(url string) error {

	domain, err := dynamo.GetDomainWithScheme(url)
	if err != nil {
		return err //Failed to fetch data from DynamoDB
	}

	isExists, err := dynamo.ExistsStructureBySiteDomainAndURL(client, pageStructureTable, domain, url)
	if err != nil {
		return err //Failed to fetch data from DynamoDB
	}

	if isExists {

		err := dynamo.IncrementStructureCountByURL(client, pageStructureTable, domain, url)
		if err != nil {
			return err //Failed to fetch data from DynamoDB
		}

	} else {

		structureItem := dynamo.PageStructureItem{
			SiteDomain: domain,
			URL:        url,
			Count:      1,
		}

		PutStructureErr := dynamo.PutStructure(client, pageStructureTable, structureItem)
		if PutStructureErr != nil {
			return err //Failed to write data to DynamoDB
		}
	}

	isGlobalStructureExists, err := dynamo.ExistsGlobalStructureBySiteDomainAndURL(client, pageGlobalStructureTable, domain)
	if err != nil {
		return err //Failed to fetch data from DynamoDB
	}

	if isGlobalStructureExists {

		err := dynamo.IncrementGlobalStructureCountByURL(client, pageGlobalStructureTable, domain)
		if err != nil {
			return err //Failed to fetch data from DynamoDB
		}

	} else {

		globalStructureItem := dynamo.PageGlobalStructureItem{
			Global:     "GLOBAL",
			SiteDomain: domain,
			Count:      1,
		}

		PutStructureErr := dynamo.PutGlobalStructure(client, pageGlobalStructureTable, globalStructureItem)
		if PutStructureErr != nil {
			return err //Failed to write data to DynamoDB
		}
	}

	return nil
}
