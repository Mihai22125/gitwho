package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/flaviostutz/gitwho/ownership"
	"github.com/flaviostutz/gitwho/utils"
	"github.com/sirupsen/logrus"
)

func RunOwnershipTimeline(osArgs []string) {
	opts := ownership.OwnershipTimelineOptions{}
	cliOpts := CliOpts{}

	flags := flag.NewFlagSet("ownership-timeline", flag.ExitOnError)
	flags.StringVar(&opts.RepoDir, "repo", ".", "Repository path to analyse")
	flags.StringVar(&opts.Branch, "branch", "main", "Branch name to analyse")
	flags.StringVar(&opts.FilesRegex, "files", ".*", "Regex for selecting which file paths to include in analysis")
	flags.StringVar(&opts.FilesNotRegex, "files-not", "", "Regex for filtering out files from analysis")
	flags.StringVar(&opts.AuthorsRegex, "authors", ".*", "Regex for selecting which authors to include in analysis")
	flags.StringVar(&opts.AuthorsNotRegex, "authors-not", "", "Regex for filtering out authors from analysis")
	flags.StringVar(&opts.Since, "since", "3 months ago", "Starting date for historical analysis. Eg: '1 year ago'")
	flags.StringVar(&opts.Until, "until", "now", "Ending date for historical analysis. Eg: 'now'")
	flags.StringVar(&opts.Period, "period", "2 weeks", "Show ownership data each [period] in the range [since]-[until]. Eg.: '7 days', '1 month'")
	flags.IntVar(&opts.MinDuplicateLines, "min-dup-lines", 4, "Min number of similar lines in a row to be considered a duplicate")
	flags.StringVar(&cliOpts.Format, "format", "full", "Output format. 'full' (more details) or 'short' (lines per author)")
	flags.StringVar(&cliOpts.GoProfileFile, "profile-file", "", "Profile file to dump golang runtime data to")
	flags.BoolVar(&cliOpts.Verbose, "verbose", false, "Show verbose logs during processing")

	flags.Parse(osArgs[2:])

	progressChan := setupBasic(cliOpts)
	defer close(progressChan)

	_, err := utils.ExecCommitsInRange(opts.RepoDir, opts.Branch, "", "")
	if err != nil {
		fmt.Printf("Branch %s not found\n", opts.Branch)
		os.Exit(1)
	}

	logrus.Debugf("Starting analysis of code ownership")
	ownershipResults, err := ownership.TimelineCodeOwnership(opts, progressChan)
	if err != nil {
		fmt.Println("Failed to perform ownership-timeline analysis. err=", err)
		os.Exit(2)
	}

	switch cliOpts.Format {
	case "full":
		ownership.PrintTimelineOwnershipResults(ownershipResults, true)
	case "short":
		ownership.PrintTimelineOwnershipResults(ownershipResults, false)
	case "graph":
		url := ownership.ServeOwnershipTimeline(ownershipResults)
		_, err := utils.ExecShellf("", "open %s", url)
		if err != nil {
			fmt.Printf("Couldn't open browser automatically. See results at %s\n", url)
		}
		fmt.Printf("Serving graph at %s\n", url)
		select {}
	}
}
