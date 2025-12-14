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

func TestUnixTime(t *testing.T) {
	type testTime struct {
		Time                 time.Time  `json:"time"`
		TimePointer          *time.Time `json:"timePointer"`
		TimeOmitZero         time.Time  `json:"timeOmitZero,omitzero"`
		TimePointerOmitEmpty *time.Time `json:"timePointerOmitEmpty,omitempty"`
	}

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
		{testTime{}, `{"time":-62135596800,"timePointer":null}`},
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
