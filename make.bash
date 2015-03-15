#!/usr/bin/env bash
# Copyright 2014 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

export GOROOT=$GOROOT_ANDROID

mkdir -p jni/armeabi
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=7 \
	$GOROOT/bin/go build -ldflags="-shared" -o jni/armeabi/libpongo.so .
ndk-build NDK_DEBUG=1
ant debug
