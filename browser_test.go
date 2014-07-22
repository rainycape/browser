package browser

import "testing"

func TestURL(t *testing.T) {
	if err := Open("http://google.com"); err != nil {
		t.Fatal(err)
	}
}
