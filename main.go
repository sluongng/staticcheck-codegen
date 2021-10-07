// staticcheck-codegen generates a set of packages from 'staticcheck' that can be
// consumed by rules_go's nogo static analysis framework.
//
// Check README.md for usage instructions.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Templates
const (
	goTpl = `package {{ .Key }}

import (
	"honnef.co/go/tools/{{ .CheckType }}"
)

var Analyzer = {{ .CheckType }}.Analyzers[{{ .Index }}].Analyzer`
)

// Some constants
const (
	PrefixStyleCheck  = "ST"
	PrefixStaticCheck = "SA"
	PrefixQuickFix    = "QF"
	CodeGenDir        = "_gen"
)

type config struct {
	// onlyFiles is a list of regular expressions that match files an analyzer
	// will emit diagnostics for. When empty, the analyzer will emit diagnostics
	// for all files.
	OnlyFiles map[string]string `json:"only_files,omitempty"`

	// excludeFiles is a list of regular expressions that match files that an
	// analyzer will not emit diagnostics for.
	ExcludeFiles map[string]string `json:"exclude_files,omitempty"`
}

func main() {
	tpl := template.Must(template.New("source").Parse(goTpl))

	err := os.RemoveAll(CodeGenDir)
	if err != nil {
		log.Fatalf("os.RemoveAll: %v", err)
	}
	err = os.Mkdir(CodeGenDir, os.ModePerm)
	if err != nil {
		log.Fatalf("os.Mkdir: %v", err)
	}

	keys := []string{}
	for _, analyzers := range [][]*lint.Analyzer{
		staticcheck.Analyzers,
		stylecheck.Analyzers,
		simple.Analyzers,
		quickfix.Analyzers,
	} {
		sort.Slice(analyzers, func(i, j int) bool {
			return analyzers[i].Analyzer.Name > analyzers[j].Analyzer.Name
		})
		for i, v := range analyzers {
			k := v.Analyzer.Name
			keys = append(keys, k)
			kUpper := k
			k = strings.ToLower(k)

			err = os.Chdir(CodeGenDir)
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
			checkType := "simple"
			if strings.HasPrefix(k, strings.ToLower(PrefixStaticCheck)) {
				checkType = "staticcheck"
			}
			if strings.HasPrefix(k, strings.ToLower(PrefixStyleCheck)) {
				checkType = "stylecheck"
			}
			if strings.HasPrefix(k, strings.ToLower(PrefixQuickFix)) {
				checkType = "quickfix"
			}

			tplFiles := []struct {
				tmplate  *template.Template
				fileName string
			}{
				{tpl, "analyzer.go"},
			}

			data := struct {
				Key       string
				KeyUpper  string
				CheckType string
				Index     int
			}{
				Index:     i,
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
	}

	sort.Strings(keys)

	govetChecks := []string{
		"asmdecl",
		"assign",
		"atomic",
		"bools",
		"buildtag",
		"composite",
		"copylock",
		"httpresponse",
		"loopclosure",
		"lostcancel",
		"nilfunc",
		"printf",
		"shift",
		"stdmethods",
		"structtag",
		"tests",
		"unreachable",
		"unsafeptr",
		"unusedresult",
	}
	temp := make(map[string]*config)
	for _, v := range append(govetChecks, keys...) {
		c := &config{ExcludeFiles: map[string]string{"external/": "third_party"}}

		if v == "ST1000" {
			c.ExcludeFiles["rules_go_work"] = "rules_go generated code"
		}

		if v == "composite" {
			v = "composites"
		}

		if v == "copylock" {
			v = "copylocks"
		}

		temp[v] = c
	}
	b, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		log.Fatalf("json marshal: %v", err)
	}

	err = os.Remove("nogo_config.json")
	if err != nil {
		log.Fatalf("rm nogo config: %v", err)
	}
	cf, err := os.Create("nogo_config.json")
	if err != nil {
		log.Fatalf("create nogo config: %v", err)
	}
	fmt.Fprint(cf, string(b))
}
