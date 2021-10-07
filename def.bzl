load("@io_bazel_rules_go//go:def.bzl", "TOOLS_NOGO", "nogo")

# TOOLS_NOGO_EXCLUDES export a list of checks that should be excluded from rules_go's TOOLS_NOGO
# to get the default behavior of 'go vet'
#
# This is exported mostly for reference/documentation and is not used
TOOLS_NOGO_EXCLUDES = [
    "@org_golang_x_tools//go/analysis/passes/{}:go_default_library".format(package)
    for package in [
        "atomicalign",
        "buildssa",
        "ctrlflow",
        "deepequalerrors",
        "errorsas",
        "fieldalignment",
        "findcall",
        "framepointer",
        "ifaceassert",
        "inspect",
        "internal/analysisutil",
        "nilness",
        "pkgfact",
        "reflectvaluecompare",
        "shadow",
        "sigchanyzer",
        "sortslice",
        "stringintconv",
        "testinggoroutine",
        "unmarshal",
        "unusedwrite",
    ]
]

# GO_VET is a subset of rules_go TOOLS_NOGO
#
# It defines only the checks that happening when you run 'go vet'
# according to 'go tool vet help'.
GO_VET = [
    "@org_golang_x_tools//go/analysis/passes/{}:go_default_library".format(package)
    for package in [
        "asmdecl",
        "assign",
        "atomic",
        "bools",
        "buildtag",
        # disabled by rules_go bug
        # "cgocall",
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
    ]
]

STATICCHECK_ANALYZERS = [
    "@com_github_sluongng_staticcheck_codegen//_gen/{}:go_default_library".format(package)
    for package in [
        "sa1000",
        "sa1001",
        "sa1002",
        "sa1003",
        "sa1004",
        "sa1005",
        "sa1006",
        "sa1007",
        "sa1008",
        "sa1010",
        "sa1011",
        "sa1012",
        "sa1013",
        "sa1014",
        "sa1015",
        "sa1016",
        "sa1017",
        "sa1018",
        "sa1019",
        "sa1020",
        "sa1021",
        "sa1023",
        "sa1024",
        "sa1025",
        "sa1026",
        "sa1027",
        "sa1028",
        "sa1029",
        "sa1030",
        "sa2000",
        "sa2001",
        "sa2002",
        "sa2003",
        "sa3000",
        "sa3001",
        "sa4000",
        "sa4001",
        "sa4003",
        "sa4004",
        "sa4005",
        "sa4006",
        "sa4008",
        "sa4009",
        "sa4010",
        "sa4011",
        "sa4012",
        "sa4013",
        "sa4014",
        "sa4015",
        "sa4016",
        "sa4017",
        "sa4018",
        "sa4019",
        "sa4020",
        "sa4021",
        "sa4022",
        "sa4023",
        "sa4024",
        "sa4025",
        "sa4026",
        "sa4027",
        "sa5000",
        "sa5001",
        "sa5002",
        "sa5003",
        "sa5004",
        "sa5005",
        "sa5007",
        "sa5008",
        "sa5009",
        "sa5010",
        "sa5011",
        "sa5012",
        "sa6000",
        "sa6001",
        "sa6002",
        "sa6003",
        "sa6005",
        "sa9001",
        "sa9002",
        "sa9003",
        "sa9004",
        "sa9005",
        "sa9006",
    ]
]

STYLECHECK_ANALYZERS = [
    "@com_github_sluongng_staticcheck_codegen//_gen/{}:go_default_library".format(package)
    for package in [
        "st1000",
        "st1001",
        "st1003",
        "st1005",
        "st1006",
        "st1008",
        "st1011",
        "st1012",
        "st1013",
        "st1015",
        "st1016",
        "st1017",
        "st1018",
        "st1019",
        "st1020",
        "st1021",
        "st1022",
        "st1023",
    ]
]

# Provide a sensible default for users to get started easily
SENSIBLE_ANALYZERS = GO_VET + STATICCHECK_ANALYZERS + STYLECHECK_ANALYZERS