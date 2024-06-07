load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "render.go",
        "root.go",
    ],
    importpath = "asset-tracker/src/alm/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//src/alm/label",
        "//src/core/asset",
        "@com_github_spf13_cobra//:cobra",
    ],
)
