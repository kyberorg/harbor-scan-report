package main

import "os"

func main() {
	//get args
	println("HH: " + os.Getenv("HARBOR_HOST"))
	println("GH event type:" + os.Getenv("GITHUB_EVENT_TYPE"))
	println("GH issue comment URL: " + os.Getenv("GITHUB_ISSUE_COMMENT_URL"))
	println("GH PR comment URL" + os.Getenv("GITHUB_PR_COMMENT_URL"))

	//kontroll args
	//get scan results
	//write comment
}
