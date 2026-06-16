package abctlx

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kr/pretty"
)

func Error(message string, e error) {
	errorLabel := color.HiRedString("[ ERROR ]")
	debugLabel := color.HiRedString("[ DEBUG ]")
	fmt.Printf("\n [ %s ] from %s", errorLabel, message)
	fmt.Printf("\n [ %s ] %v", debugLabel, e)
	os.Exit(1)
}

func LogHttp(method string, statusCode string, url string) {
	formattedMethod := color.HiGreenString("[ " + method + " ]")
	pretty.Printf("\n%s  [ %s ] - %s ", formattedMethod, statusCode, url)
}

func Log(title string, message string) {
	formattedTitle := color.HiCyanString("[ " + title + " ]")
	pretty.Printf("\n%s %s", formattedTitle, message)
}

func Table(header table.Row, body []table.Row) {
	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(body)
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()
}

func ToRow(data []string) (table.Row, error) {
	var output table.Row

	if len(data) == 0 {
		return nil, fmt.Errorf("Empty data")
	}

	for _, data := range data {
		output = append(output, data)
	}

	return output, nil
}
