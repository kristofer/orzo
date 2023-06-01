package pieces

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

type Buffer struct {
	T     *Table
	Point int
	Mark  int
}
type PieceSource int

const (
	Content PieceSource = 0
	Add     PieceSource = 1
)

type Table struct {
	Content string
	Add     string
	Mods    []*Piece
}

// Mods ends up being a series of 'runs' which contain the current
// value of the document

type Piece struct {
	Source PieceSource
	Start  int
	Run    int
}

func NewBuffer(c string) *Buffer {
	b := &Buffer{}
	b.T = &Table{Content: c, Add: "", Mods: []*Piece{}}
	b.T.Mods = append(b.T.Mods, NewPiece(Content, 0, len(c)))
	b.Point = 0
	b.Mark = 0
	return b
}

func (buf *Buffer) AddRune(r rune) {
	buf.T.Insert(string(r), buf.Point)
	buf.Point += 1
}

func NewTable(c string) *Table {
	t := &Table{Content: c, Add: "", Mods: []*Piece{}}
	t.Mods = append(t.Mods, NewPiece(Content, 0, len(c)))
	return t
}

func NewPiece(s PieceSource, pt, r int) *Piece {
	return &Piece{Source: s, Start: pt, Run: r}
}

func (t *Table) Size() int {
	i := 0
	for _, n := range t.Mods {
		i += n.Run
	}
	return i
}

func (t *Table) source(ps PieceSource) string {
	if ps == Content {
		return t.Content
	} else {
		return t.Add
	}
}

func (t *Table) RunForMod(index int) string {
	p := t.Mods[index]
	return t.source(p.Source)[p.Start : p.Start+p.Run]
}

func (t *Table) head(p *Piece, idx int) string {
	return t.source(p.Source)[:idx]
}
func (t *Table) tail(p *Piece, idx int) string {
	return t.source(p.Source)[idx:]
}

func (t *Table) AllContents() string {
	//return t.contents(0, t.size())
	// need golang's stringbuilder
	s := ""

	//t.dump()
	for i := 0; i < len(t.Mods); i++ {
		s += t.RunForMod(i)
	}
	//log.Println("ac: ", s)
	return string(s)
}

// need to add observance of start, end
func (t *Table) Contents(start, end int) string {
	// need golang's stringbuilder
	s := ""
	si, ss := t.pieceAt(start)
	sp := t.Mods[si]
	ei, ee := t.pieceAt(end)
	ep := t.Mods[ei]
	startFrag := t.head(sp, ss)
	endFrag := t.tail(ep, ee)
	s += startFrag
	for i := si + 1; i < ei; i++ {
		s += t.RunForMod(i)
	}
	s += endFrag
	return s
}

func (t *Table) IndexOf(idx int) byte {
	e, i := t.pieceAt(idx)
	return t.RunForMod(e)[i]
}

func (t *Table) appendPiece(p *Piece) {
	t.insertPieceAt(len(t.Mods), p)
}

func (t *Table) deletePieceAt(index int) {
	t.Mods = append(t.Mods[:index], t.Mods[index+1:]...)
}

func (t *Table) insertPieceAt(index int, p *Piece) {
	if index >= len(t.Mods) {
		t.Mods = append(t.Mods, p)
		return
	}
	t.Mods = append(t.Mods[:index+1], t.Mods[index:]...)
	t.Mods[index] = p
	//t.dump()
}

func (t *Table) pieceAt(idx int) (int, int) {
	i := idx
	for j, p := range t.Mods {
		if i <= p.Run {
			return j, i
		} else {
			i = i - p.Run
		}
	}
	return 0, 0
}

func (t *Table) Insert(s string, pt int) error {
	e, i := t.pieceAt(pt)
	p := t.Mods[e]
	// Appending characters to the "add file" buffer, and
	if len(t.Add) > 0 && len(t.Add) == (p.Start+p.Run) {
		t.Add += s
		p.Run += len(s)
		//log.Println("adding a rune to existing slice")
		return nil
	}
	np := NewPiece(Add, len(t.Add), len(s))
	t.Add += s
	//np.dump(log.New(os.Stderr, "np ", 0))
	// Updating the entry in piece table (breaking an entry into two or three)
	if i == 0 {
		// insert at p
		t.insertPieceAt(e, np)
		return nil
	}
	if i == p.Run {
		// insert np at p+1
		//t.insertPieceAt(e+1, np)
		t.appendPiece(np)
		return nil
	}
	// else split the piece and make { left, np, right }
	left, right := p.splitAt(i)
	t.deletePieceAt(e)
	t.insertPieceAt(e, left)
	t.insertPieceAt(e+1, right)
	t.insertPieceAt(e+1, np)
	return nil
}

func (t *Table) DeleteRune(idx int) {
	which, i := t.pieceAt(idx)
	p := t.Mods[which]
	if i == 0 || i == p.Run {
		p.trimRune(i)
	}
	// else split into two pieces
	left, right := p.splitAt(i)
	//left.Run -= 1 // trim last rune of left (orig idx above)
	right.Start += 1
	right.Run -= 1
	t.deletePieceAt(which)
	t.insertPieceAt(which, right)
	t.insertPieceAt(which, left)
}

// load a text file from a filename string
func LoadFile(filename string) *Table {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	return NewTable(string(content))
}

func (t *Table) SaveToFile(filename string) error {
	s := t.AllContents()
	f, err := os.Create("test.txt")
	if err != nil {
		return err
	}
	_, err = f.WriteString(s)
	if err != nil {
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Piece) size() int {
	return p.Run
}

func (p *Piece) trimRune(idx int) {
	if idx == 0 {
		p.Start += 1
	}
	if idx == p.Run {
		p.Run -= 1
	}
}

// splits a piece into two
func (p *Piece) splitAt(idx int) (left, right *Piece) {
	if idx == 0 || idx == p.Run {
		return nil, nil
	}
	left = NewPiece(p.Source, p.Start, idx)
	right = NewPiece(p.Source, p.Start+idx, p.Run-idx)
	return
}

func (t *Table) Dump() {
	l := log.New(os.Stderr, "", 0)
	_, filename, line, _ := runtime.Caller(1)
	l.Println(">> Table Dump", "@", filename, line)

	l.Println("Content", &t.Content, t.Content)
	l.Println("Add    ", &t.Add, t.Add)
	for i := 0; i < len(t.Mods); i++ {
		p := t.Mods[i]
		p.dump(l)
		l.Println(p, t.RunForMod(i))
	}
	l.Println(">> End")

}

func (p *Piece) dump(f *log.Logger) {
	_, filename, line, _ := runtime.Caller(1)
	f.Println(p.Source, p.Start, p.Run, "->", filename, line)
}
