//
// Copyright (c) 2019 Konstantin Ivanov <kostyarin.ivanov@gmail.com>.
// All rights reserved. This program is free software. It comes without
// any warranty, to the extent permitted by applicable law. You can
// redistribute it and/or modify it under the terms of the Do What
// The Fuck You Want To Public License, Version 2, as published by
// Sam Hocevar. See LICENSE file for more details or see below.
//

//
//        DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//                    Version 2, December 2004
//
// Copyright (C) 2004 Sam Hocevar <sam@hocevar.net>
//
// Everyone is permitted to copy and distribute verbatim or modified
// copies of this license document, and changing it is allowed as long
// as the name is changed.
//
//            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION
//
//  0. You just DO WHAT THE FUCK YOU WANT TO.
//

package srbt

import (
	"fmt"
	"math/rand"
	"testing"
)

var globalTree *Tree

func newNaturalRange(from, to int) (t *Tree) {
	t = newNatural()
	for ; from < to; from++ {
		t.Ins(from)
	}
	return
}

func randomRange(from, to int) (is []int) {
	is = make([]int, 0, to-from)
	for i := 0; from+i < to; i++ {
		is = append(is, from+i)
	}
	rand.Shuffle(len(is), func(i, j int) {
		is[i], is[j] = is[j], is[i]
	})
	return
}

func BenchmarkNew(b *testing.B) {
	// New(less LessFunc, equal EqualFunc, zero ZeroFunc) (t *Tree)
	for i := 0; i < b.N; i++ {
		globalTree = newNatural()
	}
	b.ReportAllocs()
}

func BenchmarkTree_Get(b *testing.B) {
	// Get(item interface{}) bool

	for _, n := range []int{
		1,
		10,
		100,
		1 * 1000,
		10 * 1000,
		100 * 1000,
		1000 * 1000,
	} {
		var (
			tree = newNaturalRange(0, n)
			ns   = fmt.Sprintf(" [0, %d]", n)
		)
		b.Run("missing successively"+ns, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if tree.Get(i+n) == true {
					b.Error("got item not exists")
				}
			}
			b.ReportAllocs()
		})
		b.Run("missing random"+ns, func(b *testing.B) {
			var rr = randomRange(b.N, n+b.N)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if tree.Get(rr[i]) == true {
					b.Error("got item not exists")
				}
			}
			b.ReportAllocs()
		})
		b.Run("found successively"+ns, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if tree.Get(i) == true {
					b.Error("got item not exists")
				}
			}
			b.ReportAllocs()
		})
		b.Run("missing random"+ns, func(b *testing.B) {
			var rr = randomRange(0, n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if tree.Get(rr[i]) == true {
					b.Error("got item not exists")
				}
			}
			b.ReportAllocs()
		})
	}

}

func BenchmarkTree_Ins(b *testing.B) {
	// Ins(item interface{}) (ok bool)
}

func BenchmarkTree_InsNx(b *testing.B) {
	// InsNx(item interface{}) (ok bool)
}

func BenchmarkTree_InsEx(b *testing.B) {
	// InsEx(item interface{}) (ok bool)
}

func BenchmarkTree_Del(b *testing.B) {
	// Del(item interface{}) (ok bool)
}

func BenchmarkTree_Size(b *testing.B) {
	// Size() int
}

func BenchmarkTree_Clear(b *testing.B) {
	// Clear()
}

func BenchmarkTree_Min(b *testing.B) {
	// Min() (interface{}, bool)
}

func BenchmarkTree_Max(b *testing.B) {
	// Max() (interface{}, bool)
}

func BenchmarkTree_Walk(b *testing.B) {
	// Walk(walkFunc WalkFunc)
}

func BenchmarkTree_Ascend(b *testing.B) {
	// Ascend(from, to interface{}, ascendFunc WalkFunc)
}

func BenchmarkTree_Descend(b *testing.B) {
	// Descend(from, to interface{}, descendFunc WalkFunc)
}
