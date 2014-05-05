#!/bin/sh
if [ -f src ]; then
 rm src -f
fi
echo Welcome to use GoMvc!
echo GoMvc  is a lightweight web framework ,QQ Group: 184572648
PWD=`pwd`/..
echo $GOPATH
export GOPATH=${PWD}:${GOPATH}
echo $GOPATH
echo Building ...
go build .
if [ -f src ]; then 
 echo Build succeed!!!!
fi

echo Press Enter to continue...
read -n 1
