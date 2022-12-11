ROOT=$(pwd)

# TODO: replace platform and os information to be dynamic
PROTOC_VERSION="21.9"
PROTO_ARTIFACT="protoc-${PROTOC_VERSION}-osx-universal_binary.zip"
PROTOC_URL="https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTO_ARTIFACT}"

# download and unzip protoc
mkdir -p bin && mkdir -p bin/protoc-bin
(cd bin/protoc-bin && curl -LO $PROTOC_URL && unzip -o $PROTO_ARTIFACT)
mv -f ./bin/protoc-bin/bin/protoc ./bin/protoc
rm -rf ./protos/google && mv -f ./bin/protoc-bin/include/google ./protos/google
rm -rf bin/protoc-bin

GOOGLE_IMPORT_PATH=$ROOT/protos
PROTOC_PATH=$ROOT/bin/protoc

# compile statspb
$PROTOC_PATH --proto_path=$ROOT/protos/statspb --go_out=$ROOT/protos/statspb $ROOT/protos/statspb/stats.proto