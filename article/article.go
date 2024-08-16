package article

import (
	"bytes"
	"embed"
	"fmt"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed assets/*.md
var Assets embed.FS

type Parser struct {
	markdown goldmark.Markdown
}

func NewParser() Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.TabWidth(4),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	return Parser{markdown: md}
}

func (p Parser) Parse(filename string) (string, error) {
	file, err := Assets.ReadFile(fmt.Sprintf("assets/%v.md", filename))
	if err != nil {
		return "", err
	}

	var htmlOutput bytes.Buffer
	if err := p.markdown.Convert(file, &htmlOutput); err != nil {
		return "", err
	}

	return htmlOutput.String(), nil
}

func GetAssets() ([]string, error) {
	files, err := Assets.ReadDir("assets")
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}
