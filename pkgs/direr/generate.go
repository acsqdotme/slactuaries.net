package direr

import (
	"os"
	"path/filepath"
	"strings"
)

// GenerateTree is gonna be my wrapper for all my little Tree making funcs once
// they're all functional
func GenerateTree(pathToRoot string, extFilter string, urlBase string) (*Tree, error) {
	t, err := scanTree(pathToRoot, extFilter)
	if err != nil {
		return nil, err
	}

	makeURLPaths(t, urlBase, true)
	makeFilePaths(t, pathToRoot, true)

	// lastLeaf is the variable that setPrevNext uses to flatten out files and
	// point them to each other
	var lastLeaf **Tree

	// setPrevNext preemptively declared as it's recursively called in its
	// definition
	var setPrevNext func(*Tree)
	setPrevNext = func(t *Tree) {
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
	setPrevNext(t)

	return t, nil
}

// scanTree makes a pointer to a tree based on the relative path to the root
// dir given in pathToRoot and recursively fills out FileName, and IsDir for all nodes.
// (joining pathPadding to pathToRoot to make absolute paths)
func scanTree(pathToRoot string, extFilter string) (t *Tree, err error) {
	info, err := os.Lstat(pathToRoot)
	t = &Tree{
		FileName: strings.TrimSuffix(info.Name(), extFilter),
		IsDir:    info.IsDir(),
	}

	if !t.IsDir && !strings.HasSuffix(info.Name(), extFilter) {
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
// according to the base string
func makeURLPaths(t *Tree, base string, first bool) {
	t.URLPath = filepath.Join(base, t.FileName)
	if first {
		t.URLPath = filepath.Join(base)
	}

	if t.IsDir {
		for _, ch := range t.Children {
			makeURLPaths(ch, t.URLPath, false)
		}
	}
}

// makeFilePaths walks through a tree and fills in each node's File Path
// according to the base string
// TODO don't be lazy; figure out how to combine with makeURLPaths
func makeFilePaths(t *Tree, base string, first bool) {
	t.FilePath = filepath.Join(base, t.FileName)
	if first {
		t.FilePath = filepath.Join(base)
	}

	if t.IsDir {
		for _, ch := range t.Children {
			makeFilePaths(ch, t.FilePath, false)
		}
	}
}
