package collector

import (
	_ "embed"
	"encoding/json"
	"github.com/huantt/acc/src/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/article.json
var articleData string

func TestWriteArticle(t *testing.T) {
	var article *model.Article
	err := json.Unmarshal([]byte(articleData), &article)
	if err != nil {
		panic(err)
	}
	err = writeArticle(".test", *article)
	assert.NoError(t, err)
}

func TestBuildArticleFolderPath(t *testing.T) {
	path := buildArticleFolderPath(model.Article{
		Url: "https://dev.to/jacktt/creating-dynamic-readmemd-file-388o",
	}, ".test")
	assert.NotNil(t, path)
	assert.Equal(t, ".test/dev.to/jacktt/creating-dynamic-readmemd-file-388o", path)

}
