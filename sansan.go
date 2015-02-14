package sansan

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const heapsize = 30000

type program struct {
	code []byte
	heap []int32
}

func NewProg(c []byte) *program {
	return &program{
		c,
		make([]int32, heapsize),
	}
}

func (p program) Run() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	p.run(0, &wg)
}

func (p program) run(heapPos int, wg *sync.WaitGroup) int {
	defer wg.Done()

	var childWG sync.WaitGroup
	defer childWG.Wait()

	for i := 0; i < len(p.code); i++ {
		switch ins := p.code[i]; ins {
		case '>':
			heapPos++
			heapPos = (heapPos%heapsize + heapsize) % heapsize
		case '<':
			heapPos--
			heapPos = (heapPos%heapsize + heapsize) % heapsize
		case '+':
			atomic.AddInt32(&p.heap[heapPos], 1)
		case '-':
			atomic.AddInt32(&p.heap[heapPos], -1)
		case '[':
			end := i + findClosing(p.code[i:])
			if atomic.LoadInt32(&p.heap[heapPos]) != 0 {
				// enter loop
				childWG.Add(1) // TODO: remove this on loops
				heapPos = program{p.code[i+1:], p.heap}.run(heapPos, &childWG)
			}
			i = end // goto end
		case ']':
			if atomic.LoadInt32(&p.heap[heapPos]) == 0 {
				return heapPos
			}
			i = -1

		case '{':
			end := i + findClosing(p.code[i:])
			childWG.Add(1)
			go program{p.code[i+1:], p.heap}.run(heapPos, &childWG)

			i = end // continue parrent thread
		case '}':
			return -1 // kill thread

		case '.':
			fmt.Printf("%c", atomic.LoadInt32(&p.heap[heapPos]))
		case ',':
			var n int32
			if _, err := fmt.Scanf("%d\n", &n); err != nil {
				panic(err)
			}
			atomic.SwapInt32(&p.heap[heapPos], n)
		case ' ', '\t', '\n':
		default:
			panic("instruction not implemented")
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
