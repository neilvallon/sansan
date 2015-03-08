package sansan

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const heapsize = 30000

type heap []int32

type Machine struct {
	mem heap
	wg  *sync.WaitGroup

	stdin  io.Reader
	stdout io.Writer
}

func NewMachine() Machine {
	return Machine{
		mem: make(heap, heapsize),
		wg:  &sync.WaitGroup{},

		stdin:  os.Stdin,
		stdout: os.Stdout,
	}
}

func (m Machine) Run(p program) {
	defer m.wg.Wait()
	m.run(p, new(state))
}

type state struct {
	pnt    int16
	atomic bool
}

func (m Machine) run(p program, s *state) {
	for i := 0; i < len(p); i++ {
		switch ins := p[i]; ins.Action {
		case Move:
			s.pnt += ins.Val
			s.pnt = (s.pnt%heapsize + heapsize) % heapsize
		case Modify:
			if s.atomic {
				atomic.AddInt32(&m.mem[s.pnt], int32(ins.Val))
			} else {
				m.mem[s.pnt] += int32(ins.Val)
			}
		case LStart:
			end := i + int(ins.Val)

			if (s.atomic && atomic.LoadInt32(&m.mem[s.pnt]) != 0) || (!s.atomic && m.mem[s.pnt] != 0) {
				m.run(p[i+1:end+1], s) // enter loop
			}

			i = end // goto end
		case LEnd:
			if (s.atomic && atomic.LoadInt32(&m.mem[s.pnt]) == 0) || (!s.atomic && m.mem[s.pnt] == 0) {
				return
			}
			i = -1

		case TStart:
			end := i + int(ins.Val)

			m.wg.Add(1)
			go m.runThread(p[i+1:end+1], *s)

			i = end // continue parrent thread
		case TEnd:
			return // kill thread
		case Toggle:
			// toggle atomic operations on current thread
			s.atomic = s.atomic != true

		case Print:
			var v int32
			if s.atomic {
				v = atomic.LoadInt32(&m.mem[s.pnt])
			} else {
				v = m.mem[s.pnt]
			}

			fmt.Fprintf(m.stdout, "%c", v)
		case Read:
			var n int32
			if _, err := fmt.Fscanf(m.stdin, "%d\n", &n); err != nil {
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
func (m Machine) runThread(p program, s state) {
	defer m.wg.Done()
	m.run(p, &s)
}

func (m *Machine) SetInput(r io.Reader) {
	m.stdin = r
}

func (m *Machine) SetOutput(w io.Writer) {
	m.stdout = w
}
