package hashnode

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/huantt/article-as-code/model"
	"log/slog"
	"net/http"
	"time"
)

const Endpoint = "https://api.hashnode.com"

type Service struct {
	httpClient *resty.Client
}

func NewService(APIEndpoint string, authToken string) *Service {
	httpClient := resty.New()
	httpClient.
		SetHeader("Authorization", authToken).
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
	return &Service{httpClient}
}

func (s *Service) SubmitArticle(ctx context.Context, article model.Article) error {
	if article.BodyMarkdown == nil {
		return errors.New("article has no body")
	}
	requestBody := NewCreateArticleRequest()
	requestBody.Variables.Input.Title = article.Title
	requestBody.Variables.Input.ContentMarkdown = *article.BodyMarkdown
	requestBody.Variables.Input.CoverImageURL = article.Thumbnail
	requestBody.Variables.Input.IsPartOfPublication.PublicationId = "649ba336b4f390887392535c"
	requestBody.Variables.Input.Tags = []Tag{}
	resp, err := s.httpClient.R().SetContext(ctx).
		SetBody(requestBody).
		Post("")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body())
	}
	return nil
}

func (s *Service) Exists(ctx context.Context, slug, hostname string) (bool, error) {
	requestBody := NewGetPostRequest()
	requestBody.Variables.Slug = slug
	requestBody.Variables.Hostname = hostname
	var result struct {
		Data struct {
			Post *Post `json:"post"`
		} `json:"data"`
		Errors *any `json:"errors"`
	}
	resp, err := s.httpClient.R().SetContext(ctx).
		SetBody(requestBody).
		SetResult(&result).
		Post("")
	if err != nil {
		return false, err
	}
	if resp.IsError() {
		return false, fmt.Errorf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body())
	}
	return result.Data.Post != nil, nil
}
