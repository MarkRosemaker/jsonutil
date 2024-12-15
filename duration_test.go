package jsonutil_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/go-json-experiment/json"
)

func TestDuration(t *testing.T) {
	type testDuration struct {
		Duration                 time.Duration  `json:"duration"`
		DurationPointer          *time.Duration `json:"durationPointer"`
		DurationOmitZero         time.Duration  `json:"durationOmitZero,omitzero"`
		DurationPointerOmitEmpty *time.Duration `json:"durationPointerOmitEmpty,omitempty"`
	}

	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalFuncV2(jsonutil.DurationMarshalIntSeconds)),
		json.WithUnmarshalers(json.UnmarshalFuncV2(jsonutil.DurationUnmarshalIntSeconds)),
	)

	t.Run("errors", func(t *testing.T) {
		tpDuration := reflect.PointerTo(reflect.TypeFor[time.Duration]())

		for _, data := range []string{
			`{"duration":`,      // EOF
			`{"duration": "3"}`, // not an int
		} {
			t.Run(data, func(t *testing.T) {
				out := &testDuration{}
				errSem := &json.SemanticError{}
				if err := json.Unmarshal([]byte(data), out, jsonOpts); err == nil {
					t.Fatalf("expected error")
				} else if !errors.As(err, &errSem) {
					t.Fatalf("expected error to be a semantic error, got: %v", err)
				} else if errSem.GoType != tpDuration {
					t.Fatalf("expected semantic error to have type %s, got: %s", tpDuration, errSem.GoType)
				}
			})
		}
	})

	t.Run("null", func(t *testing.T) {
		out := &testDuration{}
		if err := json.Unmarshal([]byte(`{"duration":null,"durationPointer":null,"durationOmitZero":null,"durationPointerOmitEmpty":null}`), out, jsonOpts); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	u := 30 * time.Second

	for i, tc := range []struct {
		in  testDuration
		out string
	}{
		{testDuration{}, `{"duration":0,"durationPointer":null}`},
		{
			testDuration{Duration: u, DurationPointer: &u, DurationOmitZero: u, DurationPointerOmitEmpty: &u},
			`{"duration":30,"durationPointer":30,"durationOmitZero":30,"durationPointerOmitEmpty":30}`,
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

			var out testDuration
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got, want := out.Duration.String(), tc.in.Duration.String(); got != want {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			if got, want := out.DurationOmitZero.String(), tc.in.DurationOmitZero.String(); got != want {
				t.Fatalf("want: %s, got: %s", want, got)
			}

			for _, tt := range []struct{ got, want *time.Duration }{
				{out.DurationPointer, tc.in.DurationPointer},
				{out.DurationPointerOmitEmpty, tc.in.DurationPointerOmitEmpty},
			} {
				if tt.got == nil && tt.want != nil {
					t.Fatalf("expected non-nil URLPointer")
				} else if tt.got != nil && tt.want == nil {
					t.Fatalf("expected nil URLPointer")
				} else if tt.got != nil && tt.want != nil {
					if tt.got.String() != tt.want.String() {
						t.Fatalf("want: %s, got: %s", tt.want, tt.got)
					}
				}
			}
		})
	}
}