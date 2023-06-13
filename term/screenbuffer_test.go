package term

import "testing"

func TestBufCreate(t *testing.T) {
	c, r := 6, 4

	scr := NewScreen(c, r)
	scr.Fill('+')
	t.Errorf("\n***\n%s***", scr.String())

	sr := scr.GetBytes()

	if len(sr) == c*r {
		t.Errorf("result |%s|(%d) %d, %d", sr, len(sr), c, r)
	}
}
func TestBufSet1(t *testing.T) {
	c, r := 6, 4

	scr := NewScreen(c, r)
	scr.Fill('-')

	for i := 0; i < c; i++ {
		scr.Set(i, 0, 'X')
	}
	t.Errorf("\n***\n%s***", scr.String())

	sr := scr.GetBytes()

	if len(sr) == c*r {
		t.Errorf("result |%s|(%d) %d, %d", sr, len(sr), c, r)
	}
}
func TestBufSet2(t *testing.T) {
	c, r := 6, 4

	scr := NewScreen(c, r)
	scr.Fill('-')
	for i := 0; i < r; i++ {
		scr.Set(0, i, 'X')
	}
	t.Errorf("\n***\n%s***", scr.String())
	sr := scr.GetBytes()

	if len(sr) == c*r {
		t.Errorf("result |%s|(%d) %d, %d", sr, len(sr), c, r)
	}
}
func TestBufSet3(t *testing.T) {
	c, r := 6, 4

	scr := NewScreen(c, r)
	scr.Fill('-')
	for i := 0; i < c; i++ {
		for j := 0; j < r; j++ {
			if i == j {
				scr.Set(i, j, 'X')
			}
		}
	}
	t.Errorf("\n***\n%s***", scr.String())
	sr := scr.GetBytes()

	if len(sr) == c*r {
		t.Errorf("result |%s|(%d) %d, %d", sr, len(sr), c, r)
	}
}
