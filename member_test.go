package jsonutil_test

import (
	"encoding/json/v2"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/jsonutil"
)

func TestIntUnixFor(t *testing.T) {
	t.Parallel()

	type meta struct {
		Created time.Time  `json:"created"`
		Expires civil.Date `json:"expires,omitzero"`
	}

	// some dates and times as unix time, others as strings, also through an embedded struct and in an array
	type wrapper struct {
		meta
		Raw  civil.Date   `json:"raw"`
		Fmt  civil.Date   `json:"fmt"`
		At   time.Time    `json:"at"`
		Days []civil.Date `json:"days"`
	}

	unix := []string{"expires", "raw", "at", "days"}
	opts := json.JoinOptions(
		json.WithMarshalers(json.JoinMarshalers(
			json.MarshalToFunc(jsonutil.DateMarshalIntUnixFor(unix...)),
			json.MarshalToFunc(jsonutil.TimeMarshalIntUnixFor(unix...)),
		)),
		json.WithUnmarshalers(json.JoinUnmarshalers(
			json.UnmarshalFromFunc(jsonutil.DateUnmarshalIntUnixFor(unix...)),
			json.UnmarshalFromFunc(jsonutil.TimeUnmarshalIntUnixFor(unix...)),
		)),
	)

	day := civil.Date{Year: 2025, Month: time.September, Day: 30}
	at := time.Unix(1759190400, 0).UTC()
	in := wrapper{Created: at, Expires: day, Raw: day, Fmt: day, At: at, Days: []civil.Date{day}}

	b, err := json.Marshal(in, opts)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"created":"2025-09-30T00:00:00Z","expires":1759190400,"raw":1759190400,"fmt":"2025-09-30",` +
		`"at":1759190400,"days":[1759190400]}`
	if string(b) != want {
		t.Fatalf("got %s, want %s", b, want)
	}

	var out wrapper
	if err := json.Unmarshal(b, &out, opts); err != nil {
		t.Fatal(err)
	}

	if !out.At.Equal(in.At) || !out.Created.Equal(in.Created) {
		t.Errorf("got times %s and %s, want %s", out.At, out.Created, at)
	}

	out.At, out.Created = in.At, in.Created
	if !reflect.DeepEqual(out, in) {
		t.Errorf("got %+v, want %+v", out, in)
	}
}
