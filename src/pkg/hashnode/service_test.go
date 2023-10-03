package hashnode

import (
	"context"
	"github.com/huantt/acc/src/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSubmitArticle(t *testing.T) {
	authToken := os.Getenv("HASHNODE_TOKEN")
	if authToken == "" {
		t.Skip()
	}
	service := NewService(Endpoint, authToken)
	body := "# Title \nHello world!"
	article := model.Article{
		Title:        "Test article",
		BodyMarkdown: &body,
	}
	err := service.SubmitArticle(context.Background(), article)
	assert.NoError(t, err)
}

func TestExists(t *testing.T) {
	service := NewService(Endpoint, "NO_NEED")
	tests := []struct {
		slug, hostname string
		exists         bool
	}{
		{"custom-resource-definitions-crds-extending-kubernetes", "sumanprasad", true},
		{"does-not-exist", "sumanprasad", false},
	}
	for _, test := range tests {
		exists, err := service.Exists(context.Background(), test.slug, test.hostname)
		assert.NoError(t, err)
		assert.Equal(t, test.exists, exists)
	}
}
