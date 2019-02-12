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
	"flag"
)

type rbTree struct {
	stacked    bool    // track parent reference on stack
	keyValue   bool    // key value pairs
	unque      bool    // unique (single value per node)
	threadSafe bool    // thread safe tree
	typ        string  // type of item
	valueType  string  // type of value
	comparable string  // type is comparable
	less       string  // less format
	equal      string  // equal format
	prefix     string  // name space prefix
	imports    Strings // add imports
	treeName   string  // tree type name
	printer    bool    // implement printer interface
	pkgName    string  // package name
	output     string  // output file name
}

func genRBTree(args []string) {

	var tree rbTree

	set := flag.NewFlagSet("rbtree", flag.ExitOnError)

	set.BoolVar(&tree.stacked,
		"stacked",
		false,
		"track parent references on stack")
	set.BoolVar(&tree.leftLeaning,
		"ll",
		false,
		"left-leaning Red-black tree")
	set.BoolVar(&tree.unque,
		"unique",
		false,
		"don't allow many values per node")
	set.BoolVar(&tree.threadSafe,
		"thread-safe",
		false,
		"thread-safe tree with")
	set.StringVar(&tree.typ,
		"type",
		"",
		"type of item or key")
	set.StringVar(&tree.valueType,
		"value",
		"",
		"type of value, produce tree with key-value pairs")
	set.BoolVar(&tree.comparable,
		"comparable",
		false,
		"the type is comparable using '<' and '==' operators")
	set.StringVar(&tree.less,
		"less",
		"%s < %s",
		"format of less comparison, like '%s.Less(%s)' or 'less(%s, %s)', etc")
	set.StringVar(&tree.equal,
		"equal",
		"%s == %s",
		"format of equal comparison, like '%s.Eq(%s)' or 'equal(%s, %s)', etc")
	set.StringVar(&tree.printer,
		"print",
		false,
		"add Print tree method")
	set.Var(&tree.imports,
		"import",
		"import package (reuse flag for list of packages)")
	set.StringVar(&tree.prefix,
		"prefix",
		"",
		"name space prefix")
	set.StringVar(&tree.treeName,
		"tree",
		"Tree",
		"tree type name")
	set.StringVar(&tree.treeName,
		"tree",
		"Tree",
		"tree type name")
	set.StringVar(&tree.pkgName,
		"package",
		"",
		"package name")
	set.StringVar(&tree.output,
		"o",
		"",
		"output file name")
	set.Parse(args)

}

const rbTreeTemplate = `package {{ .pkgName }}

{{ if .imports }}
import (
	{{ for .imports }}
	"{{ . }}"
	{{ end }}
)
{{ end }}

type {{ .treeName }} struct {
	{{ if .threadSafe }}
	sync.Mutex
	{{ end }}
}


`
