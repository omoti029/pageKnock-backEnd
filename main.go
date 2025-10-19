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
	client             *dynamodb.Client
	commentTable       string
	pageStructureTable string
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	region := os.Getenv("AWS_REGION")
	commentTable = os.Getenv("DYNAMO_TABLE_NAME_COMMENT")
	pageStructureTable = os.Getenv("DYNAMO_TABLE_NAME_PAGESTRUCTURE")

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

	comment := dynamo.CommentItem{
		URL:      req.URL,
		UnixTime: dynamo.GetUnixMillsecound(),
		UserID:   "0",
		Comment:  req.Comment,
	}

	err := dynamo.PutComment(client, commentTable, comment)
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

		comment := dynamo.PageStructureItem{
			URL:        url,
			SiteDomain: domain,
			Count:      1,
		}

		PutStructureErr := dynamo.PutStructure(client, pageStructureTable, comment)
		if PutStructureErr != nil {
			return err //Failed to write data to DynamoDB
		}
	}

	return nil
}
