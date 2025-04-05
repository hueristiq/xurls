package commands

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/hueristiq/hq-go-url/extractor"
	"github.com/hueristiq/xurls/internal/configuration"
	"github.com/spf13/cobra"
	hqgologger "github.com/hueristiq/hq-go-logger"
)

func Extract() (cmd *cobra.Command) {
	var concurrency int

	cmd = &cobra.Command{
		Use:     "extract",
		Aliases: []string{"e"},
		Short:   "Command for extracting URLs from text.",
		Long:    configuration.BANNER(au),
		Run: func(_ *cobra.Command, _ []string) {
			ex := extractor.New()

			regex := ex.CompileRegex()

			lines := make(chan string, concurrency)

			go func() {
				defer close(lines)

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

				if err := scanner.Err(); err != nil {
					hqgologger.Error().Msg(err.Error())
				}
			}()

			wg := &sync.WaitGroup{}

			for range concurrency {
				wg.Add(1)

				go func() {
					defer wg.Done()

					for line := range lines {
						URLs := regex.FindAllString(line, -1)

						for _, URL := range URLs {
							hqgologger.Print().Msg(URL)
						}
					}
				}()
			}

			wg.Wait()
		},
	}

	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", 30, "concurrency")

	return cmd
}
