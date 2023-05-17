package server

import "strconv"

var curr = 2

func Next() int {
	curr++
	if curr > 254 {
		curr = 2
	}
	return curr
}

func NextIP() string {
	return "10.10.0." + strconv.Itoa(Next())
}
