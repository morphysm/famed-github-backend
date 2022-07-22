package assets_test

import (
	"github.com/morphysm/famed-github-backend/assets"
	"strings"
	"testing"
)

// TestBanner test if Banner contains "Go Backend"
func TestBanner(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		t.Parallel()

		got := assets.Banner
		want := "Go Backend"

		if !strings.Contains(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
}
