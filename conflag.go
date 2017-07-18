package conflag

import (
	"bufio"
	"flag"
	"os"
)

// Conflag .
type Conflag struct {
	*flag.FlagSet

	app     string
	osArgs  []string
	cfgFile string
	args    []string

	// TODO: add shorthand? of just use pflag?
	// shorthand map[byte]string
}

// New ...
func New(args ...string) *Conflag {
	if args == nil {
		args = os.Args
	}

	c := &Conflag{}

	c.app = args[0]
	c.osArgs = args[1:]
	c.FlagSet = flag.NewFlagSet(c.app, flag.ExitOnError)
	c.FlagSet.StringVar(&c.cfgFile, "config", "", "config file path")

	return c
}

// Parse ...
func (c *Conflag) Parse() (err error) {
	// parse 1st time and see whether there is a conf file.
	err = c.FlagSet.Parse(c.osArgs)
	if err != nil || c.cfgFile == "" {
		return err
	}

	fargs, err := parseFile(c.cfgFile)
	if err != nil {
		return err
	}
	c.args = fargs
	c.args = append(c.args, c.osArgs...)

	return c.FlagSet.Parse(c.args)
}

func parseFile(cfgFile string) ([]string, error) {
	var s []string

	fp, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 || line[:1] == "#" {
			continue
		}

		s = append(s, "-"+line)
	}

	return s, nil
}
