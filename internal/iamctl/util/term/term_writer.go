// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package term

import (
	"io"
	"os"

	"github.com/mitchellh/go-wordwrap"
	"github.com/moby/term"
)

type wordWrapWriter struct {
	limit  uint
	writer io.Writer
}

// NewResponsiveWriter creates a Writer that detects the column width of the
// terminal we are in, and adjusts every line width to fit and use recommended
// terminal sizes for better readability. Does proper word wrapping automatically.
//
//	if terminal width >= 120 columns		use 120 columns
//	if terminal width >= 100 columns		use 100 columns
//	if terminal width >=  80 columns		use  80 columns
//
// In case we're not in a terminal or if it's smaller than 80 columns width,
// doesn't do any wrapping.
func NewResponsiveWriter(w io.Writer) io.Writer {
	file, ok := w.(*os.File)
	if !ok {
		return w
	}
	fd := file.Fd()
	if !term.IsTerminal(fd) {
		return w
	}

	terminalSize := GetSize(fd)
	if terminalSize == nil {
		return w
	}

	var limit uint
	switch {
	case terminalSize.Width >= 120:
		limit = 120
	case terminalSize.Width >= 100:
		limit = 100
	case terminalSize.Width >= 80:
		limit = 80
	}

	return NewWordWrapWriter(w, limit)
}

// NewWordWrapWriter is a Writer that supports a limit of characters on every line
// and does auto word wrapping that respects that limit.
func NewWordWrapWriter(w io.Writer, limit uint) io.Writer {
	return &wordWrapWriter{
		limit:  limit,
		writer: w,
	}
}

func (w wordWrapWriter) Write(p []byte) (nn int, err error) {
	if w.limit == 0 {
		return w.writer.Write(p)
	}
	original := string(p)
	wrapped := wordwrap.WrapString(original, w.limit)
	return w.writer.Write([]byte(wrapped))
}

// NewPunchCardWriter is a NewWordWrapWriter that limits the line width to 80 columns.
func NewPunchCardWriter(w io.Writer) io.Writer {
	return NewWordWrapWriter(w, 80)
}

type maxWidthWriter struct {
	maxWidth     uint
	currentWidth uint
	written      uint
	writer       io.Writer
}

// NewMaxWidthWriter is a Writer that supports a limit of characters on every
// line, but doesn't do any word wrapping automatically.
func NewMaxWidthWriter(w io.Writer, maxWidth uint) io.Writer {
	return &maxWidthWriter{
		maxWidth: maxWidth,
		writer:   w,
	}
}

func (m maxWidthWriter) Write(p []byte) (nn int, err error) {
	for _, b := range p {
		if m.currentWidth == m.maxWidth {
			_, _ = m.writer.Write([]byte{'\n'})
			m.currentWidth = 0
		}
		if b == '\n' {
			m.currentWidth = 0
		}
		_, err := m.writer.Write([]byte{b})
		if err != nil {
			return int(m.written), err
		}
		m.written++
		m.currentWidth++
	}
	return len(p), nil
}
