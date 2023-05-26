"""
Toolchains used by gaia
"""


# buildifier: disable=unnamed-macro
def register_cc_toolchains():
    native.register_toolchains(
        "//toolchains:linux_x86_64_cc",
        "//toolchains:linux_aarch64_cc",
    )

# buildifier: disable=unnamed-macro
def register_gaia_toolchains():
    register_cc_toolchains()
