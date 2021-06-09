package youdao

import (
	"context"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestClient_TextTranslation(t *testing.T) {
	r, err := recorder.New("fixtures/TestClient_TextTranslation")
	require.Nil(t, err)
	defer func() {
		_ = r.Stop()
	}()
	c := NewClient(Config{AppKey: os.Getenv("APP_KEY"), AppSecret: os.Getenv("APP_SECRET"), Client: &http.Client{Transport: r}})
	ctx := context.Background()
	resp, err := c.TextTranslation(ctx, TextTranslationReq{
		FromLang: "auto",
		ToLang:   "auto",
		Q:        "hello world",
	})
	require.Nil(t, err)
	require.Equal(t, []string{"你好世界"}, resp.Translation)
}
