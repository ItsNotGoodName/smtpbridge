package version

type Version struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

var Current Version
