#!/usr/bin/env bash

set -Eeuo pipefail

readonly INSTALL_DIR="${HOME:-}/.local/bin"
readonly INSTALL_NAME="calculate-pi-cli"
readonly PATH_START='# >>> Calculate Pi Scripts CLI PATH >>>'
readonly PATH_END='# <<< Calculate Pi Scripts CLI PATH <<<'

profile_stage=""

cleanup() {
    local status=$?
    [[ -z "$profile_stage" ]] || rm -f -- "$profile_stage"
    exit "$status"
}
trap cleanup EXIT

fail() {
    printf 'Error: %s\n' "$1" >&2
    exit 1
}

[[ -n "${HOME:-}" && "$HOME" == /* && -d "$HOME" ]] || \
    fail 'A valid, absolute HOME directory is required to locate the CLI installation.'

install_path="$INSTALL_DIR/$INSTALL_NAME"
if [[ ! -f "$install_path" || -L "$install_path" ]]; then
    fail "CLI executable not found at '$install_path'. Nothing was uninstalled."
fi

for utility in awk cp grep mktemp mv rm; do
    command -v "$utility" >/dev/null 2>&1 || \
        fail "Required system utility '$utility' was not found. The CLI was left unchanged."
done

remove_path_setup() {
    local shell_name profile_path

    shell_name=${SHELL:-}
    shell_name=${shell_name##*/}

    case "$shell_name" in
        bash) profile_path="$HOME/.bashrc" ;;
        zsh) profile_path="$HOME/.zshrc" ;;
        fish) profile_path="$HOME/.config/fish/config.fish" ;;
        *) profile_path="$HOME/.profile" ;;
    esac

    if [[ -L "$profile_path" ]]; then
        command -v readlink >/dev/null 2>&1 || \
            fail "The readlink utility is required to resolve '$profile_path'. The CLI was left unchanged."
        profile_path=$(readlink -f -- "$profile_path") || \
            fail "Could not resolve shell startup file '$profile_path'. The CLI was left unchanged."
    fi

    if [[ ! -e "$profile_path" ]]; then
        printf 'No shell startup file found at %s; no PATH setup was removed.\n' "$profile_path"
        return
    fi

    [[ -f "$profile_path" && -w "$profile_path" ]] || \
        fail "Shell startup file '$profile_path' is not a writable regular file. The CLI was left unchanged."

    if ! grep -Fq "$PATH_START" "$profile_path" && ! grep -Fq "$PATH_END" "$profile_path"; then
        printf 'No installer-managed PATH setup found in %s.\n' "$profile_path"
        return
    fi

    profile_stage=$(mktemp "${profile_path}.calculate-pi-uninstall.XXXXXX") || \
        fail "Could not create a temporary file beside '$profile_path'. The CLI was left unchanged."

    cp -p -- "$profile_path" "$profile_stage" || \
        fail "Could not preserve '$profile_path'. The CLI was left unchanged."

    awk -v start="$PATH_START" -v end="$PATH_END" '
        BEGIN { inside = 0; malformed = 0 }
        $0 == start {
            if (inside) {
                malformed = 1
                exit
            }
            inside = 1
            next
        }
        $0 == end {
            if (!inside) {
                malformed = 1
                exit
            }
            inside = 0
            next
        }
        !inside { print }
        END {
            if (inside) malformed = 1
            if (malformed) exit 2
        }
    ' "$profile_path" > "$profile_stage" || \
        fail "The PATH setup markers in '$profile_path' are incomplete or malformed. The CLI was left unchanged."

    mv -f -- "$profile_stage" "$profile_path" || \
        fail "Could not update '$profile_path'. The CLI was left unchanged."
    profile_stage=""
    printf 'Removed installer-managed PATH setup from %s.\n' "$profile_path"
}

remove_path_setup

rm -f -- "$install_path" || \
    fail "Could not remove the CLI executable at '$install_path'."

printf 'Removed CLI executable: %s\n' "$install_path"
printf 'Calculate Pi Scripts CLI uninstallation complete.\n'
