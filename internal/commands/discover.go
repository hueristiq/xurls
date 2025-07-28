package commands

import (
	"bufio"
	"os"
	"strings"
	"sync"

	hqgologger "github.com/hueristiq/hq-go-logger"
	hqgourlextractor "github.com/hueristiq/hq-go-url/extractor"
	"github.com/hueristiq/xurls/internal/configuration"
	"github.com/spf13/cobra"
)

func Discover() (cmd *cobra.Command) {
	var (
		concurrency       int
		withScheme        bool
		withSchemePattern string
		withHost          bool
		withHostPattern   string
	)

	cmd = &cobra.Command{
		Use:   "discover",
		Short: "Command for extracting URLs from text",
		Long:  configuration.BANNER(au),
		Run: func(_ *cobra.Command, _ []string) {
			options := []hqgourlextractor.Option{}

			if withScheme {
				options = append(options, hqgourlextractor.WithScheme())
			}

			if withSchemePattern != "" {
				options = append(options, hqgourlextractor.WithSchemePattern(withSchemePattern))
			}

			if withHost {
				options = append(options, hqgourlextractor.WithHost())
			}

			if withHostPattern != "" {
				options = append(options, hqgourlextractor.WithHostPattern(withHostPattern))
			}

			e := hqgourlextractor.New(options...)
			r := e.CompileRegex()

			lines := make(chan string, concurrency)
			wg := &sync.WaitGroup{}

			for range concurrency {
				wg.Add(1)

				go func() {
					defer wg.Done()

					for line := range lines {
						URLs := r.FindAllString(line, -1)

						for _, URL := range URLs {
							hqgologger.Print(URL)
						}
					}
				}()
			}

			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				line := scanner.Text()

				replacer := strings.NewReplacer(
					"*", "",
					`\u002f`, "/",
					`\u0026`, "&",
				)

				line = replacer.Replace(line)
				if line != "" {
					lines <- line
				}
			}

			close(lines)

			if err := scanner.Err(); err != nil {
				hqgologger.Fatal("input scanning failed", hqgologger.WithError(err))
			}

			wg.Wait()
		},
	}

	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", 30, "number of concurrent workers")
	cmd.Flags().BoolVar(&withScheme, "with-scheme", false, "match URLs with schemes")
	cmd.Flags().StringVar(&withSchemePattern, "with-scheme-pattern", "", "match URLs with scheme pattern")
	cmd.Flags().BoolVar(&withHost, "with-host", false, "match URLs with hosts")
	cmd.Flags().StringVar(&withHostPattern, "with-host-pattern", "", "match URLs with host pattern")

	return cmd
}
