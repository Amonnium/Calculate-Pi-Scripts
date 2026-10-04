#!/usr/bin/env bash

set -Eeuo pipefail
umask 077

readonly PROJECT_NAME="Calculate Pi Scripts"
readonly CLI_VERSION="1.0"
readonly RELEASE_TAG="calculate-pi-cli-v1.0"
readonly ASSET_NAME="calculate-pi-cli"
readonly ASSET_SIZE="11193952"
readonly ASSET_SHA256="2598df1967985bbf31a486da65d8a73e6060961c7620a3414bfec46cb88b1aa3"
readonly RELEASE_BASE="https://github.com/Amonnium/Calculate-Pi-Scripts/releases/download/${RELEASE_TAG}"
readonly ASSET_URL="${RELEASE_BASE}/${ASSET_NAME}"
readonly INSTALL_NAME="calculate-pi-cli"
readonly INSTALL_DIR="${HOME:-}/.local/bin"

work_dir=""
stage_path=""
backup_path=""
profile_stage=""
install_path=""
rollback_install=0

cleanup() {
    local status=$?

    if (( rollback_install )); then
        if [[ -n "$backup_path" && -f "$backup_path" ]]; then
            mv -f -- "$backup_path" "$install_path" 2>/dev/null || true
        elif [[ -n "$install_path" ]]; then
            rm -f -- "$install_path" 2>/dev/null || true
        fi
    fi

    [[ -z "$stage_path" ]] || rm -f -- "$stage_path" 2>/dev/null || true
    [[ -z "$backup_path" ]] || rm -f -- "$backup_path" 2>/dev/null || true
    [[ -z "$profile_stage" ]] || rm -f -- "$profile_stage" 2>/dev/null || true
    [[ -z "$work_dir" ]] || rm -rf -- "$work_dir" 2>/dev/null || true

    exit "$status"
}
trap cleanup EXIT

if [[ -t 1 ]]; then
    printf '\033[2J\033[H'
else
    clear 2>/dev/null || true
fi

printf '%s\n' \
    '   ____      _            _       ____  _' \
    '  / ___|__ _| | ___ _   _| | __ _|  _ \\(_)' \
    ' | |   / _` | |/ __| | | | |/ _` | |_) | |' \
    ' | |__| (_| | | (__| |_| | | (_| |  __/| |' \
    '  \\____\\__,_|_|\\___|\\__,_|_|\\__,_|_|   |_|' \
    '' \
    'Calculate Pi Scripts CLI Installer'
printf 'Version %s | Linux x86-64 | User installation\n\n' "$CLI_VERSION"
printf 'Project: %s\nLicense: MIT\nAuthor: Amonnium\n\n' "$PROJECT_NAME"

color_enabled=0
if [[ -t 1 ]]; then
    color_enabled=1
fi

stage() {
    if (( color_enabled )); then
        printf '\033[1;36m[%s/6]\033[0m %s\n' "$1" "$2"
    else
        printf '[%s/6] %s\n' "$1" "$2"
    fi
}

success() {
    if (( color_enabled )); then
        printf '\033[1;32m%s\033[0m\n' "$1"
    else
        printf '%s\n' "$1"
    fi
}

warn() {
    printf 'Warning: %s\n' "$1" >&2
}

fail() {
    printf 'Error: %s\n' "$1" >&2
    exit 1
}

stage 1 'Checking system...'
[[ -n "${HOME:-}" && "$HOME" == /* && -d "$HOME" && -w "$HOME" ]] || \
    fail 'A writable, absolute HOME directory is required for a user-level installation.'

if ! command -v uname >/dev/null 2>&1; then
    fail 'The uname utility is required to detect this system.'
fi

if [[ "$(uname -s)" != 'Linux' ]]; then
    fail "This installer supports Linux only; detected $(uname -s)."
fi

for utility in mktemp mkdir cp mv chmod rm wc od tr grep; do
    command -v "$utility" >/dev/null 2>&1 || \
        fail "Required system utility '$utility' was not found. No packages were installed."
done

command -v sha256sum >/dev/null 2>&1 || \
    fail 'The sha256sum utility is required to verify the release. No packages were installed.'

stage 2 'Detecting architecture...'
machine=$(uname -m)
case "$machine" in
    x86_64|amd64)
        ;;
    *)
        fail "No released Linux CLI artifact is available for architecture '$machine'. The current release provides x86-64 only."
        ;;
esac

printf 'Detected: Linux %s\n' "$machine"
printf 'Install location: %s/%s\n\n' "$INSTALL_DIR" "$INSTALL_NAME"

verify_sha256() {
    local file=$1
    local output actual

    output=$(sha256sum -- "$file") || return 1
    actual=${output%%[[:space:]]*}
    [[ "$actual" == "$ASSET_SHA256" ]]
}

verify_elf_x86_64() {
    local file=$1
    local magic elf_class elf_data machine_bytes machine_low machine_high

    magic=$(od -An -tx1 -N4 "$file" | tr -d '[:space:]') || return 1
    elf_class=$(od -An -tu1 -j4 -N1 "$file" | tr -d '[:space:]') || return 1
    elf_data=$(od -An -tu1 -j5 -N1 "$file" | tr -d '[:space:]') || return 1
    machine_bytes=$(od -An -tu1 -j18 -N2 "$file") || return 1
    read -r machine_low machine_high <<< "$machine_bytes"

    [[ "$magic" == '7f454c46' && "$elf_class" == '2' && "$elf_data" == '1' && \
       "$machine_low" == '62' && "$machine_high" == '0' ]]
}

stage 3 'Downloading Calculate Pi Scripts...'
if command -v curl >/dev/null 2>&1; then
    download_tool='curl'
elif command -v wget >/dev/null 2>&1; then
    download_tool='wget'
else
    fail 'Neither curl nor wget is available. Install one, then rerun this installer. No package manager was invoked.'
fi

work_dir=$(mktemp -d "${TMPDIR:-/tmp}/calculate-pi-cli.XXXXXX") || \
    fail 'Could not create a temporary download directory.'

download_path="$work_dir/$ASSET_NAME"

if [[ "$download_tool" == 'curl' ]]; then
    if ! curl --fail --location --silent --show-error \
        --proto '=https' --proto-redir '=https' \
        --retry 2 --connect-timeout 15 --max-time 180 \
        --output "$download_path" "$ASSET_URL"; then
        fail 'The CLI download failed. Check your network connection and try again.'
    fi
else
    if ! wget --https-only --tries=3 --timeout=30 \
        --output-document="$download_path" "$ASSET_URL"; then
        fail 'The CLI download failed. Check your network connection and try again.'
    fi
fi

[[ -f "$download_path" ]] || fail 'The download did not produce a file.'

download_size=$(wc -c < "$download_path" | tr -d '[:space:]')
[[ "$download_size" == "$ASSET_SIZE" ]] || \
    fail "The downloaded file has an unexpected size (${download_size} bytes); expected ${ASSET_SIZE}."

verify_sha256 "$download_path" || \
    fail 'The downloaded file failed SHA-256 verification. It will not be installed.'

verify_elf_x86_64 "$download_path" || \
    fail 'The downloaded release is not a valid Linux x86-64 ELF executable. It will not be installed.'

stage 4 'Installing CLI...'
install_path="$INSTALL_DIR/$INSTALL_NAME"

mkdir -p -- "$INSTALL_DIR" || \
    fail "Could not create the user install directory '$INSTALL_DIR'. Check its permissions."

[[ -w "$INSTALL_DIR" ]] || \
    fail "The user install directory '$INSTALL_DIR' is not writable. No elevated privileges were requested."

if [[ -L "$install_path" ]]; then
    fail "Refusing to replace the symlink '$install_path'. Remove or rename it yourself, then rerun the installer."
fi

if [[ -e "$install_path" && ! -f "$install_path" ]]; then
    fail "The install target '$install_path' exists and is not a regular file. It was left unchanged."
fi

stage_path=$(mktemp "$INSTALL_DIR/.calculate-pi-cli.XXXXXX") || \
    fail "Could not create a staging file in '$INSTALL_DIR'."

cp -- "$download_path" "$stage_path" || \
    fail 'Could not stage the verified CLI. The previous installation was left unchanged.'

chmod 0755 -- "$stage_path" || \
    fail 'Could not mark the staged CLI executable. The previous installation was left unchanged.'

verify_sha256 "$stage_path" && verify_elf_x86_64 "$stage_path" || \
    fail 'The staged executable failed verification. The previous installation was left unchanged.'

if [[ -e "$install_path" ]]; then
    backup_path=$(mktemp "$INSTALL_DIR/.calculate-pi-cli.previous.XXXXXX") || \
        fail 'Could not preserve the existing installation before updating it.'
    cp -p -- "$install_path" "$backup_path" || \
        fail 'Could not preserve the existing installation; it was left unchanged.'
fi

mv -f -- "$stage_path" "$install_path" || \
    fail 'Could not move the verified CLI into place. The previous installation was left unchanged.'
stage_path=''
rollback_install=1

configure_path() {
    local shell_name profile_path path_marker

    shell_name=${SHELL:-}
    shell_name=${shell_name##*/}

    case "$shell_name" in
        bash) profile_path="$HOME/.bashrc" ;;
        zsh) profile_path="$HOME/.zshrc" ;;
        fish)
            profile_path="$HOME/.config/fish/config.fish"
            mkdir -p -- "$HOME/.config/fish" || return 1
            ;;
        *) profile_path="$HOME/.profile" ;;
    esac

    if [[ -L "$profile_path" ]]; then
        command -v readlink >/dev/null 2>&1 || return 1
        profile_path=$(readlink -f -- "$profile_path" 2>/dev/null) || return 1
    fi

    if [[ -e "$profile_path" && ( ! -f "$profile_path" || ! -w "$profile_path" ) ]]; then
        return 1
    fi

    path_marker='# >>> Calculate Pi Scripts CLI PATH >>>'
    if [[ -f "$profile_path" ]] && grep -Fq "$path_marker" "$profile_path"; then
        printf 'PATH setup is already present in %s\n' "$profile_path"
        return 0
    fi

    profile_stage=$(mktemp "${profile_path}.calculate-pi.XXXXXX") || return 1

    if [[ -f "$profile_path" ]]; then
        cp -p -- "$profile_path" "$profile_stage" || return 1
    else
        chmod 0600 -- "$profile_stage" || return 1
    fi

    if [[ "$shell_name" == 'fish' ]]; then
        printf '\n%s\nif not contains -- "$HOME/.local/bin" $PATH\n    set -gx PATH "$HOME/.local/bin" $PATH\nend\n%s\n' \
            "$path_marker" '# <<< Calculate Pi Scripts CLI PATH <<<' >> "$profile_stage" || return 1
    else
        printf '\n%s\ncase ":${PATH:-}:" in\n  *":$HOME/.local/bin:"*) ;;\n  *) export PATH="$HOME/.local/bin:$PATH" ;;\nesac\n%s\n' \
            "$path_marker" '# <<< Calculate Pi Scripts CLI PATH <<<' >> "$profile_stage" || return 1
    fi

    mv -f -- "$profile_stage" "$profile_path" || return 1
    profile_stage=''
    printf 'Added PATH setup to %s\n' "$profile_path"
}

stage 5 'Configuring PATH...'
path_configured=1
case ":${PATH:-}:" in
    *":$INSTALL_DIR:"*)
        printf 'PATH already includes %s\n' "$INSTALL_DIR"
        ;;
    *)
        if ! configure_path; then
            path_configured=0
            warn "Could not safely update a shell startup file. Add '$INSTALL_DIR' to PATH manually."
        fi
        ;;
esac

export PATH="$INSTALL_DIR:${PATH:-}"

stage 6 'Verifying installation...'
[[ -f "$install_path" && -x "$install_path" ]] || \
    fail 'The installed CLI is missing or is not executable.'

verify_sha256 "$install_path" || \
    fail 'The installed CLI failed SHA-256 verification.'

verify_elf_x86_64 "$install_path" || \
    fail 'The installed CLI failed its Linux x86-64 executable check.'

# Do not launch the TUI automatically here. The CLI is interactive, and an
# installer must not block waiting for user input or enter a test loop.
rollback_install=0

if [[ -n "$backup_path" ]]; then
    rm -f -- "$backup_path"
    backup_path=''
fi

printf '\n'
success 'Installation complete!'
printf '\nRun:\n\n    %s\n\n' "$INSTALL_NAME"

if (( ! path_configured )); then
    printf 'If the command is not found, add this directory to PATH:\n\n    %s\n\n' "$INSTALL_DIR"
else
    printf 'Open a new shell if the new PATH entry is not available yet.\n\n'
fi

printf '%s is ready.\n' "$PROJECT_NAME"
printf 'Released CLI version: %s | Project author: Amonnium | License: MIT\n' "$CLI_VERSION"
