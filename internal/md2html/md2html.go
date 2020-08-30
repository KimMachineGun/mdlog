package md2html

import (
	"bytes"
	"fmt"
	"html"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type HTML struct {
	Content string
	Meta    map[string]interface{}
}

func Convert(source []byte) (*HTML, error) {
	gm := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			emoji.Emoji,
		),
	)

	var (
		buf     = bytes.NewBuffer(nil)
		context = parser.NewContext()
	)
	err := gm.Convert(source, buf, parser.WithContext(context))
	if err != nil {
		return nil, err
	}

	buf.WriteString(fmt.Sprintf("<details hidden id=\"blogger-raw\">%s</details>\n", html.EscapeString(string(source))))

	return &HTML{
		Content: buf.String(),
		Meta:    meta.Get(context),
	}, nil
}
