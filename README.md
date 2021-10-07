# staticcheck codegen

Currently Bazel `rules_go`'s linter runner nogo does not support multiple analysis.Analyzer
but require 1 analyzer per package.

This repo takes the analyzer exported as a batch from staticcheck and generate dummy packages
so nogo can use it.

## Usages

Minimum rules_go v0.29.0 is required

```starlark
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "2b1641428dff9018f9e85c0384f03ec6c10660d935b750e3fa1492a281a53b0f",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
    ],
)
```

To use this repo, you can create a `hack/tools/tools.go` file with a specific build tag `tools` so that it is not compiled by default.  In there, you can import this repository as follow.

```golang
//go:build tools
// +build tools

// Package tools is a place holder for all golang binary toolings
// needed to maintain the repository health but is not a compilation dependency.
//
// The package is kept away from standard compilations by using a specific
// 'tools' build tag. Imported packages therefore could be 'main' package of
// other modules as this tools package will never be compiled.
package tools

import (
	_ "github.com/sluongng/staticcheck-codegen"
)
```

This file could be coupled with an empty `hack/tools/BUILD.bazel` file as follow so that gazelle will skip generating golang target for this `tools` package.

```starlark
# gazelle:ignore
```

Then in your workspace, populate your `go.mod` an `go.sum` file as follow

```
> bazel run @go_sdk//:bin/go mod tidy
```

Then update your `go_repository` targets with Gazelle

```
> bazel run //:gazelle -- update-repos -from_file=go.mod <extra flags here>
```

This should give you 2 important `go_repository` targets: `co_honnef_go_tools` and `com_github_sluongng_staticcheck_codegen` as follow

```starlark
# Note that this repo uses several malfunction packages as testdata
# so we need to tell Gazelle to skip those packages when generate BUILD files
#
# Don't forget to update this to latest 'version' and matching 'sum'
go_repository(
    name = "co_honnef_go_tools",
    build_directives = [
        "gazelle:exclude **/testdata/**",  # keep
    ],
    importpath = "honnef.co/go/tools",
    sum = "<some-hash>",
    version = "<staticcheck-version>",
)
go_repository(
    name = "com_github_sluongng_staticcheck_codegen",
    importpath = "github.com/sluongng/staticcheck-codegen",
    sum = "<some-hash>",
    version = "<latest-version>",
)
```

Then you can setup your `nogo` target in your build file as follow

```starlark
load("@io_bazel_rules_go//go:def.bzl", "nogo")
load("@com_github_sluongng_staticcheck_codegen//:def.bzl", "SENSIBLE_ANALYZERS")

nogo(
    name = "nogo",
    config = "nogo_config.json",
    visibility = ["//visibility:public"],
    deps = SENSIBLE_ANALYZERS,
)
```

Note that you need to provides your own `nogo_config.json` file.
This repo does provide a default config file but it's not recommended to use it
from this repo directly as `@com_github_sluongng_staticcheck_codegen//:nogo_config.json` because
there will be needs to add in your own config/allowlist in the future.

For this reason, go to [nogo_config.json](./nogo_config.json) file and copy it to your repository.
Then you can customize it to your liking (i.e. allow excluding some paths from certain checks).

## TODO

- [ ] Implement automatic generation/update of def.bzl by using buildtools library.

- [ ] Tidy up the code

