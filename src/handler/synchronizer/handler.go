package synchronizer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huantt/acc/src/handler/collector"
	"github.com/huantt/acc/src/model"
	"log/slog"
	"os"
)

type Handler struct {
	serviceByDestination map[string]ArticleService
}

func NewHandler() *Handler {
	return &Handler{serviceByDestination: map[string]ArticleService{}}
}

func (h *Handler) RegisterService(destination string, service ArticleService) *Handler {
	h.serviceByDestination[destination] = service
	return h
}

func (h *Handler) Sync(ctx context.Context, articleFolderPath, destination string) error {
	service, found := h.serviceByDestination[destination]
	if !found {
		return errors.New("unknown destination")
	}
	articles, err := readArticles(articleFolderPath)
	if err != nil {
		return err
	}
	for _, article := range articles {
		exists, err := service.Exists(ctx, article.Slug)
		if err != nil {
			return err
		}
		if exists {
			slog.Debug(fmt.Sprintf("Article existed: %s", article.Slug))
			continue
		}

		if err = service.SubmitArticle(ctx, article); err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("Article %s submitted", article.Slug))
	}
	return nil

}

func readArticles(articleFolderPath string) ([]model.Article, error) {
	dirEntities, err := os.ReadDir(articleFolderPath)
	if err != nil {
		return nil, err
	}
	articles := make([]model.Article, 0)
	for _, dirEntity := range dirEntities {
		article, err := readArticle(fmt.Sprintf("%s/%s", articleFolderPath, dirEntity.Name()))
		if err != nil {
			return nil, err
		}
		articles = append(articles, *article)
	}
	return articles, nil
}

func readArticle(path string) (*model.Article, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, collector.MetaFileName))
	if err != nil {
		return nil, err
	}
	var article model.Article
	if err = json.Unmarshal(data, &article); err != nil {
		return nil, err
	}
	return &article, nil
}
