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

package srb

import (
	"testing"
)

var (
	globalTree      *Tree
	globalOK        bool
	globalSize      int
	globalInterface interface{}
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		globalTree = newNatiral()
	}
	b.ReportAllocs()
}

func BenchmarkTree_Ins(b *testing.B) {
	var tr = newNatiral()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, globalOK = tr.Ins(i, i)
	}
	b.ReportAllocs()
}

func BenchmarkTree_InsNx(b *testing.B) {
	var tr = newNatiral()
	b.Run("does not exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.InsNx(i, i)
		}
		b.ReportAllocs()
	})
	b.Run("exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.InsNx(i, i)
		}
		b.ReportAllocs()
	})
}

func BenchmarkTree_InsEx(b *testing.B) {
	var tr = newNatiral()
	b.Run("does not exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.InsEx(i, i)
		}
		b.ReportAllocs()
	})
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.Run("exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.InsEx(i, i)
		}
		b.ReportAllocs()
	})
}

func BenchmarkTree_Add(b *testing.B) {
	var tr = newNatiral()
	b.Run("does not exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			globalOK = tr.Add(i, i)
		}
		b.ReportAllocs()
	})
	b.Run("exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			globalOK = tr.Add(i, i)
		}
		b.ReportAllocs()
	})
}

func BenchmarkTree_Get(b *testing.B) {
	var tr = newNatiral()
	b.Run("does not exist (blank)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.Get(i)
		}
		b.ReportAllocs()
	})
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.Run("exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, globalOK = tr.Get(i)
		}
		b.ReportAllocs()
	})
	b.Run("does not exist (full)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if i%2 != 0 {
				_, globalOK = tr.Get(i + b.N)
			} else {
				_, globalOK = tr.Get(-i)
			}
		}
		b.ReportAllocs()
	})
}

func BenchmarkTree_Del(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, globalOK = tr.Del(i)
	}
	b.ReportAllocs()
}

func BenchmarkTree_Min(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, globalOK = tr.Min()
	}
	b.ReportAllocs()
}

func BenchmarkTree_Max(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, globalOK = tr.Max()
	}
	b.ReportAllocs()
}

func BenchmarkTree_Walk(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Walk(func(k, v interface{}) bool {
			globalInterface, globalInterface = k, v
			return true
		})
	}
	b.ReportAllocs()
}

func BenchmarkTree_Ascend(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Ascend(0, 0, func(k, v interface{}) bool {
			globalInterface, globalInterface = k, v
			return true
		})
	}
	b.ReportAllocs()
}

func BenchmarkTree_Descend(b *testing.B) {
	var tr = newNatiral()
	for i := 0; i < b.N; i++ {
		tr.Ins(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Descend(0, 0, func(k, v interface{}) bool {
			globalInterface, globalInterface = k, v
			return true
		})
	}
	b.ReportAllocs()
}
