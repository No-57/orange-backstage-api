package crawler

import (
	"context"
	"io"
	"testing"
)

func TestMomo_fetch(t *testing.T) {
	momo := NewMomo()
	body, err := momo.fetch(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch: %v", err)
	}
	defer body.Close()

	bs, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	t.Log(string(bs))
}

func TestMomoFetch(t *testing.T) {
	momo := NewMomo()
	products, err := momo.Fetch(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch: %v", err)
	}
	t.Logf("%+v", products)
}
