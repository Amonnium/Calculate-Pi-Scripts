# Dependencies
This README explains how to install the dependencies needed by the languages used in this repository.

## Manual Installation
The commands to install the software needed by the languages in this repo are shown per Linux distribution family below.

> Notes:
> - Package names can vary between distribution versions — adjust as needed.
> - Several tools (Julia installer, SDKMAN, rustup, nvm) use distro-agnostic installers and are shown in the "Cross-distro installers" section.

---

## Debian / Ubuntu / other Debian-based
Update packages:
```
sudo apt update && sudo apt upgrade -y
```
Install packages:
```
sudo apt install -y git python3 python-is-python3 python3-pip build-essential golang-go default-jdk lua5.4 perl ruby-full rakudo
```
(Optional) Node.js using nvm (recommended):
```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install --lts
```

---

## RHEL / CentOS / Fedora (dnf or yum)
Update packages (Fedora and modern RHEL/CentOS use dnf; older CentOS may use yum):
```
sudo dnf update -y
# or
sudo yum update -y
```
Install packages (using dnf; replace with yum on older systems):
```
sudo dnf install -y git python3 python3-pip gcc gcc-c++ make golang java-17-openjdk-devel lua perl ruby rakudo
```
Install development tools group (provides build tools):
```
sudo dnf groupinstall -y "Development Tools"
# or with yum
sudo yum groupinstall -y "Development Tools"
```
Node.js (package) or nvm:
```
sudo dnf install -y nodejs npm
# or use nvm (cross-distro) - see Cross-distro installers
```

---

## Arch Linux / Manjaro (pacman)
Update and install packages:
```
sudo pacman -Syu
sudo pacman -S --noconfirm git python python-pip base-devel go jdk-openjdk lua perl ruby rakudo nodejs npm julia zig
```
Notes:
- base-devel provides the common build tools (gcc, make, etc.).
- AUR packages may be needed for some languages/tools not in the official repos.

---

## openSUSE (Leap / Tumbleweed) - zypper
Refresh and update:
```
sudo zypper refresh && sudo zypper update -y
```
Install packages:
```
sudo zypper install -y git python3 python3-pip gcc gcc-c++ make golang java-17-openjdk-devel lua perl ruby rakudo nodejs npm julia
```
Install development pattern (optional):
```
sudo zypper install -t pattern devel_basis
```

---

## Cross-distro installers (work on most distributions)
Some projects provide their own install scripts or language-specific version managers that work across distros.

Julia (official script):
```
curl -fsSL https://install.julialang.org | sh
```
Kotlin (via SDKMAN):
```
curl -s "https://get.sdkman.io" | bash
source "$HOME/.sdkman/bin/sdkman-init.sh"
sdk install kotlin
```
Rust (rustup):
```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```
Node.js (nvm):
```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install --lts
```
Zig (snap or download releases):
- If snap is available:
```
sudo snap install zig --classic --edge
```
- Or download prebuilt binaries from the Zig releases page: https://ziglang.org/download/

---

## Quick reference (one-liners per family)
Debian/Ubuntu:
```
sudo apt update && sudo apt upgrade -y && sudo apt install -y git python3 python-is-python3 python3-pip build-essential golang-go default-jdk lua5.4 perl ruby-full rakudo
```
RHEL/Fedora:
```
sudo dnf update -y && sudo dnf install -y git python3 python3-pip gcc gcc-c++ make golang java-17-openjdk-devel lua perl ruby rakudo
```
Arch:
```
sudo pacman -Syu && sudo pacman -S --noconfirm git python python-pip base-devel go jdk-openjdk lua perl ruby rakudo
```
openSUSE:
```
sudo zypper refresh && sudo zypper update -y && sudo zypper install -y git python3 python3-pip gcc gcc-c++ make golang java-17-openjdk-devel lua perl ruby rakudo
```

---

## Using Dependencies Installer
If you don't want to copy/paste the commands every time, you can use the Dependencies Installer:
```
bash <(curl -fsSL https://raw.githubusercontent.com/Amonnium/Calculate-Pi-Scripts/refs/heads/readme.md/Dependencies/install.sh)
```
