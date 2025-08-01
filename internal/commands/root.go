package commands

import (
	hqgologger "github.com/hueristiq/hq-go-logger"
	hqgologgerformatter "github.com/hueristiq/hq-go-logger/formatter"
	hqgologgerlevels "github.com/hueristiq/hq-go-logger/levels"
	"github.com/hueristiq/xurls/internal/configuration"
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
			hqgologger.Info(configuration.BANNER(au), hqgologger.WithLabel(""))
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		hqgologger.DefaultLogger.SetFormatter(hqgologgerformatter.NewConsoleFormatter(&hqgologgerformatter.ConsoleFormatterConfiguration{
			Colorize: !monochrome,
		}))

		if silent {
			hqgologger.DefaultLogger.SetLevel(hqgologgerlevels.LevelSilent)
		}

		if verbose {
			hqgologger.DefaultLogger.SetLevel(hqgologgerlevels.LevelDebug)
		}

		au = aurora.New(aurora.WithColors(!monochrome))
	})

	rootCMD.AddCommand(Discover())
	rootCMD.AddCommand(Dissect())

	rootCMD.PersistentFlags().BoolVar(&monochrome, "monochrome", false, "disable colored console output")
	rootCMD.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "disable logging output, only results")
	rootCMD.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable detailed debug logging output")
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		hqgologger.Fatal("failed!", hqgologger.WithError(err))
	}
}
