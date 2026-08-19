#!/bin/sh

WORKDIR="$(cd "$(dirname "$0")" && pwd)"
USER="$(whoami)"

cd "$WORKDIR" || exit 1

CMD="$1"
[ -z "$CMD" ] && CMD=opencode

if [ -z "$SSH_AUTH_SOCK" ]; then
  echo "Attenzione: SSH_AUTH_SOCK non impostato. Avvia ssh-agent e fai ssh-add prima di continuare." >&2
fi

podman run --rm -it \
  --userns=keep-id \
  --network=slirp4netns \
  -v "$(pwd)":/workspace:Z \
  -v "$SSH_AUTH_SOCK":/ssh-agent:Z \
  -e SSH_AUTH_SOCK=/ssh-agent \
  -v "$HOME/.gitconfig":/$USER/.gitconfig:ro,Z \
  ai-occ "$CMD"

# podman run --rm -it \
#   --userns=keep-id \
#   --network=slirp4netns \
#   --tmpfs /home/opencode:rw,exec,uid=$(id -u),gid=$(id -g),mode=700 \
#   -e HOME=/home/opencode \
#   -v "$(pwd)":/workspace:Z \
#   -v "$SSH_AUTH_SOCK":/ssh-agent:Z \
#   -e SSH_AUTH_SOCK=/ssh-agent \
#   -v "$HOME/.gitconfig":/home/opencode/.gitconfig:ro,Z \
#   ai-occ "$CMD"