package scan

type Status struct {
	Failed         bool
	ImageFound     bool
	ScanCompleted  bool
	ScanResultsUrl string
}
