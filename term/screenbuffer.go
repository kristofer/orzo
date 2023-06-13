package term

import (
	"fmt"
	"log"
	"unicode/utf8"
)

type Screen struct {
	data []rune
	Rows int
	Cols int
}

func NewScreen(c, r int) *Screen {
	scr := &Screen{}
	scr.Rows = r
	scr.Cols = c
	// make([]rune, C*R)
	scr.data = make([]rune, c*r)
	scr.Blank()
	log.Println("created ScreenBuf size", len(scr.data))
	return scr
}

func (scr *Screen) Blank() {
	for i, _ := range scr.data {
		scr.data[i] = ' '
	}
}
func (scr *Screen) Fill(ru rune) {
	for i, _ := range scr.data {
		scr.data[i] = ru
	}
}

func (scr *Screen) checkRange(c, r int) bool {
	if c > scr.Cols {
		return false
	}
	if r > scr.Rows {
		return false
	}
	return true
}

func (scr *Screen) rowOrder(c, r int) int {
	return (r * scr.Cols) + c
}

func (scr *Screen) Set(c, r int, ch rune) {
	// board[c*C + r] = "abc" // like board[i][j] = "abc"
	if scr.checkRange(c, r) {
		//	scr.data[(r*scr.Rows)+c] = ch
		scr.data[scr.rowOrder(c, r)] = ch
	}
}

func (scr *Screen) Get(c, r int) rune {
	if scr.checkRange(c, r) {
		// return scr.data[(r*scr.Rows)+c]
		return scr.data[scr.rowOrder(c, r)]
	}
	return ' '
}

func (scr *Screen) String() string {
	s := "  |"
	for c := 0; c < scr.Cols; c++ {
		if c%10 != 0 {
			s += "-"
		} else {
			s += "+"
		}
	}
	s += "|\n"
	for r := 0; r < scr.Rows; r++ {
		s += fmt.Sprintf("%2d|", r)
		for c := 0; c < scr.Cols; c++ {
			s += string(scr.Get(c, r))
		}
		s += "|\n"
	}
	return s
}
func (scr *Screen) GetBytes() []byte {
	buf := make([]byte, len(scr.data)*utf8.UTFMax)

	count := 0
	for r := 0; r < scr.Rows; r++ {
		for c := 0; c < scr.Cols; c++ {
			count += utf8.EncodeRune(buf[count:], scr.data[scr.rowOrder(c, r)])
		}

	}
	buf = buf[:count]

	return buf
}
