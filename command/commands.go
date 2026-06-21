package command

import (
	"github.com/dihedron/excel/command/load"
	"github.com/dihedron/excel/command/query"
	"github.com/dihedron/excel/command/reports"
	"github.com/dihedron/excel/command/version"
)

// Commands is the set of root command groups.
type Commands struct {
	// Load is the command that loads an Excel file into the database.
	Load load.Load `command:"load" alias:"l" description:"Load an Excel file into the database."`
	// Query is the command that queries the database.
	Query query.Query `command:"query" alias:"q" description:"Query the database."`
	// Reports is the command that generates reports.
	Reports reports.Reports `command:"reports" alias:"r" description:"Generate reports."`
	// Version prints overlay version information and exits.
	Version version.Version `command:"version" alias:"v" description:"Show the command version and exit."`
}
