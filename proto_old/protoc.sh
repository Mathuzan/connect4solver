get_protoc () {
    PROTOBUF_VERSION=$1
    PROTOC=""

    if command -v protoc >/dev/null 2>&1 && [[ $(protoc --version) == "libprotoc $PROTOBUF_VERSION" ]]; then
        PROTOC=`which protoc`
        echo "Found protobuf "${PROTOBUF_VERSION}" installed"
    else
        platform=`uname`
        if [[ "$platform" == 'Linux' ]]; then
            platform='linux'
        elif [[ "$platform" == 'Darwin' ]]; then
            platform='osx'
        else
            echo "Unsupported OS $platform"
            exit 1
        fi

        TARGET="protoc-$platform"
        if ! command -v $TARGET/bin/protoc >/dev/null 2>&1; then

            ARCHIVE=protoc-${PROTOBUF_VERSION}-${platform}-x86_64.zip

            rm -rf $TARGET
            mkdir $TARGET
            wget --no-verbose --timeout=60 --tries=3 https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/${ARCHIVE}
            unzip ${ARCHIVE} -d $TARGET
            rm ${ARCHIVE}

            echo "Installed protobuf "${PROTOBUF_VERSION}" locally under `pwd`/$TARGET"
        fi

        PROTOC="$TARGET/bin/protoc"
        echo "Using local protobuf installation under `pwd`/$TARGET"
    fi

}
