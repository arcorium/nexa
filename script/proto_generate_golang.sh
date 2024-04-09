#!/bin/bash

WORKDIR=..
PROTODIR=$WORKDIR/proto
OUTDIR=$PROTODIR/generated/golang

#SUBDIR_COUNT=find $PROTODIR/$1 -maxdepth 1 -type d | wc -l
#
#echo $SUBDIR_COUNT
#
#if [$SUBDIR_COUNT == 1]
#then
# protoc --proto_path=$PROTODIR \
#   --go_out=$OUTDIR --go_opt=paths=source_relative\
#   --go-grpc_out=$OUTDIR --go-grpc_opt=paths=source_relative\
#   $PROTODIR/$1/*.proto
#fi

for entry in "$PROTODIR/$1"/*
do
#  echo "$enty"
 protoc --proto_path=$PROTODIR \
   --go_out=$OUTDIR --go_opt=paths=source_relative\
   --go-grpc_out=$OUTDIR --go-grpc_opt=paths=source_relative\
   $entry/*.proto
done
