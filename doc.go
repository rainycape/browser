// Package browser implements functions for opening a web browser.
//
// This package also allows opening an URL back at the original host
// when run from an SSH session. In order for this to work, SSH agent
// forwarding must be enabled, the initial host must have SSH running
// on port 22 and the local and remote usernames must match. Otherwise
// a local browser is opened.
package browser
