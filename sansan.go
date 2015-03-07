package sansan

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const heapsize = 30000

type heap []int32

type Machine struct {
	mem heap
	wg  *sync.WaitGroup
}

func newMachine() Machine {
	return Machine{
		mem: make(heap, heapsize),
		wg:  &sync.WaitGroup{},
	}
}

type state struct {
	pnt    int16
	atomic bool
}

func (m Machine) run(p Program, s *state) {
	for i := 0; i < len(p); i++ {
		switch ins := p[i]; ins {
		case '>':
			s.pnt++
			s.pnt = (s.pnt%heapsize + heapsize) % heapsize
		case '<':
			s.pnt--
			s.pnt = (s.pnt%heapsize + heapsize) % heapsize
		case '+':
			if s.atomic {
				atomic.AddInt32(&m.mem[s.pnt], 1)
			} else {
				m.mem[s.pnt]++
			}
		case '-':
			if s.atomic {
				atomic.AddInt32(&m.mem[s.pnt], -1)
			} else {
				m.mem[s.pnt]--
			}
		case '[':
			end := i + findClosing(p[i:])

			if (s.atomic && atomic.LoadInt32(&m.mem[s.pnt]) != 0) || (!s.atomic && m.mem[s.pnt] != 0) {
				m.run(p[i+1:end+1], s) // enter loop
			}

			i = end // goto end
		case ']':
			if (s.atomic && atomic.LoadInt32(&m.mem[s.pnt]) == 0) || (!s.atomic && m.mem[s.pnt] == 0) {
				return
			}
			i = -1

		case '{':
			end := i + findClosing(p[i:])

			m.wg.Add(1)
			go m.runThread(p[i+1:end+1], *s)

			i = end // continue parrent thread
		case '}':
			return // kill thread
		case '!':
			// toggle atomic operations on current thread
			s.atomic = s.atomic != true

		case '.':
			var v int32
			if s.atomic {
				v = atomic.LoadInt32(&m.mem[s.pnt])
			} else {
				v = m.mem[s.pnt]
			}

			fmt.Printf("%c", v)
		case ',':
			var n int32
			if _, err := fmt.Scanf("%d\n", &n); err != nil {
				panic(err)
			}

			if s.atomic {
				atomic.SwapInt32(&m.mem[s.pnt], n)
			} else {
				m.mem[s.pnt] = n
			}
		}
	}
}

// runThread runs the given program with a local copy of the
// heap pointer and decrements waitgroup when finished.
func (m Machine) runThread(p Program, s state) {
	defer m.wg.Done()
	m.run(p, &s)
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
