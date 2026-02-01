package color

// Public Constants
const (
	Reset  = "\x1b[0m"
	Green  = "\x1b[32m"
	Red    = "\x1b[31m"
	Blue   = "\x1b[34m"
	Orange = "\x1b[33m"
	Bold   = "\x1b[1m"
)

// Public Functions
func Paint(colorCode string, text string) string {
	return colorCode + text + Reset
}
