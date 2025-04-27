// Code example for lexer

package examples

import "fmt"

func main() {
	fmt.Println("Hello, World!")

	x := 42 + 0x42 + 0o42 + 0b101
	y := 3.14e-2
	z := true

	fmt.Println(x, y, z)

	for i := 0; i < 10; i++ {
		if i == 1 {
			break
		}

		fmt.Println(i)
	}
}
