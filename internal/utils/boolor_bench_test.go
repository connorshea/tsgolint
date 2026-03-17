package utils

import (
	"testing"

	"github.com/go-json-experiment/json"
)

// BenchmarkBoolOrUnmarshalBool measures unmarshaling a boolean value into BoolOr.
func BenchmarkBoolOrUnmarshalBool(b *testing.B) {
	data := []byte(`true`)
	b.ReportAllocs()
	for b.Loop() {
		var v BoolOr[struct{ Foo string }]
		_ = json.Unmarshal(data, &v)
	}
}

// BenchmarkBoolOrUnmarshalObject measures unmarshaling an object value into BoolOr.
func BenchmarkBoolOrUnmarshalObject(b *testing.B) {
	type detail struct {
		AllowAny     bool   `json:"allowAny"`
		AllowUnknown bool   `json:"allowUnknown"`
		Severity     string `json:"severity"`
	}
	data := []byte(`{"allowAny":true,"allowUnknown":false,"severity":"warn"}`)
	b.ReportAllocs()
	for b.Loop() {
		var v BoolOr[detail]
		_ = json.Unmarshal(data, &v)
	}
}

// BenchmarkBoolOrUnmarshalNull measures unmarshaling null into BoolOr.
func BenchmarkBoolOrUnmarshalNull(b *testing.B) {
	data := []byte(`null`)
	b.ReportAllocs()
	for b.Loop() {
		var v BoolOr[struct{ Foo string }]
		_ = json.Unmarshal(data, &v)
	}
}

// BenchmarkBoolOrMarshalBool measures marshaling a boolean BoolOr value.
func BenchmarkBoolOrMarshalBool(b *testing.B) {
	v := BoolOrValue[struct{ Foo string }](true)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(v)
	}
}

// BenchmarkBoolOrMarshalObject measures marshaling an object BoolOr value.
func BenchmarkBoolOrMarshalObject(b *testing.B) {
	type detail struct {
		AllowAny     bool   `json:"allowAny"`
		AllowUnknown bool   `json:"allowUnknown"`
		Severity     string `json:"severity"`
	}
	v := BoolOr[detail]{isSet: true, boolVal: true, objectVal: &detail{AllowAny: true, Severity: "warn"}}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(v)
	}
}

// BenchmarkBoolOrMarshalNull measures marshaling an unset BoolOr value.
func BenchmarkBoolOrMarshalNull(b *testing.B) {
	var v BoolOr[struct{ Foo string }]
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(v)
	}
}
