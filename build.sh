rm -rf ./_out/assets
cp -r ./assets ./_out/assets

go build -o ./_out/main.exe
