#! /bin/bash
# Place any shell commands needed to initialize your server here.
source env.sh
find plugins -path "*/server_side/install.sh" | xargs /bin/bash
