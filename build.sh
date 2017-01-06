#! /bin/bash
make && make gz
for p in $(ls plugins); do
    if [[ -e "plugins/$p/build.sh" ]]; then
        /bin/bash "plugins/$p/build.sh"
    fi
done
