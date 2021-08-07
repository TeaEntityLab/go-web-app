#!/usr/bin/env bash

GOHOME="${GOHOME:-$HOME/go}"
"$GOHOME/bin/swag" init --parseDependency --parseInternal --parseDepth 3
