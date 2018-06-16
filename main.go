package main

func main() {
	a := App{}
	a.Initialize("username", "password", "nexthalt")
	a.Run(":9001")
}
