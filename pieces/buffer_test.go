package pieces

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	s := "this is a test string."

	if len(s) <= 0 {
		t.Errorf("len() doesn't work!?")
	}
}

func TestPieceSplit1(t *testing.T) {
	s := "this is a test string."
	tt := NewTable(s)
	tt.Dump()
	//    0123456789123456789212
	p := NewPiece(Content, 0, len(s))

	if p.size() != len(s) {
		t.Errorf("Piece isn't the right size")
	}

	at := 1
	left, right := p.splitAt(at)

	if left.size() != at && (p.size()-right.size()) != len(s)-at {
		t.Errorf("p.splitAt(5) doesn't work!?")
		t.Errorf("%v", p)
		t.Errorf("%v, %v", left, right)
	}

}

func TestPieceSplit2(t *testing.T) {
	s := "this is a test string."
	tt := NewTable(s)
	tt.Dump()
	//    0123456789123456789212
	p := NewPiece(Content, 0, len(s))

	if p.size() != len(s) {
		t.Errorf("Piece isn't the right size")
	}

	at := 5
	left, right := p.splitAt(at)

	if left.size() != at && (p.size()-right.size()) != len(s)-at {
		t.Errorf("p.splitAt(1) doesn't work!?")
		t.Errorf("%v", p)
		t.Errorf("%v, %v", left, right)
	}
}

func TestPieceSplit3(t *testing.T) {
	s := "this is a test string."
	tt := NewTable(s)
	tt.Dump()
	//    0123456789123456789212
	p := NewPiece(Content, 0, len(s))

	if p.size() != len(s) {
		t.Errorf("Piece isn't the right size")
	}

	at := 16
	left, right := p.splitAt(at)

	if left.size() != at && (p.size()-right.size()) != len(s)-at {
		t.Errorf("p.splitAt(1) doesn't work!?")
		t.Errorf("%v", p)
		t.Errorf("%v, %v", left, right)
	}
}

func TestLoadFile(t *testing.T) {
	fname := "testtext.txt"

	tt := LoadFile(fname)
	if tt.Size() <= 0 {
		t.Errorf("File Size is bad")
	}

}

func TestSize1(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test. "

	tt := NewTable(s)

	b := tt.Size()
	if b != len(s) {
		t.Errorf("TestSize1 failed.")
	}
}

// func TestAddToBuffer(t *testing.T) {
// 	fname := "testtext.txt"
// 	s := "this is a test. "

// 	tt := LoadFile(fname)
// 	if tt.size() <= 0 {
// 		t.Errorf("File Size is bad")
// 	}

// 	err := tt.add(s, 5)
// 	if err != nil {
// 		t.Errorf("Add failed.")
// 	}

// }

func TestIndex1(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test. "

	tt := NewTable(s)

	b := tt.IndexOf(0)
	if b != 't' {
		t.Errorf("Index failed.")
	}
	//t.Errorf("b is %s", string(b))
	b1 := tt.IndexOf(11)
	if b1 != 'e' {
		t.Errorf("Index failed.")
	}
	//t.Errorf("b is %s", string(b))
	b2 := tt.IndexOf(5)
	if b2 != 'i' {
		t.Errorf("Index failed.")
	}
	//t.Errorf("b is %s", string(b))
}

func TestInsertPiece(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test. "

	tt := NewTable(s)

	foo := "foo"
	p := NewPiece(Content, 0, len(foo))
	tt.insertPieceAt(0, p)
	foo2 := "foo2"
	q := NewPiece(Content, 0, len(foo2))
	tt.insertPieceAt(0, q)

	for i, p := range tt.Mods {
		fmt.Printf("%+v %+v\n", i, p)
	}
	//t.Errorf("tt.Mods %+v", tt.Mods)
}

func TestAppendPiece(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test. "

	tt := NewTable(s)

	foo := "foo"
	p := NewPiece(Content, 0, len(foo))
	tt.appendPiece(p)
	foo2 := "foo2"
	q := NewPiece(Content, 0, len(foo2))
	tt.appendPiece(q)

	for i, p := range tt.Mods {
		fmt.Printf("%+v %+v\n", i, p)
	}
	//t.Errorf("tt.Mods %+v", tt.Mods)
}
func TestAppendInsertPiece(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test. "

	tt := NewTable(s)

	foo := "foo"
	p := NewPiece(Content, 0, len(foo))
	tt.appendPiece(p)
	foo2 := "foo2"
	q := NewPiece(Content, 0, len(foo2))
	tt.insertPieceAt(1, q)

	for i, p := range tt.Mods {
		fmt.Printf("%+v %+v\n", i, p)
	}
	//t.Errorf("tt.Mods %+v", tt.Mods)
}

func TestAllContents(t *testing.T) {
	//fname := "testtext.txt"
	s := "this is a test.  "

	tt := NewTable(s)

	s0 := tt.AllContents()
	if s != s0 {
		t.Errorf("s0 |%+v|", s0)
	}
}
func TestPieceHead(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)
	e, _ := tt.pieceAt(3)
	p := tt.Mods[e]
	s0 := tt.head(p, 3)
	if s0 != "012" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.head(p, 4)
	if s0 != "0123" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.head(p, 1)
	if s0 != "0" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.head(p, 0)
	if s0 != "" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.head(p, 8)
	if s0 != "01234567" {
		t.Errorf("s0 %+v", s0)
	}
}
func TestPieceTail(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)
	e, _ := tt.pieceAt(3)
	p := tt.Mods[e]
	s0 := tt.tail(p, 3)
	if s0 != "3456789" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.tail(p, 4)
	if s0 != "456789" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.tail(p, 1)
	if s0 != "123456789" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.tail(p, 0)
	if s0 != "0123456789" {
		t.Errorf("s0 %+v", s0)
	}
	s0 = tt.tail(p, 8)
	if s0 != "89" {
		t.Errorf("s0 %+v", s0)
	}
}
func TestAdd0(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)
	tt.Dump()

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 0)
	tt.Dump()

	c = tt.AllContents()

	if c != "xxx"+s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", tt.Size())
	tt.Dump()

	c = tt.AllContents()

	if c != "xxx"+s+"xxx" {
		t.Errorf("%+v != %+v", c, s)
	}

}

func TestAdd1(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 0)

	c = tt.AllContents()

	if c != "xxx"+s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", tt.Size())

	c = tt.AllContents()

	if c != "xxx"+s+"xxx" {
		t.Errorf("%+v != %+v", c, s)
	}

}
func TestAdd2(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 1)
	//tt.Dump()
	c = tt.AllContents()

	if c != "0xxx123456789" {
		t.Errorf("%+v != %+v", c, "0xxx123456789")
	}

}

func TestAdd3(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 5)

	c = tt.AllContents()

	if c != "01234xxx56789" {
		t.Errorf("%+v != %+v", c, s)
	}

}

func TestAdd4(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 5)

	c = tt.AllContents()

	if c != "01234xxx56789" {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("abc", 6)

	c = tt.AllContents()

	if c != "01234xabcxx56789" {
		t.Errorf("%+v != %+v", c, s)
	}

}
func TestAddDeleteRune1(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 3)

	c = tt.AllContents()

	if c != "012xxx3456789" {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Dump()
	tt.DeleteRune(8)
	tt.Dump()
	c = tt.AllContents()

	if c != "012xxx346789" {
		t.Errorf("%+v != %+v", c, s)
	}

}

func TestAddDeleteRune4(t *testing.T) {
	//fname := "testtext.txt"
	s := "0123456789"

	tt := NewTable(s)

	c := tt.AllContents()

	if c != s {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.Insert("xxx", 3)

	c = tt.AllContents()

	if c != "012xxx3456789" {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.DeleteRune(8)
	c = tt.AllContents()

	if c != "012xxx346789" {
		t.Errorf("%+v != %+v", c, s)
	}

	tt.DeleteRune(4)
	tt.Dump()
	c = tt.AllContents()

	if c != "012xx346789" {
		t.Errorf("%+v != %+v", c, s)
	}

}
