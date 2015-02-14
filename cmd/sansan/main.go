package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	sansan.NewProg(prog).Run(0, &wg)
	wg.Wait()

	fmt.Println("\n\nProgram exited in:", time.Now().Sub(t1))
}
