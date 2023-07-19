package build

type Build struct {
	BuiltBy string
	Commit  string
	Date    string
	Version string
}

var Current Build
