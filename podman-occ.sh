#!/bin/sh

WORKDIR="$(cd "$(dirname "$0")" && pwd)"

mkdir -p "$HOME/.opencode-home"

cd "$WORKDIR" || exit 1

CMD="$1"
[ -z "$CMD" ] && CMD=opencode

if [ -z "$SSH_AUTH_SOCK" ]; then
  echo "Attenzione: SSH_AUTH_SOCK non impostato. Avvia ssh-agent e fai ssh-add prima di continuare." >&2
fi

podman run --rm -it \
  --userns=keep-id \
  --network=slirp4netns \
  -v "$HOME/.opencode-home":/home/opencode:Z \
  -e GH_TOKEN="$GH_TOKEN" \
  -e HOME=/home/opencode \
  -e GIT_CONFIG_GLOBAL=/tmp/.gitconfig \
  -v "$(pwd)":/workspace:Z \
  -v "$SSH_AUTH_SOCK":/ssh-agent:Z \
  -e SSH_AUTH_SOCK=/ssh-agent \
  -v "$HOME/.gitconfig":/tmp/.gitconfig:ro,Z \
  ai-occ "$CMD"