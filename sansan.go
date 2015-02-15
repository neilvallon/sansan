package sansan

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Program []byte

func (p Program) Run() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	newHeap().run(p, &wg)
}

const heapsize = 30000

type heap struct {
	mem []int32
	pnt int
}

func newHeap() heap {
	return heap{mem: make([]int32, heapsize)}
}

func (h heap) run(p Program, wg *sync.WaitGroup) int {
	defer wg.Done()

	var childWG sync.WaitGroup
	defer childWG.Wait()

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
				childWG.Add(1) // TODO: remove this on loops
				h.pnt = h.run(p[i+1:], &childWG)
			}
			i = end // goto end
		case ']':
			if atomic.LoadInt32(&h.mem[h.pnt]) == 0 {
				return h.pnt
			}
			i = -1

		case '{':
			end := i + findClosing(p[i:])
			childWG.Add(1)
			go h.run(p[i+1:], &childWG)

			i = end // continue parrent thread
		case '}':
			return -1 // kill thread

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
	return -1
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
