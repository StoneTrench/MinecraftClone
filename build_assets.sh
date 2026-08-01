#!/bin/bash
mkdir -p "./_out/"

for file in ./include/mods/*/; do
	echo "$file"
	(cd "$file" && sh build.sh)
done

cp -r ./include/* ./_out/