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
	"github.com/logrusorgru/aurora"
)

type color bool

const (
	red   color = true
	black color = false
)

type LessFunc func(a, b interface{}) bool

type EqualFunc func(a, b interface{}) bool

type ZeroFunc func(a interface{}) bool

type node struct {
	l, r *node
	c    color
	k    interface{}
	v    interface{}
}

func newNode(k, v interface{}) (n *node) {
	n = new(node)
	n.c = red
	n.k = k
	n.v = v
	return
}

func (n *node) color() color {
	if n == nil {
		return black
	}
	return n.c
}

func (n *node) isBlack() bool {
	return n.color() == black
}

func (n *node) isRed() bool {
	return n.color() == red
}

func (n *node) setBlack() {
	if n != nil {
		n.c = black
	}
	return
}

func (n *node) setRed() {
	if n != nil {
		n.c = red
	}
	return
}

// n becomes red, its children becomes black
func (n *node) pushBlack() {
	n.setRed()
	n.l.setBlack()
	n.r.setBlack()
}

// left -> right, right, right,...
func (n *node) successor() (ss []*node, r *node) {
	ss = append(ss, n)
	if n.l != nil {
		for r = n.l; r.r != nil; r = r.r {
			ss = append(ss, r)
		}
	} else if n.r != nil {
		for r = n.r; r.l != nil; r = r.l {
			ss = append(ss, r)
		}
	}
	return
}

func (n *node) opposite(c *node) *node {
	if n.l == c {
		return n.r
	}
	return n.l
}

func (n *node) replaceChild(old, new *node) {
	if n.l == old {
		n.l = new
	} else {
		n.r = new
	}
}

// node points to at least one black
func (n *node) hasRedChild() bool {
	return n != nil && (n.l.isRed() || n.r.isRed())
}

func (n *node) copy(x *node) {
	n.k, n.v = x.k, x.v
}

type Tree struct {
	r *node

	less  LessFunc
	equal EqualFunc
	zero  ZeroFunc

	size int
}

func New(less LessFunc, equal EqualFunc, zero ZeroFunc) (t *Tree) {
	t = new(Tree)
	t.less = less
	t.equal = equal
	t.zero = zero
	return
}

var verbose bool

func print(args ...interface{}) {
	if verbose {
		fmt.Println(args...)
	}
}

// findInsertNode finds node to insert to
func (t *Tree) findInsertNode(st []*node, k interface{}) ([]*node, *node) {
	var (
		less = t.less
		p    *node
	)
	st, p = pop(st)
	for p != nil { // p - place
		st = append(st, p)
		if less(k, p.k) {
			p = p.l // left side
		} else {
			p = p.r // right side
		}
	}
	// TODO (kostyarin): avoid the pop
	st, p = pop(st)
	return st, p
}

// findNode and its dad
func (t *Tree) findNode(k interface{}) (st []*node, n *node) {
	var (
		less  = t.less
		equal = t.equal
	)
	for n = t.r; n != nil; {
		switch {
		case equal(k, n.k):
			return
		case less(k, n.k):
			st = append(st, n)
			n = n.l
		default:
			st = append(st, n)
			n = n.r
		}
	}
	return
}

func (t *Tree) isRoot(n *node) bool {
	return t.r == n
}

func (t *Tree) rightRotate(st []*node, n *node) {
	println("rightRotate", len(st), n.k.(int))
	var (
		pivot = n.l
		d     *node
	)
	st, d = pop(st)
	if d == nil {
		println("rightRotate d is nil")
		t.r = pivot
		pivot.c = black
	} else {
		println("rightRotate d:", d.k.(int))
		if d.l == n {
			d.l = pivot
		} else {
			d.r = pivot
		}
	}
	n.l = pivot.r
	pivot.r = n
}

func (t *Tree) leftRotate(st []*node, n *node) {
	println("leftRotate", len(st), n.k.(int))
	var (
		pivot = n.r
		d     *node
	)
	st, d = pop(st)
	if d == nil {
		println("leftRotate d is nil")
		t.r = pivot
		pivot.c = black
	} else {
		println("leftRotate d:", d.k.(int))
		if d.l == n {
			d.l = pivot
		} else {
			d.r = pivot
		}
	}
	n.r = pivot.l
	pivot.l = n
}

func (t *Tree) insertLeftLeftBalancing(st []*node, g, d *node) {
	d.c, g.c = g.c, d.c // swap colors
	t.rightRotate(st, g)
}

func (t *Tree) insertLeftRightBalancing(st []*node, g, d, n *node) {
	t.leftRotate(append(st, g), d)
	// the n becomes d after the leftRotate(d)
	t.insertLeftLeftBalancing(st, g, n)
}

func (t *Tree) insertRightRightBalancing(st []*node, g, d *node) {
	d.c, g.c = g.c, d.c // swap colors
	t.leftRotate(st, g)
}

func (t *Tree) insertRightLeftBalancing(st []*node, g, d, n *node) {
	t.rightRotate(append(st, g), d)
	// the n becomes d after the rightRotate(d)
	t.insertRightRightBalancing(st, g, n)
}

// balance tree after insert, the d is red
func (t *Tree) insertBalancing(st []*node, d, n *node) {
	var g, u *node
	for !t.isRoot(n) {
		if !d.isRed() {
			return
		}
		st, g = pop(st)
		if u = g.opposite(d); u.isRed() {
			g.pushBlack()
			n = g
			st, d = pop(st) // d = g.dad()
			continue
		} else { // us is black (or nil)
			// not the loop
			if g.l == d {
				if d.l == n {
					t.insertLeftLeftBalancing(st, g, d)
				} else { // n is right
					t.insertLeftRightBalancing(st, g, d, n)
				}
			} else { // d is right
				if d.r == n {
					t.insertRightRightBalancing(st, g, d)
				} else { // n is left
					t.insertRightLeftBalancing(st, g, d, n)
				}
			}
			return // done
		}
	}
	n.setBlack() // root must be black
}

// insert node to the tree and add pointer to it
// to the d
func (t *Tree) insertNode(st []*node, d, n *node) {
	t.size++
	if d == nil {
		t.r = n     // first element of the tree
		n.c = black // root must be black
		return      // done
	}
	// n already points to the d; required branch
	// (left or right) is nil and its guarantee by
	// findInsertNode
	if t.less(n.k, d.k) {
		d.l = n // left (less)
	} else {
		d.r = n // right (greater or equal)
	}
	t.insertBalancing(st, d, n)
}

// Ins is insert or overwrite, returning
//
//     1. previous value, false
//     2. nil, true
//
// The first case where an existing value overwritten. The
// second case where created new item.
func (t *Tree) Ins(k, v interface{}) (p interface{}, ok bool) {
	var (
		st, n = t.findNode(k)
		d     *node
	)
	if n != nil {
		p, n.v = n.v, v
		return // p, false
	}
	// n is nil
	st, d = t.findInsertNode(st, k)
	t.insertNode(st, d, newNode(k, v))
	return nil, true
}

// InsNx is insert if does not exist, returning
//
//     1. existing value, false
//     2. nil, true
//
// The first case if item already exists. The second case
// if item created.
func (t *Tree) InsNx(k, v interface{}) (e interface{}, ok bool) {
	var st, n = t.findNode(k)
	if n != nil {
		return n.v, false // already exists
	}
	// n is nil
	var d *node
	st, d = t.findInsertNode(st, k)
	t.insertNode(st, d, newNode(k, v))
	return nil, true
}

// InsEx is insert if exists, returning
//
//     1. previous value, true
//     2. nil, false
//
// The first case if item already exists and has been overwritten.
// The second case if item doesn't exist.
func (t *Tree) InsEx(k, v interface{}) (p interface{}, ok bool) {
	var _, n = t.findNode(k)
	if n == nil {
		return nil, false // does not exist
	}
	p, n.v, ok = n.v, v, true
	return
}

// Add is add new node even if it already exists. The Add called
// with the same key many times makes the Tree not unique. The
// Add returns true if item with given key is first in the Tree,
// i.e. if the Tree is still unique.
func (t *Tree) Add(k, v interface{}) (ok bool) {
	var (
		st, n = t.findNode(k)
		d     *node
	)
	if n != nil {
		// found, the Tree is or becomes not unique
		st, d = t.findInsertNode(append(st, n), k)
	} else {
		ok = true
		st, d = t.findInsertNode(st, k) // not found
	}
	t.insertNode(st, d, newNode(k, v))
	return
}

func (t *Tree) fixDoubleBlack(st []*node, x *node) {
	println("fixDoubleBlack", len(st), x.k.(int))
	var s, d *node
	for {
		println("fixDoubleBlack: loop", len(st), x.k.(int))
		if t.isRoot(x) {
			return
		}
		st, d = pop(st)
		s = d.opposite(x)
		if s == nil {
			println("fixDoubleBlack: s is nil")
			x = d
			continue // no recursion
		}
		println("fixDoubleBlack: s is ", s.k.(int), s.c)
		if s.isRed() {
			println("fixDoubleBlack: s is red")
			d.c = red
			s.c = black
			if d.r == s {
				println("right is:", d.r.k.(int))
				t.leftRotate(st, d)
				st = append(st, s) // stack has been changed
			} else {
				println("left is:", d.l.k.(int))
				t.rightRotate(st, d)
				st = append(st, s) // stack has been changed
			}
			st = append(st, d) // for the pop
			println("continue (1)")
			continue // no recursion
		}
		// the s is black
		if s.hasRedChild() {
			println("fixDoubleBlack: s has red child")
			if s.r.isRed() {
				println("fixDoubleBlack: right")
				if d.l == s {
					println("fixDoubleBlack: right->left")
					s.r.c = d.c
					t.leftRotate(append(st, d), s)
					t.rightRotate(st, d)
				} else {
					println("fixDoubleBlack: right->right")
					s.r.c = s.c
					s.c = d.c
					t.leftRotate(st, d)
				}
			} else { // left is red
				println("fixDoubleBlack: left")
				if d.l == s {
					println("fixDoubleBlack: left->left")
					s.l.c = s.c
					s.c = d.c
					for i, x := range st {
						println("stack", i, x.k.(int))
					}
					t.rightRotate(st, d)
				} else {
					println("fixDoubleBlack: left->right")
					s.l.c = d.c
					t.rightRotate(append(st, d), s)
					t.leftRotate(st, d)
				}
			}
			d.c = black
			return
		}
		println("fixDoubleBlack: s is black and has no red child", s.k.(int))
		s.c = red
		println("set red", s.k.(int))
		if d.c == black {
			println("fixDoubleBlack: s is black, has no red child, d is black")
			x = d
			continue
		}
		d.c = black
		println("set black", d.k.(int))
		return
	}
}

// delete and balance the Tree
func (t *Tree) delBalancing(st []*node, v *node) {
	println("delBalancing", len(st), v.k.(int))
	var (
		d, u *node
		ss   []*node
	)
	for {
		println("delBalancing: loop", len(st), v.k.(int))
		ss, u = v.successor()
		_, d = pop(st) // don't overwrite the st
		if u == nil {
			if t.isRoot(v) {
				t.r = nil
				return
			}
			if v.isBlack() {
				println("node to delete is black")
				t.fixDoubleBlack(st, v)
			} else {
				if s := d.opposite(v); s != nil {
					s.c = red
				}
			}
			d.replaceChild(v, nil)
			return
		}
		if v.l == nil || v.r == nil {
			println("has nil child")
			if t.isRoot(v) {
				v.copy(u)
				v.l, v.r = nil, nil
				return
			}
			d.replaceChild(v, u)
			if u.isBlack() && v.isBlack() {
				println("both are black (2)")
				t.fixDoubleBlack(st, u)
				return
			}
			u.c = black
			return
		}
		println("tail")
		v.copy(u)
		v = u                  // no recursion
		st = append(st, ss...) // change st
	}
}

// Get value by key. It returns (nil, false) if the
// Tree doesn't contain element with given key. If
// the Tree is not unique, the Get return first
// element. Use the Ascend or the Descend to get all
// non-unique elements.
func (t *Tree) Get(k interface{}) (v interface{}, ok bool) {
	var _, n = t.findNode(k)
	if n != nil {
		return n.v, true // got it
	}
	return nil, false // not found
}

func (t *Tree) Del(k interface{}) (v interface{}, ok bool) {
	println("Del", k.(int))
	var st, n = t.findNode(k)
	if n == nil {
		println("Del: not found", k.(int))
		return nil, false // does not exist
	}
	println("Del: found", k.(int))
	v, ok = n.v, true
	t.size--              //reduce
	t.delBalancing(st, n) // delete & balance
	return
}

func (t *Tree) minNode() (st []*node, n *node) {
	if t.r == nil {
		return
	}
	for n = t.r; n.l != nil; n = n.l {
		st = append(st, n)
	}
	return
}

func (t *Tree) maxNode() (st []*node, n *node) {
	if t.r == nil {
		return
	}
	for n = t.r; n.r != nil; n = n.r {
		st = append(st, n)
	}
	return
}

func (t *Tree) Min() (k, v interface{}, ok bool) {
	if _, n := t.minNode(); n != nil {
		k, v, ok = n.k, n.v, true
	}
	return
}

func (t *Tree) Max() (k, v interface{}, ok bool) {
	if _, n := t.maxNode(); n != nil {
		k, v, ok = n.k, n.v, true
	}
	return
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Clear() {
	t.size, t.r = 0, nil
}

// A WalkFunc is iterator. If it
// returns false iteration stops.
type WalkFunc func(k, v interface{}) (next bool)

func pop(ns []*node) (xs []*node, n *node) {
	if len(ns) == 0 {
		return ns, nil
	}
	n, xs = ns[len(ns)-1], ns[:len(ns)-1]
	return
}

func walk(n *node, walkFunc WalkFunc) bool {
	if n == nil {
		return true
	}
	return walkFunc(n.k, n.v) && walk(n.l, walkFunc) && walk(n.r, walkFunc)
}

// Walk elements of the Tree without any order.
func (t *Tree) Walk(walkFunc WalkFunc) {
	walk(t.r, walkFunc) // recursive
}

func (t *Tree) findAscendNode(k interface{}) (st []*node, n *node) {
	var (
		less  = t.less
		equal = t.equal
	)
	for n = t.r; n != nil; {
		switch {
		case equal(k, n.k):
			return
		case less(k, n.k):
			st = append(st, n)
			n = n.l
		default:
			n = n.r
		}
	}
	return
}

// [from, +inf)
func (t *Tree) ascendFrom(from interface{}, ascendFunc WalkFunc) {
	var st, n = t.findAscendNode(from)
	if n == nil {
		if st, n = t.minNode(); n != nil && t.less(n.k, from) {
			return
		}
	}
	var d *node
	for n != nil {
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				st = append(st, n)
				n = n.l
			}
		} else {
			if st, d = pop(st); d == nil {
				return
			} else if d.r == n {
				st, n = pop(st)
			} else {
				n = d
			}
		}
	}
}

// (-inf, to]
func (t *Tree) ascendTo(to interface{}, ascendFunc WalkFunc) {
	var (
		st, n = t.minNode()
		less  = t.less
	)
	for n != nil {
		if less(to, n.k) {
			return // that's all
		}
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				st = append(st, n)
				n = n.l
			}
		} else {
			st, n = pop(st)
		}
	}
}

// [from, to]
func (t *Tree) ascendFromTo(from, to interface{}, ascendFunc WalkFunc) {
	var (
		st, n = t.findAscendNode(from)
		less  = t.less
	)
	if n == nil {
		if st, n = t.minNode(); n != nil && t.less(n.k, from) {
			return
		}
	}
	for n != nil {
		if less(to, n.k) {
			return // that's all
		}
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				st = append(st, n)
				n = n.l
			}
		} else {
			st, n = pop(st)
		}
	}
}

// (-inf, +inf)
func (t *Tree) ascend(ascendFunc WalkFunc) {
	for st, n := t.minNode(); n != nil; {
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				st = append(st, n)
				n = n.l
			}
		} else {
			st, n = pop(st)
		}
	}
}

// Ascend iterates elements of the tree ascending order. The ZeroFunc
// used to determine ascending range.
func (t *Tree) Ascend(from, to interface{}, ascendFunc WalkFunc) {
	// zero function
	switch zero := t.zero; {
	case zero(from): // (-inf, to] or (-inf, +inf)
		if zero(to) {
			t.ascend(ascendFunc) // (-inf, +inf)
		} else {
			t.ascendTo(to, ascendFunc) // (-inf, to]
		}
	case zero(to): // [from, +inf)
		t.ascendFrom(from, ascendFunc)
	default: // [from, to]
		t.ascendFromTo(from, to, ascendFunc)
	}
}

func (t *Tree) findDescendNode(k interface{}) (st []*node, n *node) {
	var (
		less  = t.less
		equal = t.equal
	)
	for n = t.r; n != nil; {
		switch {
		case equal(k, n.k):
			return
		case less(k, n.k):
			n = n.l
		default:
			st = append(st, n)
			n = n.r
		}
	}
	return
}

// [from, -inf) (reversed)
func (t *Tree) descendFrom(from interface{}, descendFunc WalkFunc) {
	var st, n = t.findDescendNode(from)
	if n == nil {
		if st, n = t.maxNode(); n != nil && t.less(from, n.k) {
			return
		}
	}
	for n != nil {
		if !descendFunc(n.k, n.v) {
			return
		}
		if n.l != nil {
			n = n.l
			for n.r != nil {
				st = append(st, n)
				n = n.r
			}
		} else {
			st, n = pop(st)
		}
	}
}

// (+inf, to] (reversed)
func (t *Tree) descendTo(to interface{}, descendFunc WalkFunc) {
	var (
		st, n = t.maxNode()
		less  = t.less
	)
	for n != nil {
		if less(n.k, to) {
			return // that's all
		}
		if !descendFunc(n.k, n.v) {
			return
		}
		if n.l != nil {
			n = n.l
			for n.r != nil {
				st = append(st, n)
				n = n.r
			}
		} else {
			st, n = pop(st)
		}
	}
}

// [from, to] (reversed)
func (t *Tree) descendFromTo(from, to interface{}, descendFunc WalkFunc) {
	var (
		st, n = t.findDescendNode(from)
		less  = t.less
	)
	if n == nil {
		if st, n = t.maxNode(); n != nil && t.less(from, n.k) {
			return
		}
	}
	for n != nil {
		if less(n.k, to) {
			return // that's all
		}
		if !descendFunc(n.k, n.v) {
			return
		}
		if n.l != nil {
			n = n.l
			for n.r != nil {
				st = append(st, n)
				n = n.r
			}
		} else {
			st, n = pop(st)
		}
	}
}

// (-inf, +inf) (reversed)
func (t *Tree) descend(descendFunc WalkFunc) {
	for st, n := t.maxNode(); n != nil; {
		if !descendFunc(n.k, n.v) {
			return
		}
		if n.l != nil {
			n = n.l
			for n.r != nil {
				st = append(st, n)
				n = n.r
			}
		} else {
			st, n = pop(st)
		}
	}
}

// Descend iterates elements of the tree descending order.
func (t *Tree) Descend(from, to interface{}, descendFunc WalkFunc) {
	// zero function
	switch zero := t.zero; {
	case zero(from): // (-inf, to] or (-inf, +inf)
		if zero(to) {
			t.descend(descendFunc) // (-inf, +inf)
		} else {
			t.descendTo(to, descendFunc) // (-inf, to]
		}
	case zero(to): // [from, +inf)
		t.descendFrom(from, descendFunc)
	default: // [from, to]
		t.descendFromTo(from, to, descendFunc)
	}
}

type Printer interface {
	Add(string) Printer
}

var already map[*node]struct{}

func (n *node) print(d *node, pr Printer) {
	if n == nil {
		return
	}
	if _, ok := already[n]; ok {
		println("already", n.k.(int), d.k.(int))
		panic("")
		return
	}
	//println("print node", n.k.(int))
	var ds string
	if d == nil {
		ds = ". "
	} else if d.l == n {
		ds = "l "
	} else {
		ds = "r "
	}
	s := fmt.Sprint(ds, n.k)
	var sub Printer
	if n.isRed() {
		sub = pr.Add(aurora.Red(s).String())
	} else {
		sub = pr.Add(aurora.Blue(s).Bold().String())
	}
	n.l.print(n, sub)
	n.r.print(n, sub)

	already[n] = struct{}{}
}

func (t *Tree) Print(pr Printer) {
	already = make(map[*node]struct{})

	s := fmt.Sprintf("[%d]", t.size)
	tree := pr.Add(s)
	t.r.print(nil, tree)
}
