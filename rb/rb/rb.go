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

package rb

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
	d, l, r *node
	c       color
	k       interface{}
	v       interface{}
}

func newNode(dad *node, k, v interface{}) (n *node) {
	n = new(node)
	n.d = dad
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

func (n *node) left() *node {
	if n == nil {
		return nil
	}
	return n.l
}

func (n *node) right() *node {
	if n == nil {
		return nil
	}
	return n.r
}

func (n *node) dad() *node {
	if n == nil {
		return nil
	}
	return n.d
}

func (n *node) sibling() *node {
	if left := n.dad().left(); left != n {
		return left
	}
	return n.dad().right()
}

func (n *node) uncle() *node {
	return n.dad().sibling()
}

func (n *node) isLeft() bool {
	return n.dad().left() == n
}

func (n *node) isRight() bool {
	return n.dad().right() == n
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
func (n *node) successor() (r *node) {
	if n.l != nil {
		for r = n.l; r.r != nil; r = r.r {
		}
	} else if n.r != nil {
		for r = n.r; r.l != nil; r = r.l {
		}
	}
	return
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
	if x == nil {
		n.k, n.v = nil, nil
		return
	}
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

// findInsertNode finds node to insert to
func (t *Tree) findInsertNode(d *node, k interface{}) *node {
	for p, less := d, t.less; p != nil; { // p - place
		if less(k, p.k) {
			p, d = p.l, p // left side
		} else {
			p, d = p.r, p // right side
		}
	}
	return d
}

// findNode and its dad
func (t *Tree) findNode(k interface{}) (d, n *node) {
	var (
		less  = t.less
		equal = t.equal
	)
	for n, d = t.r, nil; n != nil; {
		switch {
		case equal(k, n.k):
			return
		case less(k, n.k):
			n, d = n.l, n
		default:
			n, d = n.r, n
		}
	}
	return
}

func (t *Tree) isRoot(n *node) bool {
	return t.r == n
}

func (t *Tree) rightRotate(n *node) {
	var pivot = n.l
	if n.d == nil {
		t.r = pivot
		pivot.c = black
		pivot.d = nil
	} else {
		pivot.d = n.d
		if n.isLeft() {
			n.d.l = pivot
		} else {
			n.d.r = pivot
		}
	}
	n.l = pivot.r
	if pivot.r != nil {
		pivot.r.d = n
	}
	n.d = pivot
	pivot.r = n
}

func (t *Tree) leftRotate(n *node) {
	var pivot = n.r
	if n.d == nil {
		t.r = pivot
		pivot.c = black
		pivot.d = nil
	} else {
		pivot.d = n.d
		if n.isLeft() {
			n.d.l = pivot
		} else {
			n.d.r = pivot
		}
	}
	n.r = pivot.l
	if pivot.l != nil {
		pivot.l.d = n
	}
	n.d = pivot
	pivot.l = n
}

func (t *Tree) insertLeftLeftBalancing(g, d *node) {
	d.c, g.c = g.c, d.c // swap colors
	t.rightRotate(g)
}

func (t *Tree) insertLeftRightBalancing(g, d, n *node) {
	t.leftRotate(d)
	// the n becomes d after the leftRotate(d)
	t.insertLeftLeftBalancing(g, n)
}

func (t *Tree) insertRightRightBalancing(g, d *node) {
	d.c, g.c = g.c, d.c // swap colors
	t.leftRotate(g)
}

func (t *Tree) insertRightLeftBalancing(g, d, n *node) {
	t.rightRotate(d)
	// the n becomes d after the rightRotate(d)
	t.insertRightRightBalancing(g, n)
}

// balance tree after insert, the d is red
func (t *Tree) insertBalancing(d, n *node) {
	var g, u *node
	for !t.isRoot(n) {
		if !d.isRed() {
			return
		}
		g = d.dad()
		if u = n.uncle(); u.isRed() {
			g.pushBlack()
			d, n = g.dad(), g
			continue
		} else { // us is black (or nil)
			// not the loop
			if d.isLeft() {
				if n.isLeft() {
					t.insertLeftLeftBalancing(g, d)
				} else { // n is right
					t.insertLeftRightBalancing(g, d, n)
				}
			} else { // d is right
				if n.isRight() {
					t.insertRightRightBalancing(g, d)
				} else { // n is left
					t.insertRightLeftBalancing(g, d, n)
				}
			}
			return // done
		}
	}
	n.setBlack() // root must be black
}

// insert node to the tree and add pointer to it
// to the d
func (t *Tree) insertNode(d, n *node) {
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
	n.d = d
	t.insertBalancing(d, n)
}

// Ins is insert or overwrite, returning
//
//     1. previous value, false
//     2. nil, true
//
// The first case where an existing value overwritten. The
// second case where created new item.
func (t *Tree) Ins(k, v interface{}) (p interface{}, ok bool) {
	var d, n = t.findNode(k)
	if n != nil {
		p, n.v = n.v, v
		return // p, false
	}
	// n is nil
	d = t.findInsertNode(d, k)
	t.insertNode(d, newNode(d, k, v))
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
	var d, n = t.findNode(k)
	if n != nil {
		return n.v, false // already exists
	}
	// n is nil
	d = t.findInsertNode(d, k)
	t.insertNode(d, newNode(d, k, v))
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
	var d, n = t.findNode(k)
	if n != nil {
		d = t.findInsertNode(n, k) // found, the Tree is or becomes not unique
	} else {
		ok, d = true, t.findInsertNode(d, k) // not found
	}
	t.insertNode(d, newNode(d, k, v))
	return
}

func (t *Tree) fixDoubleBlack(x *node) {
	for {
		if t.isRoot(x) {
			return
		}
		var (
			s = x.sibling()
			d = x.d
		)
		if s == nil {
			x = d
			continue // no recursion
		}
		if s.isRed() {
			d.c = red
			s.c = black
			if s.isRight() {
				t.leftRotate(d)
			} else {
				t.rightRotate(d)
			}
			continue // no recursion
		}
		// the s is black
		if s.hasRedChild() {
			if s.r.isRed() {
				if s.isLeft() {
					s.r.c = d.c
					t.leftRotate(s)
					t.rightRotate(d)
				} else {
					s.r.c = s.c
					s.c = d.c
					t.leftRotate(d)
				}
			} else { // left is red
				if s.isLeft() {
					s.l.c = s.c
					s.c = d.c
					t.rightRotate(d)
				} else {
					s.l.c = d.c
					t.rightRotate(s)
					t.leftRotate(d)
				}
			}
			d.c = black
			return
		}
		s.c = red
		if d.c == black {
			x = d
			continue
		}
		d.c = black
		return
	}
}

// delete and balance the Tree
func (t *Tree) delBalancing(v *node) {
	for {
		var u = v.successor()
		if u == nil {
			if t.isRoot(v) {
				t.r = nil
				return
			}
			if u.isBlack() && v.isBlack() {
				t.fixDoubleBlack(v)
			} else {
				if s := v.sibling(); s != nil {
					s.c = red
				}
			}
			v.d.replaceChild(v, nil)
			return
		}
		if v.l == nil || v.r == nil {
			if t.isRoot(v) {
				v.copy(u)
				v.l, v.r = nil, nil
				return
			}
			v.d.replaceChild(v, u)
			u.d = v.d
			if u.isBlack() && v.isBlack() {
				t.fixDoubleBlack(u)
				return
			}
			u.c = black
			return
		}
		v.copy(u)
		v = u // no recursion
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
	var _, n = t.findNode(k)
	if n == nil {
		return nil, false // does not exist
	}
	v, ok = n.v, true
	t.size--          //reduce
	t.delBalancing(n) // delete & balance
	return
}

func (t *Tree) minNode() (n *node) {
	if t.r == nil {
		return
	}
	for n = t.r; n.l != nil; n = n.l {
	}
	return
}

func (t *Tree) maxNode() (n *node) {
	if t.r == nil {
		return
	}
	for n = t.r; n.r != nil; n = n.r {
	}
	return
}

func (t *Tree) Min() (k, v interface{}, ok bool) {
	if n := t.minNode(); n != nil {
		k, v, ok = n.k, n.v, true
	}
	return
}

func (t *Tree) Max() (k, v interface{}, ok bool) {
	if n := t.maxNode(); n != nil {
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

// Walk elements of the Tree without any order.
func (t *Tree) Walk(walkFunc WalkFunc) {
	var stack []*node
	for n := t.r; n != nil; {
		if !walkFunc(n.k, n.v) {
			return
		}
		// push right
		if n.r != nil {
			stack = append(stack, n.r)
		}
		if n.l != nil {
			n = n.l // walk left
			continue
		}
		// walk right
		stack, n = pop(stack)
	}
}

// [from, +inf)
func (t *Tree) ascendFrom(from interface{}, ascendFunc WalkFunc) {
	var n *node
	if _, n = t.findNode(from); n == nil {
		if n = t.minNode(); n != nil && t.less(n.k, from) {
			return
		}
	}
	for n != nil {
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				n = n.l
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.l == n {
				n = n.d
				break
			}
		}
	}
}

// (-inf, to]
func (t *Tree) ascendTo(to interface{}, ascendFunc WalkFunc) {
	var (
		n    = t.minNode()
		less = t.less
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
				n = n.l
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.l == n {
				n = n.d
				break
			}
		}
	}
}

// [from, to]
func (t *Tree) ascendFromTo(from, to interface{}, ascendFunc WalkFunc) {
	var (
		n    *node
		less = t.less
	)
	if _, n = t.findNode(from); n == nil {
		if n = t.minNode(); n != nil && t.less(n.k, from) {
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
				n = n.l
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.l == n {
				n = n.d
				break
			}
		}
	}
}

// (-inf, +inf)
func (t *Tree) ascend(ascendFunc WalkFunc) {
	for n := t.minNode(); n != nil; {
		if !ascendFunc(n.k, n.v) {
			return
		}
		if n.r != nil {
			n = n.r
			for n.l != nil {
				n = n.l
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.l == n {
				n = n.d
				break
			}
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

// [from, -inf) (reversed)
func (t *Tree) descendFrom(from interface{}, descendFunc WalkFunc) {
	var n *node
	if _, n = t.findNode(from); n == nil {
		if n = t.maxNode(); n != nil && t.less(from, n.k) {
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
				n = n.r
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.r == n {
				n = n.d
				break
			}
		}
	}
}

// (+inf, to] (reversed)
func (t *Tree) descendTo(to interface{}, descendFunc WalkFunc) {
	var (
		n    = t.maxNode()
		less = t.less
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
				n = n.r
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.r == n {
				n = n.d
				break
			}
		}
	}
}

// [from, to] (reversed)
func (t *Tree) descendFromTo(from, to interface{}, descendFunc WalkFunc) {
	var (
		n    *node
		less = t.less
	)
	if _, n = t.findNode(from); n == nil {
		if n = t.maxNode(); n != nil && t.less(from, n.k) {
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
				n = n.r
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.r == n {
				n = n.d
				break
			}
		}
	}
}

// (-inf, +inf) (reversed)
func (t *Tree) descend(descendFunc WalkFunc) {
	for n := t.maxNode(); n != nil; {
		if !descendFunc(n.k, n.v) {
			return
		}
		if n.l != nil {
			n = n.l
			for n.r != nil {
				n = n.r
			}
			continue
		}
		for ; ; n = n.d {
			if n.d == nil {
				return
			}
			if n.d.r == n {
				n = n.d
				break
			}
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

func (n *node) print(pr Printer) {
	if n == nil {
		return
	}
	var d string
	if n.isLeft() == true {
		d = "l "
	} else {
		d = "r "
	}
	s := fmt.Sprint(d, n.k)
	var sub Printer
	if n.c == red {
		sub = pr.Add(aurora.Red(s).String())
	} else {
		sub = pr.Add(aurora.Blue(s).Bold().String())
	}
	n.l.print(sub)
	n.r.print(sub)
}

func (t *Tree) Print(pr Printer) {
	s := fmt.Sprintf("[%d]", t.size)
	tree := pr.Add(s)
	t.r.print(tree)
}
