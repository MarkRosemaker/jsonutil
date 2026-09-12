package jsonutil

import (
	"encoding/json/v2"
	"testing"
	"time"
)

func TestTimeUnmarshalStringOrIntUnixLayouts(t *testing.T) {
	jsonOpts := json.JoinOptions(
		json.WithUnmarshalers(json.UnmarshalFromFunc(
			timeUnmarshalStringOrIntUnix([]string{time.RFC3339, time.DateOnly}),
		)),
	)

	t.Run("multiple layouts", func(t *testing.T) {
		var out time.Time
		if err := json.Unmarshal([]byte(`"2024-01-02"`), &out, jsonOpts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if want := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC); !out.Equal(want) {
			t.Fatalf("want: %s, got: %s", want, out)
		}
	})

	t.Run("no layout matches", func(t *testing.T) {
		var out time.Time
		if err := json.Unmarshal([]byte(`"not a time"`), &out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("unix seconds", func(t *testing.T) {
		var out time.Time
		if err := json.Unmarshal([]byte(`0`), &out, jsonOpts); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !out.IsZero() {
			t.Fatalf("want zero time, got: %s", out)
		}
	})

}
