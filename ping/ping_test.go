package ping

import (
	"fmt"
	"math"
	"testing"
)

func TestGetScore(t *testing.T) {
	const epsilon float64 = 1e-9
	var tests = []struct {
		count, failCount int64
		averagePing      int64
		want             float64
	}{
		{0, 0, 0, 0.0},
		{10, 0, 1, 4.995},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d(%d)%d", tt.count, tt.failCount, tt.averagePing)
		t.Run(testname, func(t *testing.T) {
			ans := getScore(tt.count, tt.failCount, tt.averagePing)
			if math.Abs(ans-tt.want) > epsilon {
				t.Errorf("getScore(%d, %d, %d) = %v, want %v", tt.count, tt.failCount, tt.averagePing, ans, tt.want)
			}
		})
	}
}
