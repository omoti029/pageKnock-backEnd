package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"pageknock-backend/dynamo"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

var (
	client                  *dynamodb.Client
	commentRepo             *dynamo.CommentRepository
	commentLogRepo          *dynamo.CommentLogRepository
	pageGlobalStructureRepo *dynamo.PageGlobalStructureRepository
	pageStructureRepo       *dynamo.PageStructureRepository
	recentDomainCommentRepo *dynamo.RecentDomainCommentRepository
	recentGlobalCommentRepo *dynamo.RecentGlobalCommentRepository
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	client = dynamodb.NewFromConfig(cfg)
	commentRepo = dynamo.NewCommentRepository(client, os.Getenv("DYNAMO_TABLE_NAME_COMMENT"))
	commentLogRepo = dynamo.NewCommentLogRepository(client, os.Getenv("DYNAMO_TABLE_NAME_COMMENTLOG"))
	pageGlobalStructureRepo = dynamo.NewPageGlobalStructureRepository(client, os.Getenv("DYNAMO_TABLE_NAME_PAGEGLOBALSTRUCTURE"))
	pageStructureRepo = dynamo.NewPageStructureRepository(client, os.Getenv("DYNAMO_TABLE_NAME_PAGESTRUCTURE"))
	recentDomainCommentRepo = dynamo.NewRecentDomainCommentRepository(client, os.Getenv("DYNAMO_TABLE_NAME_RECENTDOMAINCOMMENT"))
	recentGlobalCommentRepo = dynamo.NewRecentGlobalCommentRepository(client, os.Getenv("DYNAMO_TABLE_NAME_RECENTGLOBALCOMMENT"))
}

func main() {
	http.HandleFunc("/comment", handlePostComment)
	http.HandleFunc("/getPageGlobalStructure", handleGetPageGlobalStructure)
	http.HandleFunc("/getPageStructureBySiteDomain", handleGetPageStructureBySiteDomain)
	http.HandleFunc("/getRecentGlobalCommnet", handleGetRecentGlobalCommnet)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.Host, "localhost") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func handleGetPageStructureBySiteDomain(w http.ResponseWriter, r *http.Request) {

	var req struct {
		SiteDomain string `json:"siteDomain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.SiteDomain == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	records, err := pageStructureRepo.GetStructureBySiteDomain(req.SiteDomain)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch global structure: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]dynamo.PageStructureResponse, 0, len(records))
	for _, rec := range records {
		response = append(response, dynamo.PageStructureResponse{
			Url:            rec.Url,
			CommentCount:   rec.CommentCount,
			LatestUnixTime: rec.LatestUnixTime,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetPageGlobalStructure(w http.ResponseWriter, r *http.Request) {

	enableCORS(w, r)
	records, err := pageGlobalStructureRepo.GetGlobalStructure()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch global structure: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]dynamo.PageGlobalStructureResponse, 0, len(records))
	for _, rec := range records {
		response = append(response, dynamo.PageGlobalStructureResponse{
			SiteDomain: rec.SiteDomain,
			UrlCount:   rec.UrlCount,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetRecentGlobalCommnet(w http.ResponseWriter, r *http.Request) {

	enableCORS(w, r)
	records, err := recentGlobalCommentRepo.GetRecentGlobalComment()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch global structure: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]dynamo.RecentGlobalCommentResponse, 0, len(records))
	for _, rec := range records {
		response = append(response, dynamo.RecentGlobalCommentResponse{
			UnixTime:  rec.UnixTime,
			Comment:   rec.Comment,
			CommentId: rec.CommentId,
			Url:       rec.Url,
			UserID:    rec.UserID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handlePostComment(w http.ResponseWriter, r *http.Request) {

	enableCORS(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Url     string `json:"url"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Url == "" || req.Comment == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	domain, err := dynamo.GetDomainWithScheme(req.Url)
	if err != nil {
		http.Error(w, fmt.Sprintf("URL変換処理失敗: %v", err), http.StatusInternalServerError)
		return
	}

	commentId := dynamo.GenerateCommentId()
	nowUnix := dynamo.GetUnixMillsecound()

	baseFieldDatas := dynamo.BaseFieldDatas{
		Comment:    req.Comment,
		CommentId:  commentId,
		SiteDomain: domain,
		Now:        nowUnix,
		Req:        r,
		Url:        req.Url,
		UserId:     "1",
	}

	tableRecords := dynamo.GenerateAllTableRecords(baseFieldDatas)

	err = commentRepo.PutComment(tableRecords.CommentItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	err = recentGlobalCommentRepo.PutRecentGlobalComment(tableRecords.RecentGlobalCommentItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	err = recentDomainCommentRepo.PutRecentDomainComment(tableRecords.RecentDomainCommentItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	err = commentLogRepo.PutCommentLog(tableRecords.CommentLogItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	err = handleStructureProcess(tableRecords, baseFieldDatas)
	if err != nil {
		http.Error(w, fmt.Sprintf("DynamoDB書き込み失敗: %v", err), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"message": "Insert succeeded!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleStructureProcess(tableRecords dynamo.AllTableRecords, baseFieldDatas dynamo.BaseFieldDatas) error {
	isExists, err := pageStructureRepo.ExistsStructureBySiteDomainAndURL(baseFieldDatas.SiteDomain, baseFieldDatas.Url)
	if err != nil {
		return err //Failed to fetch data from DynamoDB
	}

	if isExists {

		err := pageStructureRepo.IncrementStructureCommentCountByURL(baseFieldDatas.SiteDomain, baseFieldDatas.Url, baseFieldDatas.Now)
		if err != nil {
			return err //Failed to fetch data from DynamoDB
		}

	} else {
		PutStructureErr := pageStructureRepo.PutStructure(tableRecords.PageStructureItem)
		if PutStructureErr != nil {
			return err //Failed to write data to DynamoDB
		}
	}

	isGlobalStructureExists, err := pageGlobalStructureRepo.ExistsGlobalStructureBySiteDomainAndURL(baseFieldDatas.SiteDomain)
	if err != nil {
		return err //Failed to fetch data from DynamoDB
	}

	if isGlobalStructureExists {

		err := pageGlobalStructureRepo.IncrementGlobalStructureUrlCountByURL(baseFieldDatas.SiteDomain)
		if err != nil {
			return err //Failed to fetch data from DynamoDB
		}

	} else {
		PutStructureErr := pageGlobalStructureRepo.PutGlobalStructure(tableRecords.PageGlobalStructureItem)
		if PutStructureErr != nil {
			return err //Failed to write data to DynamoDB
		}
	}
	return nil
}
