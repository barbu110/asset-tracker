load("@buildifier_prebuilt//:rules.bzl", "buildifier")
load("@gazelle//:def.bzl", "gazelle")

BUILDIFIER_EXCLUDE_PATTERNS = [
    "./.git/*",
    "./.idea/*",
]

buildifier(
    name = "buildifier.check",
    exclude_patterns = BUILDIFIER_EXCLUDE_PATTERNS,
    lint_mode = "warn",
    mode = "diff",
)

buildifier(
    name = "buildifier",
    exclude_patterns = BUILDIFIER_EXCLUDE_PATTERNS,
    mode = "fix",
)

# gazelle:prefix github.com/barbu110/asset-tracker

gazelle(name = "gazelle")
