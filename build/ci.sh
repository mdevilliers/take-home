#!/bin/bash
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. go vet (http://golang.org/cmd/vet)
# 2. race detector (http://blog.golang.org/race-detector)
# 3. test coverage (http://blog.golang.org/cover)
# 4. build the main entry points

set -e

go vet ./...
go test -race ./...

# Run test coverage on each subdirectories and merge the coverage profile. 
echo "mode: count" > profile.cov
 
# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
go test -covermode=count -coverprofile=$dir/profile.tmp $dir
if [ -f $dir/profile.tmp ]
then

cat $dir/profile.tmp | tail -n +2 >> profile.cov
rm $dir/profile.tmp
fi
fi
done
 
go tool cover -func profile.cov

echo "building main applications"

echo "building server"
go build github.com/mdevilliers/take-home/cmd/server/ 

echo "finished"