FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    clang libclang-14-dev libc++-dev g++ util-linux ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY scan.cpp /tmp/scan.cpp
RUN clang++ -std=c++17 /tmp/scan.cpp -I/usr/lib/llvm-14/include \
    -L/usr/lib/llvm-14/lib -lclang -Wl,-rpath,/usr/lib/llvm-14/lib \
    -o /usr/local/bin/scan && rm /tmp/scan.cpp
USER 65534:65534
WORKDIR /work
