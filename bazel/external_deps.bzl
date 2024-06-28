load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

def _external_deps_impl(_):
    git_repository(
        name = "libwebp",
        remote = "https://chromium.googlesource.com/webm/libwebp",
        commit = "9ce982fdf21764ef7b273f91d6d72721656c3e03",
        build_file = "//:bazel/libwebp/BUILD.bazel",
    )

external_deps = module_extension(implementation = _external_deps_impl)
