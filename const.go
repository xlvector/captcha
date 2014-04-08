package captcha

const SCALE_HEIGHT = 32
const (
	NUMBER, ALPHA, MIX = 0, 1, 2
)

func GetChType(buf string) int {
	if buf == "number" {
		return NUMBER
	} else if buf == "alpha" {
		return ALPHA
	} else {
		return MIX
	}
}
