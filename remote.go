package browser

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/agent"
)

var (
	errNoRemote = errors.New("can't find remote address")
)

func shellJoin(cmd string, args []string) string {
	if len(args) > 0 {
		values := make([]string, len(args))
		for ii, v := range args {
			if strings.Contains(v, " ") {
				values[ii] = fmt.Sprintf("\"%s\"", strings.Replace(v, "\"", "\\\"", -1))
			} else {
				values[ii] = v
			}
		}
		added := strings.Join(values, " ")
		if added != "" {
			return cmd + " " + added
		}
	}
	return cmd

}

func runRemoteCmd(client *ssh.Client, cmd string, args ...string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()
	var out bytes.Buffer
	session.Stdout = &out
	session.Stderr = &out
	fullCmd := shellJoin(cmd, args)
	if err := session.Run(fullCmd); err != nil {
		return "", fmt.Errorf("failed to run %s: %s", fullCmd, err)
	}
	return out.String(), nil
}

func openRemoteBrowser(url string) error {
	sshClient := os.Getenv("SSH_CLIENT")
	parts := strings.Split(sshClient, " ")
	if parts[0] == "" {
		return errNoRemote
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	var auth []ssh.AuthMethod
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		conn, _ := net.Dial("unix", sock)
		if conn != nil {
			sshAgent := agent.NewClient(conn)
			signers, _ := sshAgent.Signers()
			if len(signers) > 0 {
				auth = append(auth, ssh.PublicKeys(signers...))
			}
		}
	}
	config := &ssh.ClientConfig{
		User: usr.Username,
		Auth: auth,
	}
	client, err := ssh.Dial("tcp", parts[0]+":22", config)
	if err != nil {
		return err
	}
	defer client.Close()
	out, err := runRemoteCmd(client, "uname", "-s")
	if err != nil {
		return err
	}
	cmds, err := goosCmds(strings.ToLower(strings.TrimSpace(out)))
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
			_, err = runRemoteCmd(client, cmd[0], args...)
		}
		if err == nil {
			// Browser did open
			break
		}
	}

	return nil
}
