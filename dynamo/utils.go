package dynamo

import (
	"fmt"
	"net/http"
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

func GetIpAddress(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return ip
}

func GetUserAgent(req *http.Request) string {
	userAgent := req.Header.Get("User-Agent")
	return userAgent
}

func GenerateAllTableRecords(Datas BaseFieldDatas) AllTableRecords {
	return AllTableRecords{
		CommentItem: CommentItem{
			Url:       Datas.Url,
			UnixTime:  Datas.Now,
			Comment:   Datas.Comment,
			CommentId: Datas.CommentId,
			UserID:    Datas.UserId,
		},
		CommentLogItem: CommentLogItem{
			GlobalKey: "GLOBAL", // 生成関数などで作る
			UnixTime:  Datas.Now,
			CommentId: Datas.CommentId,
			Ip:        GetIpAddress(Datas.Req),
			UserAgent: GetUserAgent(Datas.Req),
		},
		PageGlobalStructureItem: PageGlobalStructureItem{
			GlobalKey:  "GLOBAL",
			SiteDomain: Datas.SiteDomain,
			Count:      1,
		},
		PageStructureItem: PageStructureItem{
			SiteDomain: Datas.SiteDomain,
			Url:        Datas.Url,
			Count:      1,
		},
		RecentDomainCommentItem: RecentDomainCommentItem{
			SiteDomain: Datas.SiteDomain,
			UnixTime:   Datas.Now,
			Comment:    Datas.Comment,
			CommentId:  Datas.CommentId,
			Url:        Datas.Url,
			UserID:     Datas.UserId,
		},
		RecentGlobalCommentItem: RecentGlobalCommentItem{
			GlobalKey: "GLOBAL",
			UnixTime:  Datas.Now,
			Comment:   Datas.Comment,
			CommentId: Datas.CommentId,
			Url:       Datas.Url,
			UserID:    Datas.UserId,
		},
	}
}
