# Toolchains and Runtimes

An interpreter/runtime runs a script or bytecode. A compiler turns source code into a program. You only need tools for the examples you choose; this repository does not need every language installed.

## CLI-managed runtimes

The CLI manages Python 3, Java, and Julia because its selected Python/Julia examples are scripts and its Java example is JVM bytecode. It checks for a compatible runtime first and asks before installing one.

From the CLI, choose **Manage Toolchains** to view status or explicitly request installation/removal. A missing runtime can also prompt during a run. The CLI may use Winget on Windows or apt-get on Debian/Ubuntu Linux. It does not install compilers, and it records package ownership before offering uninstall. A runtime already installed outside the CLI is left alone.

| Tool | Minimum used by this project | Official installation information |
| --- | --- | --- |
| Go | 1.24 for the CLI | [Go downloads](https://go.dev/dl/) |
| Python | Python 3; CLI package is Python 3.13 on Windows | [Python downloads](https://www.python.org/downloads/) |
| Java | 21 for the supplied `Pi.class` bytecode | [Eclipse Temurin downloads](https://adoptium.net/temurin/releases/) |
| Julia | Julia 1.x | [Julia downloads](https://julialang.org/downloads/) |

On Windows, the managed package IDs are `Python.Python.3.13`, `Microsoft.OpenJDK.21`, and `JuliaLang.Julia`. On Debian/Ubuntu, the package names are `python3`, `openjdk-21-jre`, and `julia`. Package availability depends on the OS release and configured repositories. On other Linux distributions, use the linked official installation instructions; the CLI does not guess a package manager.

## Tools for compiling examples

These are needed only if you want to compile source rather than run a suitable supplied artifact or script.

| Language(s) | Tool | Official installation information |
| --- | --- | --- |
| C | GCC or another C compiler | [GCC](https://gcc.gnu.org/), [MSYS2 on Windows](https://www.msys2.org/) |
| C++ | G++ or another C++ compiler | [GCC](https://gcc.gnu.org/), [MSYS2 on Windows](https://www.msys2.org/) |
| Rust | `rustc` | [Install Rust](https://www.rust-lang.org/tools/install) |
| Go | Go compiler | [Go downloads](https://go.dev/dl/) |
| Java | JDK (includes `javac`) | [Eclipse Temurin downloads](https://adoptium.net/temurin/releases/) |
| C# | .NET SDK | [.NET downloads](https://dotnet.microsoft.com/download) |
| JavaScript | Node.js | [Node.js downloads](https://nodejs.org/en/download) |
| Kotlin | JVM plus Kotlin compiler, only to rebuild the source | [Kotlin command-line compiler](https://kotlinlang.org/docs/command-line.html) |
| Lua | Lua interpreter/compiler | [Lua downloads](https://www.lua.org/download.html) |
| Nim | Nim compiler | [Nim installation](https://nim-lang.org/install.html) |
| Perl | Perl | [Perl installation](https://www.perl.org/get.html) |
| PowerShell | PowerShell | [Install PowerShell](https://learn.microsoft.com/powershell/scripting/install/installing-powershell) |
| Raku | Rakudo | [Rakudo downloads](https://rakudo.org/downloads) |
| Ruby | Ruby | [Ruby downloads](https://www.ruby-lang.org/en/downloads/) |
| Zig | Zig compiler | [Zig downloads](https://ziglang.org/download/) |

Installation steps differ by operating system and release. Follow the official guide for your platform, then verify a tool is on `PATH` with its version command (for example, `gcc --version`, `rustc --version`, or `julia --version`). The [language guide](languages.md) contains the commands for this repository's files.
