package hashnode

import (
	"context"
	"github.com/huantt/acc/model"
	hashnodesrv "github.com/huantt/acc/pkg/hashnode"
)

type Service struct {
	srv      *hashnodesrv.Service
	username string
}

func NewService(username, authToken string) *Service {
	return &Service{
		srv:      hashnodesrv.NewService(hashnodesrv.Endpoint, authToken),
		username: username,
	}
}

func (s *Service) Exists(ctx context.Context, slug string) (bool, error) {
	return s.srv.Exists(ctx, slug, s.username)
}

func (s *Service) SubmitArticle(ctx context.Context, article model.Article) error {
	return s.srv.SubmitArticle(ctx, article)
}
