package ping

import (
	"fmt"
	"testing"
)

func TestGetScore(t *testing.T) {
	var tests = []struct {
		count, failCount, averagePing int
		want                          int
	}{
		{0, 0, 0, 0},
		{10, 0, 0, 0},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d(%d)%d", tt.count, tt.failCount, tt.averagePing)
		t.Run(testname, func(t *testing.T) {
			ans := getScore(tt.count, tt.failCount, tt.averagePing)
			if ans != tt.want {
				t.Errorf("got %f, want %f", ans, tt.want)
			}
		})
	}
}
