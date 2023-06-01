package orzo

func (e *Editor) setMark() {
	e.cb.mark.r, e.cb.mark.c = e.cb.point.r, e.cb.point.c
	e.cb.mark.ro, e.cb.mark.co = e.cb.point.ro, e.cb.point.co
	e.cb.markSet = true
	e.EditorSetStatusMessage("Mark Set (%d,%d)", e.cb.mark.r+e.cb.mark.ro, e.cb.mark.c+e.cb.mark.co)
}
func (e *Editor) noMark() bool {
	e.cb.markSet = false
	return e.cb.markSet
}
func swapCursorsMaybe(mr, mc, cr, cc int) (sr, sc, er, ec int, f bool) {
	if mr == cr {
		if mc > cc {
			return cr, cc, mr, mc, true
		}
		return mr, mc, cr, cc, false
	}
	if mr > cr {
		return cr, cc, mr, mc, true
	}
	return mr, mc, cr, cc, false
}

func (e *Editor) cutCopy(del bool) {
	if e.noMark() == true {
		return
	}
	e.pasteBuffer = ""
	sr, sc, er, ec, reverse := swapCursorsMaybe(e.cb.mark.r+e.cb.mark.ro, e.cb.mark.c+e.cb.mark.co, e.cb.point.ro+e.cb.point.r, e.cb.point.co+e.cb.point.c)
	for i := sr; i <= er; i++ {
		if sr == er {
			l := e.cb.rows[i]
			for j := sc; j < ec; j++ {
				e.pasteBuffer = e.pasteBuffer + string(l.runes[j])
			}
		} else if i == sr {
			l := e.cb.rows[i]
			for j := sc; j < l.size; j++ {
				e.pasteBuffer = e.pasteBuffer + string(l.runes[j])
			}
			e.pasteBuffer = e.pasteBuffer + "\n"
		} else if i == er {
			if i < e.cb.numrows {
				l := e.cb.rows[i]
				for j := 0; j < ec; j++ {
					e.pasteBuffer = e.pasteBuffer + string(l.runes[j])
				}
			}
		} else {
			l := e.cb.rows[i]
			for j := 0; j < l.size; j++ {
				e.pasteBuffer = e.pasteBuffer + string(l.runes[j])
			}
			e.pasteBuffer = e.pasteBuffer + "\n"
		}
	}
	// if del == true, remove all the chars
	if del == true {
		c2Remove := len(e.pasteBuffer)
		if reverse { // move cursor to delete the right runes
			e.cb.point.r, e.cb.point.c = e.cb.mark.r, e.cb.mark.c
			e.cb.point.ro, e.cb.point.co = e.cb.mark.ro, e.cb.mark.co
		}
		for k := 0; k < c2Remove; k++ {
			e.EditorDelChar()
		}
	}
	e.noMark()
}

func (e *Editor) paste() {
	for _, rch := range e.pasteBuffer {
		if rch == '\n' {
			e.EditorInsertNewline()
		} else {
			e.EditorInsertChar(rch)
		}
	}
}
