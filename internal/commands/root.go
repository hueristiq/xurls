package commands

import (
	"github.com/hueristiq/xurls/internal/configuration"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
	"go.source.hueristiq.com/logger"
	"go.source.hueristiq.com/logger/formatter"
	"go.source.hueristiq.com/logger/levels"
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
			logger.Info().Label("").Msg(configuration.BANNER(au))
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		logger.DefaultLogger.SetFormatter(formatter.NewConsoleFormatter(&formatter.ConsoleFormatterConfiguration{
			Colorize: !monochrome,
		}))

		if verbose {
			logger.DefaultLogger.SetMaxLogLevel(levels.LevelDebug)
		}

		if silent {
			logger.DefaultLogger.SetMaxLogLevel(levels.LevelSilent)
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
		logger.Fatal().Msg(err.Error())
	}
}
