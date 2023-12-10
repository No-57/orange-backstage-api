package crawler

import (
	"context"
	"testing"
)

func TestPchomeFetch(t *testing.T) {
	pchome := NewPchome()
	products, err := pchome.Fetch(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch: %v", err)
	}
	t.Logf("%+v", products)
}
