package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	hqgourl "github.com/hueristiq/hq-go-url"
	"github.com/hueristiq/hqgolog"
	"github.com/hueristiq/hqgolog/formatter"
	"github.com/hueristiq/hqgolog/levels"
	"github.com/hueristiq/xurlbits/internal/configuration"
	"github.com/hueristiq/xurlbits/pkg/stdio"
	"github.com/spf13/pflag"
)

// Extractor is a type alias for a function that processes a URL and extracts certain
// components based on a format string. It returns a slice of strings, which allows
// for multiple values to be returned from a URL, such as subdomains, query keys,
// or specific parts of the URL.
//
// Parameters:
//   - URL (*hqgourl.URL): A pointer to the URL that is being processed.
//   - format (string): A string that defines how to format or process the URL.
//
// Returns:
//   - []string: A slice of strings containing the extracted components from the URL.
type Extractor func(URL *hqgourl.URL, format string) []string

var (
	URLsListFilePath string
	unique           bool
	monochrome       bool
	verbosity        string

	// mode   string
	// fmtStr string.
)

func init() {
	// Initialize and handle CLI arguments and flags using pflag.
	pflag.StringVar(&URLsListFilePath, "urls", "", "")
	pflag.BoolVarP(&unique, "unique", "u", false, "")
	pflag.BoolVarP(&monochrome, "monochrome", "m", false, "")
	pflag.StringVarP(&verbosity, "verbosity", "v", string(levels.LevelInfo), "")

	// Custom usage/help message for the utility.
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, configuration.BANNER)

		h := "\nUSAGE:\n"
		h += "  xurlbits [MODE] [FORMATSTRING] [OPTIONS]\n"

		h += "\nINPUT:\n"
		h += "     --urls string         target URLs list file path\n"

		h += "\nOUTPUT:\n"
		h += "     --monochrome bool     display no color output\n"
		h += " -u, --unique bool         output unique values\n"
		h += fmt.Sprintf(" -v, --verbosity           debug, info, warning, error, fatal or silent (default: %s)\n", string(levels.LevelInfo))

		h += "\nMODES:\n"
		h += " domains                   the hostname (e.g. sub.example.com)\n"
		h += " apexes                    the apex domain (e.g. example.com from sub.example.com)\n"
		h += " paths                     the request path (e.g. /users)\n"
		h += " query                     `key=value` pairs from the query string (one per line)\n"
		h += " params                    keys from the query string (one per line)\n"
		h += " values                    query string values (one per line)\n"
		h += " format                    custom format (see below)\n"

		h += "\nFORMAT DIRECTIVES:\n"
		h += "  %%                a literal percent character\n"
		h += "  %s                the request scheme (e.g. https)\n"
		h += "  %u                the user info (e.g. user:pass)\n"
		h += "  %d                the domain (e.g. sub.example.com)\n"
		h += "  %S                the subdomain (e.g. sub)\n"
		h += "  %r                the root of domain (e.g. example)\n"
		h += "  %t                the TLD (e.g. com)\n"
		h += "  %P                the port (e.g. 8080)\n"
		h += "  %p                the path (e.g. /users)\n"
		h += "  %e                the path's file extension (e.g. jpg, html)\n"
		h += "  %q                the raw query string (e.g. a=1&b=2)\n"
		h += "  %f                the page fragment (e.g. page-section)\n"
		h += "  %@                inserts an @ if user info is specified\n"
		h += "  %:                inserts a colon if a port is specified\n"
		h += "  %?                inserts a question mark if a query string exists\n"
		h += "  %#                inserts a hash if a fragment exists\n"
		h += "  %a                authority (alias for %u%@%d%:%P)\n"

		fmt.Fprint(os.Stderr, h)
	}

	pflag.Parse()

	// Initialize logger with the specified verbosity and colorization options.
	hqgolog.DefaultLogger.SetMaxLevel(levels.LevelStr(verbosity))
	hqgolog.DefaultLogger.SetFormatter(formatter.NewCLI(&formatter.CLIOptions{
		Colorize: !monochrome,
	}))
}

func main() {
	// Set the mode (e.g., "domains", "paths") and format string from CLI arguments.
	mode := pflag.Arg(0)
	fmtStr := pflag.Arg(1)

	// Map each mode to its corresponding extractor function.
	procFn, ok := map[string]Extractor{
		"domains": Domains,
		"apexes":  Apexes,
		"paths":   Paths,
		"query":   Query,
		"params":  Parameters,
		"values":  Values,
		"format":  Format,
	}[mode]

	// If an unknown mode is provided, log a fatal error and exit
	if !ok {
		hqgolog.Fatal().Msgf("unknown mode: %s", mode)
	}

	URLs := make(chan string)

	go func() {
		defer close(URLs)

		// Read URLs from a file if provided.
		if URLsListFilePath != "" {
			var file *os.File

			file, err := os.Open(URLsListFilePath)
			if err != nil {
				hqgolog.Error().Msg(err.Error())
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				URL := scanner.Text()

				if URL != "" {
					URLs <- URL
				}
			}

			if err := scanner.Err(); err != nil {
				hqgolog.Error().Msg(err.Error())
			}
		}

		// Alternatively, read URLs from stdin.
		if stdio.HasStdIn() {
			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				URL := scanner.Text()

				if URL != "" {
					URLs <- URL
				}
			}

			if err := scanner.Err(); err != nil {
				hqgolog.Error().Msg(err.Error())
			}
		}
	}()

	wg := &sync.WaitGroup{}

	parser := hqgourl.NewParser()
	seen := &sync.Map{}

	for URL := range URLs {
		wg.Add(1)

		go func(URL string) {
			defer wg.Done()

			// Parse the URL using the parser.
			parsed, err := parser.Parse(URL)
			if err != nil {
				hqgolog.Error().Msgf("parse failure: %s", err.Error())

				return
			}

			// Process and print the extracted values using the appropriate mode function.
			for _, value := range procFn(parsed, fmtStr) {
				if value == "" {
					continue
				}

				if unique {
					_, loaded := seen.LoadOrStore(value, struct{}{})
					if loaded {
						continue
					}
				}

				hqgolog.Print().Msg(value)
			}
		}(URL)
	}

	wg.Wait()
}

// Format takes a URL and a format string and returns a formatted string.
// The format directives allow extracting specific parts of the URL, such as the scheme, domain, path, etc.
func Format(u *hqgourl.URL, f string) []string {
	out := &bytes.Buffer{}

	inFormat := false

	for _, r := range f {
		if r == '%' && !inFormat {
			inFormat = true

			continue
		}

		if !inFormat {
			out.WriteRune(r)

			continue
		}

		switch r {
		// a literal percent rune
		case '%':
			out.WriteByte('%')

		// the scheme; e.g. http
		case 's':
			out.WriteString(u.Scheme)

		// the userinfo; e.g. user:pass
		case 'u':
			if u.User != nil {
				out.WriteString(u.User.String())
			}

		// the domain; e.g. sub.example.com
		case 'd':
			out.WriteString(u.Hostname())

		// the port; e.g. 8080
		case 'P':
			out.WriteString(u.Port())

		// the subdomain; e.g. www
		case 'S':
			out.WriteString(u.Domain.Subdomain)

		// the root; e.g. example
		case 'r':
			out.WriteString(u.Domain.SLD)

		// the tld; e.g. com
		case 't':
			out.WriteString(u.Domain.TLD)

		// the path; e.g. /users
		case 'p':
			out.WriteString(u.EscapedPath())

		// the paths's file extension
		case 'e':
			paths := strings.Split(u.EscapedPath(), "/")
			if len(paths) > 1 {
				parts := strings.Split(paths[len(paths)-1], ".")
				if len(parts) > 1 {
					out.WriteString(parts[len(parts)-1])
				}
			} else {
				parts := strings.Split(u.EscapedPath(), ".")
				if len(parts) > 1 {
					out.WriteString(parts[len(parts)-1])
				}
			}

		// the query string; e.g. one=1&two=2
		case 'q':
			out.WriteString(u.RawQuery)

		// the fragment / hash value; e.g. section-1
		case 'f':
			out.WriteString(u.Fragment)

		// an @ if user info is specified
		case '@':
			if u.User != nil {
				out.WriteByte('@')
			}

		// a colon if a port is specified
		case ':':
			if u.Port() != "" {
				out.WriteByte(':')
			}

		// a question mark if there's a query string
		case '?':
			if u.RawQuery != "" {
				out.WriteByte('?')
			}

		// a hash if there is a fragment
		case '#':
			if u.Fragment != "" {
				out.WriteByte('#')
			}

		// the authority; e.g. user:pass@example.com:8080
		case 'a':
			out.WriteString(Format(u, "%u%@%d%:%P")[0])

		// default to literal
		default:
			// output untouched
			out.WriteByte('%')
			out.WriteRune(r)
		}

		inFormat = false
	}

	return []string{out.String()}
}

// hostnames returns the domain portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"sub.example.com"}.
func Domains(u *hqgourl.URL, _ string) []string {
	return Format(u, "%d")
}

// Apexes return the apex portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"example.com"}.
func Apexes(u *hqgourl.URL, _ string) []string {
	return Format(u, "%r.%t")
}

// domains returns the path portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"/path"}.
func Paths(u *hqgourl.URL, _ string) []string {
	return Format(u, "%p")
}

// keyPairs returns all the key=value pairs in
// the query string portion of the URL. E.g for
// /?one=1&two=2&three=3 it will return
// []string{"one=1", "two=2", "three=3"}.
func Query(u *hqgourl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for key, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return out
}

// Parameters returns all of the keys used in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"one", "two", "three"}.
func Parameters(u *hqgourl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for key := range u.Query() {
		out = append(out, key)
	}

	return out
}

// values returns all of the values in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"1", "2", "3"}.
func Values(u *hqgourl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for _, value := range u.Query() {
		// value: [string]{items...}
		out = append(out, value...)
	}

	return out
}
