package forem

import (
	"context"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/huantt/acc/src/model"
	forem2 "github.com/huantt/acc/src/pkg/forem"
	"log/slog"
	"strings"
	"sync"
)

type Service struct {
	foremService *forem2.Service
	uniqueSlugs  map[string]bool
	username     *string
}

func NewService(APIEndpoint string, maxRps int) *Service {
	return &Service{foremService: forem2.NewService(APIEndpoint, maxRps)}
}

func NewAuthenticatedService(APIEndpoint string, maxRps int, username string, authToken string) *Service {
	return &Service{
		foremService: forem2.NewAuthenticatedService(APIEndpoint, maxRps, authToken),
		username:     &username,
	}
}

func (s *Service) GetArticles(ctx context.Context, username string, page, perPage int) ([]model.Article, error) {
	articles, err := s.foremService.GetArticles(ctx, forem2.GetArticlesPrams{
		UserName: username,
		Page:     page,
		PerPage:  perPage,
	})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(articles))
	for i := range articles {
		go func(index int) {
			defer wg.Done()
			if fillBodyErr := s.fillBody(ctx, &articles[index]); fillBodyErr != nil {
				err = fillBodyErr
			}
		}(i)
	}
	wg.Wait()

	return toModels(articles), nil
}

func (s *Service) fillBody(ctx context.Context, article *forem2.Article) error {
	fullArticle, err := s.foremService.GetArticleById(ctx, fmt.Sprintf("%d", article.Id))
	if err != nil {
		return err
	}
	if fullArticle == nil {
		slog.Debug("Article not found: %d", article.Id)
		return nil
	}
	article.BodyMarkdown = fullArticle.BodyMarkdown
	article.BodyHtml = fullArticle.BodyHtml
	slog.Debug(fmt.Sprintf("Filled body for article: %s", article.Url))
	return nil
}

func toModels(articles []forem2.Article) []model.Article {
	var result []model.Article
	for _, article := range articles {
		result = append(result, toModel(article))
	}
	return result
}

func toModel(article forem2.Article) model.Article {
	return model.Article{
		Url:          article.Url,
		Slug:         article.Slug,
		Title:        article.Title,
		Description:  article.Description,
		Thumbnail:    article.CoverImage,
		BodyMarkdown: article.BodyMarkdown,
		BodyHtml:     article.BodyHtml,
		Author: model.Author{
			Id:   article.User.Username,
			Name: article.User.Name,
		},
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.EditedAt,
		Tags:      article.GetTags(),
	}
}

func (s *Service) Exists(ctx context.Context, slug string) (bool, error) {
	if s.uniqueSlugs == nil {
		if err := s.loadSlugs(ctx); err != nil {
			return false, err
		}
	}

	return s.uniqueSlugs[slug], nil
}

func (s *Service) loadSlugs(ctx context.Context) error {
	s.uniqueSlugs = map[string]bool{}
	if s.username == nil {
		return errors.New("username has not been specified")
	}
	page := 1
	pageSize := 20
	for {
		articles, err := s.foremService.GetArticles(ctx, forem2.GetArticlesPrams{
			UserName: *s.username,
			Page:     page,
			PerPage:  pageSize,
		})
		if err != nil {
			return err
		}
		for _, article := range articles {
			s.uniqueSlugs[extractRealSlug(article.Slug)] = true
			s.uniqueSlugs[slug.Make(article.Title)] = true
		}
		slog.Debug(fmt.Sprintf("[Page %d - Page size %d] Loaded %d slug to uniqueSlugs", page, pageSize, len(articles)))

		if len(articles) < pageSize {
			break
		}
		page++
	}
	return nil
}

func extractRealSlug(slug string) string {
	parts := strings.Split(slug, "-")
	return strings.Join(parts[0:len(parts)-1], "-")
}

func (s *Service) SubmitArticle(ctx context.Context, article model.Article) error {
	return s.foremService.SubmitArticle(ctx, article)
}
