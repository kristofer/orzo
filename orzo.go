package orzo

import (
	"fmt"
	"syscall"
	"time"

	termbox "github.com/nsf/termbox-go"
)

/* orzo -- A very simple Editor in less than 1000 lines of Go code .
 *
 * -----------------------------------------------------------------------
 *
 * Copyright (c) 2018, Kristofer Younger <kryounger at gmail dot com>
 * (who ripped out all the highlighting stuff)
 * based on the work of https://github.com/antirez/kilo
 * which was marked
 * Copyright (C) 2016 Salvatore Sanfilippo <antirez at gmail dot com>
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 *
 *  *  Redistributions of source code must retain the above copyright
 *     notice, this list of conditions and the following disclaimer.
 * met:
 *
 *  *  Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTabILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES LOSS OF USE,
 * DATA, OR PROFITS OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

const orzoVersion = "2.0"
const orzoListBuffers = "*orzo Buffers*"
const tabWidth = 4

// orzo is the top-level exported type
type Orzo struct {
	Orzo *Editor
}

/* This structure represents a single line of the file we are editing. */
type erow struct {
	idx    int    /* Row index in the file, zero-based. */
	size   int    /* Size of the row, excluding the null term. */
	rsize  int    /* Size of the rendered row. */
	runes  []rune //string /* Row content. */
	render []rune /* Row content "rendered" for screen (for Tabs). */
}

type cursor struct {
	r  int /* Cursor row position */
	c  int /* Cursor column position */
	ro int /* Cursor rowoffset */
	co int /* Cursor coloffset */
}

// cursor rowoffset is the number of rows into the file of the 0-th
// screenrow. so if ro is 500, then the 0th row of the screen is
// looking at the 500th row in the file.
// It's the "scroll" offset

type Editor struct {
	events        chan termbox.Event
	Buffers       []*Buffer
	cb            *Buffer
	screenrows    int /* Number of rows that we can show */
	screencols    int /* Number of cols that we can show */
	pasteBuffer   string
	quitTimes     int
	done          bool
	statusmsg     string
	statusmsgTime time.Time
	fgcolor       termbox.Attribute
	bgcolor       termbox.Attribute
}

type TextContainer interface {
	SetPoint(pt int)
	SetMark(pt int)
	Insert(s string, pt int) error
	AddRune(r rune)
	IndexOf(idx int) byte
	DeleteRune(idx int)
	Size() int
}

type Terminal interface {
	
}

type Buffer struct {
	point    cursor
	mark     cursor
	markSet  bool
	numrows  int     /* Number of rows in file */
	rows     []*erow /* Rows */
	dirty    bool    /* File modified but not saved. */
	readonly bool
	filename string /* Currently open filename */
}

func (e *Editor) checkErr(er error) {
	if er != nil {
		e.EditorSetStatusMessage("%s", er)
	}
}

// KEY constants
const (
	KeyNull   = 0   /* NULL ctrl-space set mark */
	CtrlA     = 1   /* Ctrl-a BOL */
	CtrlB     = 2   /* Ctrl-b list Buffers */
	CtrlC     = 3   /* Ctrl-c  copy */
	CtrlE     = 5   /* Ctrl-e  EOL */
	CtrlD     = 4   /* Ctrl-d del forward? */
	CtrlF     = 6   /* Ctrl-f find */
	CtrlH     = 8   /* Ctrl-h del backward*/
	Tab       = 9   /* Tab */
	CtrlK     = 11  /* Ctrl+k killToEOL */
	CtrlL     = 12  /* Ctrl+l redraw */
	Enter     = 13  /* Enter */
	CtrlN     = 14  /* Ctrl-n nextBuffer */
	CtrlO     = 15  /* Ctrl-o load(open) file */
	CtrlQ     = 17  /* Ctrl-q quit*/
	CtrlS     = 19  /* Ctrl-s save*/
	CtrlU     = 21  /* Ctrl-u number of times??*/
	CtrlV     = 22  /* Ctrl-v paste */
	CtrlW     = 23  /* Ctrl-w kill Buffer */
	CtrlX     = 24  /* Ctrl-x cut */
	CtrlY     = 25  /* Help */
	CtrlZ     = 26  /* ?? */
	Esc       = 27  /* Escape */
	Space     = 32  /* Space */
	Backspace = 127 /* Backspace */
)

// Cursor movement keys
const (
	ArrowLeft  = termbox.KeyArrowLeft
	ArrowRight = termbox.KeyArrowRight
	ArrowUp    = termbox.KeyArrowUp
	ArrowDown  = termbox.KeyArrowDown
	DelKey     = termbox.KeyDelete
	HomeKey    = termbox.KeyHome
	EndKey     = termbox.KeyEnd
	PageUp     = termbox.KeyPgup
	PageDown   = termbox.KeyPgdn
)

// State contains the state of a terminal.
type State struct {
	termios syscall.Termios
}

// NewEditor generates a new Editor for use
func (e *Editor) initEditor() {
	e.done = false
	e.Buffers = []*Buffer{}
	e.addNewBuffer()
	e.cb.point.c, e.cb.point.r = 0, 0
	e.cb.point.ro, e.cb.point.co = 0, 0
	e.cb.numrows = 0
	e.cb.rows = []*erow{}
	e.cb.dirty = false
	e.screencols, e.screenrows = termbox.Size()
	e.screenrows -= 2 /* Get room for status bar. */
	e.quitTimes = 3
	e.fgcolor, e.bgcolor = termbox.ColorDefault, termbox.ColorDefault // retrieved from environment
}

func (e *Editor) resize() {
	e.screencols, e.screenrows = termbox.Size()
	e.screenrows -= 2 /* Get room for status bar. */
}

func (e *Editor) readOnly() {
	e.EditorSetStatusMessage("Buffer %s is Read Only", e.cb.filename)
}

// Start runs an Editor
func (z *Orzo) Start(filename string) {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	z.Orzo = &Editor{}
	e := z.Orzo
	e.initEditor()

	err = e.EditorOpen(filename)
	if err != nil {
		termbox.Close()
		fmt.Printf("orzo: error %s", err)
	}
	termbox.SetOutputMode(termbox.OutputNormal)
	termbox.SetInputMode(termbox.InputAlt | termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	e.EditorSetStatusMessage("CTRL-Y = HELP | Ctrl-S = save | Ctrl-Q = quit | Ctrl-F = find ")

	e.events = make(chan termbox.Event, 20)
	go func() {
		for {
			e.events <- termbox.PollEvent()
		}
	}()
	for e.done != true {
		e.EditorRefreshScreen(true)
		select {
		case ev := <-e.events:
			e.EditorProcessEvent(ev)
			termbox.Flush()
		}
	}

}
