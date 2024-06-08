//go:build amd64
package main

import (
    "C"
    "fmt"
)

//export TestFunc
func TestFunc() {
    fmt.Println("Hello World!")
}

func main() {

}
