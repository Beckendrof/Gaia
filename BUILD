load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix beckendrof/gaia
gazelle(name = "gazelle")

go_library(
    name = "gaia_lib",
    srcs = ["gaia.go"],
    data = [".env"],
    importpath = "beckendrof/gaia",
    visibility = ["//visibility:private"],
    deps = [
        "//src/controllers/grpc/apostolis",
        "//src/services/grpc/apostolis",
        "//src/services/metrics",
        "//src/services/mjolnir",
        "//src/utils",
        "@com_github_joho_godotenv//:godotenv",
        "@org_golang_google_grpc//:grpc",
    ],
)

go_binary(
    name = "gaia",
    embed = [":gaia_lib"],
    visibility = ["//visibility:public"],
)

# gazelle:proto disable_global
gazelle(
    name = "update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
        "-build_file_proto_mode=disable_global",
    ],
    command = "update-repos",
)
