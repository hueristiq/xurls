package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	hqgologger "github.com/hueristiq/hq-go-logger"
	"github.com/hueristiq/hq-go-url/parser"
	"github.com/hueristiq/xurls/internal/configuration"
	"github.com/hueristiq/xurls/internal/input"
	"github.com/spf13/cobra"
)

type Extractor func(URL *parser.URL, format string) []string

func Parse() (cmd *cobra.Command) {
	var (
		inputURLs             []string
		inputURLsListFilePath string
		unique                bool
	)

	cmd = &cobra.Command{
		Use:     "parse",
		Aliases: []string{"p"},
		Short:   "Command for parsing URLs into granular components.",
		Long:    configuration.BANNER(au),
		Run: func(_ *cobra.Command, args []string) {
			var err error

			var mode string

			if len(args) >= 1 {
				mode = args[0]
			}

			var fmtStr string

			if len(args) >= 2 {
				fmtStr = args[1]
			}

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
				hqgologger.Fatal().Msgf("unknown mode: %s", mode)
			}

			URLs := make(chan string)

			go func() {
				defer close(URLs)

				if len(inputURLs) > 0 {
					for _, URL := range inputURLs {
						URLs <- URL
					}
				}

				if inputURLsListFilePath != "" {
					var file *os.File

					file, err = os.Open(inputURLsListFilePath)
					if err != nil {
						hqgologger.Error().Msg(err.Error())
					}

					scanner := bufio.NewScanner(file)

					for scanner.Scan() {
						URL := scanner.Text()

						if URL != "" {
							URLs <- URL
						}
					}

					if err = scanner.Err(); err != nil {
						hqgologger.Error().Msg(err.Error())
					}

					file.Close()
				}

				if input.HasStdin() {
					scanner := bufio.NewScanner(os.Stdin)

					for scanner.Scan() {
						URL := scanner.Text()

						if URL != "" {
							URLs <- URL
						}
					}

					if err = scanner.Err(); err != nil {
						hqgologger.Error().Msg(err.Error())
					}
				}
			}()

			wg := &sync.WaitGroup{}

			p := parser.New(parser.WithDefaultScheme("http"))

			seen := &sync.Map{}

			for URL := range URLs {
				wg.Add(1)

				go func(URL string) {
					defer wg.Done()

					parsed, err := p.Parse(URL)
					if err != nil {
						hqgologger.Error().Msgf("parse failure: %s", err.Error())

						return
					}

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

						hqgologger.Print().Msg(value)
					}
				}(URL)
			}

			wg.Wait()
		},
	}

	cmd.Flags().StringSliceVarP(&inputURLs, "url", "u", []string{}, "target URL")
	cmd.Flags().StringVarP(&inputURLsListFilePath, "list", "l", "", "target URLs list file path")
	cmd.Flags().BoolVar(&unique, "unique", false, "output unique values")

	h := "Usage:\n"
	h += fmt.Sprintf(" %s parse [MODE] [FORMATSTRING] [OPTIONS]\n", configuration.NAME)

	h += "\nAliases:\n"
	h += "  parse, p\n"

	h += "\nModes:\n"
	h += "  domains                   the hostname (e.g. sub.example.com)\n"
	h += "  apexes                    the apex domain (e.g. example.com from sub.example.com)\n"
	h += "  paths                     the request path (e.g. /users)\n"
	h += "  query                     `key=value` pairs from the query string (one per line)\n"
	h += "  params                    keys from the query string (one per line)\n"
	h += "  values                    query string values (one per line)\n"
	h += "  format                    custom format (see below)\n"

	h += "\nFormat Directives:\n"
	h += "   %%                       a literal percent character\n"
	h += "   %s                       the request scheme (e.g. https)\n"
	h += "   %u                       the user info (e.g. user:pass)\n"
	h += "   %d                       the domain (e.g. sub.example.com)\n"
	h += "   %S                       the subdomain (e.g. sub)\n"
	h += "   %r                       the root of domain (e.g. example)\n"
	h += "   %t                       the TLD (e.g. com)\n"
	h += "   %P                       the port (e.g. 8080)\n"
	h += "   %p                       the path (e.g. /users)\n"
	h += "   %e                       the path's file extension (e.g. jpg, html)\n"
	h += "   %q                       the raw query string (e.g. a=1&b=2)\n"
	h += "   %f                       the page fragment (e.g. page-section)\n"
	h += "   %@                       inserts an @ if user info is specified\n"
	h += "   %:                       inserts a colon if a port is specified\n"
	h += "   %?                       inserts a question mark if a query string exists\n"
	h += "   %#                       inserts a hash if a fragment exists\n"
	h += "   %a                       authority (alias for %u%@%d%:%P)\n"

	h += "\nFlags:\n"
	h += "  -u, --url string[]        target URL\n"
	h += "  -l, --list string         target URLs list file path\n"

	h += "\nTIP: For multiple input URLs use comma(,) separated value with `-u`,\n"
	h += "     specify multiple `-u`, load from file with `-l` or load from stdin.\n"

	h += "\nGlobal Flags::\n"
	h += "      --unique bool         output unique values\n"
	h += "      --monochrome bool     display no color output\n"
	h += "  -s, --silent bool         stdout values only output\n"
	h += "  -v, --verbose bool        stdout verbose output\n\n"

	cmd.SetUsageTemplate(h)

	return cmd
}

func Format(u *parser.URL, f string) []string {
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
			out.WriteString(u.Domain.SecondLevelDomain)

		// the tld; e.g. com
		case 't':
			out.WriteString(u.Domain.TopLevelDomain)

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

func Domains(u *parser.URL, _ string) []string {
	return Format(u, "%d")
}

func Apexes(u *parser.URL, _ string) []string {
	return Format(u, "%r.%t")
}

func Paths(u *parser.URL, _ string) []string {
	return Format(u, "%p")
}

func Query(u *parser.URL, _ string) []string {
	out := make([]string, 0)

	for key, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return out
}

func Parameters(u *parser.URL, _ string) []string {
	out := make([]string, 0)

	for key := range u.Query() {
		out = append(out, key)
	}

	return out
}

func Values(u *parser.URL, _ string) []string {
	out := make([]string, 0)

	for _, value := range u.Query() {
		out = append(out, value...)
	}

	return out
}
