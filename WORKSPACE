workspace(name = "asset_tracker")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("//:bazel/external_deps.bzl", "register_build_system_external_repos", "register_deps_external_repos")

register_build_system_external_repos()

# Setup RulesForeignCc
load("@rules_foreign_cc//foreign_cc:repositories.bzl", "rules_foreign_cc_dependencies")

rules_foreign_cc_dependencies()

# Setup Golang
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.22.4")

# gazelle:repo bazel_gazelle

# Setup Gazelle
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//:bazel/go_deps.bzl", "go_deps")

# gazelle:repository_macro bazel/go_deps.bzl%go_deps
go_deps()

gazelle_dependencies()

# Setup Protocol Buffers Rules
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

load("@rules_proto//proto:setup.bzl", "rules_proto_setup")

rules_proto_setup()

load("@rules_proto//proto:toolchains.bzl", "rules_proto_toolchains")

rules_proto_toolchains()

# Setup Protocol Buffers required by Build Tools
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "ae34c344514e08c23e90da0e2d6cb700fcd28e80c02e23e4d5715dddcb42f7b3",
    strip_prefix = "buildtools-4.2.2",
    urls = [
        "https://github.com/bazelbuild/buildtools/archive/refs/tags/4.2.2.tar.gz",
    ],
)

register_deps_external_repos()

go_deps()
