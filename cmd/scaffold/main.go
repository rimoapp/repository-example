package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const templatePath = "./templates/scaffold/"

type TemplateData struct {
	ModelName string
	BasePath  string
	TableName string
}

var rootCmd = &cobra.Command{
	Use:   "[camel case model name]",
	Short: "go run cmd/scaffold/main.go CamelCaseModelName",
	Args:  cobra.MinimumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		camel := args[0]
		snake := camelToSnake(camel)
		plural := pluralize(snake)
		data := TemplateData{
			ModelName: camel,
			BasePath:  plural,
			TableName: plural,
		}
		generateFromTemplate("model.go.tpl", fmt.Sprintf("./model/%s.go", snake), data)
		generateFromTemplate("service.go.tpl", fmt.Sprintf("./service/%s.go", snake), data)
		generateFromTemplate("service_test.go.tpl", fmt.Sprintf("./service/%s_test.go", snake), data)
		generateFromTemplate("handler.go.tpl", fmt.Sprintf("./handler/%s.go", snake), data)
		// generateFromTemplate("handler_test.go.tpl", fmt.Sprintf("./handler/%s_test.go", snake), data)
		generateFromTemplate("repository.go.tpl", fmt.Sprintf("./repository/%s.go", snake), data)
	},
}

func camelToSnake(s string) string {
	// Convert to snake_case
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	} else if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	} else {
		return word + "s"
	}
}

func generateFromTemplate(templateFileName, outputFileName string, data TemplateData) {
	tpl, err := template.ParseFiles(templatePath + templateFileName)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(outputFileName); err == nil {
		fmt.Printf("File %s already exists. Overwrite? (y/N): ", outputFileName)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(response)) != "y" {
			fmt.Println("Skipping...")
			return
		}
	}

	outFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = tpl.Execute(outFile, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
