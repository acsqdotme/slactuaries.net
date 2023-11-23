package direr

import (
	"fmt"
	"strings"
)

// Tree is the data struct for files and directories in a defined root.
// FileName, Path, and IsDir are exactly what they sound like. Prev and Next
// are meant only for going from file end node to file end node, so no dir will
// have those filled in. And obviously, only dirs can have Children.
type Tree struct {
	FileName string
	IsDir    bool
	URLPath  string
	Prev     *Tree
	Next     *Tree
	Children []*Tree
}



// PrintTree is for testing; TODO turn into real method
func PrintTree(t Tree, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}

	if t.IsDir {
		fmt.Print("📁")
	} else {
		fmt.Print("📃")
	}

	prev := "<nil>"
	next := "<nil>"
	if t.Prev != nil {
		prev = t.Prev.FileName
	}
	if t.Next != nil {
		next = t.Next.FileName
	}

	fmt.Printf("%s ; \t\tprev: %s \t\tnext: %s \t\t %s\n", t.FileName, prev, next, t.URLPath)
	if t.IsDir {
		for _, ch := range t.Children {
			PrintTree(*ch, indent+1)
		}
	}
}
