package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	hqgologger "github.com/hueristiq/hq-go-logger"
	hqgourlparser "github.com/hueristiq/hq-go-url/parser"
	"github.com/hueristiq/xurls/internal/configuration"
	"github.com/hueristiq/xurls/internal/input"
	"github.com/spf13/cobra"
)

type Extractor func(URL *hqgourlparser.URL, format string) []string

func Dissect() (cmd *cobra.Command) {
	var (
		inputURLs             []string
		inputURLsListFilePath string
		concurrency           int
		unique                bool
	)

	cmd = &cobra.Command{
		Use:   "dissect",
		Short: "Command for pulling out bits of URLs",
		Long:  configuration.BANNER(au),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				hqgologger.Fatal("mode argument is required")
			}

			mode := args[0]

			fmtStr := ""

			if len(args) > 1 {
				fmtStr = args[1]
			}

			procFn, ok := map[string]Extractor{
				"domains": Domains,
				"apexes":  Apexes,
				"paths":   Paths,
				"query":   Query,
				"params":  Parameters,
				"values":  Values,
				"format":  Format,
			}[mode]

			if !ok {
				hqgologger.Fatal("unknown mode", hqgologger.WithString("mode", mode))
			}

			URLs := make(chan string)

			wg := &sync.WaitGroup{}
			seen := &sync.Map{}

			p := hqgourlparser.New(hqgourlparser.WithDefaultScheme("http"))

			for range concurrency {
				wg.Add(1)

				go func() {
					defer wg.Done()

					for URL := range URLs {
						parsed, err := p.Parse(URL)
						if err != nil {
							hqgologger.Error("parsing failed", hqgologger.WithString("url", URL), hqgologger.WithError(err))

							continue
						}

						for _, value := range procFn(parsed, fmtStr) {
							if value == "" {
								continue
							}

							if unique {
								if _, exists := seen.LoadOrStore(value, struct{}{}); exists {
									continue
								}
							}

							hqgologger.Print(value)
						}
					}
				}()
			}

			go func() {
				defer close(URLs)

				for _, URL := range inputURLs {
					URLs <- URL
				}

				if inputURLsListFilePath != "" {
					file, err := os.Open(inputURLsListFilePath)
					if err != nil {
						hqgologger.Fatal("file open failed", hqgologger.WithString("path", inputURLsListFilePath), hqgologger.WithError(err))
					}

					scanner := bufio.NewScanner(file)

					for scanner.Scan() {
						if url := scanner.Text(); url != "" {
							URLs <- url
						}
					}

					if err := scanner.Err(); err != nil {
						hqgologger.Fatal("file read failed", hqgologger.WithString("path", inputURLsListFilePath), hqgologger.WithError(err))
					}

					file.Close()
				}

				if input.HasStdin() {
					scanner := bufio.NewScanner(os.Stdin)

					for scanner.Scan() {
						if url := scanner.Text(); url != "" {
							URLs <- url
						}
					}

					if err := scanner.Err(); err != nil {
						hqgologger.Fatal("stdin read failed", hqgologger.WithError(err))
					}
				}
			}()

			wg.Wait()
		},
	}

	cmd.Flags().StringSliceVarP(&inputURLs, "url", "u", []string{}, "target URLs (comma separated)")
	cmd.Flags().StringVarP(&inputURLsListFilePath, "list", "l", "", "file containing list of target URLs")
	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", 30, "number of concurrent workers")
	cmd.Flags().BoolVar(&unique, "unique", false, "output only unique values")

	helpTemplate := `Usage:
  {{.CommandPath}} [MODE] [FORMATSTRING] [OPTIONS]

Aliases:
  {{.Aliases}}

Modes:
  domains                   the hostname (e.g. sub.example.com)
  apexes                    the apex domain (e.g. example.com from sub.example.com)
  paths                     the request path (e.g. /users)
  query                     key=value pairs from query string
  params                    keys from query string
  values                    values from query string
  format                    custom format (see below)

Format Directives:
  %%                        literal percent character
  %s                        request scheme (e.g. https)
  %u                        user info (e.g. user:pass)
  %d                        domain (e.g. sub.example.com)
  %S                        subdomain (e.g. sub)
  %r                        root domain (e.g. example)
  %t                        TLD (e.g. com)
  %P                        port (e.g. 8080)
  %p                        path (e.g. /users)
  %e                        file extension (e.g. jpg, html)
  %q                        raw query string (e.g. a=1&b=2)
  %f                        page fragment (e.g. page-section)
  %@                        inserts @ if user info exists
  %:                        inserts colon if port exists
  %?                        inserts ? if query string exists
  %#                        inserts # if fragment exists
  %a                        authority (%u%@%d%:%P)

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Examples:
  {{.CommandPath}} domains -u https://sub.example.com
  {{.CommandPath}} format "Root: %r, TLD: %t" -u https://sub.example.com
  {{.CommandPath}} params -l urls.txt --unique

TIP: Use comma-separated values with -u, multiple -u flags, 
     -l for file input, or stdin for multiple URLs.
`

	cmd.SetUsageTemplate(helpTemplate)

	return cmd
}

func Format(u *hqgourlparser.URL, f string) []string {
	if f == "" {
		return []string{}
	}

	var out bytes.Buffer

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
		case '%':
			out.WriteByte('%')
		case 's':
			out.WriteString(u.Scheme)
		case 'u':
			if u.User != nil {
				out.WriteString(u.User.String())
			}
		case 'd':
			out.WriteString(u.Hostname())
		case 'P':
			out.WriteString(u.Port())
		case 'S':
			out.WriteString(u.Domain.Subdomain)
		case 'r':
			out.WriteString(u.Domain.SecondLevelDomain)
		case 't':
			out.WriteString(u.Domain.TopLevelDomain)
		case 'p':
			out.WriteString(u.EscapedPath())
		case 'e':
			if ext := extractExtension(u.EscapedPath()); ext != "" {
				out.WriteString(ext)
			}
		case 'q':
			out.WriteString(u.RawQuery)
		case 'f':
			out.WriteString(u.Fragment)
		case '@':
			if u.User != nil {
				out.WriteByte('@')
			}
		case ':':
			if u.Port() != "" {
				out.WriteByte(':')
			}
		case '?':
			if u.RawQuery != "" {
				out.WriteByte('?')
			}
		case '#':
			if u.Fragment != "" {
				out.WriteByte('#')
			}
		case 'a':
			out.WriteString(Format(u, "%u%@%d%:%P")[0])
		default:
			out.WriteByte('%')
			out.WriteRune(r)
		}

		inFormat = false
	}

	return []string{out.String()}
}

func extractExtension(path string) string {
	if path == "" {
		return ""
	}

	lastSegment := path

	if idx := strings.LastIndex(path, "/"); idx != -1 {
		lastSegment = path[idx+1:]
	}

	if dotIndex := strings.LastIndex(lastSegment, "."); dotIndex != -1 && dotIndex < len(lastSegment)-1 {
		return lastSegment[dotIndex+1:]
	}

	return ""
}

func Domains(u *hqgourlparser.URL, _ string) []string {
	return []string{u.Hostname()}
}

func Apexes(u *hqgourlparser.URL, _ string) []string {
	if u.Domain.SecondLevelDomain == "" || u.Domain.TopLevelDomain == "" {
		return []string{}
	}

	return []string{u.Domain.SecondLevelDomain + "." + u.Domain.TopLevelDomain}
}

func Paths(u *hqgourlparser.URL, _ string) []string {
	return []string{u.EscapedPath()}
}

func Query(u *hqgourlparser.URL, _ string) []string {
	query := u.Query()

	total := 0

	for _, values := range query {
		total += len(values)
	}

	pairs := make([]string, 0, total)

	for key, vals := range query {
		for _, val := range vals {
			pairs = append(pairs, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return pairs
}

func Parameters(u *hqgourlparser.URL, _ string) []string {
	query := u.Query()

	params := make([]string, 0, len(query))

	for key := range query {
		params = append(params, key)
	}

	return params
}

func Values(u *hqgourlparser.URL, _ string) []string {
	query := u.Query()

	total := 0

	for _, values := range query {
		total += len(values)
	}

	values := make([]string, 0, total)

	for _, vals := range query {
		values = append(values, vals...)
	}

	return values
}
