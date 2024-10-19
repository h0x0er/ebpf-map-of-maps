#!/bin/bash

set -eo pipefail

VMLINUX="./bpf/includes/vmlinux.h"

log() {
    printf "$(date) [%s] %s\n" "$1" "$2"
}

vmlinux() {
    if [ ! -e $VMLINUX ]; then
        bpftool btf dump file /sys/kernel/btf/vmlinux format c >$VMLINUX
        log "vmlinux" "created"
    fi

}

generate() {

    # make sure vmlinux exists
    vmlinux

    log "generate" "generating files"

    go generate .

    log "generate" "done"

}

build-main() {

    log "build" "building main"
    go build -o maps
    log "build" "done"

}

all() {
    generate
    build-main
}

logk() {
    sudo cat /sys/kernel/debug/tracing/trace_pipe
}

# run cmd
$1
