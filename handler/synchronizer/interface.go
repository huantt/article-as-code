package synchronizer

import (
	"context"
	"github.com/huantt/article-as-code/model"
)

type ArticleService interface {
	Exists(ctx context.Context, slug string) (bool, error)
	SubmitArticle(ctx context.Context, article model.Article) error
}
