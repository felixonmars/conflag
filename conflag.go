package conflag

import (
	"bufio"
	"flag"
	"os"
	"strings"
)

// Conflag .
type Conflag struct {
	*flag.FlagSet

	app     string
	osArgs  []string
	cfgFile string
	args    []string

	// TODO: add shorthand? or just use pflag?
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

// NewFromFile ...
func NewFromFile(app, cfgFile string) *Conflag {
	c := &Conflag{}

	if app != "" {
		c.app = app
	} else {
		c.app = os.Args[0]
	}

	c.cfgFile = cfgFile
	c.FlagSet = flag.NewFlagSet(c.app, flag.ExitOnError)
	return c
}

// Parse ...
func (c *Conflag) Parse() (err error) {
	// parse 1st time and see whether there is a conf file.
	err = c.FlagSet.Parse(c.osArgs)
	if err != nil {
		return err
	}

	// if there is no args, just try to load the app.conf file.
	if c.cfgFile == "" && len(c.osArgs) == 0 {
		// trim app exetension
		for i := len(c.app) - 1; i >= 0 && c.app[i] != '/' && c.app[i] != '\\'; i-- {
			if c.app[i] == '.' {
				c.cfgFile = c.app[:i]
				break
			}
		}

		if c.cfgFile == "" {
			c.cfgFile = c.app
		}

		c.cfgFile += ".conf"
	}

	if c.cfgFile == "" {
		return nil
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
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[:1] == "#" {
			continue
		}
		s = append(s, "-"+line)
	}

	return s, nil
}
