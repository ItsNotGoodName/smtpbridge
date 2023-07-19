package envelope

import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

var extensionRe *regexp.Regexp = regexp.MustCompile("^\\.[a-zA-Z0-9]+$")

func validExtension(str string) bool {
	return extensionRe.MatchString(str)
}

func fileExtension(name string, mimeT string) string {
	extension := strings.ToLower(filepath.Ext(name))
	if validExtension(extension) {
		return extension
	}

	extension = ""
	extensions, err := mime.ExtensionsByType(mimeT)
	if err == nil && extensions != nil {
		extension = extensions[0]
		// Use extension from name if it is valid
		unknownExt := filepath.Ext(name)
		for _, ext := range extensions {
			if ext == unknownExt {
				extension = ext
				break
			}
		}
	}
	return extension
}
