#!/bin/bash

read -p "Enter Application Name: " APP_NAME
if [ -z "$APP_NAME" ]; then
    echo "Error: Application name is required."
    exit 1
fi

read -p "Enter Project Root (default: .): " PROJECT_ROOT
PROJECT_ROOT=${PROJECT_ROOT:-.}

read -p "Enter Output Directory (default: bin): " OUTPUT_DIR
OUTPUT_DIR=${OUTPUT_DIR:-bin}

read -p "Enter Version (default: 1.0.0): " VERSION
VERSION=${VERSION:-1.0.0}

read -p "Enter Commit Hash (default: unknown): " COMMIT_HASH
COMMIT_HASH=${COMMIT_HASH:-unknown}

read -p "Zip binaries into $OUTPUT_DIR/compressed? (y/n): " DO_ZIP
read -p "Generate checksums? (y/n): " DO_CHECKSUM

PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64" "windows/arm64")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

if [[ "$DO_ZIP" =~ ^[Yy](es)?$ ]]; then
    mkdir -p "$OUTPUT_DIR/compressed"
fi

echo "Starting build for $APP_NAME version $VERSION"

for PLATFORM in "${PLATFORMS[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$PLATFORM"
    
    BASENAME="${APP_NAME}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        BINARY_NAME="${BASENAME}.exe"
    else
        BINARY_NAME="${BASENAME}"
    fi

    TARGET_PATH="$OUTPUT_DIR/$BINARY_NAME"

    echo "Building $BINARY_NAME..."

    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-s -w -X main.version=$VERSION -X main.buildDate=$BUILD_DATE -X main.commitHash=$COMMIT_HASH" \
        -o "$TARGET_PATH" \
        "$PROJECT_ROOT"

    if [ $? -ne 0 ]; then
        echo "Build failed for $GOOS/$GOARCH"
        exit 1
    fi

    if [[ "$DO_ZIP" =~ ^[Yy](es)?$ ]]; then
        pushd "$OUTPUT_DIR" > /dev/null
        
        if [ "$GOOS" = "windows" ]; then
            ZIP_NAME="${BASENAME}.zip"
            
            if command -v zip &> /dev/null; then
                zip -q -j "compressed/$ZIP_NAME" "$BINARY_NAME"
                echo "  > Created compressed/$ZIP_NAME"
            elif command -v powershell &> /dev/null; then
                echo "  > Zipping via PowerShell..."
                powershell -Command "Compress-Archive -Path '${BINARY_NAME}' -DestinationPath 'compressed/${ZIP_NAME}' -Force"
            else
                echo "  > Error: No zip tool found."
            fi
        else
            TAR_NAME="${BASENAME}.tar.gz"
            tar -czf "compressed/$TAR_NAME" "$BINARY_NAME"
            echo "  > Created compressed/$TAR_NAME"
        fi
        
        popd > /dev/null
    fi
done

if [[ "$DO_CHECKSUM" =~ ^[Yy](es)?$ ]]; then
    echo "Generating checksums..."
    cd "$OUTPUT_DIR" || exit
    touch checksums.txt

    hash_file() {
        local f="$1"
        if command -v sha256sum &> /dev/null; then
            sha256sum "$f"
        elif command -v shasum &> /dev/null; then
            shasum -a 256 "$f"
        fi
    }

    for file in *; do
        if [ -f "$file" ] && [ "$file" != "checksums.txt" ]; then
            hash_file "$file" >> checksums.txt
        fi
    done

    if [[ -d "compressed" ]]; then
        cd compressed || exit
        for file in *; do
            if [ -f "$file" ]; then
                hash_file "$file" >> ../checksums.txt
            fi
        done
        cd ..
    fi
    
    echo "Checksums generated."
fi

echo "Done."