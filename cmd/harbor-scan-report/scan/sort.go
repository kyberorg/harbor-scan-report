package scan

type BySeverity []Vulnerability
type ByScore []Vulnerability

func (a BySeverity) Len() int {
	return len(a)
}

func (a BySeverity) Less(i, j int) bool {
	return a[i].Severity.IsMoreCriticalThen(a[j].Severity)
}

func (a BySeverity) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByScore) Len() int {
	return len(a)
}

func (a ByScore) Less(i, j int) bool {
	return a[i].Score > a[j].Score
}

func (a ByScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
