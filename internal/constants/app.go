package constants

var (
	debug       = "true"
	Debug       = true
	ServiceName = ""
)

func init() {
	if debug != "true" {
		Debug = false
	}
}
