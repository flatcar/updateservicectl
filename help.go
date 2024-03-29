package main

import (
	"flag"
	"fmt"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/flatcar/updateservicectl/client/update/v1"
	"github.com/flatcar/updateservicectl/version"
)

var (
	cmdHelp = &Command{
		Name:        "help",
		Summary:     "Show a list of commands or help for one command",
		Usage:       "[COMMAND]",
		Description: "Show a list of commands or detailed help for one command",
		Run:         runHelp,
	}

	globalUsageTemplate  *template.Template
	commandUsageTemplate *template.Template
	templFuncs           = template.FuncMap{
		"descToLines": func(s string) []string {
			// trim leading/trailing whitespace and split into slice of lines
			return strings.Split(strings.Trim(s, "\n\t "), "\n")
		},
		"printOption": func(name, defvalue, usage string) string {
			prefix := "--"
			if len(name) == 1 {
				prefix = "-"
			}
			return fmt.Sprintf("\t%s%s=%s\t%s", prefix, name, defvalue, usage)
		},
	}
)

func init() {
	globalUsageTemplate = template.Must(template.New("global_usage").Funcs(templFuncs).Parse(`
NAME:
{{printf "\t%s - %s" .Executable .Description}}

USAGE:
{{printf "\t%s" .Executable}} [global options] <command> [command options] [arguments...]

VERSION:
{{printf "\t%s" .Version}}

COMMANDS:{{range .Commands}}
{{printf "\t%s\t%s" .Name .Summary}}{{end}}

GLOBAL OPTIONS:{{range .Flags}}
{{printOption .Name .DefValue .Usage}}{{end}}

Run "{{.Executable}} <command> --help" for more details on a specific command.
`[1:]))
	commandUsageTemplate = template.Must(template.New("command_usage").Funcs(templFuncs).Parse(`
NAME:
{{printf "\t%s - %s" .Cmd.Name .Cmd.Summary}}

USAGE:
{{printf "\t%s %s %s" .Executable .Cmd.Name .Cmd.Usage}}

DESCRIPTION:
{{range $line := descToLines .Cmd.Description}}{{printf "\t%s" $line}}
{{end}}
{{if .Cmd.Subcommands}}COMMANDS:{{range .Cmd.Subcommands}}
{{printf "\t%s\t%s" .Name .Summary}}{{end}}
{{end}}
{{if .CmdFlags}}OPTIONS:{{range .CmdFlags}}
{{printOption .Name .DefValue .Usage}}{{end}}
{{end}}For help on global options run "{{.Executable}} --help"
`[1:]))
}

func runHelp(args []string, service *update.Service, out *tabwriter.Writer) int {
	if len(args) < 1 {
		printGlobalUsage()
		return OK
	}

	var cmd *Command

	for _, c := range commands {
		if c.Name == args[0] {
			cmd = c
			break
		}
	}

	if cmd == nil {
		fmt.Println("Unrecognized command:", args[0])
		return ERROR_NO_COMMAND
	}

	printCommandUsage(cmd)
	return OK
}

func printGlobalUsage() {
	globalUsageTemplate.Execute(out, struct {
		Executable  string
		Commands    []*Command
		Flags       []*flag.Flag
		Description string
		Version     string
	}{
		cliName,
		commands,
		getAllFlags(),
		cliDescription,
		version.Version,
	})
}

func printCommandUsage(cmd *Command) {
	commandUsageTemplate.Execute(out, struct {
		Executable string
		Cmd        *Command
		CmdFlags   []*flag.Flag
	}{
		cliName,
		cmd,
		getFlags(&cmd.Flags),
	})
}
