#! /bin/bash
source env.sh
find plugins -path "*/server_side/install.sh" | xargs /bin/bash
