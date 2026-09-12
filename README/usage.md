To install the library, use the following command:

```shell
go get github.com/MarkRosemaker/jsonutil
```


### Custom Marshaling and Unmarshaling for `url.URL`

To use the custom marshaler and unmarshaler for `url.URL`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"
	"net/url"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.MarshalToFunc(jsonutil.URLMarshal)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.URLUnmarshal)),
)

type MyStruct struct {
	Foo  string  `json:"foo"`
	Link url.URL `json:"link"`
}

func main() {
	out := &MyStruct{}
	if err := json.Unmarshal([]byte(`{"foo":"bar","link":"https://example.com/"}`), out, jsonOpts); err != nil {
		panic(err)
	}

	if out.Link.String() != "https://example.com/" {
		panic("something went wrong")
	}

	res, err := json.Marshal(MyStruct{
		Foo:  "baz",
		Link: url.URL{Scheme: "http", Host: "example.com", Path: "/some-path"},
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"foo":"baz","link":"http://example.com/some-path"}
}
```

### Custom Marshaling and Unmarshaling for `time.Duration`

To use the custom marshaler and unmarshaler for `time.Duration`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"
	"time"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.MarshalToFunc(jsonutil.DurationMarshalIntSeconds)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.DurationUnmarshalIntSeconds)),
)

type MyStruct struct {
	Foo      string        `json:"foo"`
	Duration time.Duration `json:"duration"`
}

func main() {
	out := &MyStruct{}
	if err := json.Unmarshal([]byte(`{"foo":"bar","duration":3}`), out, jsonOpts); err != nil {
		panic(err)
	}

	if out.Duration != 3*time.Second {
		panic("something went wrong")
	}

	res, err := json.Marshal(MyStruct{
		Foo:      "baz",
		Duration: time.Minute,
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"foo":"baz","duration":60}
}
```

### Custom Marshaling for Maps with Ordered Keys

To use the custom marshaler for `M` where `[M ~map[K]V, K cmp.Ordered, V any]`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"

	"github.com/MarkRosemaker/jsonutil"
)

type myMap map[string]int

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.MarshalToFunc(jsonutil.OrderedMapMarshal[myMap])),
)

func main() {
	res, err := json.Marshal(myMap{
		"foo": 1,
		"bar": 2,
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"bar":2,"bar":1}
}
```

### Custom Marshaling and Unmarshaling for `http.Header`

To use the custom marshaler and unmarshaler for `http.Header`, you can import the package and use the provided functions:

```go
package main

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/MarkRosemaker/jsonutil"
)

var jsonOpts = json.JoinOptions(
	json.WithMarshalers(json.MarshalToFunc(jsonutil.HTTPHeaderMarshal)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.HTTPHeaderUnmarshal)),
)

func main() {
	out := &http.Header{}
	if err := json.Unmarshal([]byte(`{"foo":"bar","baz":"quux"}`), out, jsonOpts); err != nil {
		panic(err)
	}

	res, err := json.Marshal(http.Header{
		"foo": []string{"bar"},
		"baz": []string{"quux"},
	}, jsonOpts)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
	// Output: {"Foo":"bar","Baz":"quux"}
}
```
