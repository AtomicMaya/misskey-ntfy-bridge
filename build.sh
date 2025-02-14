echo "Building artifacts in ./build"
VERSION=$1

build () {
	GOOS=$1
	GOARCH=$2
	EXT=""
	if [ "$1" = "windows" ]
	then
		EXT=".exe"
	fi

	GCO_ENABLED=0 go build -tags netgo -a -ldflags "-w" -o build/misskey-ntfy-bridge-v$VERSION-$GOOS-$GOARCH$EXT ./app
}

#build darwin    	amd64
#build darwin    	arm64
#build dragonfly 	amd64
#build freebsd   	386
#build freebsd   	amd64
#build freebsd   	arm
#build linux     	386
build linux			amd64
build linux			arm
build linux			arm64
# build linux	 		ppc64
# build linux	 		ppc64le
# build linux	 		mips
# build linux	 		mipsle
# build linux	 		mips64
# build linux	 		mips64le
# build netbsd		386
# build netbsd		amd64
# build netbsd		arm
# build openbsd		386
# build openbsd		amd64
# build openbsd		arm
# build plan9			386
# build plan9			amd64
# build solaris		amd64
build windows 	386
build windows 	amd64

echo "Done!"