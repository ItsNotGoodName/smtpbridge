package build

type Build struct {
	BuiltBy string
	Commit  string
	Date    string
	Version string
	RepoURL string
}

var Current Build

var RepoURL string = "https://github.com/ItsNotGoodName/smtpbridge"
