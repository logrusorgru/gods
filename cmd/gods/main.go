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
	"flags"
	"fmt"
	"io"
	"os"
	"strings"
)

const version = "1.0"

func showHelp(out io.Writer, code int) {
	fmt.Fprintf(out, `The Gods is a generator of Golang data structures

Supported data structures:

    rbtree     Red-black tree
    avltree    AVL-tree
    version    show generator version

Use '%s help [data structure]' for details.
`, os.Args[0])
	os.Exit(code)
}

func main() {

	if len(os.Args) < 2 {
		showHelp(os.Stderr, 1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "rbtree":
		genRBTree(os.Args[1:])
	case "avltree":
		genAVLTree(os.Args[1:])
	case "version":
		fmt.Println("gods", version)
	case "help":
		showHelp(os.Stdout, 0)
	}

}

// Strings represents list of strings
type Strings []string

// String implements flag.Value interface
func (s Strings) String() (ss string) {
	for _, x := range s {
		ss += x + ","
	}
	if len(ss) > 0 {
		ss = ss[len(ss)-1:] // trim trailing comma
	}
	return
}

// Set implements flag.Value interface
func (s *Strings) Set(val string) (_ error) {
	*s = append(*s, val)
	return
}

// Contain returns true if the
// Strings contain given element
func (s Strings) Contain(name string) bool {
	for _, x := range s {
		if x == name {
			return true
		}
	}
	return false
}
