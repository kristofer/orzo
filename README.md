# orzo


orzo is a small terminal/screen text Editor in less than 2K lines of Go code. It uses Go's rune (unicode) machinery. Has multiple Buffers. But still single window. hope this makes it into the vault  

### Usage: orzo `<filename>`

To build: clone repo into gopath;
 ```
  $ cd cmd
  $ go get -d ./...
  $ go build -o orzo
 ```

 Copy to someplace in your PATH

be sure you have https://godoc.org/github.com/nsf/termbox-go
so, you may need to `go get github.com/nsf/termbox-go`

### Key Commands:

#### Movement
* CTRL-Y: HELP
* Use ArrowKeys to move Cursor around.
* Home, End, and PageUp & PageDown should work
* CTRL-A: Move to beginning of current line
* CTRL-E: Move to end of current line
* on mac keyboards:
  * FN+ArrowUp: PageUp (screen full)
  * FN+ArrowDown: PageDown (screen full)

#### File/Buffer 
* CTRL-S: Save the file
* CTRL-Q: Quit the Editor
* CTRL-F: Find string in file 
	(ESC to exit search mode, arrows to navigate to next/prev find)
* CTRL-N: Next Buffer
* CTRL-B: List all the Buffers
* CTRL-O: (control oh) Open File into new Buffer
* CTRL-W: kill Buffer

#### Cut/Copy/Paste & Deletion
* CTRL-Space: Set Mark
* CTRL-X: Cut region from Mark to Cursor into paste Buffer
* CTRL-C: Copy region from Mark to Cursor into paste Buffer
* CTRL-V: Paste copied/cut region into file at Cursor
_Once you've set the Mark, as you move the cursor,
you should be getting underlined text showing the current
selection/region._
* Delete: to delete a rune backward
* CTRL-K: killtoEndOfLine (once) removeLine (twice)


Setting the cursor with a mouse click should work. (and so,
it should work to set the selection. but hey, you MUST SetMark
for a selection to start... sorry, it's not a real mouse based Editor.)
    
### Implementation Notes
orzo was based on Kilo, a project by Salvatore Sanfilippo <antirez at gmail dot com> at  https://github.com/antirez/kilo.

It's a very simple Editor, with kinda-"Mac-Emacs"-like key bindings. It uses `go get github.com/nsf/termbox-go" for simple termio junk.

The central data structure is an array of lines (type erow struct). Each line in the file has a struct, which contains an array of rune. (If you're not familiar with Go's _runes_, they are Go's unicode code points (or characters))

Multiple Buffers, but no window splits. Two _mini modes_,  one for the search modal operations, and
one for opening files.

Notice the goroutine attached to events coming from termbox-go, that is pretty cool. Yet another real reason that Go routines are handy.

_orzo was written in Go by K Younger and is released
under the BSD 2 clause license._


## Orzo interface with WindowServer

- a Pty interface, where vt100/ansi terminal codes rule
- a Web interface where browser/JS events flow, but only on input.

- a light go based window server