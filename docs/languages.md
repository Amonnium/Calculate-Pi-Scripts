# Languages and Commands

Run commands from the repository root unless stated otherwise. Install the required runtime or compiler first; official download links are in [Toolchains](toolchains.md). Commands that produce a local build use a `local` output name so they do not overwrite the supplied binaries.

## CLI-supported implementations

| Language | Source | Method | Run from the repository root |
| --- | --- | --- | --- |
| Python | `Pi_in_Python/Chudnovsky.py`, `Leibniz.py`, `Monte_Carlo.py`, `Math_Pi.py` | Chudnovsky, Leibniz, Monte Carlo, and `math.pi` | `python3 Pi_in_Python/Chudnovsky.py` (Linux/macOS); `py -3 Pi_in_Python\Chudnovsky.py` (Windows) |
| C | `Pi_in_C_Folder/Pi.c` | Leibniz series | Use the supplied Windows binary or compile; see below |
| C++ | `Pi_in_C++_Folder/Pi.cpp` | Leibniz series | Use the supplied Windows binary or compile; see below |
| Rust | `Pi_in_Rust_Folder/Pi.rs`, `Math_Pi.rs` | Leibniz series and standard-library constant | Use the supplied Windows binary or compile; see below |
| Go | `Pi_in_Go_Folder/Pi.go`, `Math_Pi.go` | Leibniz series and standard-library constant | `go run Pi_in_Go_Folder/Pi.go` |
| Java | `Pi_in_Java/Pi.java`, `Math_Pi.java` | Machin's formula with `BigDecimal` and `Math.PI` | `java -cp Pi_in_Java Pi` |
| Julia | `Pi_in_Julia/Pi.jl`, `Math_Pi.jl` | Chudnovsky with binary splitting and `pi` | `julia Pi_in_Julia/Pi.jl` |

The Java command runs the checked-in `Pi.class`, which requires Java 21 or later. To compile from source without replacing the provided class file, use a separate output directory.

Windows PowerShell:

```powershell
New-Item -ItemType Directory -Force build\java
javac -d build\java Pi_in_Java\Pi.java
java -cp build\java Pi
```

Linux/macOS:

```sh
mkdir -p build/java
javac -d build/java Pi_in_Java/Pi.java
java -cp build/java Pi
```

## Compile native CLI examples

The existing C, C++, Go, and Rust binaries are Windows x86-64 only. On another platform, the repository has source but no matching prebuilt binary. The CLI does not compile these sources automatically; you can build them yourself with a compiler for your platform.

### C

Windows, using GCC:

```powershell
gcc -O2 Pi_in_C_Folder\Pi.c -o Pi_in_C_Folder\Pi-local.exe
.\Pi_in_C_Folder\Pi-local.exe
```

Linux/macOS:

```sh
cc -O2 Pi_in_C_Folder/Pi.c -o Pi_in_C_Folder/pi-local
./Pi_in_C_Folder/pi-local
```

### C++

Windows, using GCC's C++ compiler:

```powershell
g++ -O2 Pi_in_C++_Folder\Pi.cpp -o Pi_in_C++_Folder\Pi-local.exe
.\Pi_in_C++_Folder\Pi-local.exe
```

Linux/macOS:

```sh
g++ -O2 Pi_in_C++_Folder/Pi.cpp -o Pi_in_C++_Folder/pi-local
./Pi_in_C++_Folder/pi-local
```

### Rust

Windows:

```powershell
rustc -O Pi_in_Rust_Folder\Pi.rs -o Pi_in_Rust_Folder\Pi-local.exe
.\Pi_in_Rust_Folder\Pi-local.exe
```

Linux/macOS:

```sh
rustc -O Pi_in_Rust_Folder/Pi.rs -o Pi_in_Rust_Folder/pi-local
./Pi_in_Rust_Folder/pi-local
```

### Go

Run the single-file program without building a separate executable:

```sh
go run Pi_in_Go_Folder/Pi.go
```

Build it for the current platform:

Windows PowerShell: `go build -o Pi_in_Go_Folder\Pi-local.exe Pi_in_Go_Folder\Pi.go`

Linux/macOS: `go build -o Pi_in_Go_Folder/Pi-local Pi_in_Go_Folder/Pi.go`

## Other source examples

These implementations remain part of the project but are not selectable in the CLI.

| Language | File(s) | Run command from repository root | Notes |
| --- | --- | --- | --- |
| C# | `Pi.cs` | See the project setup below | Uses .NET `BigInteger` and Machin's formula. |
| JavaScript | `Pi_in_JavaScript/Pi.js`, `Pi_in_JavaScript/Math_Pi.js` | `node Pi_in_JavaScript/Pi.js` or `node Pi_in_JavaScript/Math_Pi.js` | Requires Node.js. |
| Kotlin | `Pi_in_Kotlin/Pi.kt` | `java -jar Pi_in_Kotlin/Pi_in_Kotlin.jar` | The supplied JAR includes the Kotlin runtime. |
| Lua | `Pi.lua` | `lua Pi.lua` | Leibniz series. |
| Nim | `Pi_in_Nim_Folder/Pi.nim`, `Pi_in_Nim_Folder/Math_Pi.nim` | Use the platform-specific compiler command below | Windows binaries are also supplied. |
| Perl | `Pi.pl` | `perl Pi.pl` | Uses Perl's `bignum` pragma. |
| PowerShell | `Pi_in_PowerShell/Pi.ps1`, `Math_Pi.ps1` | `pwsh -File Pi_in_PowerShell/Pi.ps1` | Alternating reciprocal triple-products and `[math]::Pi`; PowerShell 7 works across supported platforms. |
| Raku | `Pi.raku` | `raku Pi.raku` | Uses exact rational arithmetic in the BBP series. |
| Ruby | `Pi.rb` | `ruby Pi.rb` | Nilakantha series. |
| Zig | `Pi.zig` | `zig run Pi.zig` | Prints the standard-library Pi constant. |
| Windows batch | `Pi.bat` | `cmd /c Pi.bat` | Windows-only; uses Windows Script Host and VBScript. |

### Run the C# source

`Pi.cs` is a standalone source file and this repository does not include a `.csproj`. Create a temporary console project, replace its generated entry-point source, and run it.

Windows PowerShell:

```powershell
dotnet new console --output build\csharp
Remove-Item build\csharp\Program.cs
Copy-Item Pi.cs build\csharp\Program.cs
dotnet run --project build\csharp
```

Linux/macOS:

```sh
dotnet new console --output build/csharp
rm build/csharp/Program.cs
cp Pi.cs build/csharp/Program.cs
dotnet run --project build/csharp
```

## Provided artifacts

| Artifact | Platform | How to run |
| --- | --- | --- |
| `Pi_in_C_Folder/Pi_in_C.exe` | Windows x86-64 | `& .\Pi_in_C_Folder\Pi_in_C.exe` in PowerShell, or run it from Command Prompt |
| `Pi_in_C++_Folder/Pi_in_C++.exe` | Windows x86-64 | `& '.\Pi_in_C++_Folder\Pi_in_C++.exe'` in PowerShell |
| `Pi_in_Go_Folder/Pi_in_Go.exe` | Windows x86-64 | `& .\Pi_in_Go_Folder\Pi_in_Go.exe` |
| `Pi_in_Go_Folder/Math_Pi_in_Go.exe` | Windows x86-64 | `& .\Pi_in_Go_Folder\Math_Pi_in_Go.exe` |
| `Pi_in_Rust_Folder/Pi_in_Rust.exe` | Windows x86-64 | `& .\Pi_in_Rust_Folder\Pi_in_Rust.exe` |
| `Pi_in_Rust_Folder/Math_Pi_in_Rust.exe` | Windows x86-64 | `& .\Pi_in_Rust_Folder\Math_Pi_in_Rust.exe` |
| `Pi_in_Nim_Folder/Pi_in_Nim.exe` and `Math_Pi_in_Nim.exe` | Windows x86-64 | Run the selected `.exe` from PowerShell |
| `Pi_in_Java/Pi.class`, `Math_Pi.class` | JVM bytecode | `java -cp Pi_in_Java Pi` or `java -cp Pi_in_Java Math_Pi` |
| `Pi_in_Kotlin/Pi_in_Kotlin.jar` | JVM; Kotlin runtime bundled | `java -jar Pi_in_Kotlin/Pi_in_Kotlin.jar` |

The non-`.exe` native binaries referred to by older README instructions are not present in this checkout. Do not rename a Windows executable and expect it to become a Linux binary.

### Nim

On Windows, the repository's `nim.cfg` sets cache and output directories under `C:/temp`, so this command follows that configuration:

```powershell
nim c -r Pi_in_Nim_Folder\Pi.nim
```

On Linux, override those repository-level Windows paths with writable Linux paths:

```sh
nim c -r --nimcache:/tmp/pi-nimcache --outdir:/tmp/pi-nimout Pi_in_Nim_Folder/Pi.nim
```
