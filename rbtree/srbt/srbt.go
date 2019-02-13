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

// Package SRRB represents stacked Red-black tree.
// A node of the Tree doesn't contain parent
// reference. This way any insert, delete and walk
// operation requires some stack space to keep way
// of references to root.
package rbtree

// LessFunc
type LessFunc func(a, b interface{}) bool

// EqualFunc returns true if a is equal to b.
type EqualFunc func(a, b interface{}) bool

// ZeroFunc return true if given item is zero.
type ZeroFunc func(a interface{}) bool

type color bool

const (
	black color = false
	red   color = true
)

type node struct {
	left  *node
	right *node
	color color

	item interface{}
}

func (n *node) isRed() bool {
	return n.color == red
}

func (n *node) isBlack() bool {
	return n.color == black
}

func (n *node) isSentinel() bool {
	return n == &sentinel
}

func (n *node) replaceChild(old, new *node) {
	if n.left == old {
		n.left = new
	} else {
		n.right = new
	}
}

func (n *node) oppositeChild(x *node) *node {
	if n.left == x {
		return n.right
	}
	return n.left
}

// successor for this deleted node with
// path to it; the n is not the sentinel
func (n *node) successor(br []*node) (sr []*node, x *node) {
	var end = &sentinel
	for x = n.right; x.left != end; x = x.left {
		sr = append(sr, x)
	}
	return
}

// BST-replacement for this deleted node with
// relative path; then n is not the sentinel
func (n *node) replacement(br []*node) (sr []*node, x *node) {
	var end = &sentinel
	br = push(br, n) // this
	switch {
	case n.left != end && n.right != end:
		return n.successor(br) // min of the right
	case n.left != end:
		return br, n.left // has the left only
	case n.right != end:
		return br, n.right // has the right only
	default:
		// no successors
	}
	return br, end
}

var sentinel node

// initialize the sentinel
func init() {
	sentinel.left = &sentinel
	sentinel.right = &sentinel
	sentinel.color = black
	sentinel.item = nil
}

func push(br []*node, n *node) []*node {
	return append(br, n)
}

func pop(br []*node) (nb []*node, n *node) {
	if len(br) == 0 {
		return // nil, nil
	}
	n, nb = br[len(br)-1], br[:len(br)-1]
	return
}

type Tree struct {
	root *node

	less  LessFunc
	equal EqualFunc
	zero  ZeroFunc

	size int // number of items
}

func New(less LessFunc, equal EqualFunc, zero ZeroFunc) (t *Tree) {
	t = new(Tree)
	t.root = &sentinel
	t.less = less
	t.equal = equal
	t.zero = zero
	return
}

func (t *Tree) isRoot(n *node) bool {
	return t.root == n
}

func (t *Tree) findNode(item interface{}) *node {
	var (
		n     = t.root
		less  = t.less
		equal = t.equal

		end = &sentinel
	)

	for n != end {
		switch {
		case equal(item, n.item):
			return n // found
		case less(item, n.item):
			n = n.left
		default:
			n = n.right
		}
	}

	return nil // not found
}

// Get returns true if the Tree contains given item.
func (t *Tree) Get(item interface{}) bool {
	if n := t.findNode(item); n != nil {
		return true
	}
	return false
}

func (t *Tree) replaceChild(n *node, old, new *node) {
	if n == nil {
		t.root = new
	} else {
		n.replaceChild(old, new)
	}
}

func (t *Tree) rotateLeft(br []*node, x *node) {
	var d, pivot *node
	br, d = pop(br)
	pivot = x.right
	t.replaceChild(d, x, pivot) // the pivot might become root
	if d != nil {
		d.replaceChild(x, pivot)
	}
	x.right = pivot.left
	pivot.left = x
}

// the x already popped from the branch, i.e. last
// element of the branch is parent of the x; if the
// branch is empty, then right node of the x becomes
// root of the Tree;
func (t *Tree) rotateRight(br []*node, x *node) {
	var d, pivot *node
	br, d = pop(br)
	pivot = x.left
	t.replaceChild(d, x, pivot) // the pivot might become root
	if d != nil {
		d.replaceChild(x, pivot)
	}
	x.left = pivot.right
	pivot.right = x
}

// d is red, u is black
func (t *Tree) insert4(br []*node, d, g, u, x *node) {
	// (not the loop)
	if g.left == d {
		if d.left == x {
			// 1. left left case
			t.rotateRight(br, g)                //
			g.color, d.color = d.color, g.color // swap colors
		} else {
			// 2. left right case
			t.rotateLeft(push(br, g), d)
			// and apply the left left case
			t.rotateRight(br, g)
			g.color, x.color = x.color, g.color // swap colors
		}
	} else {
		if d.right == x {
			// 3. right right case
			t.rotateLeft(br, g)                 //
			g.color, d.color = d.color, g.color // swap colors
		} else {
			// 4. right left case
			t.rotateRight(push(br, g), d)
			// and apply the left left case
			t.rotateLeft(br, g)
			g.color, x.color = x.color, g.color // swap colors
		}
	}
}

func (t *Tree) insertNode(br []*node, x *node) {
	// if the Tree is empty
	if len(br) == 0 {
		t.root = x
		x.color = black // root must be black
		return
	}
	var d, g, u *node // dad, granddad, uncle
	for t.isRoot(x) == false {
		br, d = pop(br)
		if d.isRed() == true {
			br, g = pop(br)
			u = g.oppositeChild(d)
			if u.isRed() == true {
				d.color, u.color = black, black
				g.color, x = red, g
				continue
			} else { // the u is black
				t.insert4(br, d, g, u, x)
				return // done
			}
		} else { // d is black
			return // done
		}
	}
	x.color = black // change root to black
	return
}

// find node and track branch
func (t *Tree) findNodeBranch(item interface{}) (br []*node, n *node) {
	var (
		less  = t.less
		equal = t.equal

		end = &sentinel
	)

	for n = t.root; n != end; {
		if equal(item, n.item) == true {
			return // found
		}
		br = push(br, n)
		if less(item, n.item) == true {
			n = n.left
		} else {
			n = n.right
		}
	}

	n = nil // not found
	return
}

// Ins inserts given item to the Tree. The Ins
// returns true if item not overwritten.
func (t *Tree) Ins(item interface{}) (ok bool) {
	var br, n = t.findNodeBranch(item)

	// found
	if n != nil {
		n.item = item
		return // false, no new node created
	}

	// not found
	var end = &sentinel // }
	t.size++            // } new node will be created
	ok = true           // }

	// create new red node and insert it using the branch
	t.insertNode(br, &node{
		left:  end,  // nil
		right: end,  // nil
		color: red,  // new is red
		item:  item, // item
	})

	return
}

// InsNx creates new item and returns true. If item already
// exists, then the InsNx returns false without touching the
// Tree. Other words the InsNx is create. Mnemonic is
// 'insert if not exists'.
func (t *Tree) InsNx(item interface{}) (ok bool) {
	var br, n = t.findNodeBranch(item)

	// found
	if n != nil {
		return // false, already exists
	}

	// not found, create
	var end = &sentinel // }
	t.size++            // } new node will be created
	ok = true           // }

	// create new red node and insert it using the branch
	t.insertNode(br, &node{
		left:  end,  // nil
		right: end,  // nil
		color: red,  // new is red
		item:  item, // item
	})

	return
}

// InsEx overwrites existing item and returns true. If item
// doesn't exist in the Tree, then the InsEx returns false
// without touching the Tree. Other words, the InsEx is
// overwrite. Mnemonic is 'insert if exists'.
func (t *Tree) InsEx(item interface{}) (ok bool) {
	var (
		br    []*node
		n     = t.root
		less  = t.less
		equal = t.equal

		end = &sentinel
	)

	for n != end {
		if equal(item, n.item) == true {
			n.item = item
			return true // overwritten
		}
		br = push(br, n)
		if less(item, n.item) == true {
			n = n.left
		} else {
			n = n.right
		}
	}

	return // false, not found
}

func (t *Tree) fixDoubleBlack(br []*node, n *node) {
	if t.isRoot(n) == true {
		return
	}
	var (
		d, s *node
		end  = &sentinel
	)
	// the n is not root, then the d is not nil
	br, d = pop(br)
	if s = d.oppositeChild(n); s == end {
		t.fixDoubleBlack(br, d) // push up
		return
	}
	if s.isRed() == true {
		d.color, s.color = red, black
		if d.left == s {
			t.rotateRight(br, d)
		} else {
			t.rotateLeft(br, d)
		}
		// rotated around d, then we have to fix the br content;
		// now parent of the d points to the n; but if the d is root
		// of the tree, then n becomes root; since, the d already
		// popped, thus we can leave br as it
		t.fixDoubleBlack(br, n)
		return
	}
	// s is black
	if s.left.isRed() || s.right.isRed() {
		if s.left.isRed() == true {
			if d.left == s {
				s.left.color, s.color = s.color, d.color
				t.rotateRight(br, d)
			} else {
				s.left.color = d.color
				t.rotateRight(push(br, d), s)
				t.rotateLeft(br, d)
			}
		} else {
			// the s.right is red (is not the sentinel)
			if d.left == s {
				s.right.color = d.color
				t.rotateLeft(push(br, d), s)
				t.rotateRight(br, d)
			} else {
				s.right.color, s.color = s.color, d.color
				t.rotateLeft(br, d)
			}
		}
		d.color = black
		return
	}
	s.color = red
	if d.isBlack() == true {
		t.fixDoubleBlack(br, d)
	} else {
		d.color = black
	}
}

// delete given node with branch of its ancestors
func (t *Tree) delete(br []*node, n *node) {
	var (
		sr   []*node
		d, r *node
		end  = &sentinel
	)
	br, d = pop(br)
	sr, r = n.replacement(br)
	if r == end {
		if d == nil {
			t.root = end // the n in root, the Tree becomes empty
			return
		}
		// the n is not root, then the d is not nil
		if r.isBlack() && n.isBlack() == true {
			t.fixDoubleBlack(br, n)
		} else {
			// sibling
			if s := d.oppositeChild(n); s != end {
				s.color = red
			}
		}
		d.replaceChild(n, end)
		return
	}
	// the r is not sentinel
	if n.left == end || n.right == end {
		if t.isRoot(n) == true {
			n.item = r.item
			n.left, n.right = end, end
		} else {
			// the n is not root, then d is not nil
			d.replaceChild(n, r)
			if n.isBlack() && r.isBlack() == true {
				t.fixDoubleBlack(sr, r)
			} else {
				r.color = black
			}
		}
		return
	}
	r.item, n.item = n.item, r.item
	t.delete(sr, r)
	return
}

func (t *Tree) Del(item interface{}) (ok bool) {
	var br, n = t.findNodeBranch(item)
	if n == nil {
		return // false, no such item
	}
	ok = true       // found
	t.size--        // reduce
	t.delete(br, n) // delete
	return
}

// Size returns number of items in the tree
func (t *Tree) Size() int {
	return t.size
}

// Clear the Tree.
func (t *Tree) Clear() {
	t.root, t.size = &sentinel, 0
}

func (t *Tree) min() (n *node) {
	var end = &sentinel
	for n = t.root; n.left != end; n = n.left {
	}
	return
}

// Min item of the Tree. It returns (nil, false)
// if the Tree is empty.
func (t *Tree) Min() (interface{}, bool) {
	var n = t.min()
	if n == nil {
		return nil, false // the Tree is empty
	}
	return n.item, true
}

func (t *Tree) max() (n *node) {
	var end = &sentinel
	for n = t.root; n.right != end; n = n.right {
	}
	return
}

// Max item of the Tree. It returns (nil, false)
// if the Tree if empty.
func (t *Tree) Max() (interface{}, bool) {
	var n = t.max()
	if n == nil {
		return nil, false // the Tree is empty
	}
	return n.item, true
}

type WalkFunc func(item interface{}) (next bool)

// Walk over all items of the Tree without any order.
// The given WalkFunc must not change the Tree.
//
// TODO (kostyarin): allow changes (?)
func (t *Tree) Walk(walkFunc WalkFunc) {
	var (
		rs []*node  // rights
		n  = t.root //

		end = &sentinel
	)
	for n != end && len(rs) > 0 {
		for n != end {
			if walkFunc(n) == false {
				return
			}
			if n.right != end {
				rs = push(rs, n.right)
			}
			n = n.left
		}
		if len(rs) > 0 {
			rs, n = pop(rs)
		}
	}
}

// min node and track branch
func (t *Tree) minBranch() (br []*node, n *node) {
	var end = &sentinel
	for n = t.root; n.left != end; n = n.left {
		br = push(br, n)
	}
	return
}

// max node and track branch
func (t *Tree) maxBranch() (br []*node, n *node) {
	var end = &sentinel
	for n = t.root; n.right != end; n = n.right {
		br = push(br, n)
	}
	return
}

// [from, +inf)
func (t *Tree) ascendFrom(from interface{}, ascendFunc WalkFunc) {
	var (
		br, n = t.findNodeBranch(from)
		end   = &sentinel
	)
	if n == nil {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		if ascendFunc(n) == false {
			return
		}
		if n.right != end {
			n = n.right
			for n.left != end {
				br = push(br, n)
				n = n.left
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// (-inf, to]
func (t *Tree) ascendTo(to interface{}, ascendFunc WalkFunc) {
	var (
		br, n = t.minBranch()

		less  = t.less
		equal = t.equal

		end = &sentinel
	)
	if n == end {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		// TODO (kostyarin): the less && equal
		if less(to, n.item) == false && equal(to, n.item) == false {
			return // that's all
		}
		if ascendFunc(n) == false {
			return
		}
		if n.right != end {
			n = n.right
			for n.left != end {
				br = push(br, n)
				n = n.left
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// [from, to]
func (t *Tree) ascendFromTo(from, to interface{}, ascendFunc WalkFunc) {
	var (
		br, n = t.findNodeBranch(from)

		less  = t.less
		equal = t.equal

		end = &sentinel
	)
	if n == nil {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		if less(to, n.item) && equal(to, n.item) == false {
			return // that's all
		}
		if ascendFunc(n) == false {
			return
		}
		if n.right != end {
			n = n.right
			for n.left != end {
				br = push(br, n)
				n = n.left
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// (-inf, +inf)
func (t *Tree) ascend(ascendFunc WalkFunc) {
	var (
		br, n = t.minBranch()
		end   = &sentinel
	)
	if n == end {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		if ascendFunc(n) == false {
			return
		}
		if n.right != end {
			n = n.right
			for n.left != end {
				br = push(br, n)
				n = n.left
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

func (t *Tree) Ascend(from, to interface{}, ascendFunc WalkFunc) {
	// zero function
	switch zero := t.zero; {
	case zero(from) == true: // (-inf, to] or (-inf, +inf)
		if zero(to) == true {
			t.ascend(ascendFunc) // (-inf, +inf)
		} else {
			t.ascendTo(to, ascendFunc) // (-inf, to]
		}
	case zero(to) == true: // [from, +inf)
		t.ascendFrom(from, ascendFunc)
	default: // [from, to]
		t.ascendFromTo(from, to, ascendFunc)
	}
}

// [from, -inf) (reversed)
func (t *Tree) descendFrom(from interface{}, descendFunc WalkFunc) {
	var (
		br, n = t.findNodeBranch(from)
		end   = &sentinel
	)
	if n == nil {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		if descendFunc(n) == false {
			return
		}
		if n.left != end {
			n = n.left
			for n.right != end {
				br = push(br, n)
				n = n.right
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// (+inf, to] (reversed)
func (t *Tree) descendTo(to interface{}, descendFunc WalkFunc) {
	var (
		br, n = t.maxBranch()

		less  = t.less
		equal = t.equal

		end = &sentinel
	)
	if n == end {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		// TODO (kostyarin): the less && equal
		if less(to, n.item) == false && equal(to, n.item) == false {
			return // that's all
		}
		if descendFunc(n) == false {
			return
		}
		if n.left != end {
			n = n.left
			for n.right != end {
				br = push(br, n)
				n = n.right
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// [from, to] (reversed)
func (t *Tree) descendFromTo(from, to interface{}, descendFunc WalkFunc) {
	var (
		br, n = t.findNodeBranch(from)

		less  = t.less
		equal = t.equal

		end = &sentinel
	)
	if n == nil {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		// TODO (kostyarin): the less && equal
		if less(to, n.item) == false && equal(to, n.item) == false {
			return // that's all
		}
		if descendFunc(n) == false {
			return
		}
		if n.left != end {
			n = n.left
			for n.right != end {
				br = push(br, n)
				n = n.right
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

// (-inf, +inf) (reversed)
func (t *Tree) descend(descendFunc WalkFunc) {
	var (
		br, n = t.maxBranch()
		end   = &sentinel
	)
	if n == end {
		if len(br) == 0 {
			return
		}
		br, n = pop(br)
	}
	for n != end {
		if descendFunc(n) == false {
			return
		}
		if n.left != end {
			n = n.left
			for n.right != end {
				br = push(br, n)
				n = n.right
			}
			continue // the n pushed and popped virtually
		}
		br, n = pop(br)
	}
}

func (t *Tree) Descend(from, to interface{}, descendFunc WalkFunc) {
	// zero function
	switch zero := t.zero; {
	case zero(from) == true: // (-inf, to] or (-inf, +inf)
		if zero(to) == true {
			t.descend(descendFunc) // (-inf, +inf)
		} else {
			t.descendTo(to, descendFunc) // (-inf, to]
		}
	case zero(to) == true: // [from, +inf)
		t.descendFrom(from, descendFunc)
	default: // [from, to]
		t.descendFromTo(from, to, descendFunc)
	}
}
