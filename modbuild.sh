#!/bin/bash
set -e

for file in ./include/mods/*/; do
	echo \[Mod Build\] Building $file
	(cd "$file" && go mod tidy && sh build.sh)
	echo \[Mod Build\] Finished
done

cp -r ./include/* ./_out/