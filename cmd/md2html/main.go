package main

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/toc"
)

const (
	tmplFileExt = ".tmpl.html"
	htmlDir     = "./html"
)

// Front is Frontmatter metadata from markdown docs;
// only exported to interact with go.abhg.dev/goldmark/frontmatter
type Front struct {
	Title string `yaml:"title"`
	Desc  string `yaml:"description"`
}

// Valid is to make sure frontmatter isn't empty in a file
func (m Front) Valid() bool {
	if len(m.Desc) > 0 && len(m.Title) > 0 {
		return true
	}
	return false
}

type hashes struct{}

func (*hashes) AnchorText(h *anchor.HeaderInfo) []byte {
	return bytes.Repeat([]byte{'#'}, h.Level)
}

func md2HTML(md []byte) (html []byte, m Front, err error) {
	md2html := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
			extension.GFM,
			extension.Footnote,
			figure.Figure,
			&toc.Extender{
				Title:   "Contents",
				TitleID: "toc-title",
				ListID:  "toc",
				Compact: true,
			},
			&anchor.Extender{
				Position: anchor.Before,
				Texter:   &hashes{},
			},
			mathjax.MathJax, // TODO move to something a pkg lets me do $$ delimiters on same line as display math
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldhtml.WithUnsafe(),
		),
	)

	context := parser.NewContext()
	var buf bytes.Buffer
	err = md2html.Convert(md, &buf, parser.WithContext(context))
	if err != nil {
		return []byte{}, Front{}, err
	}

	data := frontmatter.Get(context)
	if err = data.Decode(&m); err != nil {
		return []byte{}, Front{}, err
	}

	html = buf.Bytes()
	return html, m, nil
}

func html2TMPL(html string, m Front) (tmpl []byte, err error) {
	if !m.Valid() {
		return []byte{}, errors.New("Frontmatter not set")
	}
	return []byte(`{{ define "title" }}` + m.Title + `{{ end }}
{{ define "description"}}` + m.Desc + `{{ end }}

{{ define "article" }}
` + html + `
{{ end }}`), nil
}

func convertFile(pathToMD string) (err error) {
	pathToMD = filepath.Clean(pathToMD)
	md, err := os.ReadFile(pathToMD)
	if err != nil {
		return err
	}

	html, m, err := md2HTML(md)
	if err != nil {
		return err
	}

	tmpl, err := html2TMPL(string(html), m)
	if err != nil {
		return err
	}

	pathToHTML := strings.TrimSuffix(pathToMD, ".md") + tmplFileExt

	err = os.WriteFile(pathToHTML, tmpl, 0644)
	if err != nil {
		return err
	}
	return nil
}

func convertDir(pathToDir string) (err error) {
	err = filepath.WalkDir(pathToDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			err = convertFile(path)
			if err != nil {
				return err
			}
		}
		return nil

	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// TODO actually scan through to find lessons
	err := convertDir("./html/topics/p/lessons")
	if err != nil {
		log.Println(err.Error())
	}
	err = convertDir("./html/topics/fm/lessons")
	if err != nil {
		log.Println(err.Error())
	}
}
