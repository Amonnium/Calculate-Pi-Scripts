# Troubleshooting and FAQ

## Common problems

### `go` is not recognized or the CLI will not build

Install Go 1.24 or later and open a new terminal so its `bin` directory is on `PATH`. Run `go version`, then from the repository root retry `go -C CLI run .`. The first run needs internet access to download the CLI's Go dependencies.

### The CLI cannot download a script or binary

Check internet access and allow HTTPS to `raw.githubusercontent.com`. The CLI downloads on demand, checks the pinned artifact hash, and keeps verified files in the OS user cache. A failed or interrupted download is not treated as a valid cached artifact; retry after network access is restored.

### I selected C, C++, Rust, or Go on Linux/macOS and no binary is available

The supplied native executables are Windows x86-64 only. The CLI does not compile from source. Install the appropriate compiler and use the source-build commands in [Languages and commands](languages.md), or run a different implementation with an available runtime.

### Java says the class file has an unsupported version

The checked-in `Pi.class` uses Java class-file version 65, which requires Java 21 or newer. Check `java -version`; install a current JRE/JDK or compile the Java source with a JDK into a separate output directory.

### A runtime is missing or too old

Use **Manage Toolchains** to inspect Python, Java, and Julia. The CLI only offers a system package installation after you choose it or after a run needs that runtime. On Windows it needs Winget; on Linux, automatic package support is limited to Debian/Ubuntu with apt-get. On other Linux distributions, use the official links in [Toolchains](toolchains.md).

### Package installation fails or asks for elevated permissions

The CLI does not bypass system security prompts. Winget may ask you to accept package agreements; apt-get requires root privileges and may use `sudo`. Confirm network access, permissions, and package availability for your OS release. If the package manager is missing, install the package manager using your OS's official instructions or install the runtime manually.

### PowerShell refuses to run a script

Use PowerShell 7 (`pwsh`) where available and run the file explicitly, for example `pwsh -File Pi_in_PowerShell/Pi.ps1`. Follow your organization's execution-policy guidance rather than changing system policy blindly.

### A calculation is slow or the answer has fewer correct digits than expected

Some examples deliberately use slow series. C, C++, and Rust use 100 million Leibniz terms; Go uses 10 million. More iterations improve that series slowly, and floating-point rounding eventually limits accuracy. See [Algorithms](algorithms.md) before comparing results.

### `Pi.bat` cannot find `cscript`

The batch example calls Windows Script Host and VBScript. It only applies to Windows systems where those components are enabled. It is not a cross-platform script.

## Frequently asked questions

### Do I need every language installed?

No. Install only the runtime/compiler for examples you want to run or build. The CLI manages only the Python, Java, and Julia runtimes it needs, and it prefers suitable supplied binaries.

### Does the CLI support every language in the repository?

No. Its implementation menu contains only Python, C, C++, Rust, Go, Java, and Julia. Other examples remain in the repository for direct use and learning.

### Does the CLI compile source automatically?

No. It uses compatible supplied artifacts or runs the selected source with a required runtime. If a native executable is unavailable on your OS, compile the source yourself using the documented commands.

### Can I run the checked-in `.exe` files on Linux?

No. They are Windows x86-64 binaries. Use source with a compiler for your platform instead.

### Will the CLI remove a runtime that I installed myself?

No. Uninstall is offered only for a system package recorded as installed by this CLI, and it asks for confirmation before removal.
