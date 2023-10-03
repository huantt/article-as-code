package collector

import (
	"context"
	"github.com/huantt/article-as-code/model"
)

type ArticleService interface {
	GetArticles(ctx context.Context, username string, page, perPage int) ([]model.Article, error)
}
