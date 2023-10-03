package synchronizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadArticles(t *testing.T) {
	articles, err := readArticles("./testdata/articles")
	assert.NoError(t, err)
	assert.NotEmpty(t, articles)
}

func TestReadArticle(t *testing.T) {
	article, err := readArticle("./testdata/articles/advanced-go-build-techniques")
	assert.NoError(t, err)
	assert.NotNil(t, article)
	assert.NotNil(t, article.BodyMarkdown)
	assert.NotNil(t, article.Title)
}
