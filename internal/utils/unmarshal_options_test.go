package utils

import (
	"testing"

	"github.com/go-json-experiment/json"
)

// testOptions mimics a typical rule options struct with defaults applied via UnmarshalJSON.
type testOptions struct {
	AllowAny     bool     `json:"allowAny,omitempty"`
	AllowUnknown bool     `json:"allowUnknown,omitempty"`
	Severity     string   `json:"severity,omitempty"`
	Patterns     []string `json:"patterns,omitempty"`
}

func (j *testOptions) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain testOptions
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["allowAny"]; !ok || v == nil {
		plain.AllowAny = true
	}
	if v, ok := raw["allowUnknown"]; !ok || v == nil {
		plain.AllowUnknown = true
	}
	if v, ok := raw["severity"]; !ok || v == nil {
		plain.Severity = "error"
	}
	if v, ok := raw["patterns"]; !ok || v == nil {
		plain.Patterns = []string{}
	}
	*j = testOptions(plain)
	return nil
}

func TestUnmarshalOptionsCaching(t *testing.T) {
	ResetOptionsCache()

	// Create options map (simulates what headless mode provides)
	opts := map[string]interface{}{
		"allowAny": false,
		"severity": "warn",
	}

	// First call should unmarshal
	result1 := UnmarshalOptions[testOptions](opts, "test-rule")
	if result1.AllowAny != false {
		t.Errorf("expected AllowAny=false, got %v", result1.AllowAny)
	}
	if result1.AllowUnknown != true {
		t.Errorf("expected AllowUnknown=true (default), got %v", result1.AllowUnknown)
	}
	if result1.Severity != "warn" {
		t.Errorf("expected Severity=warn, got %v", result1.Severity)
	}

	// Second call with same map should return cached result
	result2 := UnmarshalOptions[testOptions](opts, "test-rule")
	if result2.AllowAny != result1.AllowAny || result2.Severity != result1.Severity {
		t.Error("cached result differs from first result")
	}

	// Nil options should also be cacheable (standard mode)
	resultNil1 := UnmarshalOptions[testOptions](nil, "test-rule")
	if resultNil1.AllowAny != true {
		t.Errorf("expected AllowAny=true (default for nil), got %v", resultNil1.AllowAny)
	}
	resultNil2 := UnmarshalOptions[testOptions](nil, "test-rule")
	if resultNil2.AllowAny != resultNil1.AllowAny {
		t.Error("cached nil result differs from first nil result")
	}

	// Different map should produce different result
	opts2 := map[string]interface{}{
		"severity": "off",
	}
	result3 := UnmarshalOptions[testOptions](opts2, "test-rule")
	if result3.Severity != "off" {
		t.Errorf("expected Severity=off, got %v", result3.Severity)
	}
}

// BenchmarkUnmarshalOptionsUncached measures the cost of UnmarshalOptions without caching.
func BenchmarkUnmarshalOptionsUncached(b *testing.B) {
	opts := map[string]interface{}{
		"allowAny":     false,
		"allowUnknown": true,
		"severity":     "warn",
		"patterns":     []string{"foo", "bar"},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		ResetOptionsCache() // clear cache each iteration to force re-unmarshal
		_ = UnmarshalOptions[testOptions](opts, "test-rule")
	}
}

// BenchmarkUnmarshalOptionsCached measures the cost of UnmarshalOptions with cache hits.
func BenchmarkUnmarshalOptionsCached(b *testing.B) {
	opts := map[string]interface{}{
		"allowAny":     false,
		"allowUnknown": true,
		"severity":     "warn",
		"patterns":     []string{"foo", "bar"},
	}

	ResetOptionsCache()
	// Warm up the cache
	_ = UnmarshalOptions[testOptions](opts, "test-rule")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = UnmarshalOptions[testOptions](opts, "test-rule")
	}
}

// BenchmarkUnmarshalOptionsNilUncached measures the nil-options path without caching.
func BenchmarkUnmarshalOptionsNilUncached(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		ResetOptionsCache()
		_ = UnmarshalOptions[testOptions](nil, "test-rule")
	}
}

// BenchmarkUnmarshalOptionsNilCached measures the nil-options path with caching.
func BenchmarkUnmarshalOptionsNilCached(b *testing.B) {
	ResetOptionsCache()
	_ = UnmarshalOptions[testOptions](nil, "test-rule")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = UnmarshalOptions[testOptions](nil, "test-rule")
	}
}

// BenchmarkUnmarshalOptionsSimulated simulates headless mode: 30 rules × 1000 files
// with shared options maps (as would happen when files share the same config).
func BenchmarkUnmarshalOptionsSimulated(b *testing.B) {
	const numRules = 30
	const numFiles = 1000

	// Create 30 different options maps (one per rule)
	ruleOptions := make([]map[string]interface{}, numRules)
	for i := range ruleOptions {
		ruleOptions[i] = map[string]interface{}{
			"allowAny": i%2 == 0,
			"severity": "warn",
		}
	}

	b.Run("no_cache", func(b *testing.B) {
		// Simulate no caching: create a fresh options map for every file×rule call
		// to prevent pointer-identity cache hits.
		b.ReportAllocs()
		for b.Loop() {
			ResetOptionsCache()
			for file := range numFiles {
				_ = file
				for i := range numRules {
					// Fresh map each time = different pointer = no cache hit
					opts := map[string]interface{}{
						"allowAny": i%2 == 0,
						"severity": "warn",
					}
					_ = UnmarshalOptions[testOptions](opts, "test-rule")
				}
			}
		}
	})

	b.Run("with_cache", func(b *testing.B) {
		// Simulate caching: reuse the same options maps (as headless mode does).
		// Cache is populated on the first file, remaining 999 files hit the cache.
		b.ReportAllocs()
		for b.Loop() {
			ResetOptionsCache()
			for file := range numFiles {
				_ = file
				for _, opts := range ruleOptions {
					_ = UnmarshalOptions[testOptions](opts, "test-rule")
				}
			}
		}
	})
}
