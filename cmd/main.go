package main

import "os"

func main() {
	//get args
	println("HH: " + os.Getenv("HARBOR_HOST"))
	println("GH URL: " + os.Getenv("GITHUB_URL"))

	//kontroll args
	//get scan results
	//write comment
}
