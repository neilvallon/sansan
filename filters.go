package sansan

import "errors"

type filter func([]byte) ([]byte, error)

// Brainfuck removes all sansan specific instructions
func Brainfuck(b []byte) ([]byte, error) {
	for i := range b {
		switch b[i] {
		case '{', '}', '!':
			b[i] = ' '
		}
	}

	return b, nil
}

// NoRead will turn any read instructions into noops
// and return an error.
func NoRead(b []byte) ([]byte, error) {
	var err error
	for i := range b {
		if b[i] == ',' {
			b[i] = ' '
			err = errors.New("sansan: input has been disabled")
		}
	}

	return b, err
}
