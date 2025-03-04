package ownership

import (
	"fmt"
	"strconv"

	"github.com/flaviostutz/gitwho/ownership"
)

type authorLinesDate struct {
	date        string
	authorLines ownership.AuthorLines
}

func FormatCodeOwnershipResults(ownershipResult ownership.OwnershipResult, full bool) string {
	text := fmt.Sprintf("\nTotal authors: %d\n", len(ownershipResult.AuthorsLines))
	text += fmt.Sprintf("Total files: %d\n", ownershipResult.TotalFiles)
	if full {
		text += fmt.Sprintf("Avg line age: %s\n", avgLineAgeStr(ownershipResult.LinesAgeDaysSum, ownershipResult.TotalLines))
		text += fmt.Sprintf("Duplicated lines: %d (%d%%)\n", ownershipResult.TotalLinesDuplicated, int(100*float64(ownershipResult.TotalLinesDuplicated)/float64(ownershipResult.TotalLines)))
	}
	text += fmt.Sprintf("Total lines: %d\n", ownershipResult.TotalLines)
	for _, authorLines := range ownershipResult.AuthorsLines {
		mailStr := ""
		additional := ""
		if full {
			additional = fmt.Sprintf(" avg-days:%d dup:%d orig:%d dup-others:%d",
				int((authorLines.OwnedLinesAgeDaysSum / float64(authorLines.OwnedLinesTotal))),
				authorLines.OwnedLinesDuplicate,
				authorLines.OwnedLinesDuplicateOriginal,
				authorLines.OwnedLinesDuplicateOriginalOthers)
			mailStr = fmt.Sprintf(" %s", authorLines.AuthorMail)
		}
		text += fmt.Sprintf("  %s%s: %d (%s%%)%s\n",
			authorLines.AuthorName,
			mailStr,
			authorLines.OwnedLinesTotal,
			strconv.FormatFloat(float64(100)*(float64(authorLines.OwnedLinesTotal)/float64(ownershipResult.TotalLines)), 'f', 1, 32),
			additional)
	}

	return text
}

func FormatDuplicatesResults(ownershipResult ownership.OwnershipResult, full bool) string {
	text := fmt.Sprintf("Total lines: %d\n", ownershipResult.TotalLines)
	text += fmt.Sprintf("Duplicated lines: %d (%d%%)\n", ownershipResult.TotalLinesDuplicated, int(100*float64(ownershipResult.TotalLinesDuplicated)/float64(ownershipResult.TotalLines)))
	counter := 0
	for _, lineGroup := range ownershipResult.DuplicateLineGroups {
		text += fmt.Sprintf("%s:%d - %d\n", lineGroup.FilePath, lineGroup.LineNumber, lineGroup.LineNumber+lineGroup.LineCount)
		for _, relatedGroup := range lineGroup.RelatedLinesGroup {
			text += fmt.Sprintf("  %s:%d - %d\n", relatedGroup.FilePath, relatedGroup.LineNumber, relatedGroup.LineNumber+relatedGroup.LineCount)
			if !full {
				counter++
				if counter > 20 {
					text += "...(use --format \"full\" for more results)\n"
					return text
				}
			}
		}
	}
	return text
}

func avgLineAgeStr(linesAgeDaysSum float64, totalLines int) string {
	return fmt.Sprintf("%1.f days", (linesAgeDaysSum / float64(totalLines)))
}
