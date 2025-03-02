package changes

import (
	"fmt"
	"sort"

	"github.com/flaviostutz/gitwho/changes"
	"github.com/flaviostutz/gitwho/utils"
)

func FormatFullTextResults(cResult changes.ChangesResult) string {
	if cResult.TotalCommits == 0 {
		return "No changes found"
	}

	text := fmt.Sprintf("Total authors active: %d\n", len(cResult.AuthorsLines))
	text += fmt.Sprintf("Total files touched: %d\n", cResult.TotalFiles)
	if cResult.TotalLinesTouched.Changes > 0 {
		text += fmt.Sprintf("Average line age when changed: %d days\n", (int(cResult.TotalLinesTouched.AgeDaysSum / float64(cResult.TotalLinesTouched.Changes))))
	}
	text += formatLinesTouched(cResult.TotalLinesTouched, changes.LinesTouched{})

	for _, authorLines := range cResult.AuthorsLines {
		if authorLines.LinesTouched.New+authorLines.LinesTouched.Changes == 0 {
			continue
		}
		mailStr := fmt.Sprintf(" %s", authorLines.AuthorMail)
		text += fmt.Sprintf("\nAuthor: %s%s\n", authorLines.AuthorName, mailStr)
		text += formatLinesTouched(authorLines.LinesTouched, cResult.TotalLinesTouched)
		text += formatTopTouchedFiles(authorLines.FilesTouched)
	}
	return text
}

func FormatTopTextResults(changes changes.ChangesResult) string {
	if changes.TotalCommits == 0 {
		return "No changes found"
	}

	// top coders
	sort.Slice(changes.AuthorsLines, func(i, j int) bool {
		ai := changes.AuthorsLines[i].LinesTouched
		aj := changes.AuthorsLines[j].LinesTouched
		return calcTopCoderScore(ai) > calcTopCoderScore(aj)
	})
	text := "\nTop Coders (new+refactor-churn)\n"
	for i := 0; i < len(changes.AuthorsLines) && i < 3; i++ {
		al := changes.AuthorsLines[i]
		mailStr := fmt.Sprintf(" %s", al.AuthorMail)
		text += fmt.Sprintf("  %s%s: %d%s\n", al.AuthorName, mailStr, calcTopCoderScore(al.LinesTouched), utils.CalcPercStr(calcTopCoderScore(al.LinesTouched), calcTopCoderScore(changes.TotalLinesTouched)))
	}

	// top new liners
	sort.Slice(changes.AuthorsLines, func(i, j int) bool {
		ai := changes.AuthorsLines[i].LinesTouched
		aj := changes.AuthorsLines[j].LinesTouched
		return ai.New > aj.New
	})
	text += "\nTop New Liners\n"
	for i := 0; i < len(changes.AuthorsLines) && i < 3; i++ {
		al := changes.AuthorsLines[i]
		text += fmt.Sprintf("  %s: %d%s\n", al.AuthorName, al.LinesTouched.New, utils.CalcPercStr(al.LinesTouched.New, changes.TotalLinesTouched.New))
	}

	// top refactorers
	sort.Slice(changes.AuthorsLines, func(i, j int) bool {
		ai := changes.AuthorsLines[i].LinesTouched
		aj := changes.AuthorsLines[j].LinesTouched
		return ai.RefactorOther+ai.RefactorOwn > aj.RefactorOther+aj.RefactorOwn
	})
	text += "\nTop Refactorers\n"
	for i := 0; i < len(changes.AuthorsLines) && i < 3; i++ {
		al := changes.AuthorsLines[i]
		text += fmt.Sprintf("  %s: %d%s\n", al.AuthorName, al.LinesTouched.RefactorOther+al.LinesTouched.RefactorOwn, utils.CalcPercStr(al.LinesTouched.RefactorOther+al.LinesTouched.RefactorOwn, changes.TotalLinesTouched.RefactorOther+changes.TotalLinesTouched.RefactorOwn))
	}

	// top helpers
	sort.Slice(changes.AuthorsLines, func(i, j int) bool {
		ai := changes.AuthorsLines[i].LinesTouched
		aj := changes.AuthorsLines[j].LinesTouched
		return ai.ChurnOther > aj.ChurnOther
	})
	text += "\nTop Helpers\n"
	for i := 0; i < len(changes.AuthorsLines) && i < 3; i++ {
		al := changes.AuthorsLines[i]
		text += fmt.Sprintf("  %s: %d%s\n", al.AuthorName, al.LinesTouched.ChurnOther, utils.CalcPercStr(al.LinesTouched.ChurnOther, changes.TotalLinesTouched.ChurnOther))
	}

	// top churners
	sort.Slice(changes.AuthorsLines, func(i, j int) bool {
		ai := changes.AuthorsLines[i].LinesTouched
		aj := changes.AuthorsLines[j].LinesTouched
		return ai.ChurnReceived+ai.ChurnOwn > aj.ChurnReceived+aj.ChurnOwn
	})
	text += "\nTop Churners\n"
	for i := 0; i < len(changes.AuthorsLines) && i < 3; i++ {
		al := changes.AuthorsLines[i]
		text += fmt.Sprintf("  %s: %d%s\n", al.AuthorName, al.LinesTouched.ChurnOwn+al.LinesTouched.ChurnReceived, utils.CalcPercStr(al.LinesTouched.ChurnOwn+al.LinesTouched.ChurnReceived, changes.TotalLinesTouched.ChurnOwn+changes.TotalLinesTouched.ChurnReceived))
	}

	return text
}

func formatTopTouchedFiles(filesTouched []changes.FileTouched) string {
	text := fmt.Sprintf("  - Top files:\n")
	sort.Slice(filesTouched, func(i, j int) bool {
		return filesTouched[i].Lines > filesTouched[j].Lines
	})
	for i := 0; i < len(filesTouched) && i < 5; i++ {
		text += fmt.Sprintf("    - %s (%d)\n", filesTouched[i].Name, filesTouched[i].Lines)
	}
	return text
}

func calcTopCoderScore(ai changes.LinesTouched) int {
	return ai.New + 3*ai.RefactorOther + 2*ai.RefactorOwn - 2*ai.ChurnOwn - 4*ai.ChurnReceived
}

func formatLinesTouched(changes changes.LinesTouched, totals changes.LinesTouched) string {
	totalTouched := changes.New + changes.Changes
	text := fmt.Sprintf("- Total lines touched: %d%s\n", totalTouched, utils.CalcPercStr(changes.New+changes.Changes, totals.New+totals.Changes))
	text += fmt.Sprintf("  - New lines: %d%s\n", changes.New, utils.CalcPercStr(changes.New, totalTouched))
	text += fmt.Sprintf("  - Changed lines: %d%s\n", changes.Changes, utils.CalcPercStr(changes.Changes, totalTouched))
	text += fmt.Sprintf("    - Refactor: %d%s\n", changes.RefactorOwn+changes.RefactorOther, utils.CalcPercStr(changes.RefactorOwn+changes.RefactorOther, changes.Changes))
	text += fmt.Sprintf("      - Refactor of own lines: %d%s\n", changes.RefactorOwn, utils.CalcPercStr(changes.RefactorOwn, changes.RefactorOwn+changes.RefactorOther))
	text += fmt.Sprintf("      - Refactor of other's lines: %d%s\n", changes.RefactorOther, utils.CalcPercStr(changes.RefactorOther, changes.RefactorOwn+changes.RefactorOther))
	text += fmt.Sprintf("      * Refactor done by others to own lines (help received): %d\n", changes.RefactorReceived)
	text += fmt.Sprintf("    - Churn: %d%s\n", changes.ChurnOwn+changes.ChurnOther, utils.CalcPercStr(changes.ChurnOwn+changes.ChurnOther, changes.Changes))
	text += fmt.Sprintf("      - Churn of own lines: %d%s\n", changes.ChurnOwn, utils.CalcPercStr(changes.ChurnOwn, changes.ChurnOwn+changes.ChurnOther))
	text += fmt.Sprintf("      - Churn of other's lines (help given): %d%s\n", changes.ChurnOther, utils.CalcPercStr(changes.ChurnOther, changes.ChurnOwn+changes.ChurnOther))
	text += fmt.Sprintf("      * Churn done by others to own lines (help received): %d\n", changes.ChurnReceived)
	return text
}
