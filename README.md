# sansan
A Concurrent Brainfuck Dialect.


## Instruction Set

| op |                 Function                  |
|----|-------------------------------------------|
| +  | Add one to current cell                   |
| -  | Subtract one from current cell            |
| >  | Move right one cell                       |
| <  | Move left one cell                        |
| [  | Enter loop if cell is not zero            |
| ]  | Continue loop untill current cell is zero |
| {  | Start new thread                          |
| }  | Stops current thread                      |
| .  | Print value of current cell               |
| ,  | Read integer into current cell            |

## Usage

#### Install
	go get vallon.me/sansan/cmd/sansan

#### Run
	sansan filename.san
