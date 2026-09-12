* Custom marshaler and unmarshaler for `url.URL`:
  * `URLMarshal` marshals `url.URL` as a string.
  * `URLUnmarshal` unmarshals `url.URL` from a string.
* Custom marshaler and unmarshaler for `time.Duration`:
  * `DurationMarshalIntSeconds` marshals `time.Duration` as an integer representing seconds.
  * `DurationUnmarshalIntSeconds` unmarshals `time.Duration` from an integer assuming it represents seconds.
* Custom marshaler for maps with ordered keys:
  * `OrderedMapMarshal[M ~map[K]V, K cmp.Ordered, V any]` marshals `M` so that the keys are sorted.
* Custom marshaler for `http.Header`:
  * `HTTPHeaderMarshal` marshals the values of `http.Header` as single strings.
  * `HTTPHeaderUnmarshal` unmarshals the values of `http.Header` from single strings.
