load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "encoder",
    srcs = [
        "encoder.go",
        "options.go",
    ],
    cgo = True,
    cdepts = ["//external/libwebp"],
    importpath = "github.com/kolesa-team/go-webp/encoder",
    visibility = ["//visibility:public"],
)

alias(
    name = "go_default_library",
    actual = ":encoder",
    visibility = ["//visibility:public"],
)

go_test(
    name = "encoder_test",
    srcs = [
        "encoder_test.go",
        "options_test.go",
    ],
    embed = [":encoder"],
    deps = [
        "@com_github_stretchr_testify//assert:go_default_library",
        "@com_github_stretchr_testify//require:go_default_library",
    ],
)
