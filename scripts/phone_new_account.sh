#!/bin/sh
set -eux

DATA_DIR="/User/Containers/Data/Application/E1ED6F9C-FFCE-44D7-9FA9-63AE60B479DC"

rm -rf "$DATA_DIR"
cp -rv "~/FIFA_INIT" "$DATA_DIR"
curl "http://setreplace/$VENDOR_ID"

open com.ea.ios.fifaultimate