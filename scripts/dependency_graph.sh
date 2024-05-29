#!/bin/bash
# This creates a go dependency graph of the project.
# You need to have godepgraph and graphviz installed.

set -e
command -v dot -- > /dev/null && echo "graphviz Found In \$PATH" || echo "graphviz Not Found in \$PATH"
command -v godepgraph -- > /dev/null && echo "godepgraph Found In \$PATH" || echo "godepgraph Not Found in \$PATH"
echo "some common offenders(standard lib, req,...) are excluded from the graph"
godepgraph -s -p golang.org,github.com/imroc/req,github.com/riverqueue/river/internal . | dot -Tsvg > output.svg && echo "see output.svg"