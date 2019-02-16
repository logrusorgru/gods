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

package main

import (
	"fmt"
	"github.com/disiqueira/gotree"
	"github.com/kr/pretty"
	"github.com/logrusorgru/gods/rb/rb"
)

func newNatural() *rb.Tree {
	return rb.New(
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

type Tree struct {
	gotree.Tree
}

func (t *Tree) Add(text string) rb.Printer {
	return t.Tree.Add(text)
}

func main() {
	var tree = newNatural()

	for i := 0; i <= 20; i++ {
		tree.Ins(i, 0)
	}

	var pt = Tree{gotree.New("rb-tree")}
	tree.Print(pt)

	fmt.Println(pt.Print())
	// pretty.Println(tree)

	for _, i := range []int{
		1,
		2,
		4,
		5,
		7,
		8,
		10,
		11,
		// 13,
		// 14,
		// 16,
		// 17,
		// 19,
		// 20,
	} {
		fmt.Println("DELETE", i)
		tree.Del(i)
	}

	for i := 0; i <= 20; i++ {
		if i%3 == 0 {
			continue
		}
		_ = i
	}

	{

		var pt = Tree{gotree.New("rb-tree")}
		tree.Print(pt)

		fmt.Println(pt.Print())
		// pretty.Println(tree)
	}

	_ = pretty.Print

}
