package orzo

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Functions which manipulate Buffers.

func (e *Editor) indexOfBuffer(element *Buffer) int {
	for k, v := range e.Buffers {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func (e *Editor) anyBufferDirty() bool {
	for _, v := range e.Buffers {
		if v.dirty {
			return true
		}
	}
	return false
}

func (e *Editor) indexOfBufferNamed(name string) (int, error) {
	for k, v := range e.Buffers {
		if strings.Compare(v.filename, name) == 0 {
			return k, nil
		}
	}
	return -1, errors.New("not found") //not found.
}

func (e *Editor) removeBuffer(element *Buffer) {
	i := e.indexOfBuffer(element)
	e.Buffers = append(e.Buffers[:i], e.Buffers[i+1:]...)
}

func (e *Editor) nextBuffer() {
	idx := e.indexOfBuffer(e.cb)
	if (idx + 1) >= len(e.Buffers) {
		e.cb = e.Buffers[0]
		return
	}
	e.cb = e.Buffers[idx+1]
}

func (e *Editor) addNewBuffer() {
	nb := &Buffer{}
	e.Buffers = append(e.Buffers, nb)
	e.cb = nb
	e.cb.filename = "*scratch*"
}

func (e *Editor) listBuffers() {
	found, err := e.indexOfBufferNamed(orzoListBuffers)
	if err == nil {
		e.cb = e.Buffers[found]
	} else {
		e.addNewBuffer()
		e.cb.filename = orzoListBuffers
	}
	Bufferlist := " ** Current Buffers\n\n"
	for k, v := range e.Buffers {
		Bufferlist = Bufferlist + fmt.Sprintf("%d - %s - %d lines \n", k, v.filename, v.numrows)
	}
	e.cb.rows = nil
	e.cb.numrows = 0
	e.cb.readonly = false
	scanner := bufio.NewScanner(strings.NewReader(Bufferlist))
	for scanner.Scan() {
		line := scanner.Text()
		e.EditorInsertRow(e.cb.numrows, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	e.cb.dirty = false
	e.cb.readonly = true
}

func (e *Editor) killBuffer() {
	if e.cb.dirty {
		e.EditorSetStatusMessage("Buffer %s is modified.", e.cb.filename)
		return
	}
	e.removeBuffer(e.cb)
	if len(e.Buffers) > 0 {
		e.cb = e.Buffers[0]
	}
	if len(e.Buffers) == 0 {
		e.done = true
	}
}
