#!/bin/bash

# This script expects as first argument the path to a gitinspector source folder
# and as second argument the path to a git repository

gitinspector_src=$(readlink -f $1)
repo=$(readlink -f $2)

docker run --user inspector --rm -it -v $repo:/repo -v $gitinspector_src:/gitinspector -w /repo gitinspector /bin/bash
