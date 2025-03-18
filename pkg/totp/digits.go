package totp

import "fmt"

type Digits int

const (
	DigitsSix   Digits = 6
	DigitsEight Digits = 8
)

func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}
