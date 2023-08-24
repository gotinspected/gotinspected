#!/bin/bash

# This script expects as first argument the path to a gitinspector source folder
# and as second argument the path to a git repository

gitinspector_src=$1
repo=$2

docker run --user inspector --rm -it -v $repo:/repo -v $gitinspector_src:/gitinspector -w /repo gitinspector /bin/bash
