# staticcheck codegen

Currently Bazel `rules_go`'s linter runner nogo does not support multiple analysis.Analyzer
but require 1 analyzer per package.

This repo takes the analyzer exported as a batch from staticcheck and generate dummy packages
so nogo can use it.

## Usages

*NOTE*: this repo depends on [a fix in upstream rules_go](https://github.com/bazelbuild/rules_go/commit/63dfd99403076331fef0775d52a8039d502d4115) that has yet been released by the time of writing (latest rules_go release is v0.28.0). So if you want to try this out, please import `rules_go` via 'master' branch or with a specific git commit

```starlark
    http_archive(
        name = "io_bazel_rules_go",
        strip_prefix = "rules_go-master",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/master.zip",
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/archive/master.zip",
        ],
    )
```

To use this repo, you need to have this and staticcheck imported as `go_repository`

```starlark
    # Note that this repo uses several malfunction packages as testdata
    # so we need to tell Gazelle to skip those packages when generate BUILD files
    go_repository(
        name = "co_honnef_go_tools",
        build_directives = [
            "gazelle:exclude **/testdata/**",  # keep
        ],
        importpath = "honnef.co/go/tools",
        sum = "h1:UoveltGrhghAA7ePc+e+QYDHXrBps2PqFZiHkGR/xK8=",
        version = "v0.0.1-2020.1.4",
    )
```

The code gen binary will generate a list of dependencies to add in your BUILD file that define the `nogo` target

```starlark
govet = [
    "@org_golang_x_tools//go/analysis/passes/asmdecl:go_tool_library",
    "@org_golang_x_tools//go/analysis/passes/assign:go_tool_library",
    "@org_golang_x_tools//go/analysis/passes/atomic:go_tool_library",
    ...
    "@org_golang_x_tools//go/analysis/passes/unreachable:go_tool_library",
    "@org_golang_x_tools//go/analysis/passes/unsafeptr:go_tool_library",
    "@org_golang_x_tools//go/analysis/passes/unusedresult:go_tool_library",
]

staticcheck = [
    "//projects/staticcheck-codegen/_gen/sa1000:go_tool_library",
    ...
    "//projects/staticcheck-codegen/_gen/sa9005:go_tool_library",
    "//projects/staticcheck-codegen/_gen/st1000:go_tool_library",
    ...
    "//projects/staticcheck-codegen/_gen/st1022:go_tool_library",
]

nogo(
    name = "nogo",
    config = "nogo_config.json",
    visibility = ["//visibility:public"],
    deps = govet + staticcheck,
)
```

The code gen binary also provides you with a list of JSON config for you to copy paste to your nogo_config.json file

```json
  ...
  "SA1005": {
    "exclude_files": {
      "external/": "third_party"
    }
  },
  "SA1006": {
    "exclude_files": {
      "external/": "third_party"
    }
  },
  "SA1007": {
    "exclude_files": {
      "external/": "third_party"
    }
  },
  ...
```

To obtain both the deps and the json config, just run `go run .` in this repo to have them printed in stdout.

## TODO

- [ ] Make the code compatible with staticcheck 2021.1.1(v0.2.1): In a recent refactoring, staticcheck stopped exporting the analyzer as a mapped of `checkName -> analyzer` but instead just export them as a slice of analyzer instead. I need to check with staticcheck author to see whether that change can be reversed. If not, we can build a reverse map of `docTitle -> checkName` and get the checkName(i.e. `SA1007`) via matching docTitle.

- [ ] Make a `deps.bzl` and `def.bzl` so that people can use this repository as a bazel dependency out of the box. Preferablly with the list of nogo dependencies exported as constants.
