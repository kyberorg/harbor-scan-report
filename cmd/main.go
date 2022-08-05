package main

import "os"

func main() {
	//get args
	println("HH: " + os.Getenv("HARBOR_HOST"))
	println("GH URL: " + os.Getenv("GITHUB_URL"))
	println("Fail level: " + os.Getenv("FAIL_LEVEL"))

	//kontroll args
	//get scan results
	//write comment
}
