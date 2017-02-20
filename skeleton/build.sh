#! /bin/bash
for p in $(ls local/plugins); do
    if [[ -e "local/plugins/$p/build.sh" ]]; then
        /bin/bash "local/plugins/$p/build.sh"
    fi
done
