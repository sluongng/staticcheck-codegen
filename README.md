# staticcheck codegen

Currently Bazel `rules_go`'s linter runner nogo does not support multiple analysis.Analyzer
but require 1 analyzer per package.

This repo takes the analyzer exported as a batch from staticcheck and generate dummy packages
so nogo can use it.

This is very much a WIP.

