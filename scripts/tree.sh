#!/bin/bash

# Generate directory structure using tree and convert to markdown format
echo "# 项目目录结构" > dir.md
echo '```' >> dir.md
tree -I 'node_modules|vendor|*.git' >> dir.md
echo '```' >> dir.md