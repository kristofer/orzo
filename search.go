package orzo

import (
	termbox "github.com/nsf/termbox-go"
)

/* =============================== Find mode ================================ */

func (e *Editor) runeAt(r, c int) rune {
	return e.cb.rows[r].runes[c]
}

func (e *Editor) searchForward(startr, startc int, stext string) (fr, rc int) {
	if len(stext) == 0 {
		return -1, -1
	}
	s := []rune(stext)
	ss := 0

	c := startc
	r := startr
	for r < e.cb.numrows {
		for c < e.cb.rows[r].size {
			rch := e.runeAt(r, c)
			if ss < len(s) && s[ss] == rch {
				ss++
			} else {
				if ss == len(s) {
					//log.Println("found", r, c)
					return r, c
				}
				if s[ss] != rch {
					ss = 0
				}
			}
			c++
		}
		c = 0
		r++
		ss = 0
	}
	return -1, -1
}

func reverse(s string) []rune {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}

func (e *Editor) searchBackwards(startr, startc int, stext string) (fr, fc int) {
	if len(stext) == 0 {
		return -1, -1
	}
	s := reverse(stext)
	ss := 0

	c := startc
	r := startr
	//log.Printf("before lop r %d c %d\n", r, c)
	for r > -1 {
		//rintf("r is %d\n", r)
		for c >= 0 && c < e.cb.rows[r].size {
			rch := e.runeAt(r, c)
			if ss < len(s) && s[ss] == rch {
				fr = r
				fc = c
				ss++
			} else if s[ss] != rch {
				ss = 0
			}

			if ss == len(s) {
				//log.Println("found", r, c)
				return fr, fc
			}
			c--
		}
		if r > 0 {
			r--
			c = e.cb.rows[r].size - 1
		} else {
			break
		}
		ss = 0
	}
	return -1, -1
}

func (e *Editor) EditorFind() {
	query := ""
	startrow := e.cb.point.ro + e.cb.point.r
	startcol := e.cb.point.co + e.cb.point.c
	lastLineMatch := startrow /* Last line where a match was found. -1 for none. */
	lastColMatch := startcol  // last column where a match was found
	findDirection := 1        /* if 1 search next, if -1 search prev. */

	/* Save the cursor position in order to restore it later. */
	savedCx, savedCy := e.cb.point.c, e.cb.point.r
	savedColoff, savedRowoff := e.cb.point.co, e.cb.point.ro
	mesg := "Search: %s (Use Esc/Arrows/Enter)"
	for {
		e.EditorSetStatusMessage(mesg, query)
		e.EditorRefreshScreen(false)
		ev := <-e.events
		if ev.Ch != 0 {
			ch := ev.Ch
			query = query + string(ch)
			lastLineMatch, lastColMatch = startrow, startcol
			findDirection = 1
		}
		if ev.Ch == 0 {
			switch ev.Key {
			case termbox.KeyTab:
				query = query + string('\t')
				lastLineMatch, lastColMatch = startrow, startcol
				findDirection = 1
			case termbox.KeySpace:
				query = query + string(' ')
				lastLineMatch, lastColMatch = startrow, startcol
				findDirection = 1
			case termbox.KeyEnter, termbox.KeyCtrlR:
				e.EditorSetStatusMessage("")
				return
			case termbox.KeyCtrlC:
				e.EditorSetStatusMessage("killed.")
				return
			case termbox.KeyBackspace2, termbox.KeyBackspace:
				if len(query) > 0 {
					query = query[:len(query)-1]
				} else {
					query = ""
				}
				lastLineMatch, lastColMatch = startrow, startcol
				findDirection = -1
			case termbox.KeyCtrlG, termbox.KeyEsc:
				e.cb.point.c, e.cb.point.r = savedCx, savedCy
				e.cb.point.co, e.cb.point.ro = savedColoff, savedRowoff
				e.EditorSetStatusMessage("")
				return
			case termbox.KeyArrowDown, termbox.KeyArrowRight:
				findDirection = 1
			case termbox.KeyArrowLeft, termbox.KeyArrowUp:
				findDirection = -1
			default:
				e.EditorSetStatusMessage(mesg, query)
				e.EditorRefreshScreen(false)
				termbox.SetCursor(len(mesg)+len(query), e.screenrows+1)
				findDirection = 0
			}
		}
		if findDirection != 0 {
			currentLine, matchOffset := -1, -1
			if findDirection == 1 {
				currentLine, matchOffset = e.searchForward(lastLineMatch, lastColMatch, query)
			}
			if findDirection == -1 {
				currentLine, matchOffset = e.searchBackwards(lastLineMatch, lastColMatch, query)
			}
			//findDirection = 0

			if currentLine != -1 {
				lastLineMatch = currentLine
				lastColMatch = matchOffset
				e.cb.point.r = 0
				e.cb.point.c = matchOffset
				e.cb.point.ro = currentLine
				e.cb.point.co = 0
				/* Scroll horizontally as needed. */
				if e.cb.point.c > e.screencols {
					diff := e.cb.point.c - e.screencols
					e.cb.point.c -= diff
					e.cb.point.co += diff
				}
			}
		}
	}
}
