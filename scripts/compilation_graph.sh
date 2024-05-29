#!/bin/bash
go install github.com/unravelin/actiongraph@latest
sudo apt install -y -q graphviz # inkscape
echo "CLEAN BUILD"
go clean -cache
go build -debug-actiongraph=compile.json
actiongraph top --file compile.json
echo "OUTPUTTING GRAPH"
actiongraph graph --file compile.json > ./output.dot
echo "RENDERING TO SVG"
cat output.dot | dot -Tsvg > ./cache/compilestats_clean.svg

echo "REBUILD"
go build -debug-actiongraph=compile.json
actiongraph top --file compile.json
echo "OUTPUTTING GRAPH"
actiongraph graph --file compile.json > ./output.dot
echo "RENDERING TO SVG"
cat output.dot | dot -Tsvg > ./cache/cmopilestats.svg # && inkscape -z -w 32000 -h 10000 output.svg -o output.png