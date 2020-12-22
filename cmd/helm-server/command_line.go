package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/norwoodj/helm-docs/pkg/document"
	"github.com/norwoodj/helm-docs/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	rootCmd = &cobra.Command{
		Use:     "helm-docs",
		Short:   "helm-docs automatically generates markdown documentation for helm charts from requirements and values files",
		Version: version,
	}
)

func possibleLogLevels() []string {
	levels := make([]string, 0)

	for _, l := range log.AllLevels {
		levels = append(levels, l.String())
	}

	return levels
}

func initializeCli() {
	logLevelName := viper.GetString("log-level")
	logLevel, err := log.ParseLevel(logLevelName)
	if err != nil {
		log.Errorf("Failed to parse provided log level %s: %s", logLevelName, err)
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(logLevel)
}

func init() {
	serverCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the helm-doc server and listen on 0.0.0.0:8080",
		Long:  `Runs a helm-doc server for querying the loaded charts via a RESTful API`,
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}

	rootCmd.AddCommand(serverCmd)

	logLevelUsage := fmt.Sprintf("Level of logs that should printed, one of (%s)", strings.Join(possibleLogLevels(), ", "))
	rootCmd.PersistentFlags().StringP("chart-search-root", "c", ".", "directory to search recursively within for charts")
	rootCmd.PersistentFlags().StringP("ignore-file", "i", ".helmdocsignore", "The filename to use as an ignore file to exclude chart directories")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", logLevelUsage)
	rootCmd.PersistentFlags().StringP("sort-values-order", "s", document.AlphaNumSortOrder, fmt.Sprintf("order in which to sort the values table (\"%s\" or \"%s\")", document.AlphaNumSortOrder, document.FileSortOrder))
	rootCmd.PersistentFlags().StringSliceP("template-files", "t", []string{"README.md.gotmpl"}, "gotemplate file paths relative to each chart directory from which documentation will be generated")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HELM_DOCS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
