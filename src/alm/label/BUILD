load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "label",
    srcs = [
        "label.go",
        "raster_label.go",
        "utils.go",
    ],
    data = [
        ":data/fonts/CourierNew.ttf",
        ":data/fonts/CourierNewBold.ttf",
    ],
    importpath = "asset-tracker/src/alm/label",
    visibility = ["//visibility:public"],
    deps = [
        "//src/core/asset",
        "@com_github_ajstarks_svgo//:svgo",
        "@com_github_boombuler_barcode//:barcode",
        "@com_github_boombuler_barcode//datamatrix",
        "@com_github_tdewolff_canvas//:canvas",
        "@com_github_tdewolff_canvas//renderers/rasterizer",
        "@com_github_vincent_petithory_dataurl//:dataurl",
    ],
)

go_test(
    name = "label_test",
    srcs = ["utils_test.go"],
    embed = [":label"],
)
