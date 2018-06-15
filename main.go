package main

func main() {
	a := App{}
	a.Initialize("nex", "next@1", "nexthalt")
	a.Run(":9001")
}
