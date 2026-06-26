package base

// Command is the base command.
type Command struct {
	DB     string `short:"D" long:"db" description:"The database file to use." required:"true"`
	Format string `short:"F" long:"format" description:"The format of the output." optional:"true" default:"none" choice:"text" choice:"json" choice:"yaml" choice:"csv" choice:"none" env:"EXCEL_FORMAT"`
}
