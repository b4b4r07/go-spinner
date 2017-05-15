package spinner

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/sethgrid/curse"
)

var (
	Box1    = `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`
	Box2    = `⠋⠙⠚⠞⠖⠦⠴⠲⠳⠓`
	Box3    = `⠄⠆⠇⠋⠙⠸⠰⠠⠰⠸⠙⠋⠇⠆`
	Box4    = `⠋⠙⠚⠒⠂⠂⠒⠲⠴⠦⠖⠒⠐⠐⠒⠓⠋`
	Box5    = `⠁⠉⠙⠚⠒⠂⠂⠒⠲⠴⠤⠄⠄⠤⠴⠲⠒⠂⠂⠒⠚⠙⠉⠁`
	Box6    = `⠈⠉⠋⠓⠒⠐⠐⠒⠖⠦⠤⠠⠠⠤⠦⠖⠒⠐⠐⠒⠓⠋⠉⠈`
	Box7    = `⠁⠁⠉⠙⠚⠒⠂⠂⠒⠲⠴⠤⠄⠄⠤⠠⠠⠤⠦⠖⠒⠐⠐⠒⠓⠋⠉⠈⠈`
	Default = Box1
)

type Spinner struct {
	Line         int
	progressChan chan string
	frames       []rune
	length       int
	pos          int
	done         bool
}

func (s *Spinner) next() string {
	r := s.frames[s.pos/2%s.length]
	s.pos++
	return string(r)
}

type ProgressFunc func(title string)

type Screen struct {
	Spinners []*Spinner

	screenLines             int
	startingLine            int
	totalNewlines           int
	historicNewlinesCounter int

	history map[int]string
	sync.Mutex
}

func New() (*Screen, error) {
	_, lines, _ := curse.GetScreenDimensions()
	_, line, _ := curse.GetCursorPosition()

	history := make(map[int]string)

	b := &Screen{screenLines: lines, startingLine: line, history: history}
	return b, nil
}

func (b *Screen) Listen() {
	for len(b.Spinners) == 0 {
		time.Sleep(time.Millisecond * 100)
	}
	cases := make([]reflect.SelectCase, len(b.Spinners))
	for i, bar := range b.Spinners {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(bar.progressChan)}
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			b.Spinners[chosen].Done()
			cases[chosen].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}
		b.Spinners[chosen].Update(value.String())
	}
	b.Println()
}

func (b *Screen) MakeBar() ProgressFunc {
	ch := make(chan string)
	bar := &Spinner{
		progressChan: ch,
		frames:       []rune(Default),
	}

	bar.length = len(bar.frames)
	b.Spinners = append(b.Spinners, bar)
	bar.Line = b.startingLine + b.totalNewlines
	b.history[bar.Line] = ""
	// bar.Update("?")
	b.Println()

	return func(title string) { bar.progressChan <- title }
}

func (p *Spinner) Done() {
	p.done = true
}

func (p *Spinner) Update(title string) {
	c, _ := curse.New()
	c.Move(1, p.Line)
	c.EraseCurrentLine()
	if p.done {
		fmt.Printf("\r %s %s", "x", title)
	} else {
		fmt.Printf("\r %s %s", p.next(), title)
	}
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}

func (b *Screen) addedNewlines(count int) {
	b.totalNewlines += count
	b.historicNewlinesCounter += count

	if b.startingLine+b.totalNewlines > b.screenLines {
		b.totalNewlines -= count
		for _, bar := range b.Spinners {
			bar.Line -= count
		}
		b.redrawAll(count)
	}
}

func (b *Screen) redrawAll(moveUp int) {
	c, _ := curse.New()

	newHistory := make(map[int]string)
	for line, printed := range b.history {
		newHistory[line+moveUp] = printed
		c.Move(1, line)
		c.EraseCurrentLine()
		c.Move(1, line+moveUp)
		c.EraseCurrentLine()
		fmt.Print(printed)
	}
	b.history = newHistory
	c.Move(c.StartingPosition.X, c.StartingPosition.Y)
}

func (b *Screen) Print(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Print(a...)
}

func (b *Screen) Printf(format string, a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := strings.Count(format, "\n")
	newlines += countAllNewlines(a...)
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprintf(format, a...)
	return fmt.Printf(format, a...)
}

func (b *Screen) Println(a ...interface{}) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	newlines := countAllNewlines(a...) + 1
	b.addedNewlines(newlines)
	thisLine := b.startingLine + b.totalNewlines
	b.history[thisLine] = fmt.Sprint(a...)
	return fmt.Println(a...)
}

func countAllNewlines(interfaces ...interface{}) int {
	count := 0
	for _, iface := range interfaces {
		switch s := iface.(type) {
		case string:
			count += strings.Count(s, "\n")
		}
	}
	return count
}
