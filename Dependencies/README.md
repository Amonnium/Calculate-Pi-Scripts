# Dependencies
This README will explain how to install the **dependencies** needed by **all** of the **languages** on the **repo**.

## Manual Installation
**The commands** to **install the software** needed by **all** of the **languages** are:

Update the packages:
```
sudo apt update && sudo apt upgrade -y
```
Python:
```
sudo apt install -y python3 python-is-python3 python3-pip
```
GCC and G++:
```
sudo apt install -y build-essential
```
Go:
```
sudo apt install -y golang-go
```
JDK:
```
sudo apt install -y default-jdk
```
Julia:
```
curl -fsSL https://install.julialang.org | sh
```
Kotlin:
```
curl -s "https://get.sdkman.io" | bash
source "$HOME/.sdkman/bin/sdkman-init.sh"
sdk install kotlin
```
Rust:
```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```
Node.js:
```
curl -o-
https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install --lts
```
Lua:
```
sudo apt install -y lua5.4
```
Perl:
```
sudo apt install -y perl
```
Ruby:
```
sudo apt install -y ruby-full
```
Zig:
```
sudo snap install zig --classic --edge
```
Raku:
```
curl https://rakubrew.org/install-on-macos.sh | sh
rakubrew init
```
## Using Dependencies Installer
If you don't want to copy-paste the commands every time, you can use the Dependencies Installer by using this command:
```
bash <(curl -fsSL https://raw.githubusercontent.com/Amonnium/Calculate-Pi-Scripts/refs/heads/readme.md/Dependencies/install.sh)
```
Note: This method, at the moment I'm writing this, **only works** on **Debian-based systems**, like **Debian**, **Ubuntu**, **Linux Mint**, **Zorin OS**, and many others.
