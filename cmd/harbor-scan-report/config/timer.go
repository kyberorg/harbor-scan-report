package config

type Timer struct {
	SecondsLeft int
}

func (t Timer) IsTimeOver() bool {
	return t.SecondsLeft <= 0
}
