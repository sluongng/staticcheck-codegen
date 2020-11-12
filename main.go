package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
)

const analyzersTpl = `
package {{ .Key }}

import (
	"honnef.co/go/tools/staticcheck"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = staticcheck.Analyzers["{{ .Key }}"]
`

func main() {
	tpl := template.Must(template.New("source").Parse(analyzersTpl))

	err := os.RemoveAll("_gen")
	if err != nil {
		log.Fatalf("os.RemoveAll: %v", err)
	}
	err = os.Mkdir("_gen", os.ModePerm)
	if err != nil {
		log.Fatalf("os.Mkdir: %v", err)
	}

	for k, v := range staticcheck.Analyzers {
		k = strings.ToLower(k)

		err = os.Chdir("_gen")
		if err != nil {
			log.Fatalf("os.Chdir: %v", err)
		}

		log.Println("Creating " + k)
		err = os.Mkdir(k, os.ModePerm)
		if err != nil {
			log.Fatalf("os.Mkdir: %v", err)
		}
		err = os.Chdir(k)
		if err != nil {
			log.Fatalf("os.Chdir: %v", err)
		}

		outFile, err := os.Create("analyzer.go")
		if err != nil {
			log.Fatalf("os.Create: %v", err)
		}

		data := struct {
			Key      string
			Analyzer *analysis.Analyzer
		}{
			Key:      k,
			Analyzer: v,
		}
		if err = tpl.Execute(outFile, data); err != nil {
			log.Fatalf("template.Execute failed: %v", err)
		}

		os.Chdir("../../")
	}
}
