load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "cmd",
    srcs = [
        "render.go",
        "root.go",
        "save_asset.go",
    ],
    importpath = "asset-tracker/pkg/alm/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/alm/label",
        "//pkg/asset_manager",
        "//pkg/core/asset",
        "@com_github_spf13_cobra//:cobra",
    ],
)
