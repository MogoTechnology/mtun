package ping

import (
	"fmt"
	"math"
	"testing"
)

func TestGetScore(t *testing.T) {
	const epsilon float64 = 1e-3
	var tests = []struct {
		count, failCount int64
		averagePing      int64
		wantScore        float64
	}{
		// 全部成功，Ping 值很低，分值为 5
		{1, 0, 1, 6.99},
		{2, 0, 1, 6.99},
		{3, 0, 1, 6.99},
		{4, 0, 1, 6.99},
		{5, 0, 1, 6.99},
		{10, 0, 1, 6.99},
		{100, 0, 1, 6.99},

		// 4次Ping，Ping 值很低
		{4, 0, 1, 6.99},
		{4, 1, 1, 5.49},
		{4, 2, 1, 3.49},
		{4, 3, 1, 1.49},
		{4, 4, 0, -7.5},

		// 4次Ping，Ping 值较高
		{4, 0, 699, 0.01},
		{4, 1, 699, -1.49},
		{4, 2, 699, -3.49},
		{4, 3, 699, -5.49},
		{4, 4, 0, -7.5},

		// 4次Ping，全部成功，Ping 值变化
		{4, 0, 1, 6.99},
		{4, 0, 10, 6.9},
		{4, 0, 100, 6.0},
		{4, 0, 200, 5.0},
		{4, 0, 300, 4.0},
		{4, 0, 500, 2.0},
		{4, 0, 699, 0.01},
		{4, 0, 700, 1},  // avePing >= m*k 时，basicScore 保持初始值 1
		{4, 0, 701, 1},
		{4, 0, 800, 1},
		{4, 0, 1000, 1},
		{4, 0, 10000, 1},

		// 单次Ping
		{1, 0, 1, 6.99},
		{1, 0, 500, 2.0},
		{1, 0, 699, 0.01},
		{1, 0, 700, 1},  // avePing >= m*k 时，basicScore 保持初始值 1
		{1, 0, 701, 1},
		{1, 0, 1000, 1},
		{1, 0, 10000, 1},
		{1, 1, 0, -7.5},

		// 边界值，实际测试不会出现，如：测试次数为0或负数，失败次数大于测试次数，全部失败时ping值非0。
		{0, 0, 1, 6.99},
		{0, 0, 0, 0},
		{0, 1, 0, 0},
		{0, -1, 0, 0},
		{-1, 0, 1, 6.99},
		{1, 0, 0, 0},
		{1, 1, 1, -0.51},
		{1, 1, 1000, -6.5},
		{4, 10, 100, -13.5},
		{4, 0, -100, 8.0},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d(%d)%d", tt.count, tt.failCount, tt.averagePing)
		t.Run(testname, func(t *testing.T) {
			gotScore := getScore(tt.count, tt.failCount, tt.averagePing)
			if math.Abs(gotScore-tt.wantScore) > epsilon {
				t.Errorf("getScore(%d, %d, %d) = %v, want %v", tt.count, tt.failCount, tt.averagePing, gotScore, tt.wantScore)
			}
		})
	}
}
