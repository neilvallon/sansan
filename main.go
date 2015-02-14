package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	// Print 65
	// prog := `++++++ [ > ++++++++++ < - ] > +++++ .`

	// Echo Input
	// prog := `, [ > + < - ] > .`

	// Multiply input
	prog := `,>,< [ > [ >+ >+ << - ] >> [- << + >> ] <<< - ] >> .`

	NewProg([]byte(prog)).run(0)
	fmt.Println()
}

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

func (p program) run(heapPos int) {
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
				program{p.code[i+1:], p.heap}.run(heapPos)
			}
			i = end // goto end
		case ']':
			if atomic.LoadInt32(&p.heap[heapPos]) == 0 {
				return
			}
			i = -1
		case '.':
			fmt.Printf("%d", atomic.LoadInt32(&p.heap[heapPos]))
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
}

func findClosing(prog []byte) int {
	braces := 0
	for i := 0; i < len(prog); i++ {
		switch prog[i] {
		case '[':
			braces++
		case ']':
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
