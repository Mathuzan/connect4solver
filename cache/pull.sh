#!/bin/bash
cd "$(dirname "$0")"

# sudo apt install p7zip-full
wget https://github.com/igrek51/connect4endgames/raw/cache/cache/cache_7x6.protobuf.7z.001
7z x cache_7x6.protobuf.7z.001
rm cache_7x6.protobuf.7z.001
echo "Cached endgames ready"
