package command

import (
	"github.com/dihedron/excel/command/load"
	"github.com/dihedron/excel/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Login is the command that checks logins to an LDAP server.
	Load load.Load `command:"load" alias:"l" description:"Load an Excel file into the database."`
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
