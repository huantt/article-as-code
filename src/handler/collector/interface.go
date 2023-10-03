package collector

import (
	"context"
	"github.com/huantt/acc/src/model"
)

type ArticleService interface {
	GetArticles(ctx context.Context, username string, page, perPage int) ([]model.Article, error)
}
