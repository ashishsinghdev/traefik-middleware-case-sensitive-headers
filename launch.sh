#!/bin/bash
if [ ! -f traefik ]; then
  curl -L https://github.com/traefik/traefik/releases/download/v3.0.0/traefik_v3.0.0_linux_amd64.tar.gz | tar xz
  chmod +x traefik
fi

if [ ! -f traefik-config-validator ]; then
  curl -L https://github.com/otto-de/traefik-config-validator/releases/download/v0.0.2/traefik-config-validator_0.0.2_Linux_x86_64.tar.gz | tar xz
  chmod +x traefik-config-validator
fi

docker run -d --network host containous/whoami -port 5000

./traefik-config-validator -cfg traefik-config.yml -cfgdir provider

mkdir -p plugins-local/src/github.com/ashishsinghdev
rsync -av ./ plugins-local/src/github.com/ashishsinghdev/traefik-middleware-case-sensitive-headers --include="*.go" --include="go.mod" --include=".traefik.yml" --exclude="*"

./traefik --configFile traefik-config.yml
