* Custom marshaler and unmarshaler for `url.URL`:
  * `URLMarshal` marshals `url.URL` as a string.
  * `URLUnmarshal` unmarshals `url.URL` from a string.
* Custom marshalers and unmarshalers for `time.Duration`:
  * `DurationMarshalIntSeconds` marshals `time.Duration` as an integer representing seconds.
  * `DurationUnmarshalIntSeconds` unmarshals `time.Duration` from an integer assuming it represents seconds.
  * `DurationMarshalString` marshals `time.Duration` as a string in its canonical units, such as `"1h30m"` or `"500ms"`.
  * `DurationUnmarshalString` unmarshals `time.Duration` from such a string, parsed with `time.ParseDuration`.
* Custom marshalers and unmarshalers for `time.Time`:
  * `TimeMarshalIntUnix` marshals `time.Time` as an integer representing unix time.
  * `TimeUnmarshalIntUnix` unmarshals `time.Time` from an integer assuming it represents unix time.
  * `TimeUnmarshalStringOrIntUnix` unmarshals `time.Time` from an RFC 3339 string, or a string such as
    `"Mon Jan 2 2006 15:04:05 MST-0700"`, or an integer representing unix time. Null is the zero time.
  * `TimeMarshalIntUnixFor(names...)` and `TimeUnmarshalIntUnixFor(names...)` return a marshaler and an unmarshaler
    that do as `TimeMarshalIntUnix` and `TimeUnmarshalIntUnix` do, but only for the times held by an object member of
    one of the names given, or by an array such a member holds. Any other time is left to the next (un)marshaler, or to
    the default string: for an API that writes some times as unix time and others as strings.
* Custom marshalers and unmarshalers for `civil.Date`:
  * `DateMarshalIntUnix` marshals `civil.Date` as an integer representing the unix time of its midnight UTC.
  * `DateUnmarshalIntUnix` unmarshals `civil.Date` from such an integer, reading the day in UTC.
  * `DateMarshalIntUnixFor(names...)` and `DateUnmarshalIntUnixFor(names...)` do as the two above, but only for the
    dates held by an object member of one of the names given, as `TimeMarshalIntUnixFor` does for times.
* Custom unmarshaler for `float64`:
  * `Float64UnmarshalStringOrNumber` unmarshals `float64` from a number or from a string holding one, such as `"9.5"`:
    for an API that sends numbers as strings, but for the odd number. Null leaves the value as it is. With a field's
    `string` option, the value is marshaled as a string too.
* Custom marshaler for maps with ordered keys:
  * `OrderedMapMarshal[M ~map[K]V, K cmp.Ordered, V any]` marshals `M` so that the keys are sorted.
* Custom marshaler for `http.Header`:
  * `HTTPHeaderMarshal` marshals the values of `http.Header` as single strings.
  * `HTTPHeaderUnmarshal` unmarshals the values of `http.Header` from single strings.
* Reading and writing files:
  * `ReadFile[T any](name, opts...)` reads the JSON file `name` and unmarshals it into a `T`.
  * `WriteFile[T any](name, data, opts...)` marshals `data` and writes it to the JSON file `name`.
