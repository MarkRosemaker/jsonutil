package jsonutil_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"maps"
	"reflect"
	"strconv"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
)

type orderedMap map[string]int

func TestMap(t *testing.T) {
	jsonOpts := json.JoinOptions(
		json.WithMarshalers(json.MarshalToFunc(jsonutil.OrderedMapMarshal[orderedMap])),
		json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.DurationUnmarshalIntSeconds)),
	)

	t.Run("EOF", func(t *testing.T) {
		out := &orderedMap{}
		errSyn := &jsontext.SyntacticError{}

		if err := json.Unmarshal([]byte(`{"foo":`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSyn) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if want := `unexpected EOF`; errSyn.Err.Error() != want {
			t.Fatalf("expected syntactic error be %s, got: %#v", want, errSyn.Err)
		}
	})

	t.Run("not an int", func(t *testing.T) {
		out := &orderedMap{}
		errSem := &json.SemanticError{}

		if err := json.Unmarshal([]byte(`{"foo": "3"}`), out, jsonOpts); err == nil {
			t.Fatalf("expected error")
		} else if !errors.As(err, &errSem) {
			t.Fatalf("expected error to be a semantic error, got: %v", err)
		} else if tpInt := reflect.TypeFor[int](); errSem.GoType != tpInt {
			t.Fatalf("expected semantic error to have type %s, got: %s", tpInt, errSem.GoType)
		}
	})

	for i, tc := range []struct {
		in  orderedMap
		out string
	}{
		{orderedMap{"foo": 0, "bar": 3}, `{"bar":3,"foo":0}`},
		{orderedMap{"bar": 3, "foo": 0}, `{"bar":3,"foo":0}`},
		{nil, `null`},
		{orderedMap{}, `{}`},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			b, err := json.Marshal(tc.in, jsonOpts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != tc.out {
				t.Fatalf("want: %s, got: %s", tc.out, string(b))
			}

			var out orderedMap
			if err := json.Unmarshal(b, &out, jsonOpts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !maps.Equal(tc.in, out) {
				t.Fatalf("want: %v, got: %v", tc.in, out)
			}
		})
	}
}
