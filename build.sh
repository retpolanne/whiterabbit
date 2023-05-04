#!/bin/bash
# Taken from https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

mkdir out || true

package_name="whiterabbit"
platforms=("darwin/amd64" "darwin/arm64" "linux/amd64" "windows/amd64")
for platform in "${platforms[@]}" 
do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name=$package_name'-'$GOOS'-'$GOARCH

	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi
  env GOOS=$GOOS GOARCH=$GOARCH go build -o out/$output_name github.com/retpolanne/whiterabbit
done
