package main

import (
	"fmt"

	"github.com/lemon-mint/rfyy/internal/checkenv"
)

func main() {
	env, err := checkenv.Check()
	if err != nil {
		panic(err)
	}
	fmt.Println(env.String())
}
