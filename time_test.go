package jsonutil_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/MarkRosemaker/jsonutil"
)

type testTime struct {
	Time                 time.Time  `json:"time"`
	TimePointer          *time.Time `json:"timePointer"`
	TimeOmitZero         time.Time  `json:"timeOmitZero,omitzero"`
	TimePointerOmitEmpty *time.Time `json:"timePointerOmitEmpty,omitempty"`
}

func TestUnixTime(t *testing.T) {
	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalToFunc(jsonutil.TimeMarshalIntUnix)),
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.TimeUnmarshalIntUnix)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &testTime{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"time":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not an int", func(t *testing.T) {
		out := &testTime{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"time": "3"}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpInt := reflect.TypeFor[int64](); errSem.GoType != tpInt {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpInt, errSem.GoType)
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testTime{}
		if err := json.Unmarshal([]byte(`{"time":null,"timePointer":null,"timeOmitZero":null,"timePointerOmitEmpty":null}`), out, jsonOpts); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	now := time.Now().Truncate(time.Second)
	nowUnix := now.Unix()

	for i, tc := range []struct {
		in  testTime
		out string
	}{
		{testTime{}, `{"time":0,"timePointer":null}`},
		{
			testTime{Time: now, TimePointer: &now, TimeOmitZero: now, TimePointerOmitEmpty: &now},
			fmt.Sprintf(`{"time":%[1]d,"timePointer":%[1]d,"timeOmitZero":%[1]d,"timePointerOmitEmpty":%[1]d}`, nowUnix),
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tc.in, jsonOpts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != tc.out {
				t.Fatalf("want: %s, got: %s", tc.out, string(b))
			}

			var out testTime
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out.Time, tc.in.Time; !got.Equal(want) {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			if got, want := out.TimeOmitZero, tc.in.TimeOmitZero; !got.Equal(want) {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			for _, tt := range []struct{ got, want *time.Time }{
				{out.TimePointer, tc.in.TimePointer},
				{out.TimePointerOmitEmpty, tc.in.TimePointerOmitEmpty},
			} {
				if tt.got == nil && tt.want != nil {
					t.Fatalf("expected non-nil URLPointer")
				} else if tt.got != nil && tt.want == nil {
					t.Fatalf("expected nil URLPointer")
				} else if tt.got != nil && tt.want != nil {
					if !tt.got.Equal(*tt.want) {
						t.Fatalf("want: %s, got: %s", tt.want, tt.got)
					}
				}
			}
		})
	}
}

func TestStringOrUnixTime(t *testing.T) {
	const altLayout = "Mon Jan 2 2006 15:04:05 MST-0700"

	jsonOpts := json.JoinOptions(
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.TimeUnmarshalStringOrIntUnix)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &testTime{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"time":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not a string or int", func(t *testing.T) {
		out := &testTime{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"time": true}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpTime := reflect.TypeFor[*time.Time](); errSem.GoType != tpTime {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpTime, errSem.GoType)
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testTime{}
		if err := json.Unmarshal([]byte(`{"time":null,"timePointer":null,"timeOmitZero":null,"timePointerOmitEmpty":null}`), out, jsonOpts); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	now := time.Now().Truncate(time.Second)
	nowUnix := now.Unix()
	milli := now.Add(123 * time.Millisecond)

	for i, tc := range []struct {
		in   string
		want testTime
	}{
		{`{"time":0,"timePointer":null}`, testTime{}},
		{`{"time":""}`, testTime{}},
		{
			fmt.Sprintf(`{"time":%[1]d,"timePointer":%[1]d,"timeOmitZero":%[1]d,"timePointerOmitEmpty":%[1]d}`, nowUnix),
			testTime{Time: now, TimePointer: &now, TimeOmitZero: now, TimePointerOmitEmpty: &now},
		},
		{
			fmt.Sprintf(`{"time":%[1]d,"timePointer":%[1]d,"timeOmitZero":%[1]d,"timePointerOmitEmpty":%[1]d}`, milli.UnixMilli()),
			testTime{Time: milli, TimePointer: &milli, TimeOmitZero: milli, TimePointerOmitEmpty: &milli},
		},
		{`{"time":1485487350827}`, testTime{Time: time.Date(2017, 1, 27, 3, 22, 30, 827e6, time.UTC)}},
		{`{"time":99999999999}`, testTime{Time: time.Unix(99999999999, 0)}},
		{
			fmt.Sprintf(`{"time":%[1]q,"timePointer":%[1]q,"timeOmitZero":%[1]q,"timePointerOmitEmpty":%[1]q}`, now.Format(time.RFC3339)),
			testTime{Time: now, TimePointer: &now, TimeOmitZero: now, TimePointerOmitEmpty: &now},
		},
		{
			fmt.Sprintf(`{"time":%[1]q,"timePointer":%[1]q,"timeOmitZero":%[1]q,"timePointerOmitEmpty":%[1]q}`, now.Format(altLayout)),
			testTime{Time: now, TimePointer: &now, TimeOmitZero: now, TimePointerOmitEmpty: &now},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var out testTime
			if err := json.Unmarshal([]byte(tc.in), &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out.Time, tc.want.Time; !got.Equal(want) {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			if got, want := out.TimeOmitZero, tc.want.TimeOmitZero; !got.Equal(want) {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			for _, tt := range []struct{ got, want *time.Time }{
				{out.TimePointer, tc.want.TimePointer},
				{out.TimePointerOmitEmpty, tc.want.TimePointerOmitEmpty},
			} {
				if tt.got == nil && tt.want != nil {
					t.Fatalf("expected non-nil URLPointer")
				} else if tt.got != nil && tt.want == nil {
					t.Fatalf("expected nil URLPointer")
				} else if tt.got != nil && tt.want != nil {
					if !tt.got.Equal(*tt.want) {
						t.Fatalf("want: %s, got: %s", tt.want, tt.got)
					}
				}
			}
		})
	}

	t.Run("no layout matches", func(t *testing.T) {
		out := &testTime{}
		if err := json.Unmarshal([]byte(`{"time":"not a time"}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		}
	})
}

// A time is written and read as unix milliseconds, to the millisecond, and the zero time as 0; by member name, beside
// times in unix seconds and in the type's own form.
func TestUnixMilliTime(t *testing.T) {
	t.Parallel()

	type trade struct {
		At       time.Time   `json:"at"`
		Read     time.Time   `json:"read"`
		Created  time.Time   `json:"created"`
		Ticks    []time.Time `json:"ticks"`
		Settled  time.Time   `json:"settled"`
		Reported time.Time   `json:"reported"`
	}

	ms, secs := []string{"at", "ticks", "settled"}, []string{"read"}
	opts := json.JoinOptions(
		json.WithMarshalers(json.JoinMarshalers(
			json.MarshalToFunc(jsonutil.TimeMarshalIntUnixMilliFor(ms...)),
			json.MarshalToFunc(jsonutil.TimeMarshalIntUnixFor(secs...)),
		)),
		json.WithUnmarshalers(json.JoinUnmarshalers(
			json.UnmarshalFromFunc(jsonutil.TimeUnmarshalIntUnixMilliFor(ms...)),
			json.UnmarshalFromFunc(jsonutil.TimeUnmarshalIntUnixFor(secs...)),
		)),
	)

	at := time.UnixMilli(1795464914123).UTC()
	in := trade{At: at, Read: at.Truncate(time.Second), Created: at, Ticks: []time.Time{at}, Reported: at}

	b, err := json.Marshal(in, opts)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"at":1795464914123,"read":1795464914,"created":"2026-11-23T20:15:14.123Z","ticks":[1795464914123],` +
		`"settled":0,"reported":"2026-11-23T20:15:14.123Z"}`
	if string(b) != want {
		t.Fatalf("got %s, want %s", b, want)
	}

	var out trade
	if err := json.Unmarshal(b, &out, opts); err != nil {
		t.Fatal(err)
	}

	if !out.At.Equal(at) || !out.Ticks[0].Equal(at) || !out.Read.Equal(in.Read) || !out.Settled.IsZero() {
		t.Errorf("got %+v, want %+v", out, in)
	}

	// without names, every time is unix milliseconds
	var plain time.Time
	if err := json.Unmarshal([]byte(`1795464914123`), &plain, json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.TimeUnmarshalIntUnixMilli))); err != nil || !plain.Equal(at) {
		t.Errorf("got %v (%v), want %v", plain, err, at)
	}
}
