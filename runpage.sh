#!/bin/bash
set -e

cp ./_out/main.wasm ./webpage/
cd webpage
# npm install -g serve
serve