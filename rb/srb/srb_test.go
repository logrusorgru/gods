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
	"fmt"
	"github.com/disiqueira/gotree"
	"math/rand"
	"testing"
)

const (
	keyBelow = -100          //
	keyMin   = 0             //
	keyMax   = 100           //
	keyAbove = keyMax + 1000 //
)

func newNatiral() *Tree {
	return New(
		func(a, b interface{}) bool {
			return a.(int) < b.(int)
		},
		func(a, b interface{}) bool {
			return a.(int) == b.(int)
		},
		func(a interface{}) bool {
			return a.(int) == 0
		})
}

func TestNew(t *testing.T) {
	tr := newNatiral()
	if tr == nil {
		t.Fatal("new returns nil")
	}
	if tr.Size() != 0 {
		t.Error("size is not zero")
	}
}

func rs(r []int) string {
	return fmt.Sprintf("[%d, ..., %d] %d", r[0], r[len(r)-1], len(r))
}

type Rng struct {
	f, t int
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (r Rng) size() int {
	return abs(r.t-r.f) + 1
}

func (r Rng) cond(i int) bool {
	if r.f < r.t {
		return i <= r.t
	}
	return i >= r.t
}

func (r Rng) iter(i int) int {
	if r.f < r.t {
		return i + 1
	}
	return i - 1
}

func (r Rng) String() string {
	if r.f < r.t {
		return fmt.Sprintf("[%d, %d)", r.f, r.t)
	}
	return fmt.Sprintf("[%d, %d) (reversed)", r.f, r.t)
}

func Range(f, t int, random bool) (vs []int) {
	var r = Rng{f, t}
	vs = make([]int, 0, r.size())
	for i := r.f; r.cond(i); i = r.iter(i) {
		vs = append(vs, i)
	}
	if random == true {
		rand.Shuffle(len(vs), func(i, j int) {
			vs[i], vs[j] = vs[j], vs[i]
		})
	}
	return
}

var Ranges = [][]int{
	Range(keyMin, keyMax, false),
	Range(keyMax, keyMin, false),
	Range(keyMin, keyMax, true),
	Range(keyMax, keyMin, true),
}

func TestTree_Ins(t *testing.T) {
	// Ins(k, v interface{}) (p interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		verbose = false

		for _, i := range r {
			p, ok := tr.Ins(i, i)
			if ok == false {
				t.Error("ok is false")
			}
			if p != nil {
				t.Error("p is not nil")
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		for _, i := range r {
			p, ok := tr.Ins(i, 0)
			if ok == true {
				t.Error("ok is true", i)
			}
			if j, ok := p.(int); !ok {
				t.Error("p is not int", p, i)
			} else if j != i {
				t.Error("p is not", i, j)
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		if t.Failed() {
			return
		}

	}

}

func TestTree_InsNx(t *testing.T) {
	// InsNx(k, v interface{}) (e interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			p, ok := tr.InsNx(i, i)
			if ok == false {
				t.Error("ok is false")
			}
			if p != nil {
				t.Error("p is not nil")
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		for _, i := range r {
			e, ok := tr.InsNx(i, 0)
			if ok == true {
				t.Error("ok is true")
			}
			if j, ok := e.(int); !ok {
				t.Error("e is not int")
			} else if j != i {
				t.Error("e is not", i, j)
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_InsEx(t *testing.T) {
	// InsEx(k, v interface{}) (p interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			p, ok := tr.Ins(i, i)
			if ok == false {
				t.Error("ok is false")
			}
			if p != nil {
				t.Error("p is not nil")
			}
			if t.Failed() {
				return
			}
		}

		for _, i := range r {
			p, ok := tr.InsEx(i, i)
			if ok == false {
				t.Error("ok is false")
			}
			if j, ok := p.(int); !ok {
				t.Error("p is not int")
			} else if j != i {
				t.Error("p is not", i, j)
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		tr.Clear()

		for _, i := range r {
			p, ok := tr.InsEx(i, i)
			if ok == true {
				t.Error("ok is true")
			}
			if p != nil {
				t.Error("p is not nil")
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != 0 {
			t.Error("wrong size", tr.Size(), "want", 0)
		}
		if t.Failed() {
			return
		}

	}

}

func TestTree_Add(t *testing.T) {
	// Add(k, v interface{}) (ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			if tr.Add(i, i) == false {
				t.Error("Add returns false")
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		for _, i := range r {
			if tr.Add(i, i) == true {
				t.Error("Add returns true", i, tr.Size(), rs(r))
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r)*2 {
			t.Error("wrong size", tr.Size(), "want", len(r)*2)
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Get(t *testing.T) {
	// Get(k interface{}) (v interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			v, ok := tr.Get(i)
			if ok == true {
				t.Error("ok is true")
			}
			if v != nil {
				t.Error("v is not nil")
			}
			if t.Failed() {
				return
			}
		}

		for _, i := range r {
			tr.Ins(i, i)
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		for _, i := range r {
			v, ok := tr.Get(i)
			if ok == false {
				t.Error("ok is false")
			}
			if j, ok := v.(int); !ok {
				t.Error("v is not int")
			} else if j != i {
				t.Error("j is not i", j, i)
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}
		if t.Failed() {
			return
		}
	}

}

type Print struct {
	gotree.Tree
}

func (p *Print) Add(name string) Printer {
	return &Print{p.Tree.Add(name)}
}

func TestTree_Del(t *testing.T) {
	// Del(k interface{}) (v interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			tr.Ins(i, i)
		}

		if tr.Size() != len(r) {
			t.Error("wrong size", tr.Size(), "want", len(r))
		}

		for _, i := range r {
			t.Log(i)
			v, ok := tr.Del(i)
			if ok == false {
				t.Error("ok is false", i, rs(r))
			}
			if j, ok := v.(int); !ok {
				t.Error("v is not int")
			} else if j != i {
				t.Error("j is not i", j, i)
			}
			if t.Failed() {
				return
			}
			////////////////////////////////////////////////////////////////////
			var tree = &Print{gotree.New("rb")}
			tr.Print(tree)
			t.Log(tree.Print())
			t.Fatal("fatality")
			////////////////////////////////////////////////////////////////////
		}

		if tr.Size() != 0 {
			t.Error("wrong size", tr.Size(), "want", 0)
		}

		for _, i := range r {
			v, ok := tr.Del(i)
			if ok == true {
				t.Error("ok is true")
			}
			if v != nil {
				t.Error("v is not nil")
			}
			if t.Failed() {
				return
			}
		}

		if tr.Size() != 0 {
			t.Error("wrong size", tr.Size(), "want", 0)
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Min(t *testing.T) {
	// Min() (k, v interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		k, v, ok := tr.Min()
		if ok == true {
			t.Error("ok is true")
		}
		if v != nil {
			t.Error("v is not nil")
		}
		if k != nil {
			t.Error("k is not nil")
		}

		var min = keyAbove

		for _, i := range r {
			tr.Ins(i, i)
			if i < min {
				min = i
			}
			k, v, ok := tr.Min()
			if ok == false {
				t.Error("ok is false")
			}
			if v != k {
				t.Error("v is not k")
			}
			if m, ok := k.(int); !ok {
				t.Error("k is not int")
			} else if m != min {
				t.Error("m is not min", m, min)
			}
			if t.Failed() {
				return
			}
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Max(t *testing.T) {
	// Max() (k, v interface{}, ok bool)

	for _, r := range Ranges {
		tr := newNatiral()

		k, v, ok := tr.Max()
		if ok == true {
			t.Error("ok is true")
		}
		if v != nil {
			t.Error("v is not nil")
		}
		if k != nil {
			t.Error("k is not nil")
		}

		var max int = keyBelow
		for _, i := range r {
			tr.Ins(i, i)
			if i > max {
				max = i
			}
			k, v, ok := tr.Max()
			if ok == false {
				t.Error("ok is false")
			}
			if v != k {
				t.Error("v is not k")
			}
			if m, ok := k.(int); !ok {
				t.Error("k is not int")
			} else if m != max {
				t.Error("m is not 100", m, max)
			}
			if t.Failed() {
				return
			}
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Size(t *testing.T) {
	// Size() int

	for _, r := range Ranges {
		tr := newNatiral()

		if tr.Size() != 0 {
			t.Error("wrong size", tr.Size(), "want", 0)
		}

		for j, i := range r {
			tr.Ins(i, i)
			if tr.Size() != j+1 {
				t.Error("wrong size", tr.Size(), "want", j+1)
			}
			if t.Failed() {
				return
			}
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Clear(t *testing.T) {
	// Clear()

	for _, r := range Ranges {
		tr := newNatiral()

		for _, i := range r {
			tr.Ins(i, i)
		}

		tr.Clear()

		if tr.Size() != 0 {
			t.Error("wrong size", tr.Size(), "want", 0)
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Walk(t *testing.T) {
	// Walk(walkFunc WalkFunc)

	for _, r := range Ranges {
		tr := newNatiral()

		var called int
		tr.Walk(func(k, v interface{}) bool {
			called++
			return true
		})
		if called != 0 {
			t.Error("called", called)
		}

		for _, i := range r {
			tr.Ins(i, i)
		}

		called = 0
		var mp = make(map[interface{}]interface{})
		tr.Walk(func(k, v interface{}) bool {
			called++
			if v, ok := mp[k]; ok {
				t.Fatal("already", k, v)
			}
			mp[k] = v
			return true
		})

		if len(mp) != tr.Size() {
			t.Error("wrong size walked")
		}

		for _, i := range r {
			if v, ok := mp[i].(int); !ok {
				t.Fatal("wrong or missing value", i)
			} else if v != i {
				t.Fatal("wrong value", i)
			}
		}

		called = 0
		tr.Walk(func(k, v interface{}) bool {
			called++
			return false
		})
		if called != 1 {
			t.Error("wrong called", called)
		}
		if t.Failed() {
			return
		}
	}

}

func TestTree_Ascend(t *testing.T) {
	// Ascend(from, to interface{}, ascendFunc WalkFunc)

	t.Run("full", func(t *testing.T) {
		for _, r := range Ranges {
			tr := newNatiral()
			var called int
			tr.Ascend(0, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = 0
			tr.Ascend(0, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, tr.Size())
			}
			called = 0
			tr.Ascend(0, 0, func(k, v interface{}) bool {
				called++
				return false
			})
			if called != 1 {
				t.Error("wrong called", called)
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("from", func(t *testing.T) {
		for _, r := range Ranges {
			const from = 50
			tr := newNatiral()
			var called int
			tr.Ascend(from, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = from
			tr.Ascend(from, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != len(r) {
				t.Error("wrong called", called, len(r), rs(r))
			}
			// before
			called = 0
			tr.Ascend(keyBelow, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			// after
			called = 0
			tr.Ascend(keyAbove, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("to", func(t *testing.T) {
		for _, r := range Ranges {
			const to = 50
			tr := newNatiral()
			var called int
			tr.Ascend(0, to, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			tr.Ascend(0, to, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != to+1 {
				t.Error("wrong called", called, to+1, rs(r))
			}
			// before
			called = 0
			tr.Ascend(0, keyBelow, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			// after
			called = 0
			tr.Ascend(0, keyAbove, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("from to", func(t *testing.T) {
		for _, r := range Ranges {
			const from, to = 45, 55
			tr := newNatiral()
			var called int
			tr.Ascend(from, to, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = from
			tr.Ascend(from, to, func(k, v interface{}) bool {
				t.Log(k)
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called-to != 1 {
				t.Error("wrong called", called, rs(r))
			}
			// before & after
			called = 0
			tr.Ascend(keyBelow, keyAbove, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != called {
					t.Fatal("wrong j", j, called, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			if t.Failed() {
				return
			}
		}
	})

}

func TestTree_Descend(t *testing.T) {
	// Descend(from, to interface{}, descendFunc WalkFunc)

	t.Run("full", func(t *testing.T) {
		for _, r := range Ranges {
			tr := newNatiral()
			var called int
			tr.Descend(0, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = 0
			tr.Descend(0, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != len(r)-called-1 {
					t.Fatal("wrong j", j, len(r)-called-1, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called)
			}
			called = 0
			tr.Descend(0, 0, func(k, v interface{}) bool {
				called++
				return false
			})
			if called != 1 {
				t.Error("wrong called", called)
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("from", func(t *testing.T) {
		for _, r := range Ranges {
			const from = 50
			tr := newNatiral()
			var called int
			tr.Descend(from, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = 0
			tr.Descend(from, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != from-called {
					t.Fatal("wrong j", j, from-called, rs(r))
				}
				called++
				return true
			})
			if called != from+1 {
				t.Error("wrong called", called, from+1)
			}
			// before
			called = 0
			tr.Descend(keyAbove, 0, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != len(r)-called-1 {
					t.Fatal("wrong j", j, len(r)-called-1, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, tr.Size(), rs(r))
			}
			// after
			called = 0
			tr.Descend(keyBelow, 0, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("to", func(t *testing.T) {
		for _, r := range Ranges {
			const to = 50
			tr := newNatiral()
			var called int
			tr.Descend(0, to, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			tr.Descend(0, to, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != len(r)-called-1 {
					t.Fatal("wrong j", j, len(r)-called-1, rs(r))
				}
				called++
				return true
			})
			if called+to != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			// before
			called = 0
			tr.Descend(0, keyAbove, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			// after
			called = 0
			tr.Descend(0, keyBelow, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != len(r)-called-1 {
					t.Fatal("wrong j", j, len(r)-called-1, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			if t.Failed() {
				return
			}
		}
	})

	t.Run("from to", func(t *testing.T) {
		for _, r := range Ranges {
			const from, to = 55, 45
			tr := newNatiral()
			var called int
			tr.Descend(from, to, func(k, v interface{}) bool {
				called++
				return true
			})
			if called != 0 {
				t.Error("wrong called", called)
			}
			for _, i := range r {
				tr.Ins(i, i)
			}
			called = 0
			tr.Descend(from, to, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != from-called {
					t.Fatal("wrong j", j, from-called, rs(r))
				}
				called++
				return true
			})
			if called != from-to+1 {
				t.Error("wrong called", called, from-to+1, rs(r))
			}
			// before & after
			called = 0
			tr.Descend(keyAbove, keyBelow, func(k, v interface{}) bool {
				if k != v {
					t.Fatal("k is not v")
				}
				if j, ok := k.(int); !ok {
					t.Fatal("k is not int")
				} else if j != len(r)-called-1 {
					t.Fatal("wrong j", j, len(r)-called-1, rs(r))
				}
				called++
				return true
			})
			if called != tr.Size() {
				t.Error("wrong called", called, rs(r))
			}
			if t.Failed() {
				return
			}
		}
	})

}
