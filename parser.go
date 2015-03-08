package sansan

import "errors"

type program []instruction

func Program(p []byte) program {
	prog, err := Parse(p)
	if err != nil {
		panic(err)
	}

	return prog
}

func Parse(p []byte) (program, error) {
	// allocate some space for cleaned program
	b := make([]byte, 0, len(p)/2)

	for _, c := range p {
		switch c {
		case '>', '<', '+', '-', '[', ']', '{', '}', '.', ',', '!':
			b = append(b, c)
		}
	}

	return parse(b)
}

type action uint8

const (
	Modify action = iota
	Move

	Read
	Print

	LStart
	LEnd

	TStart
	TEnd

	Toggle
)

type instruction struct {
	Action action
	Val    int16
}

func parse(p []byte) (program, error) {
	prog := make(program, 0, len(p)/2)

	for len(p) != 0 {
		var i instruction
		switch p[0] {
		case '+', '-':
			i, p = parseModify(p)
		case '>', '<':
			i, p = parseMove(p)
		case ',':
			i.Action, p = Read, p[1:]
		case '.':
			i.Action, p = Print, p[1:]
		case '[':
			i.Action, p = LStart, p[1:]
		case ']':
			i.Action, p = LEnd, p[1:]
		case '{':
			i.Action, p = TStart, p[1:]
		case '}':
			i.Action, p = TEnd, p[1:]
		case '!':
			i.Action, p = Toggle, p[1:]
		}
		prog = append(prog, i)
	}

	err := findLoopEnds(prog)
	if err != nil {
		return nil, err
	}

	return prog, nil
}

func parseModify(p []byte) (i instruction, rest []byte) {
	i.Action = Modify

	for n, c := range p {
		switch c {
		case '+':
			i.Val++
		case '-':
			i.Val--
		default:
			return i, p[n:]
		}
	}
	return i, rest
}

func parseMove(p []byte) (i instruction, rest []byte) {
	i.Action = Move

	for n, c := range p {
		switch c {
		case '>':
			i.Val++
		case '<':
			i.Val--
		default:
			return i, p[n:]
		}
	}
	return i, rest
}

func findLoopEnds(p program) error {
	for i := range p {
		switch p[i].Action {
		case LStart, TStart:
			end, err := findClosing(p[i:])
			if err != nil {
				return err
			}

			p[i].Val = int16(end)
		}
	}

	return nil
}

func findClosing(prog program) (int, error) {
	braces := 0
	for i := 0; i < len(prog); i++ {
		switch prog[i].Action {
		case LStart, TStart:
			braces++
		case LEnd, TEnd:
			braces--
			if braces < 0 {
				return 0, errors.New("sansan: unbalanced braces")
			}
			if braces == 0 {
				return i, nil
			}
		}
	}

	return 0, errors.New("sansan: reached end of input inside loop or thread")
}
