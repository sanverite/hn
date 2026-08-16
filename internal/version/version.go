package version 

// These are set at build time via -ldflags
var (
	GitSHA = "dev"
	BuildTime = "unkown"
	GoVersion = "unkown"
)
