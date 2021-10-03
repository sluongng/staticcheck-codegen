// staticcheck-codegen generates a set of packages from 'staticcheck' that can be
// consumed by rules_go's nogo static analysis framework.
//
// Check README.md for usage instructions.
package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Templates
const (
	goTpl = `package {{ .Key }}

import (
	"honnef.co/go/tools/{{ .CheckType }}"
)

var Analyzer = {{ .CheckType }}.Analyzers["{{ .KeyUpper }}"]`

	buildTpl = `# gazelle:ignore
load("@io_bazel_rules_go//go:def.bzl", "go_tool_library")

go_tool_library(
    name = "go_tool_library",
    srcs = ["analyzer.go"],
    importpath = "github.com/sluongng/staticcheck-codegen/_gen/{{ .Key }}",
    visibility = ["//visibility:public"],
    deps = [
        "@co_honnef_go_tools//{{ .CheckType }}",
        "@org_golang_x_tools//go/analysis:go_tool_library",
    ],
)`
)

// Some constants
const (
	PrefixStyleCheck = "ST"
	CodeGenDir       = "_gen"
)

func main() {
	tpl := template.Must(template.New("source").Parse(goTpl))
	buildTpl := template.Must(template.New("source").Parse(buildTpl))

	err := os.RemoveAll(CodeGenDir)
	if err != nil {
		log.Fatalf("os.RemoveAll: %v", err)
	}
	err = os.Mkdir(CodeGenDir, os.ModePerm)
	if err != nil {
		log.Fatalf("os.Mkdir: %v", err)
	}

	analyzerMap := staticcheck.Analyzers
	for k, v := range stylecheck.Analyzers {
		analyzerMap[k] = v
	}

	keys := []string{}
	for k, v := range analyzerMap {
		keys = append(keys, k)
		kUpper := k
		k = strings.ToLower(k)

		err = os.Chdir("_gen")
		if err != nil {
			log.Fatalf("os.Chdir: %v", err)
		}
		err = os.Mkdir(k, os.ModePerm)
		if err != nil {
			log.Fatalf("os.Mkdir: %v", err)
		}
		err = os.Chdir(k)
		if err != nil {
			log.Fatalf("os.Chdir: %v", err)
		}

		// Use stylecheck template instead of
		checkType := "staticcheck"
		if strings.HasPrefix(k, strings.ToLower(PrefixStyleCheck)) {
			checkType = "stylecheck"
		}

		tplFiles := []struct {
			tmplate  *template.Template
			fileName string
		}{
			{tpl, "analyzer.go"},
			{buildTpl, "BUILD.bazel"},
		}

		data := struct {
			Analyzer  *analysis.Analyzer
			Key       string
			KeyUpper  string
			CheckType string
		}{
			Analyzer:  v,
			Key:       k,
			KeyUpper:  kUpper,
			CheckType: checkType,
		}

		// Render the files
		for _, tplFile := range tplFiles {
			outFile, err := os.Create(tplFile.fileName)
			if err != nil {
				log.Fatalf("os.Create: %v", err)
			}

			if err = tplFile.tmplate.Execute(outFile, data); err != nil {
				log.Fatalf("template.Execute failed: %v", err)
			}
		}

		// Return to repo's root
		os.Chdir("../../")
	}

	sort.Strings(keys)

	fmt.Printf("\n\n\tprinting deps for nogo target in BUILD files\n\n")
	for _, s := range keys {
		fmt.Printf(`    "@com_github_sluongng_staticcheck_codegen//_gen/%s:go_tool_library",
`, strings.ToLower(s))
	}

	fmt.Printf("\n\n\tprinting external exclusions for nogo json config\n\n")
	for _, s := range keys {
		fmt.Printf(`  "%s": {
    "exclude_files": {
      "external/": "third_party"
    }
  },
`, s)
	}
}
