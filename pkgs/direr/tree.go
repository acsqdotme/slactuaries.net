package direr

import (
	"fmt"
	"os"
	"path/filepath"
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

// GenerateTree is gonna be my wrapper for all my little Tree making funcs once
// they're all functional
func GenerateTree(pathToRoot string, extFilter string, urlBase string) (Tree, error) {
	t, err := scanTree(pathToRoot, extFilter)
	if err != nil {
		return Tree{}, err
	}

	makePaths(t, urlBase, extFilter, true)
	setPrevNext(t)

	return *t, nil
}

// scanTree makes a pointer to a tree based on the relative path to the root
// dir given in pathToRoot and recursively fills out FileName, and IsDir for all nodes.
// (joining pathPadding to pathToRoot to make absolute paths)
func scanTree(pathToRoot string, extFilter string) (t *Tree, err error) {
	info, err := os.Lstat(pathToRoot)
	t = &Tree{
		FileName: info.Name(),
		IsDir:    info.IsDir(),
	}

	if !t.IsDir && !strings.HasSuffix(t.FileName, extFilter) {
		return nil, nil
	}

	if t.IsDir {
		contents, err := os.ReadDir(pathToRoot)
		if err != nil {
			return nil, err
		}

		for _, c := range contents {
			child, err := scanTree(filepath.Join(pathToRoot, c.Name()), extFilter)
			if err != nil {
				return nil, err
			}

			if child != nil {
				t.Children = append(t.Children, child)
			}
		}
	}

	if t.IsDir && len(t.Children) < 1 {
		return nil, nil
	}

	return t, nil
}

// makePaths walks through a tree and fills in each node's URL Path
// according to the base string and the cut out file extension filter
// TODO make internal func and bootstrap with scanTree to a generic
// GenerateTree func
func makePaths(t *Tree, base string, cutExt string, first bool) {
	t.URLPath = filepath.Join(base, strings.TrimSuffix(t.FileName, cutExt))
	if first {
		t.URLPath = filepath.Join(base)
	}

	if t.IsDir {
		for _, ch := range t.Children {
			makePaths(ch, t.URLPath, cutExt, false)
		}
	}
}

// lastLeaf is a helper variable for setPrevNext
// TODO make unglobal
var lastLeaf **Tree

// setPrevNext flattens and sets previous and next pointers for all files in
// the tree
func setPrevNext(t *Tree) {
	if !t.IsDir {
		if lastLeaf != nil {
			(*lastLeaf).Next = t
			t.Prev = *lastLeaf
		}
		lastLeaf = &t
	}

	if t.IsDir {
		for _, child := range t.Children {
			setPrevNext(child)
		}
	}
}

// PrintTree is for testing
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
