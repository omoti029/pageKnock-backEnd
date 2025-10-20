package dynamo

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetUnixMillsecound() int64 {
	return time.Now().UnixMilli()
}

func GetDomainWithScheme(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URLの解析に失敗しました: %w", err)
	}

	scheme := parsed.Scheme
	if scheme == "" {
		scheme = "https"
	}

	domain := fmt.Sprintf("%s://%s", scheme, parsed.Host)
	return strings.TrimRight(domain, "/"), nil
}

func GenerateCommentId() string {
	id := uuid.New()
	return id.String()
}
