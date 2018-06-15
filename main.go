package main

func main() {
	a := App{}
	a.Initialize("nexthalt", "next@1234", "nexthalt")
	a.Run(":9000")
}
