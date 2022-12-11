package main

import (
	"fmt"

	"github.com/umutozd/stats-keeper/protos/statspb"
)

func main() {
	hm := statspb.HelloMessage{
		Hello: "hello",
		World: 42,
	}

	fmt.Printf("Hello: %q, World: %d\n", hm.Hello, hm.World)
}
