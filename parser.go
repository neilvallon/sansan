package sansan

type program []instruction

func Program(p []byte) program {
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

func parse(p []byte) program {
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

	findLoopEnds(prog)

	return prog
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

func findLoopEnds(p program) {
	for i := range p {
		switch p[i].Action {
		case LStart, TStart:
			p[i].Val = int16(findClosing(p[i:]))
		}
	}
}

func findClosing(prog program) int {
	braces := 0
	for i := 0; i < len(prog); i++ {
		switch prog[i].Action {
		case LStart, TStart:
			braces++
		case LEnd, TEnd:
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
