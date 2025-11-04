package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

func ShowFeatureTable(features []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Feature"})
	for i, f := range features {
		t.AppendRow(table.Row{i + 1, f})
	}
	color.Cyan("\nAvailable Vault Features:\n")
	t.Render()
}

func StartProgress(label string) *progressbar.ProgressBar {
	fmt.Printf("\n%s\n", color.BlueString("🚀 "+label))
	bar := progressbar.NewOptions(50,
		progressbar.OptionSetDescription("Processing..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetRenderBlankState(true),
	)
	return bar
}

func ShowSuccessBox(title, message string) {
	fmt.Println(color.GreenString("\n" + title))
	fmt.Println("───────────────────────────────")
	color.White(message)
	fmt.Println("───────────────────────────────\n")
}

func ShowErrorBox(title, message string) {
	fmt.Println(color.RedString("\n" + title))
	fmt.Println("───────────────────────────────")
	color.White(message)
	fmt.Println("───────────────────────────────\n")
}

func ShowFileDetails(name, hash string, size int64, mode, location string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})
	t.AppendRows([]table.Row{
		{"File Name", name},
		{"Hash", hash},
		{"Size (bytes)", size},
		{"Mode", mode},
		{"Location", location},
		{"Timestamp", time.Now().Format(time.RFC3339)},
	})
	color.Cyan("\nFile Details:\n")
	t.Render()
}
