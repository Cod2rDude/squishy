package config

// Config
const (
	Version string = "1.0.0"
)

var DefaultExpectedFileExtensions = map[string]bool{
	".squishy":  true,
	".sqy": true,
}

const (
	GRANDMASTER int = 0
	MASTER      int = 2
	FELLOWCRAFT int = 4
	APPRENTICE  int = 6
	ROYAL       int = 8
	UPPERCLASS  int = 10
	MIDCLASS    int = 12
	BOTTOMCLASS int = 14
)
