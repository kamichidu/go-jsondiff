package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kamichidu/go-jsondiff"
	"github.com/kamichidu/go-jsondiff/internal/cmdutil"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	_ "github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

type VersionDescriptor struct {
	Version, Commit, Date, BuiltBy string
}

func (v *VersionDescriptor) String() string {
	return strings.Join([]string{
		fmt.Sprintf("%s - %s", v.Version, v.Commit),
		fmt.Sprintf("built by %s at %s", v.BuiltBy, v.Date),
	}, "\n")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "json-diff [flags] {fileA} {fileB}",
	Args:  cobra.ExactArgs(2),
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		setProperties := mustGetStrings(cmd.Flags().GetStringArray("set-property"))
		setPropertiesName := mustGetString(cmd.Flags().GetString("set-property-from-file"))
		ignoreProperties := mustGetStrings(cmd.Flags().GetStringArray("ignore-property"))
		ignorePropertiesName := mustGetString(cmd.Flags().GetString("ignore-property-from-file"))
		prettyPrint := mustGetBool(cmd.Flags().GetBool("pretty-print"))
		outName := mustGetString(cmd.Flags().GetString("output"))
		verbose := mustGetInt(cmd.Flags().GetCount("verbose"))
		colorMode := cmd.Flag("color").Value.String()
		aName := args[0]
		bName := args[1]

		var logger *log.Logger
		if verbose > 0 {
			logger = log.New(os.Stderr, "> ", 0)
		} else {
			logger = log.New(ioutil.Discard, "", 0)
		}

		if setPropertiesName != "" {
			if v, err := readLines(setPropertiesName); err != nil {
				return err
			} else {
				setProperties = append(setProperties, v...)
			}
		}
		if ignorePropertiesName != "" {
			if v, err := readLines(ignorePropertiesName); err != nil {
				return err
			} else {
				ignoreProperties = append(ignoreProperties, v...)
			}
		}

		var stdout io.Writer
		switch colorMode {
		case "never":
			stdout = colorable.NewNonColorable(os.Stdout)
		case "auto":
			// is pipe or not
			if info, err := os.Stdout.Stat(); err != nil {
				logger.Printf("unable to detect stdout is or is not a pipe: %v", err)
				stdout = colorable.NewNonColorable(os.Stdout)
			} else if (info.Mode() & os.ModeCharDevice) == 0 {
				// is a pipe
				stdout = colorable.NewNonColorable(os.Stdout)
			} else {
				// is not a pipe
				stdout = colorable.NewColorable(os.Stdout)
			}
		default:
			stdout = colorable.NewColorable(os.Stdout)
		}

		var opts []jsondiff.Option
		for _, p := range setProperties {
			logger.Printf("the path %q treat as a set", p)
			opts = append(opts, jsondiff.WithSetPath(p))
		}
		for _, p := range ignoreProperties {
			logger.Printf("the path %q will be ignored", p)
			opts = append(opts, jsondiff.WithIgnorePath(p))
		}
		opts = append(opts, jsondiff.WithLogger(logger))

		var w io.Writer
		if outName == "-" {
			w = stdout
		} else {
			file, err := os.Create(outName)
			if err != nil {
				return err
			}
			defer file.Close()
			w = file
		}

		jsonFmtFn := newJSONFormatFunc(prettyPrint)

		aContent, aStat, err := readJSONFile(aName)
		if err != nil {
			return err
		}
		bContent, bStat, err := readJSONFile(bName)
		if err != nil {
			return err
		}

		hunks, err := jsondiff.Diff(aContent, bContent, opts...)
		if err != nil {
			return err
		}
		if len(hunks) == 0 {
			return nil
		}

		// want "diff -u" output
		const diffTimeLayout = "2006-01-02 15:04:05.000000000 -0700"
		fmt.Fprintf(w, "--- %q\t%s\n", aName, aStat.ModTime().Format(diffTimeLayout))
		fmt.Fprintf(w, "+++ %q\t%s\n", bName, bStat.ModTime().Format(diffTimeLayout))
		for i, hunk := range hunks {
			if i > 0 {
				fmt.Fprintln(w)
			}
			writeHunk(w, hunk, jsonFmtFn)
		}

		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func newJSONFormatFunc(prettyPrint bool) func(string, []byte) string {
	return func(pfx string, b []byte) string {
		if !prettyPrint {
			return pfx + string(b)
		}
		var out bytes.Buffer
		je := json.NewEncoder(&out)
		je.SetIndent(pfx, "  ")
		if err := je.Encode(json.RawMessage(b)); err != nil {
			panic(err)
		}
		return pfx + strings.TrimSuffix(out.String(), "\n")
	}
}

func writeHunk(w io.Writer, hunk jsondiff.Hunk, jsonFormatFn func(string, []byte) string) {
	var oldStr string
	if hunk.Old != nil {
		oldStr = jsonFormatFn("-", *hunk.Old)
	}
	var newStr string
	if hunk.New != nil {
		newStr = jsonFormatFn("+", *hunk.New)
	}
	fmt.Fprintln(w, aurora.Cyan(fmt.Sprintf("@@ %s @@", hunk.Path)))
	if oldStr != "" {
		fmt.Fprintln(w, aurora.Red(oldStr))
	}
	if newStr != "" {
		fmt.Fprintln(w, aurora.Green(newStr))
	}
}

func readLines(name string) ([]string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out, nil
}

func readJSONFile(name string) ([]byte, os.FileInfo, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	return content, stat, nil
}

func mustGetString(v string, err error) string {
	if err != nil {
		panic(err)
	}
	return v
}

func mustGetBool(v bool, err error) bool {
	if err != nil {
		panic(err)
	}
	return v
}

func mustGetStrings(v []string, err error) []string {
	if err != nil {
		panic(err)
	}
	return v
}

func mustGetInt(v int, err error) int {
	if err != nil {
		panic(err)
	}
	return v
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version VersionDescriptor) {
	log.SetOutput(os.Stderr)
	log.SetPrefix("> ")
	log.SetFlags(0)

	rootCmd.SetVersionTemplate(`{{ .Version }}`)
	rootCmd.Version = version.String()
	if err := rootCmd.Execute(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func init() {
	c := rootCmd
	c.Flags().StringArrayP("set-property", "", nil, "treat given properties as a set")
	c.Flags().StringP("set-property-from-file", "", "", "read a file as a --set-property for each line")
	c.Flags().StringArrayP("ignore-property", "", nil, "ignore given properties to compare")
	c.Flags().StringP("ignore-property-from-file", "", "", "read a file as a --ignore-property for each line")
	c.Flags().BoolP("pretty-print", "", false, "pretty print")
	c.Flags().StringP("output", "o", "-", "output diffs")
	c.Flags().CountP("verbose", "v", "verbose logging")
	c.Flags().VarP(&cmdutil.StringEnumValue{
		Choices: []string{"always", "auto", "never"},
	}, "color", "", "colorize the output; `WHEN` can be 'always' (default if omitted), 'auto', or 'never'")
}
