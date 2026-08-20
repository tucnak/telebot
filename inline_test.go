package telebot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreparedInlineMessageUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "prepared_msg_abc123",
		"expiration_date": 1700000000
	}`

	var pim PreparedInlineMessage
	err := json.Unmarshal([]byte(jsonData), &pim)
	assert.NoError(t, err)
	assert.Equal(t, "prepared_msg_abc123", pim.ID)
	assert.Equal(t, int64(1700000000), pim.ExpirationDate)
}

func TestPreparedInlineMessageMarshal(t *testing.T) {
	pim := PreparedInlineMessage{
		ID:             "prepared_xyz",
		ExpirationDate: 1750000000,
	}

	data, err := json.Marshal(pim)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "prepared_xyz", result["id"])
	assert.Equal(t, float64(1750000000), result["expiration_date"])
}

func TestArticleResultNoHideURL(t *testing.T) {
	// Test that ArticleResult can be marshaled without hide_url field
	article := ArticleResult{
		ResultBase: ResultBase{ID: "test123"},
		Title:      "Test Article",
		Text:       "This is a test article",
		URL:        "https://example.com",
	}

	data, err := json.Marshal(article)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	assert.NotContains(t, result, "hide_url")
	assert.Equal(t, "Test Article", result["title"])
	assert.Equal(t, "This is a test article", result["message_text"])
	assert.Equal(t, "https://example.com", result["url"])
}

func TestArticleResultJSONUnmarshal(t *testing.T) {
	// Test that ArticleResult can be unmarshaled from JSON
	jsonData := `{
		"id": "article123",
		"type": "article",
		"title": "Example Article",
		"message_text": "This is the article content",
		"url": "https://example.com/article",
		"description": "Article description"
	}`

	var article ArticleResult
	err := json.Unmarshal([]byte(jsonData), &article)
	assert.NoError(t, err)
	assert.Equal(t, "article123", article.ID)
	assert.Equal(t, "article", article.Type)
	assert.Equal(t, "Example Article", article.Title)
	assert.Equal(t, "This is the article content", article.Text)
	assert.Equal(t, "https://example.com/article", article.URL)
	assert.Equal(t, "Article description", article.Description)
}
