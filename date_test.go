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

	"cloud.google.com/go/civil"
	"github.com/MarkRosemaker/jsonutil"
)

func TestUnixDate(t *testing.T) {
	type testDate struct {
		Date                 civil.Date  `json:"date"`
		DatePointer          *civil.Date `json:"datePointer"`
		DateOmitZero         civil.Date  `json:"dateOmitZero,omitzero"`
		DatePointerOmitEmpty *civil.Date `json:"datePointerOmitEmpty,omitempty"`
	}

	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalToFunc(jsonutil.DateMarshalIntUnix)),
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.DateUnmarshalIntUnix)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &testDate{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"date":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not an int", func(t *testing.T) {
		out := &testDate{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"date": "3"}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpInt := reflect.TypeFor[int64](); errSem.GoType != tpInt {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpInt, errSem.GoType)
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testDate{}
		if err := json.Unmarshal([]byte(`{"date":null,"datePointer":null,"dateOmitZero":null,"datePointerOmitEmpty":null}`), out, jsonOpts); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	now := time.Now()
	today := civil.DateOf(now)
	todayUnix := now.Truncate(time.Hour * 24).Unix()

	for i, tc := range []struct {
		in  testDate
		out string
	}{
		{testDate{}, `{"date":0,"datePointer":null}`},
		{
			testDate{Date: today, DatePointer: &today, DateOmitZero: today, DatePointerOmitEmpty: &today},
			fmt.Sprintf(`{"date":%[1]d,"datePointer":%[1]d,"dateOmitZero":%[1]d,"datePointerOmitEmpty":%[1]d}`, todayUnix),
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

			var out testDate
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out.Date, tc.in.Date; got.Compare(want) != 0 {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			if got, want := out.DateOmitZero, tc.in.DateOmitZero; got.Compare(want) != 0 {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			for _, tt := range []struct{ got, want *civil.Date }{
				{out.DatePointer, tc.in.DatePointer},
				{out.DatePointerOmitEmpty, tc.in.DatePointerOmitEmpty},
			} {
				if tt.got == nil && tt.want != nil {
					t.Fatalf("expected non-nil URLPointer")
				} else if tt.got != nil && tt.want == nil {
					t.Fatalf("expected nil URLPointer")
				} else if tt.got != nil && tt.want != nil {
					if tt.got.Compare(*tt.want) != 0 {
						t.Fatalf("want: %s, got: %s", tt.want, tt.got)
					}
				}
			}
		})
	}
}
