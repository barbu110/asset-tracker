load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "label",
    srcs = ["label.go"],
    importpath = "asset-tracker/src/alm/label",
    visibility = ["//visibility:public"],
    deps = [
        "//src/alm/asset",
        "@com_github_ajstarks_svgo//:svgo",
        "@com_github_boombuler_barcode//:barcode",
        "@com_github_boombuler_barcode//datamatrix",
        "@com_github_vincent_petithory_dataurl//:dataurl",
    ],
)
