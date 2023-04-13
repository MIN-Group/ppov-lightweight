package test

import (
	"fmt"
	"testing"
)

func TestShift(t *testing.T) {
	var a []int
	var b []uint8

	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			a = append(a, 1)
		} else {
			a = append(a, -1)
		}
	}
	a = append(a, 0, 1, 0, -1, 1, 1, 0, -1, -1, -1, 0, 0, 0, -1, -1, -1)
	fmt.Println(len(a))
	var tmp uint8 = 0
	for i := 1; i <= len(a); i++ {

		if a[i-1] == 0 {
			tmp = tmp << 1
			tmp = tmp | 1
			tmp = tmp << 1
		} else if a[i-1] == 1 {
			tmp = tmp << 2
			tmp = tmp | 1
		} else if a[i-1] == -1 {
			tmp = tmp << 1
			tmp = tmp | 1
			tmp = tmp << 1
			tmp = tmp | 1
		}

		if i%3 == 0 || i == len(a) {
			b = append(b, tmp)
		}

		if i%3 == 0 {
			tmp = 0
		}
	}

	fmt.Println(b)

	var c []int
	for i := 0; i < len(b); i++ {
		var s []int
		for j := 0; j < 6; j += 2 {
			tmp1, tmp2 := int(b[i]&1), int(b[i]&2>>1)
			if tmp1 == 1 && tmp2 == 0 {
				s = append(s, 1)
			} else if tmp1 == 1 && tmp2 == 1 {
				s = append(s, -1)
			} else if tmp1 == 0 && tmp2 == 1 {
				s = append(s, 0)
			}
			b[i] = b[i] >> 2
		}
		for k := len(s) - 1; k >= 0; k-- {
			c = append(c, s[k])
		}
	}

	fmt.Println(a)
	fmt.Println(c)

	fmt.Println(len(a))
	fmt.Println(len(b))
}


