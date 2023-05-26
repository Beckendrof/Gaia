load(":cc_toolchain_config.bzl", "cc_toolchain_config")
load("@rules_cc//cc:defs.bzl", "cc_toolchain")

filegroup(
    name = "empty",
    srcs = [],
)

toolchain(
    name = "linux_x86_64_cc",
    target_compatible_with = [
        "@platforms//cpu:x86_64",
        "@platforms//os:linux",
    ],
    toolchain = ":x86_64_linux",
    toolchain_type = "@bazel_tools//tools/cpp:toolchain_type",
)

cc_toolchain(
    name = "x86_64_linux",
    all_files = ":empty",
    ar_files = ":empty",
    as_files = ":empty",
    compiler_files = ":empty",
    dwp_files = ":empty",
    linker_files = ":empty",
    objcopy_files = ":empty",
    strip_files = ":empty",
    toolchain_config = ":linux_x86_64_config",
)

cc_toolchain_config(
    name = "linux_x86_64_config",
    cpu = "local",
    cxx_builtin_include_directories = [
        "/usr/include",
        "/usr/include/c++/9",
        "/usr/include/c++/9/backward",
        "/usr/include/x86_64-linux-gnu",
        "/usr/include/x86_64-linux-gnu/c++/9",
        "/usr/lib/gcc/x86_64-linux-gnu/9/include-fixed",
        "/usr/lib/gcc/x86_64-linux-gnu/9/include",
        "/usr/local/include",
    ],
    cxx_opts = [
        "-D_DEFAULT_SOURCE",
        # Specific for glibc https://en.cppreference.com/w/cpp/types/integer
        "-D__STDC_FORMAT_MACROS",
    ],
    link_opts = [
        "-lm",
        "-fuse-ld=gold",
    ],
    tool_paths = {
        "ar": "/usr/bin/ar",
        "ld": "/usr/bin/gold",
        "cpp": "/usr/bin/gcc",
        "gcc": "/usr/bin/gcc",
        "dwp": "/usr/bin/dwp",
        "gcov": "/usr/bin/gcov",
        "nm": "/usr/bin/nm",
        "objcopy": "/usr/bin/objcopy",
        "objdump": "/usr/bin/objdump",
        "strip": "/usr/bin/strip",
    },
)

toolchain(
    name = "linux_aarch64_cc",
    target_compatible_with = [
        "@platforms//cpu:aarch64",
        "@platforms//os:linux",
    ],
    toolchain = ":aarch64_linux",
    toolchain_type = "@bazel_tools//tools/cpp:toolchain_type",
)

cc_toolchain(
    name = "aarch64_linux",
    all_files = ":empty",
    ar_files = ":empty",
    as_files = ":empty",
    compiler_files = ":empty",
    dwp_files = ":empty",
    linker_files = ":empty",
    objcopy_files = ":empty",
    strip_files = ":empty",
    toolchain_config = ":linux_aarch64_config",
)

cc_toolchain_config(
    name = "linux_aarch64_config",
    cpu = "aarch64",
    cxx_builtin_include_directories = [
        "/usr/aarch64-linux-gnu/include",
        "/usr/aarch64-linux-gnu/include/c++/9",
        "/usr/lib/gcc-cross/aarch64-linux-gnu/9/include",
        "/usr/lib/gcc-cross/aarch64-linux-gnu/9/include-fixed",
        "/usr/aarch64-linux-gnu/include/c++/9/backward",
        "/usr/aarch64-linux-gnu/include",
    ],
    cxx_opts = [
        "-D_DEFAULT_SOURCE",
        # Specific for glibc https://en.cppreference.com/w/cpp/types/integer
        "-D__STDC_FORMAT_MACROS",
    ],
    link_opts = [
        "-lm",
    ],
    tool_paths = {
        "ar": "/usr/bin/aarch64-linux-gnu-ar",
        "ld": "/usr/bin/aarch64-linux-gnu-ld",
        "cpp": "/usr/bin/gcc",
        "gcc": "/usr/bin/gcc",
        "dwp": "/usr/bin/aarch64-linux-gnu-dwp",
        "gcov": "/usr/bin/aarch64-linux-gnu-gcov",
        "nm": "/usr/bin/aarch64-linux-gnu-nm",
        "objcopy": "/usr/bin/aarch64-linux-gnu-objcopy",
        "objdump": "/usr/bin/aarch64-linux-gnu-objdump",
        "strip": "/usr/bin/aarch64-linux-gnu-strip",
    },
)
