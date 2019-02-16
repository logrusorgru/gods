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
	"testing"
)

func newNatural() *Tree {
	return New(
		func(a, b interface{}) bool {
			return a.(int) < b.(int)
		},
		func(a, b interface{}) bool {
			return a.(int) == b.(int)
		},
		func(a interface{}) bool {
			return a.(int) == 0
		},
	)
}

func testee() []int {
	return []int{
		100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		80, 81, 82, 83, 84, 85, 86, 86, 87, 88, 89, 90,
		200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210,
		50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	}
}

func TestNew(t *testing.T) {
	// New(less LessFunc, equal EqualFunc, zero ZeroFunc) (t *Tree)
	if tree := newNatural(); tree == nil {
		t.Error("New returns nil")
	}
}

func TestTree_Get(t *testing.T) {
	// Get(item interface{}) bool
	var tree = newNatural()
	for _, item := range testee() {
		if tree.Get(item) == true {
			t.Errorf("got item not exist %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Ins(item) == false {
			t.Errorf("new item inserted with false %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Get(item) == false {
			t.Errorf("missing item %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Del(item) == false {
			t.Errorf("delete item with false %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Get(item) == false {
			t.Errorf("got deleted item %d", item)
		}
	}
}

func TestTree_Ins(t *testing.T) {
	// Ins(item interface{}) (ok bool)
}

func TestTree_InsNx(t *testing.T) {
	// InsNx(item interface{}) (ok bool)
	var tree = newNatural()
	for _, item := range testee() {
		if tree.InsNx(item) == false {
			t.Errorf("new item inserted with false %d", item)
		}
	}
	for _, item := range testee() {
		if tree.InsNx(item) == true {
			t.Errorf("item overwritten %d", item)
		}
	}
}

func TestTree_InsEx(t *testing.T) {
	// InsEx(item interface{}) (ok bool)
	var tree = newNatural()
	for _, item := range testee() {
		if tree.InsEx(item) == true {
			t.Errorf("new item created, but should not %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Ins(item) == false {
			t.Errorf("new item inserted with false %d", item)
		}
	}
	for _, item := range testee() {
		if tree.InsEx(item) == false {
			t.Errorf("item not overwritten %d", item)
		}
	}
}

func TestTree_Del(t *testing.T) {
	// Del(item interface{}) (ok bool)
	var tree = newNatural()
	for _, item := range testee() {
		if tree.Del(item) == true {
			t.Errorf("new item created, but should not %d", item)
		}
	}
	for _, item := range testee() {
		if tree.Ins(item) == false {
			t.Errorf("new item inserted with false %d", item)
		}
	}
	for _, item := range testee() {
		if tree.InsEx(item) == false {
			t.Errorf("item not overwritten %d", item)
		}
	}
}

func TestTree_Size(t *testing.T) {
	// Size() int
}

func TestTree_Clear(t *testing.T) {
	// Clear()
}

func TestTree_Min(t *testing.T) {
	// Min() (interface{}, bool)
}

func TestTree_Max(t *testing.T) {
	// Max() (interface{}, bool)
}

func TestTree_Walk(t *testing.T) {
	// Walk(walkFunc WalkFunc)
}

func TestTree_Ascend(t *testing.T) {
	// Ascend(from, to interface{}, ascendFunc WalkFunc)
}

func TestTree_Descend(t *testing.T) {
	// Descend(from, to interface{}, descendFunc WalkFunc)
}
