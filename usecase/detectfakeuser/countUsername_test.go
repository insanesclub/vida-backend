package detectfakeuser

import "fmt"

func Example_atomicBool() {
	var atombool = new(atomicBool)
	atombool.set(true)
	fmt.Println(atombool.get())
	atombool.set(false)
	fmt.Println(atombool.get())
	// Output:
	// true
	// false
}

func ExampleUsernames() {
	fmt.Println(len(Usernames("gopher")))
	// Output: 10
}
