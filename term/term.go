package term

import (
	"bufio"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/gorilla/websocket"
)

type Terminal struct {
	Kind   TermType
	Input  *bufio.Reader
	Output *bufio.Writer
	Origin *syscall.Termios
	// sigwinch       = make(chan os.Signal, 1)
	// sigio          = make(chan os.Signal, 1)
	Quit     chan int
	Contents strings.Builder
	ScrBuf   *Screen
	Conn     *websocket.Conn
	CurCol   int
	CurRow   int
	NumCol   int
	NumRow   int
}

type (
	TermType int
)

const (
	Web TermType = iota
	Pty
)

// Init()
func NewTerminal(height int, width int, kind TermType) (t *Terminal) {
	t = &Terminal{}
	t.NumCol = width
	t.NumRow = height
	t.Kind = kind

	if kind == Pty {
		t.Input = bufio.NewReader(os.Stdin)
		t.Output = bufio.NewWriter(os.Stdout)
		stdin := os.Stdin.Fd()
		termios := GetTermios(stdin)

		t.Origin = termios

		SetRaw(termios)
		SetTermios(stdin, termios)
		t.Output.Write([]byte(CSIStart()))

	}

	if kind == Web {
		t.ScrBuf = NewScreen(80, 24) // termsize cols, rows
	}

	return t
}

func (t *Terminal) IsPty() bool { return t.Kind == Pty }
func (t *Terminal) IsWeb() bool { return t.Kind == Web }

func (t *Terminal) Close() {
	if t.IsPty() {
		t.Output.Write([]byte(CSIStart()))
		SetTermios(os.Stdin.Fd(), t.Origin)
	}
}

// termbox.Size()
func (t *Terminal) Size() (int, int) {
	return t.NumCol, t.NumRow
}

func (t *Terminal) SetOutputMode(m OutputMode) {}

// termbox.InputAlt | termbox.InputEsc | termbox.InputMouse
func (t *Terminal) SetInputMode(m InputMode) {}
func (t *Terminal) Clear(fg, bg Attribute)   {}
func (t *Terminal) Flush()                   {}

func (t *Terminal) SetCursor(r, c int) {
	t.CurCol = c
	t.CurRow = r
	if t.IsWeb() {
		// t.Write([]byte(CUP(t.CurCol, t.CurRow)))
		t.Write([]byte(CUP(t.CurCol, t.CurRow)))
	}
}
func (t *Terminal) SetCell(c, r int, ch rune, fgcolor, bgcolor Attribute) {
	t.ScrBuf.Set(c, r, ch)
}

func (t *Terminal) Write(b []byte) {
	if t.Kind == Pty {
		t.Output.Write(b)
		t.Output.Flush()
	}
	if t.Kind == Web {
		msgType := 1
		msg := b
		if err := t.Conn.WriteMessage(msgType, msg); err != nil {
			log.Println("unable to write message to frontend")
			return
		}
	}
}

// termios

func GetTermios(fd uintptr) *syscall.Termios {
	var t syscall.Termios
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		syscall.TIOCGETA, //TCGETS,
		uintptr(unsafe.Pointer(&t)),
		0, 0, 0)

	if err != 0 {
		panic("err")
	}

	return &t
}

func SetTermios(fd uintptr, term *syscall.Termios) {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdin.Fd(),
		syscall.TIOCSETA, //TCSETS,
		uintptr(unsafe.Pointer(term)),
		0, 0, 0)
	if err != 0 {
		panic("err")
	}
}

func SetRaw(term *syscall.Termios) {
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	term.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	// newState.Oflag &^= syscall.OPOST
	term.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	term.Cflag &^= syscall.CSIZE | syscall.PARENB
	term.Cflag |= syscall.CS8

	term.Cc[syscall.VMIN] = 1
	term.Cc[syscall.VTIME] = 0
}

// constants
