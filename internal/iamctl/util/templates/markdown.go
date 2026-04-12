// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package templates

import (
	"bytes"
	"fmt"
	"io"

	"github.com/russross/blackfriday"
)

const linebreak = "\n"

// ASCIIRenderer implements blackfriday.Renderer for v2 API.
var _ blackfriday.Renderer = &ASCIIRenderer{}

// ASCIIRenderer is a blackfriday.Renderer intended for rendering markdown
// documents as plain text, well suited for human reading on terminals.
type ASCIIRenderer struct {
	Indentation string

	listItemCount uint
	listLevel     uint
}

// RenderNode implements the v2 Renderer interface.
func (r *ASCIIRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Text:
		r.renderText(w, node)
	case blackfriday.Paragraph:
		_, _ = w.Write([]byte(linebreak))
	case blackfriday.Heading:
		if entering {
			_, _ = w.Write([]byte(linebreak))
		}
	case blackfriday.HorizontalRule:
		_, _ = w.Write([]byte(linebreak + "----------" + linebreak))
	case blackfriday.BlockQuote:
		r.renderBlockQuote(w, node, entering)
	case blackfriday.CodeBlock:
		r.renderCodeBlock(w, node)
	case blackfriday.Code:
		r.renderCode(w, node)
	case blackfriday.Emph:
		r.renderEmph(w, node, entering)
	case blackfriday.Strong:
		r.renderStrong(w, node, entering)
	case blackfriday.Link:
		r.renderLink(w, node, entering)
	case blackfriday.Image:
		r.renderImage(w, node, entering)
	case blackfriday.List:
		r.renderList(w, node, entering)
	case blackfriday.Item:
		r.renderItem(w, node, entering)
	default:
		r.fw(w, node.Literal)
	}
	return blackfriday.GoToNext
}

func (r *ASCIIRenderer) renderText(w io.Writer, node *blackfriday.Node) {
	raw := string(node.Literal)
	lines := bytes.Split([]byte(raw), []byte(linebreak))
	for i, line := range lines {
		trimmed := bytes.Trim(line, " \n\t")
		if len(trimmed) > 0 && trimmed[0] != '_' {
			_, _ = w.Write([]byte(" "))
		}
		_, _ = w.Write(trimmed)
		if i < len(lines)-1 {
			_, _ = w.Write([]byte(linebreak))
		}
	}
}

func (r *ASCIIRenderer) renderBlockQuote(w io.Writer, node *blackfriday.Node, entering bool) {
	r.fw(w, node.Literal)
}

func (r *ASCIIRenderer) renderCodeBlock(w io.Writer, node *blackfriday.Node) {
	_, _ = w.Write([]byte(linebreak))
	lines := bytes.Split(node.Literal, []byte(linebreak))
	for _, line := range lines {
		indented := append([]byte(r.Indentation), line...)
		_, _ = w.Write(indented)
		_, _ = w.Write([]byte(linebreak))
	}
}

func (r *ASCIIRenderer) renderCode(w io.Writer, node *blackfriday.Node) {
	r.fw(w, node.Literal)
}

func (r *ASCIIRenderer) renderEmph(w io.Writer, node *blackfriday.Node, entering bool) {
	r.fw(w, node.Literal)
}

func (r *ASCIIRenderer) renderStrong(w io.Writer, node *blackfriday.Node, entering bool) {
	r.fw(w, node.Literal)
}

func (r *ASCIIRenderer) renderLink(w io.Writer, node *blackfriday.Node, entering bool) {
	_, _ = w.Write([]byte(" "))
	r.fw(w, node.Destination)
}

func (r *ASCIIRenderer) renderImage(w io.Writer, node *blackfriday.Node, entering bool) {
	r.fw(w, node.Destination)
}

func (r *ASCIIRenderer) renderList(w io.Writer, node *blackfriday.Node, entering bool) {
	if entering {
		r.listLevel++
		_, _ = w.Write([]byte(linebreak))
	} else {
		r.listLevel--
	}
}

func (r *ASCIIRenderer) renderItem(w io.Writer, node *blackfriday.Node, entering bool) {
	if entering {
		indent := bytes.Repeat([]byte(r.Indentation), int(r.listLevel))
		_, _ = w.Write(indent)

		if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
			r.listItemCount++
			fmt.Fprintf(w, "%d.", r.listItemCount)
		} else {
			_, _ = w.Write([]byte("*"))
		}
		_, _ = w.Write([]byte(" "))
	}
}

// RenderHeader implements the v2 Renderer interface.
func (r *ASCIIRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
}

// RenderFooter implements the v2 Renderer interface.
func (r *ASCIIRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
}

func (r *ASCIIRenderer) fw(w io.Writer, text ...[]byte) {
	for _, t := range text {
		_, _ = w.Write(t)
	}
}
