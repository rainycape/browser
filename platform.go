package browser

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	platforms = map[string][][]string{
		"windows": [][]string{
			[]string{"cmd", "/c", "start"},
		},
		"darwin": [][]string{
			[]string{"open"},
		},
		"unix": [][]string{
			[]string{"sensible-browser"},
			[]string{"xdg-open"},
		},
		"dragonfly": [][]string{{"@unix"}},
		"freebsd":   [][]string{{"@unix"}},
		"netbsd":    [][]string{{"@unix"}},
		"openbsd":   [][]string{{"@unix"}},
		"linux":     [][]string{{"@unix"}},
		"solaris":   [][]string{{"@unix"}},
	}
)

func goosCmds(goos string) ([][]string, error) {
	cmds := platforms[goos]
	if len(cmds) == 0 {
		return nil, fmt.Errorf("unsupported platform %s", goos)
	}
	if len(cmds) == 1 && len(cmds[0]) == 1 && cmds[0][0] != "" && cmds[0][0][0] == '@' {
		return goosCmds(cmds[0][0][1:])
	}
	return cmds, nil
}

func openBrowser(url string) error {
	if !strings.HasPrefix(url, "file://") {
		if err := openRemoteBrowser(url); err != nil {
			if err != errNoRemote {
				// There's a remote connection, but the
				// browser couldn't be opened
				return err
			}
		} else {
			// Succesfully opened a remote browser
			return nil
		}
	}
	cmds, err := goosCmds(runtime.GOOS)
	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		if len(cmd) == 0 || cmd[0] == "" {
			// Should not happen, but...
			continue
		}
		_, err = exec.LookPath(cmd[0])
		if err == nil {
			var args []string
			args = append(args, cmd[1:]...)
			args = append(args, url)
			c := exec.Command(cmd[0], args...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			err = c.Run()
		}
		if err == nil {
			// Browser did open
			break
		}
	}
	return err
}
