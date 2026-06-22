package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kr/pretty"
)

func Error(message string, e error) {
	var debugLabel string
	errorLabel := labelFormat("error", "error")

	if e != nil {
		debugLabel = labelFormat("debug", "error")
	}

	fmt.Printf("\n%s %s", errorLabel, message)
	fmt.Printf("\n%s %v", debugLabel, e)
	os.Exit(1)
}

func labelFormat(label string, labelType string) string {
	str := "[ " + strings.ToUpper(label) + " ]"
	formattedStr := str

	switch strings.ToLower(labelType) {
	case "error":
		formattedStr = color.HiRedString(str)
	case "command":
		formattedStr = color.HiYellowString(str)
	case "request", "success":
		formattedStr = color.HiGreenString(str)
	case "info":
		formattedStr = color.YellowString(str)
	}

	return formattedStr
}

func LogHttp(method string, statusCode string, url string) {
	formattedMethod := labelFormat(method, "request")
	pretty.Printf("\n%s  %s - %s ", formattedMethod, statusCode, url)
}

func Log(title string, message string) {
	formattedTitle := labelFormat(title, "command")
	pretty.Printf("\n%s %s", formattedTitle, message)
}

func BoolLog(res bool) {
	var label string

	if res {
		label = labelFormat("Command Successful", "success")
	} else {
		label = labelFormat("Command Unsuccessful", "error")
	}

	pretty.Printf("\n%s", label)
}

func Table(header table.Row, body []table.Row, title string) {
	Log("command", "Building table...")
	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.SetTitle(title)
	t.AppendRows(body)
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Render()
}

func ToRow(data []string) table.Row {
	var output table.Row

	if len(data) == 0 {
		Error("Empty ToRow data", nil)
	}

	for _, data := range data {
		output = append(output, data)
	}

	return output
}
