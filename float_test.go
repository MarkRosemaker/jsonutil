package jsonutil_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
)

func TestFloat64UnmarshalStringOrNumber(t *testing.T) {
	t.Parallel()

	type option struct {
		Strike float64 `json:"strike,string"`
	}

	opts := json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.Float64UnmarshalStringOrNumber))

	for in, want := range map[string]float64{`{"strike":"9.5"}`: 9.5, `{"strike":"52"}`: 52, `{"strike":0.0}`: 0, `{"strike":3}`: 3} {
		var o option
		if err := json.Unmarshal([]byte(in), &o, opts); err != nil {
			t.Errorf("%s: %v", in, err)
		} else if o.Strike != want {
			t.Errorf("%s: got %v, want %v", in, o.Strike, want)
		}
	}

	for _, in := range []string{`{"strike":"x"}`, `{"strike":true}`} {
		var o option
		if err := json.Unmarshal([]byte(in), &o, opts); err == nil {
			t.Errorf("%s: got no error", in)
		}
	}

	// written as the string the API mostly sends
	if b, err := json.Marshal(option{Strike: 9.5}); err != nil || string(b) != `{"strike":"9.5"}` {
		t.Errorf("got %s and %v, want {\"strike\":\"9.5\"}", b, err)
	}
}
