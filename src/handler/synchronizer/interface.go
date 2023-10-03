package synchronizer

import (
	"context"
	"github.com/huantt/acc/src/model"
)

type ArticleService interface {
	Exists(ctx context.Context, slug string) (bool, error)
	SubmitArticle(ctx context.Context, article model.Article) error
}
