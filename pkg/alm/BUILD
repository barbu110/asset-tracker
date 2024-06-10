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
    importpath = "asset-tracker/pkg/alm",
    visibility = ["//visibility:private"],
    deps = ["//pkg/alm/cmd"],
)

sh_binary(
    name = "create_asset",
    srcs = ["create_asset.sh"],
    data = [
        ":alm",
    ],
    env = {
        "ALM_EXEC": "$(location :alm)",
    },
    deps = ["@bazel_tools//tools/bash/runfiles"],
)
