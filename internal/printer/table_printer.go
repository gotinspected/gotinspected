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

func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}

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
	
	// THE FIX: Header colors are set via the Color struct, not directly on the Style struct.
	tw.Style().Color.Header = text.Colors{text.FgHiWhite}
	tw.Style().Options.SeparateRows = false

	tw.SetTitle("📊 GITINSPECTOR CODE IMPACT")
	tw.AppendHeader(table.Row{"Author", "Commits", "Insertions (+)", "Deletions (-)", "Impact Share"})

	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignRight},
		{Number: 3, Align: text.AlignRight},
		{Number: 4, Align: text.AlignRight},
		{Number: 5, Align: text.AlignLeft},
	})

	cyan := text.Colors{text.FgHiCyan}
	green := text.Colors{text.FgGreen}
	red := text.Colors{text.FgRed}
	yellow := text.Colors{text.FgHiYellow}

	for _, author := range result.Authors {
		changes := author.Insertions + author.Deletions
		percent := 0.0
		if result.TotalChanges > 0 {
			percent = float64(changes) / float64(result.TotalChanges) * 100
		}

		bar := renderProgressBar(percent, 12)

		commitsColored := cyan.Sprintf("%d", author.Commits)
		insColored := green.Sprintf("+%d", author.Insertions)
		delColored := red.Sprintf("-%d", author.Deletions)
		shareColored := fmt.Sprintf("%s %s", bar, yellow.Sprintf("%.1f%%", percent))

		tw.AppendRow(table.Row{
			author.Name,
			commitsColored,
			insColored,
			delColored,
			shareColored,
		})
	}

	tw.AppendSeparator()
	tw.AppendFooter(table.Row{
		"TOTALS",
		cyan.Sprint(result.TotalCommits),
		green.Sprintf("+%d", result.TotalIns),
		red.Sprintf("-%d", result.TotalDel),
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

	blue := text.Colors{text.FgHiBlue}
	return blue.Sprint(bar.String())
}
