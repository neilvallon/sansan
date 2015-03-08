package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"vallon.me/sansan"
)

func main() {
	flag.Parse()

	prog, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Println(err)
		return
	}

	t1 := time.Now()

	sansan.NewMachine().Run(sansan.Program(prog))

	fmt.Fprintln(os.Stderr, "\nProgram exited in:", time.Now().Sub(t1))
}
