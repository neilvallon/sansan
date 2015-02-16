package sansan

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Program []byte

func (p Program) Run() {
	h := newHeap()
	defer h.wg.Wait()

	h.wg.Add(1)
	h.run(p)
}

const heapsize = 30000

type heap struct {
	mem []int32
	pnt int
	wg  *sync.WaitGroup
}

func newHeap() *heap {
	return &heap{
		mem: make([]int32, heapsize),
		wg:  &sync.WaitGroup{},
	}
}

func (h *heap) run(p Program) {
	defer h.wg.Done()

	for i := 0; i < len(p); i++ {
		switch ins := p[i]; ins {
		case '>':
			h.pnt++
			h.pnt = (h.pnt%heapsize + heapsize) % heapsize
		case '<':
			h.pnt--
			h.pnt = (h.pnt%heapsize + heapsize) % heapsize
		case '+':
			atomic.AddInt32(&h.mem[h.pnt], 1)
		case '-':
			atomic.AddInt32(&h.mem[h.pnt], -1)
		case '[':
			end := i + findClosing(p[i:])

			if atomic.LoadInt32(&h.mem[h.pnt]) != 0 {
				// enter loop
				h.wg.Add(1) // TODO: remove this on loops
				h.run(p[i+1 : end+1])
			}

			i = end // goto end
		case ']':
			if atomic.LoadInt32(&h.mem[h.pnt]) == 0 {
				return
			}
			i = -1

		case '{':
			end := i + findClosing(p[i:])

			newH := *h // copy heap mem and pointer
			h.wg.Add(1)
			go newH.run(p[i+1 : end+1])

			i = end // continue parrent thread
		case '}':
			return // kill thread

		case '.':
			fmt.Printf("%c", atomic.LoadInt32(&h.mem[h.pnt]))
		case ',':
			var n int32
			if _, err := fmt.Scanf("%d\n", &n); err != nil {
				panic(err)
			}
			atomic.SwapInt32(&h.mem[h.pnt], n)
		case ' ', '\t', '\n':
		default:
		}
	}
}

func findClosing(prog []byte) int {
	braces := 0
	for i := 0; i < len(prog); i++ {
		switch prog[i] {
		case '[', '{':
			braces++
		case ']', '}':
			braces--
			if braces < 0 {
				panic("invalid program: unbalanced braces")
			}
			if braces == 0 {
				return i
			}
		}
	}
	panic("invalid program: could not find closing ']'")
}
