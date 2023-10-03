package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/huantt/acc/model"
	"log/slog"
	"os"
	"strings"
)

type Handler struct {
	articleService ArticleService
	articleFolder  string
}

func NewHandler(articleService ArticleService, articleFolder string) *Handler {
	return &Handler{
		articleService: articleService,
		articleFolder:  articleFolder,
	}
}

const (
	ReadmeFileName = "README.md"
	MetaFileName   = "metadata.json"
)

func (h *Handler) Collect(ctx context.Context, username string) error {
	page := 1
	pageSize := 20
	for {
		articles, err := h.articleService.GetArticles(ctx, username, page, pageSize)
		if err != nil {
			return err
		}
		slog.Info(fmt.Sprintf("Got %d articles - page: %d, perPage: %d", len(articles), page, pageSize))
		for i := range articles {
			articles[i].Slug = slug.Make(articles[i].Title)
		}

		if err := writeArticles(h.articleFolder, articles); err != nil {
			return err
		}
		slog.Info("Wrote successfully")
		if len(articles) < pageSize {
			break
		}
		page++
	}
	return nil
}

func writeArticles(rootFolder string, articles []model.Article) error {
	for _, article := range articles {
		if err := writeArticle(rootFolder, article); err != nil {
			return err
		}
	}
	return nil
}

func writeArticle(rootFolder string, article model.Article) error {
	if article.BodyMarkdown == nil {
		slog.Warn(fmt.Sprintf("Article %s has no body", article.Url))
		return nil
	}

	articleFolderPath := buildArticleFolderPath(article, rootFolder)

	if err := mkdirIfNotExists(articleFolderPath); err != nil {
		return errors.Join(err, errors.New(articleFolderPath))
	}

	articleFilePath := fmt.Sprintf("%s/%s", articleFolderPath, ReadmeFileName)
	if err := os.WriteFile(articleFilePath, []byte(*article.BodyMarkdown), 0644); err != nil {
		return err
	}
	metadataFile := fmt.Sprintf("%s/%s", articleFolderPath, MetaFileName)
	metadata, err := json.Marshal(article)
	if err != nil {
		return err
	}
	if err := os.WriteFile(metadataFile, metadata, 0644); err != nil {
		return err
	}
	return nil
}

func buildArticleFolderPath(article model.Article, rootFolder string) string {
	path := fmt.Sprintf("%s/%s", rootFolder, article.Slug)
	return path
}

func mkdirIfNotExists(folderPath string) error {
	for i := range strings.Split(folderPath, "/") {
		path := strings.Join(strings.Split(folderPath, "/")[:i+1], "/")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}
