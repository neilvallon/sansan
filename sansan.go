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

	h.run(p.Clean())
}

// Clean returns a new program with invalid instructions and
// whitespace removed.
func (p Program) Clean() Program {
	// allocate some space for cleaned program
	np := make(Program, 0, len(p)/2)

	for _, c := range p {
		switch c {
		case '>', '<', '+', '-', '[', ']', '{', '}', '.', ',', '!':
			np = append(np, c)
		}
	}

	return np
}

const heapsize = 30000

type heap struct {
	mem    []int32
	pnt    int
	atomic bool
	wg     *sync.WaitGroup
}

func newHeap() *heap {
	return &heap{
		mem: make([]int32, heapsize),
		wg:  &sync.WaitGroup{},
	}
}

func (h *heap) run(p Program) {
	for i := 0; i < len(p); i++ {
		switch ins := p[i]; ins {
		case '>':
			h.pnt++
			h.pnt = (h.pnt%heapsize + heapsize) % heapsize
		case '<':
			h.pnt--
			h.pnt = (h.pnt%heapsize + heapsize) % heapsize
		case '+':
			if h.atomic {
				atomic.AddInt32(&h.mem[h.pnt], 1)
			} else {
				h.mem[h.pnt]++
			}
		case '-':
			if h.atomic {
				atomic.AddInt32(&h.mem[h.pnt], -1)
			} else {
				h.mem[h.pnt]--
			}
		case '[':
			end := i + findClosing(p[i:])

			if (h.atomic && atomic.LoadInt32(&h.mem[h.pnt]) != 0) || h.mem[h.pnt] != 0 {
				h.run(p[i+1 : end+1]) // enter loop
			}

			i = end // goto end
		case ']':
			if (h.atomic && atomic.LoadInt32(&h.mem[h.pnt]) == 0) || h.mem[h.pnt] == 0 {
				return
			}
			i = -1

		case '{':
			end := i + findClosing(p[i:])

			h.wg.Add(1)
			go h.runThread(p[i+1 : end+1])

			i = end // continue parrent thread
		case '}':
			return // kill thread
		case '!':
			// toggle atomic operations on current thread
			h.atomic = h.atomic != true

		case '.':
			var v int32
			if h.atomic {
				v = atomic.LoadInt32(&h.mem[h.pnt])
			} else {
				v = h.mem[h.pnt]
			}

			fmt.Printf("%c", v)
		case ',':
			var n int32
			if _, err := fmt.Scanf("%d\n", &n); err != nil {
				panic(err)
			}

			if h.atomic {
				atomic.SwapInt32(&h.mem[h.pnt], n)
			} else {
				h.mem[h.pnt] = n
			}
		}
	}
}

// runThread runs the given program with a local copy of the
// heap pointer and decrements waitgroup when finished.
func (h heap) runThread(p Program) {
	defer h.wg.Done()
	h.run(p)
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
