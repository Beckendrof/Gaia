# Copyright 2019 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""A Starlark cc_toolchain configuration rule"""

load(
    "@rules_cc//cc:action_names.bzl",
    _ASSEMBLE_ACTION_NAME = "ASSEMBLE_ACTION_NAME",
    _CLIF_MATCH_ACTION_NAME = "CLIF_MATCH_ACTION_NAME",
    _CPP_COMPILE_ACTION_NAME = "CPP_COMPILE_ACTION_NAME",
    _CPP_HEADER_PARSING_ACTION_NAME = "CPP_HEADER_PARSING_ACTION_NAME",
    _CPP_LINK_DYNAMIC_LIBRARY_ACTION_NAME = "CPP_LINK_DYNAMIC_LIBRARY_ACTION_NAME",
    _CPP_LINK_EXECUTABLE_ACTION_NAME = "CPP_LINK_EXECUTABLE_ACTION_NAME",
    _CPP_LINK_NODEPS_DYNAMIC_LIBRARY_ACTION_NAME = "CPP_LINK_NODEPS_DYNAMIC_LIBRARY_ACTION_NAME",
    _CPP_MODULE_CODEGEN_ACTION_NAME = "CPP_MODULE_CODEGEN_ACTION_NAME",
    _CPP_MODULE_COMPILE_ACTION_NAME = "CPP_MODULE_COMPILE_ACTION_NAME",
    _C_COMPILE_ACTION_NAME = "C_COMPILE_ACTION_NAME",
    _LINKSTAMP_COMPILE_ACTION_NAME = "LINKSTAMP_COMPILE_ACTION_NAME",
    _LTO_BACKEND_ACTION_NAME = "LTO_BACKEND_ACTION_NAME",
    _PREPROCESS_ASSEMBLE_ACTION_NAME = "PREPROCESS_ASSEMBLE_ACTION_NAME",
)
load(
    "@rules_cc//cc:cc_toolchain_config_lib.bzl",
    "action_config",
    "feature",
    "flag_group",
    "flag_set",
    "tool",
    "tool_path",
    "with_feature_set",
)

all_compile_actions = [
    _C_COMPILE_ACTION_NAME,
    _CPP_COMPILE_ACTION_NAME,
    _LINKSTAMP_COMPILE_ACTION_NAME,
    _ASSEMBLE_ACTION_NAME,
    _PREPROCESS_ASSEMBLE_ACTION_NAME,
    _CPP_HEADER_PARSING_ACTION_NAME,
    _CPP_MODULE_COMPILE_ACTION_NAME,
    _CPP_MODULE_CODEGEN_ACTION_NAME,
    _CLIF_MATCH_ACTION_NAME,
    _LTO_BACKEND_ACTION_NAME,
]

all_cpp_compile_actions = [
    _CPP_COMPILE_ACTION_NAME,
    _LINKSTAMP_COMPILE_ACTION_NAME,
    _CPP_HEADER_PARSING_ACTION_NAME,
    _CPP_MODULE_COMPILE_ACTION_NAME,
    _CPP_MODULE_CODEGEN_ACTION_NAME,
    _CLIF_MATCH_ACTION_NAME,
]

preprocessor_compile_actions = [
    _C_COMPILE_ACTION_NAME,
    _CPP_COMPILE_ACTION_NAME,
    _LINKSTAMP_COMPILE_ACTION_NAME,
    _PREPROCESS_ASSEMBLE_ACTION_NAME,
    _CPP_HEADER_PARSING_ACTION_NAME,
    _CPP_MODULE_COMPILE_ACTION_NAME,
    _CLIF_MATCH_ACTION_NAME,
]

codegen_compile_actions = [
    _C_COMPILE_ACTION_NAME,
    _CPP_COMPILE_ACTION_NAME,
    _LINKSTAMP_COMPILE_ACTION_NAME,
    _ASSEMBLE_ACTION_NAME,
    _PREPROCESS_ASSEMBLE_ACTION_NAME,
    _CPP_MODULE_CODEGEN_ACTION_NAME,
    _LTO_BACKEND_ACTION_NAME,
]

all_link_actions = [
    _CPP_LINK_EXECUTABLE_ACTION_NAME,
    _CPP_LINK_DYNAMIC_LIBRARY_ACTION_NAME,
    _CPP_LINK_NODEPS_DYNAMIC_LIBRARY_ACTION_NAME,
]

def _impl(ctx):
    tool_paths = [
        tool_path(name = name, path = path)
        for name, path in ctx.attr.tool_paths.items()
    ]

    if (ctx.attr.cpu == "local"):
        toolchain_identifier = "local_linux"
        target_system_name = "local"
        target_cpu = "local"
        target_libc = "local"
        compiler = "compiler"
        abi_version = "local"
        abi_libc_version = "local"
    elif (ctx.attr.cpu == "aarch64"):
        toolchain_identifier = "aarch64"
        target_system_name = "aarch64"
        target_cpu = "aarch64"
        target_libc = "aarch64"
        compiler = "aarch64"
        abi_version = "aarch64"
        abi_libc_version = "aarch64"
    else:
        fail("Unreachable")

    obj_copy_tool_path = None
    for t in tool_paths:
        if t.name == "objcopy":
            obj_copy_tool_path = t.path

    objcopy_embed_data_action = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        if (obj_copy_tool_path == None):
            fail("Unreachable")
        objcopy_embed_data_action = action_config(
            action_name = "objcopy_embed_data",
            enabled = True,
            tools = [tool(path = obj_copy_tool_path)],
        )

    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        action_configs = [objcopy_embed_data_action]
    else:
        fail("Unreachable")

    default_link_flags_feature = None
    if ctx.attr.cpu == "local":
        default_link_flags_feature = feature(
            name = "default_link_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = all_link_actions,
                    flag_groups = [
                        flag_group(
                            flags = [
                                "-lstdc++",
                                "-Wl,-z,relro,-z,now",
                                "-no-canonical-prefixes",
                                "-pass-exit-codes",
                                "-B/usr/bin",
                            ],
                        ),
                        flag_group(
                            flags = ctx.attr.link_opts,
                        ),
                    ],
                ),
                flag_set(
                    actions = all_link_actions,
                    flag_groups = [flag_group(flags = ["-Wl,--gc-sections"])],
                    with_features = [with_feature_set(features = ["opt"])],
                ),
            ],
        )
    elif ctx.attr.cpu == "aarch64":
        default_link_flags_feature = feature(
            name = "default_link_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = all_link_actions,
                    flag_groups = [
                        flag_group(
                            flags = [
                                "-lstdc++",
                                "-Wl,-z,relro,-z,now",
                                "-no-canonical-prefixes",
                                "-pass-exit-codes",
                            ],
                        ),
                        flag_group(
                            flags = ctx.attr.link_opts,
                        ),
                    ],
                ),
                flag_set(
                    actions = all_link_actions,
                    flag_groups = [flag_group(flags = ["-Wl,--gc-sections"])],
                    with_features = [with_feature_set(features = ["opt"])],
                ),
            ],
        )

    unfiltered_compile_flags_feature = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        unfiltered_compile_flags_feature = feature(
            name = "unfiltered_compile_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = [
                        _ASSEMBLE_ACTION_NAME,
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [
                        flag_group(
                            flags = [
                                "-no-canonical-prefixes",
                                "-fno-canonical-system-headers",
                                "-Wno-builtin-macro-redefined",
                                "-D__DATE__=\"redacted\"",
                                "-D__TIMESTAMP__=\"redacted\"",
                                "-D__TIME__=\"redacted\"",
                            ],
                        ),
                    ],
                ),
            ],
        )

    supports_pic_feature = feature(name = "supports_pic", enabled = True)

    default_compile_flags_feature = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        default_compile_flags_feature = feature(
            name = "default_compile_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = [
                        _ASSEMBLE_ACTION_NAME,
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [
                        flag_group(
                            flags = [
                                "-fstack-protector",
                                "-Wall",
                                "-Werror",
                                "-Wunused-but-set-parameter",
                                "-Wno-free-nonheap-object",
                                "-fno-omit-frame-pointer",
                            ],
                        ),
                        flag_group(
                            flags = ctx.attr.cxx_opts,
                        ),
                    ],
                ),
                flag_set(
                    actions = [
                        _ASSEMBLE_ACTION_NAME,
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [flag_group(flags = ["-g"])],
                    with_features = [with_feature_set(features = ["dbg"])],
                ),
                flag_set(
                    actions = [
                        _ASSEMBLE_ACTION_NAME,
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [
                        flag_group(
                            flags = [
                                "-g0",
                                "-O2",
                                "-U_FORTIFY_SOURCE",
                                "-D_FORTIFY_SOURCE=2",
                                "-DNDEBUG",
                                "-ffunction-sections",
                                "-fdata-sections",
                            ],
                        ),
                    ],
                    with_features = [with_feature_set(features = ["opt"])],
                ),
                flag_set(
                    actions = [
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [flag_group(flags = ["-std=c++14"])],
                ),
            ],
        )

    opt_feature = feature(name = "opt")

    supports_dynamic_linker_feature = feature(name = "supports_dynamic_linker", enabled = True)

    objcopy_embed_flags_feature = feature(
        name = "objcopy_embed_flags",
        enabled = True,
        flag_sets = [
            flag_set(
                actions = ["objcopy_embed_data"],
                flag_groups = [flag_group(flags = ["-I", "binary"])],
            ),
        ],
    )

    dbg_feature = feature(name = "dbg")

    user_compile_flags_feature = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        user_compile_flags_feature = feature(
            name = "user_compile_flags",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = [
                        _ASSEMBLE_ACTION_NAME,
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                    ],
                    flag_groups = [
                        flag_group(
                            flags = ["%{user_compile_flags}"],
                            iterate_over = "user_compile_flags",
                            expand_if_available = "user_compile_flags",
                        ),
                    ],
                ),
            ],
        )

    sysroot_feature = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        sysroot_feature = feature(
            name = "sysroot",
            enabled = True,
            flag_sets = [
                flag_set(
                    actions = [
                        _PREPROCESS_ASSEMBLE_ACTION_NAME,
                        _LINKSTAMP_COMPILE_ACTION_NAME,
                        _C_COMPILE_ACTION_NAME,
                        _CPP_COMPILE_ACTION_NAME,
                        _CPP_HEADER_PARSING_ACTION_NAME,
                        _CPP_MODULE_COMPILE_ACTION_NAME,
                        _CPP_MODULE_CODEGEN_ACTION_NAME,
                        _LTO_BACKEND_ACTION_NAME,
                        _CLIF_MATCH_ACTION_NAME,
                        _CPP_LINK_EXECUTABLE_ACTION_NAME,
                        _CPP_LINK_DYNAMIC_LIBRARY_ACTION_NAME,
                        _CPP_LINK_NODEPS_DYNAMIC_LIBRARY_ACTION_NAME,
                    ],
                    flag_groups = [
                        flag_group(
                            flags = ["--sysroot=%{sysroot}"],
                            expand_if_available = "sysroot",
                        ),
                    ],
                ),
            ],
        )

    features = None
    if (ctx.attr.cpu == "local" or ctx.attr.cpu == "aarch64"):
        features = [
            default_compile_flags_feature,
            default_link_flags_feature,
            supports_dynamic_linker_feature,
            supports_pic_feature,
            objcopy_embed_flags_feature,
            opt_feature,
            dbg_feature,
            user_compile_flags_feature,
            sysroot_feature,
            unfiltered_compile_flags_feature,
        ]

    artifact_name_patterns = []
    make_variables = []

    out = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.write(out, "Fake executable")
    return [
        cc_common.create_cc_toolchain_config_info(
            ctx = ctx,
            features = features,
            action_configs = action_configs,
            artifact_name_patterns = artifact_name_patterns,
            cxx_builtin_include_directories = ctx.attr.cxx_builtin_include_directories,
            toolchain_identifier = toolchain_identifier,
            target_system_name = target_system_name,
            target_cpu = target_cpu,
            target_libc = target_libc,
            compiler = compiler,
            abi_version = abi_version,
            abi_libc_version = abi_libc_version,
            tool_paths = tool_paths,
            make_variables = make_variables,
        ),
        DefaultInfo(
            executable = out,
        ),
    ]

cc_toolchain_config = rule(
    implementation = _impl,
    attrs = {
        "cpu": attr.string(mandatory = True, values = ["local", "aarch64"]),
        "tool_paths": attr.string_dict(mandatory = True),
        "cxx_builtin_include_directories": attr.string_list(mandatory = True),
        "cxx_opts": attr.string_list(),
        "link_opts": attr.string_list(),
    },
    provides = [CcToolchainConfigInfo],
    executable = True,
)
