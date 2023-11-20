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
	t, err := MakeTree(pathToRoot, extFilter)
	if err != nil {
		return Tree{}, err
	}

	MakePaths(t, urlBase, extFilter, true)

	SetPrevNext(t, nil)

	return *t, nil
}

// MakeTree makes a pointer to a tree based on the relative path to the root
// dir given in pathToRoot and recursively fills out FileName, and IsDir for all nodes.
// (joining pathPadding to pathToRoot to make absolute paths)
func MakeTree(pathToRoot string, extFilter string) (t *Tree, err error) {
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
			child, err := MakeTree(filepath.Join(pathToRoot, c.Name()), extFilter)
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

// TODO fill out prev and next and URLPath

// MakePaths walks through a tree and fills in each node's URL Path
// according to the base string and the cut out file extension filter
// TODO make internal func and bootstrap with MakeTree to a generic
// GenerateTree func
func MakePaths(t *Tree, base string, cutExt string, first bool) {
	t.URLPath = filepath.Join(base, strings.TrimSuffix(t.FileName, cutExt))
	if first {
		t.URLPath = filepath.Join(base)
	}

	if t.IsDir {
		for _, ch := range t.Children {
			MakePaths(ch, t.URLPath, cutExt, false)
		}
	}
}

// SetPrevNext *supposed* to set all the end nodes to point to each other in a
// really convenient way, but it doesn't work :(((
func SetPrevNext(t *Tree, lastLeaf **Tree) {
	if !t.IsDir {
		if lastLeaf != nil {
			(*lastLeaf).Next = t
			fmt.Println((*lastLeaf).FileName, ".Next is now", t.FileName) // TODO
			t.Prev = *lastLeaf
			fmt.Println(t.FileName, ".Previous is now", (*lastLeaf).FileName) // TODO
		}
		lastLeaf = &t
		fmt.Println("setting lastLeaf to", t.FileName) // TODO rm debugging
	}

	if t.IsDir {
		for _, ch := range t.Children {
			SetPrevNext(ch, lastLeaf)
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
