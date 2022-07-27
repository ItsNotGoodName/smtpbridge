package version

type Version struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

var Current Version

func SetCurrentVersion(version, commit, date, builtBy string) {
	Current = Version{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	}
}
