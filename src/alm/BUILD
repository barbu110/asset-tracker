load("@rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:build_file_name BUILD,BUILD.bazel

go_binary(
    name = "alm",
    embed = [":alm_lib"],
    importpath = "github.com/barbu110/asset-tracker",
    visibility = ["//visibility:public"],
)

go_library(
    name = "alm_lib",
    srcs = ["main.go"],
    importpath = "asset-tracker/src/alm",
    visibility = ["//visibility:private"],
    deps = ["//src/alm/cmd"],
)
