#!/bin/bash
cd "$(dirname "$0")"

# sudo apt install p7zip-full
7z -v99m a cache_7x6.protobuf.7z cache_7x6.protobuf
# git checkout --orphan cache
# git add cache_7x6.protobuf.7z.001
