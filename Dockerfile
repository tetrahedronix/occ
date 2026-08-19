FROM docker.io/library/golang:1.26-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# nvm + Node.js
ENV NVM_DIR=/usr/local/nvm
RUN mkdir -p $NVM_DIR \
    && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | PROFILE=/dev/null bash
RUN . "$NVM_DIR/nvm.sh" && nvm install --lts && nvm alias default lts/*

# Risolvi il percorso reale e fissalo esplicitamente nel PATH
RUN ln -s "$NVM_DIR/versions/node/$(. $NVM_DIR/nvm.sh && nvm version default)" /usr/local/nvm-default
ENV PATH="/usr/local/nvm-default/bin:${PATH}"

# Installazione di OpenCode via npm
RUN . "$NVM_DIR/nvm.sh" && npm install -g --allow-scripts=opencode-ai opencode-ai 

WORKDIR /workspace