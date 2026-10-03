# Getting Started

## Choose how to explore

Use the interactive CLI for the seven curated implementations, or run a source file directly if you want to inspect or experiment with a particular language. The [language guide](languages.md) lists every implementation and its commands.

## Clone the repository

Git is only needed to clone or update the source repository. Downloading an individual CLI artifact does not require Git.

```sh
git clone https://github.com/Amonnium/Calculate-Pi-Scripts.git
cd Calculate-Pi-Scripts
```

You can also download the repository as a ZIP and open its extracted directory.

## Start the CLI

Install Go 1.24 or later from [go.dev/dl](https://go.dev/dl/), then run this command from the repository root on Windows, Linux, or macOS:

```sh
go -C CLI run .
```

On its first run, Go downloads the CLI's Go module dependencies. When you choose an implementation, the CLI obtains only its required repository artifact and caches it under your operating system's user cache directory. Python, Java, or Julia is installed only if the chosen implementation needs it and you confirm the package-manager prompt.

Use the arrow keys and Enter to navigate. Press B to go back, Q to quit, and N/P to move between pages. Help is available from the main menu.

## Run a source file directly

For example, Python 3 can run the Chudnovsky script from the repository root:

Windows PowerShell:

```powershell
py -3 Pi_in_Python\Chudnovsky.py
```

Linux:

```sh
python3 Pi_in_Python/Chudnovsky.py
```

The script asks how many decimal places to calculate. It uses Python's standard library and needs no third-party Python package. See [Languages and commands](languages.md) for other examples.

## Understand the supplied binaries

The checked-in C, C++, Go, Rust, and Nim executables are Windows x86-64 programs. They are not portable to Linux or macOS. Java `.class` files are platform-independent bytecode but need a Java 21-or-later runtime. The Kotlin JAR can be launched with Java. The CLI has a smaller supported set than the full source archive; the [artifact table](languages.md#provided-artifacts) shows the current files.
