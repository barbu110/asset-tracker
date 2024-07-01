load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def register_deps_external_repos():
    git_repository(
        name = "libwebp",
        remote = "https://chromium.googlesource.com/webm/libwebp",
        commit = "9ce982fdf21764ef7b273f91d6d72721656c3e03",
        build_file = "//:bazel/libwebp/BUILD.bazel",
    )

def register_build_system_external_repos():
    git_repository(
        name = "rules_foreign_cc",
        remote = "https://github.com/bazelbuild/rules_foreign_cc.git",
        commit = "0ed9aaa68282f8a7de56ae4f96191891a75d4dfb"
    )

    http_archive(
        name = "io_bazel_rules_go",
        sha256 = "b2038e2de2cace18f032249cb4bb0048abf583a36369fa98f687af1b3f880b26",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.48.1/rules_go-v0.48.1.zip",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.48.1/rules_go-v0.48.1.zip",
        ],
    )

    http_archive(
        name = "bazel_gazelle",
        integrity = "sha256-12v3pg/YsFBEQJDfooN6Tq+YKeEWVhjuNdzspcvfWNU=",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
        ],
    )

    http_archive(
        name = "rules_proto",
        sha256 = "6fb6767d1bef535310547e03247f7518b03487740c11b6c6adb7952033fe1295",
        strip_prefix = "rules_proto-6.0.2",
        url = "https://github.com/bazelbuild/rules_proto/releases/download/6.0.2/rules_proto-6.0.2.tar.gz",
    )

    # Required by Build Tools below.
    http_archive(
        name = "com_google_protobuf",
        integrity = "sha256-5P8q63Z9pvT1JIXC5yRolg3f5SYkg4ee9q1VLlJ1enc=",
        strip_prefix = "protobuf-27.2",
        url = "https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protobuf-27.2.tar.gz",
    )
