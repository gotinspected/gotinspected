package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/yourusername/gitinspector-go/internal/analyzer"
)

type TablePrinter struct{}

func NewTablePrinter() *TablePrinter { return &TablePrinter{} }

func (t *TablePrinter) Print(result *analyzer.Result) {
	if len(result.Authors) == 0 {
		fmt.Println("No activity found for the specified criteria.")
		return
	}

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.SetStyle(table.StyleRounded)
	tw.Style().Title.Align = text.AlignCenter
	tw.Style().Title.Colors = text.Colors{text.FgHiCyan}
	tw.Style().Color.Header = text.Colors{text.FgHiWhite}

	// Add slight spacing between rows since we are using multiline author cells now
	tw.Style().Options.SeparateRows = true

	tw.SetTitle("📊 GITINSPECTOR CODE IMPACT")
	tw.AppendHeader(table.Row{"Author", "Commits", "Raw (+/-)", "Code (+/-)", "Current Lines", "Noise Breakdown", "Impact Share"})

	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft},
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter},
		{Number: 5, Align: text.AlignRight},
		{Number: 6, Align: text.AlignRight},
		{Number: 7, Align: text.AlignLeft},
	})

	cyan := text.Colors{text.FgHiCyan}
	green := text.Colors{text.FgGreen}
	red := text.Colors{text.FgRed}
	yellow := text.Colors{text.FgHiYellow}
	gray := text.Colors{text.FgHiBlack}

	for _, author := range result.Authors {
		sigChanges := author.Insertions + author.Deletions
		rawChanges := author.RawInsertions + author.RawDeletions

		impactPercent, noisePercent := 0.0, 0.0
		if result.TotalChanges > 0 {
			impactPercent = float64(sigChanges) / float64(result.TotalChanges) * 100
		}
		if rawChanges > 0 {
			noisePercent = float64(rawChanges-sigChanges) / float64(rawChanges) * 100
		}

		bar := renderProgressBar(impactPercent, 12)

		// NEW: Format the author column to show "Name \n <email>"
		authorStr := fmt.Sprintf("%s\n%s", author.Name, gray.Sprintf("<%s>", author.Email))

		commitsStr := cyan.Sprintf("%d", author.Commits)
		rawStr := gray.Sprintf("+%d / -%d", author.RawInsertions, author.RawDeletions)
		sigStr := fmt.Sprintf("%s / %s", green.Sprintf("+%d", author.Insertions), red.Sprintf("-%d", author.Deletions))
		currentStr := green.Sprintf("%d", author.CurrentLines)

		noiseStr := gray.Sprintf("%.1f%% (💬 %d | ∅ %d)", noisePercent, author.Comments, author.EmptyLines)
		shareStr := fmt.Sprintf("%s %s", bar, yellow.Sprintf("%.1f%%", impactPercent))

		tw.AppendRow(table.Row{
			authorStr, commitsStr, rawStr, sigStr, currentStr, noiseStr, shareStr,
		})
	}

	tw.AppendSeparator()

	totalNoise := 0.0
	if result.TotalRawChanges > 0 {
		totalNoise = float64(result.TotalRawChanges-result.TotalChanges) / float64(result.TotalRawChanges) * 100
	}

	tw.AppendFooter(table.Row{
		"TOTALS",
		cyan.Sprint(result.TotalCommits),
		gray.Sprintf("+%d / -%d", result.TotalRawIns, result.TotalRawDel),
		fmt.Sprintf("%s / %s", green.Sprintf("+%d", result.TotalIns), red.Sprintf("-%d", result.TotalDel)),
		green.Sprint(result.TotalCurrent),
		gray.Sprintf("%.1f%% (💬 %d | ∅ %d)", totalNoise, result.TotalComments, result.TotalEmpty),
		yellow.Sprint("100.0%"),
	})

	fmt.Println()
	tw.Render()
	fmt.Println()
}

func renderProgressBar(percent float64, width int) string {
	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	var bar strings.Builder
	for i := 0; i < width; i++ {
		if i < filled {
			bar.WriteString("█")
		} else {
			bar.WriteString("░")
		}
	}
	return text.Colors{text.FgHiBlue}.Sprint(bar.String())
}
