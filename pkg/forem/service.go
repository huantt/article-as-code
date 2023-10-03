package forem

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/huantt/acc/model"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// Service docs: https://developers.forem.com/api
type Service struct {
	httpClient  *resty.Client
	rateLimiter *time.Ticker
	authToken   *string
}

func NewService(APIEndpoint string, rps int) *Service {
	httpClient := resty.New()
	httpClient.
		SetRetryCount(12).
		SetRetryWaitTime(5 * time.Second).
		SetBaseURL(APIEndpoint).AddRetryCondition(func(response *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		if response.StatusCode() == http.StatusInternalServerError ||
			response.StatusCode() == http.StatusBadGateway ||
			response.StatusCode() == http.StatusGatewayTimeout ||
			response.StatusCode() == http.StatusServiceUnavailable {
			slog.Warn(fmt.Sprintf("Response status code is %d - Request: %s - Body: %s - Retrying...", response.StatusCode(), response.Request.URL, response.Body()))
			return true
		}

		return false
	})
	return &Service{httpClient: httpClient, rateLimiter: time.NewTicker(time.Second / time.Duration(rps))}
}

func NewAuthenticatedService(APIEndpoint string, rps int, authToken string) *Service {
	service := NewService(APIEndpoint, rps)
	service.authToken = &authToken
	return service
}

// GetArticles Docs: https://developers.forem.com/api/v1#tag/articles/operation/getArticles
func (s *Service) GetArticles(ctx context.Context, params GetArticlesPrams) ([]Article, error) {
	<-s.rateLimiter.C
	endpoint := "/api/articles"
	if params.MostRecent {
		endpoint = fmt.Sprintf("%s/latest", endpoint)
	}

	var articles []Article
	request := s.httpClient.R().SetContext(ctx).SetResult(&articles)
	if params.Page > 0 {
		request = request.SetQueryParam("page", fmt.Sprintf("%d", params.Page))
	}
	if params.PerPage > 0 {
		request = request.SetQueryParam("per_page", fmt.Sprintf("%d", params.PerPage))
	}
	if params.Tag != "" {
		request = request.SetQueryParam("tag", params.Tag)
	}
	if len(params.Tags) > 0 {
		request = request.SetQueryParam("tags", strings.Join(params.Tags, ","))
	}
	if len(params.TagsExclude) > 0 {
		request = request.SetQueryParam("tags_exclude", strings.Join(params.TagsExclude, ","))
	}
	if params.UserName != "" {
		request = request.SetQueryParam("username", params.UserName)
	}
	if params.State != "" {
		request = request.SetQueryParam("state", params.State)
	}
	if params.Top > 0 {
		request = request.SetQueryParam("top", fmt.Sprintf("%d", params.Top))
	}
	if params.CollectionID > 0 {
		request = request.SetQueryParam("collection_id", fmt.Sprintf("%d", params.CollectionID))
	}
	resp, err := request.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body()))
	}
	return articles, nil
}

func (s *Service) GetArticleById(ctx context.Context, id string) (*Article, error) {
	<-s.rateLimiter.C
	slog.Debug(fmt.Sprintf("Requesting article by id: %s", id))
	var article *Article
	resp, err := s.httpClient.R().SetContext(ctx).
		SetPathParam("id", id).
		SetResult(&article).
		Get("/api/articles/{id}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body()))
	}
	return article, nil
}

func (s *Service) SubmitArticle(ctx context.Context, article model.Article) error {
	if s.authToken == nil {
		return errors.New("auth token is not set")
	}
	<-s.rateLimiter.C
	requestBody := SubmitArticleRequest{
		Title:        article.Title,
		Published:    true,
		BodyMarkdown: *article.BodyMarkdown,
		Tags:         article.Tags,
	}
	resp, err := s.httpClient.R().SetContext(ctx).
		SetBody(map[string]SubmitArticleRequest{
			"article": requestBody,
		}).
		SetHeader("api-key", *s.authToken).
		Post("/api/articles")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New(fmt.Sprintf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body()))
	}
	return nil
}
