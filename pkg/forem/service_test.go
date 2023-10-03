package forem

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetArticles(t *testing.T) {
	userName := "jacktt"
	perPage := 10
	service := NewService(DevToEndpoint, 3)
	articles, err := service.GetArticles(context.Background(), GetArticlesPrams{
		UserName: userName,
		PerPage:  perPage,
	})
	assert.NoError(t, err)
	assert.Equal(t, 10, len(articles))
	for _, article := range articles {
		assert.Equal(t, userName, article.User.Username)
	}
}

func TestGetArticleById(t *testing.T) {
	service := NewService(DevToEndpoint, 3)
	articleId := "1594564"
	article, err := service.GetArticleById(context.Background(), articleId)
	assert.NoError(t, err)
	assert.NotNil(t, article.BodyMarkdown)
	assert.NotEmpty(t, *article.BodyMarkdown)
	assert.NotNil(t, article.BodyHtml)
	assert.NotEmpty(t, *article.BodyHtml)
}
