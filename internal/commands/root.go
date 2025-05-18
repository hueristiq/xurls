package commands

import (
	hqgologger "github.com/hueristiq/hq-go-logger"
	hqgologgerformatter "github.com/hueristiq/hq-go-logger/formatter"
	hqgologgerlevels "github.com/hueristiq/hq-go-logger/levels"
	"github.com/hueristiq/xurl/internal/configuration"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

var (
	monochrome bool
	silent     bool
	verbose    bool

	au = aurora.New(aurora.WithColors(true))

	rootCMD = &cobra.Command{
		Use:  configuration.NAME,
		Long: configuration.BANNER(au),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			hqgologger.Info().Label("").Msg(configuration.BANNER(au))
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		hqgologger.DefaultLogger.SetFormatter(hqgologgerformatter.NewConsoleFormatter(&hqgologgerformatter.ConsoleFormatterConfiguration{
			Colorize: !monochrome,
		}))

		if verbose {
			hqgologger.DefaultLogger.SetLevel(hqgologgerlevels.LevelDebug)
		}

		if silent {
			hqgologger.DefaultLogger.SetLevel(hqgologgerlevels.LevelSilent)
		}

		au = aurora.New(aurora.WithColors(!monochrome))
	})

	rootCMD.AddCommand(Extract())
	rootCMD.AddCommand(Parse())

	rootCMD.PersistentFlags().BoolVar(&monochrome, "monochrome", false, "display no color output")
	rootCMD.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "stdout values only output")
	rootCMD.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "stdout verbose output")
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		hqgologger.Fatal().Msg(err.Error())
	}
}
