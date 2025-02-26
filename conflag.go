package conflag

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// A Conflag represents the state of a conflag.
type Conflag struct {
	// embeded the standard FlagSet so wen can use all it's methods.
	*flag.FlagSet

	app     string
	osArgs  []string
	cfgFile string
	args    []string

	includes []string

	// TODO: add shorthand? or just use pflag?
	// shorthand map[byte]string
}

// New parses os args and returns a new Conflag instance.
func New(args ...string) *Conflag {
	if args == nil {
		args = os.Args
	}

	c := &Conflag{app: args[0], osArgs: args[1:]}

	c.FlagSet = flag.NewFlagSet(c.app, flag.ExitOnError)
	c.StringVar(&c.cfgFile, "config", "", "config file path")
	c.StringSliceUniqVar(&c.includes, "include", nil, "include file")

	return c
}

// NewFromFile parses cfgFile and returns a new Conflag instance.
func NewFromFile(app, cfgFile string) *Conflag {
	if app == "" {
		app = os.Args[0]
	}

	c := &Conflag{app: app, cfgFile: cfgFile}

	c.FlagSet = flag.NewFlagSet(c.app, flag.ExitOnError)
	c.StringSliceUniqVar(&c.includes, "include", nil, "include file")

	return c
}

// Parse parses config file and flags.
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

	// parse 2nd time to get the include file values
	err = c.FlagSet.Parse(c.args)
	if err != nil {
		return err
	}

	dir := filepath.Dir(c.cfgFile)

	// parse 3rd time to parse flags in include file
	for _, include := range c.includes {
		include = filepath.Join(dir, include)
		fargs, err := parseFile(include)
		if err != nil {
			return err
		}

		c.args = fargs
		c.args = append(c.args, c.osArgs...)

		err = c.FlagSet.Parse(c.args)
	}

	return err
}

func parseFile(cfgFile string) ([]string, error) {
	var s []string

	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[:1] == "#" {
			continue
		}
		s = append(s, "-"+line)
	}

	return s, nil
}

// AppDir returns the app dir.
func (c *Conflag) AppDir() string {
	return filepath.Dir(os.Args[0])
}

// ConfDir returns the config file dir.
func (c *Conflag) ConfDir() string {
	return filepath.Dir(c.cfgFile)
}
