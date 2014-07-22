package browser

import (
	"os"
	"path/filepath"
	"strings"
)

// Open opens the given file or URL in a browser. If fileOrUrl is
// a valid and existing file path, it calls OpenFile. Otherwise, it
// calls OpenURL.
func Open(fileOrURL string) error {
	if st, err := os.Stat(fileOrURL); err == nil && !st.IsDir() {
		return OpenFile(fileOrURL)
	}
	return OpenURL(fileOrURL)
}

// OpenFile opens the file at the given path in a browser.
func OpenFile(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return OpenURL("file://" + abs)
}

// OpenURL opens the given URL in a browser. If the URL does not
// have a scheme, http is used.
func OpenURL(url string) error {
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}
	return openBrowser(url)
}
