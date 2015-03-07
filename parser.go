package sansan

type Program []byte

func (p Program) Run() {
	m := newMachine()
	defer m.wg.Wait()

	m.run(p.Clean(), new(state))
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
