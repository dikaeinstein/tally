package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/dikaeinstein/tally/tournament"
)

type Config struct {
	verbose bool

	args []string
}

func Run(progName string, out io.Writer) int {
	cfg, err := ParseArgs(progName, os.Args[1:], out)
	if err != nil {
		fmt.Fprintln(out, err)
		return 2
	}

	in, err := processFile(cfg)
	if err != nil {
		fmt.Fprintln(out, err)
		return 1
	}
	defer in.Close()

	matches, err := tournament.ParseInput(in)
	if err != nil {
		fmt.Fprintln(out, err)
		return 1
	}

	table := tournament.Tally(matches)
	printTable(table, out)

	return 0
}

func ParseArgs(progName string, args []string, out io.Writer) (*Config, error) {
	flags := flag.NewFlagSet(progName, flag.ContinueOnError)
	flags.SetOutput(out)
	flags.Usage = func() { usage(flags, progName, out) }

	var conf Config
	flags.BoolVar(&conf.verbose, "verbose", false, "set verbosity")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}
	conf.args = flags.Args()
	return &conf, nil
}

func processFile(cfg *Config) (*os.File, error) {
	if len(cfg.args) == 0 {
		return os.Stdin, nil
	}

	if len(cfg.args) == 1 {
		return os.OpenFile(cfg.args[0], os.O_RDONLY, os.ModePerm)
	}

	return nil, errors.New("takes at most one input")
}

// usage is a replacement usage function for the flags package.
func usage(flags *flag.FlagSet, progName string, out io.Writer) {
	fmt.Fprintf(out, "Usage of %s\n", progName)
	fmt.Fprintf(out, "\t%s <filepath>\n", progName)
	fmt.Fprintf(out, "\tcat file.csv | %s\n", progName)
	fmt.Fprintf(out, "For more information run\n")
	fmt.Fprintf(out, "\t%s -h\n\n", progName)
	fmt.Fprintf(out, "Flags:\n")
	flags.PrintDefaults()
}

func printTable(t *tournament.Table, out io.Writer) {
	// ```text
	// Team                           | MP |  W |  D |  L |  P
	// Devastating Donkeys            |  3 |  2 |  1 |  0 |  7
	// Allegoric Alaskans             |  3 |  2 |  0 |  1 |  6
	// Blithering Badgers             |  3 |  1 |  0 |  2 |  3
	// Courageous Californians        |  3 |  0 |  1 |  2 |  1
	// ```

	fmt.Fprintln(out, "Team                            | MP |  W |  D |  L |  P")
	for _, r := range t.Rows {
		fmt.Fprintf(out, "%-31s |  %d |  %d |  %d |  %d |  %d\n",
			r.Team, r.MatchPlayed, r.MatchesWon, r.MatchesDrawn, r.MatchesLost, r.Points)
	}
}
