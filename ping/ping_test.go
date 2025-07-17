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
		{1, 0, 1, 4.995},
		{2, 0, 1, 4.995},
		{3, 0, 1, 4.995},
		{4, 0, 1, 4.995},
		{5, 0, 1, 4.995},
		{10, 0, 1, 4.995},
		{100, 0, 1, 4.995},

		// 4次Ping，Ping 值很低
		{4, 0, 1, 4.995},
		{4, 1, 1, 3.495},
		{4, 2, 1, 1.495},
		{4, 3, 1, -0.505},
		{4, 4, 0, -7.5},

		// 4次Ping，Ping 值较高
		{4, 0, 999, 0.005},
		{4, 1, 999, -1.495},
		{4, 2, 999, -3.495},
		{4, 3, 999, -5.495},
		{4, 4, 0, -7.5},

		// 4次Ping，全部成功，Ping 值变化
		{4, 0, 1, 4.995},
		{4, 0, 10, 4.95},
		{4, 0, 100, 4.5},
		{4, 0, 200, 4.0},
		{4, 0, 300, 3.5},
		{4, 0, 500, 2.5},
		{4, 0, 800, 1},
		{4, 0, 999, 0.005},
		{4, 0, 1000, 1}, // 1000及以上时跳变为1
		{4, 0, 1001, 1},
		{4, 0, 10000, 1},

		// 单次Ping
		{1, 0, 1, 4.995},
		{1, 0, 500, 2.5},
		{1, 0, 999, 0.005},
		{1, 0, 1000, 1}, // 1000及以上时跳变为1
		{1, 0, 10000, 1},
		{1, 1, 0, -7.5},

		// 边界值，实际测试不会出现，如：测试次数为0或负数，失败次数大于测试次数，全部失败时ping值非0。
		{0, 0, 1, 4.995},
		{0, 0, 0, 0},
		{0, 1, 0, 0},
		{0, -1, 0, 0},
		{-1, 0, 1, 4.995},
		{1, 0, 0, 0},
		{1, 1, 1, -2.505},
		{1, 1, 1000, -6.5},
		{4, 10, 100, -15},
		{4, 0, -100, 5.5},
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
