package main

import (
	"os"
	"path"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/norwoodj/helm-docs/pkg/document"
	"github.com/norwoodj/helm-docs/pkg/helm"
)

func retrieveInfoAndPrintDocumentation(chartDirectory string, chartSearchRoot string, templateFiles []string, waitGroup *sync.WaitGroup, dryRun bool) {
	defer waitGroup.Done()
	chartDocumentationInfo, err := helm.ParseChartInformation(path.Join(chartSearchRoot, chartDirectory))

	if err != nil {
		log.Warnf("Error parsing information for chart %s, skipping: %s", chartDirectory, err)
		return
	}

	document.PrintDocumentation(chartDocumentationInfo, chartSearchRoot, templateFiles, dryRun, version)

}

func helmDocs(cmd *cobra.Command, _ []string) {
	initializeCli()

	chartSearchRoot := viper.GetString("chart-search-root")
	var fullChartSearchRoot string

	if path.IsAbs(chartSearchRoot) {
		fullChartSearchRoot = chartSearchRoot
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Warnf("Error getting working directory: %s", err)
			return
		}

		fullChartSearchRoot = path.Join(cwd, chartSearchRoot)
	}

	chartDirs, err := helm.FindChartDirectories(fullChartSearchRoot)
	if err != nil {
		log.Errorf("Error finding chart directories: %s", err)
		os.Exit(1)
	}

	log.Infof("Found Chart directories [%s]", strings.Join(chartDirs, ", "))

	templateFiles := viper.GetStringSlice("template-files")
	log.Debugf("Rendering from optional template files [%s]", strings.Join(templateFiles, ", "))

	dryRun := viper.GetBool("dry-run")
	waitGroup := sync.WaitGroup{}

	for _, c := range chartDirs {
		waitGroup.Add(1)
		retrieveInfoAndPrintDocumentation(c, fullChartSearchRoot, templateFiles, &waitGroup, dryRun)
	}

	waitGroup.Wait()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Failed to start the CLI: %s", err)
		os.Exit(1)
	}
}
