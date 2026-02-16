package utils

import (
	"math"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64
	}{
		{
			name:     "London to Paris (~344 km)",
			lat1:     51.5074,
			lon1:     -0.1278,
			lat2:     48.8566,
			lon2:     2.3522,
			expected: 344.0,
			delta:    5.0,
		},
		{
			name:     "Same location",
			lat1:     51.5074,
			lon1:     -0.1278,
			lat2:     51.5074,
			lon2:     -0.1278,
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "New York to Los Angeles (~3944 km)",
			lat1:     40.7128,
			lon1:     -74.0060,
			lat2:     34.0522,
			lon2:     -118.2437,
			expected: 3944.0,
			delta:    50.0,
		},
		{
			name:     "Short distance (~5 km)",
			lat1:     51.5074,
			lon1:     -0.1278,
			lat2:     51.5155,
			lon2:     -0.0922,
			expected: 3.0,
			delta:    1.0,
		},
		{
			name:     "Equator crossing",
			lat1:     -10.0,
			lon1:     0.0,
			lat2:     10.0,
			lon2:     0.0,
			expected: 2222.0,
			delta:    10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)

			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("HaversineDistance() = %.2f, want %.2f (±%.2f)", result, tt.expected, tt.delta)
			}
		})
	}
}

func TestDegreesToRadians(t *testing.T) {
	tests := []struct {
		name     string
		degrees  float64
		expected float64
	}{
		{
			name:     "Zero degrees",
			degrees:  0,
			expected: 0,
		},
		{
			name:     "90 degrees",
			degrees:  90,
			expected: math.Pi / 2,
		},
		{
			name:     "180 degrees",
			degrees:  180,
			expected: math.Pi,
		},
		{
			name:     "360 degrees",
			degrees:  360,
			expected: 2 * math.Pi,
		},
		{
			name:     "Negative degrees",
			degrees:  -90,
			expected: -math.Pi / 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := degreesToRadians(tt.degrees)

			if math.Abs(result-tt.expected) > 0.000001 {
				t.Errorf("degreesToRadians(%v) = %v, want %v", tt.degrees, result, tt.expected)
			}
		})
	}
}

func BenchmarkHaversineDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HaversineDistance(51.5074, -0.1278, 48.8566, 2.3522)
	}
}