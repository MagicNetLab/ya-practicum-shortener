package shortgen

import "fmt"

func ExampleGetShortLink() {
	short := GetShortLink(7)
	fmt.Println(short)
}
