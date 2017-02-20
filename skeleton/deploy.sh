#! /bin/bash
[[ $HOST =~ localhost ]] && DST=$DIR || DST="$HOST:$DIR"
cd local
for p in $(ls plugins); do
    if [[ -e "plugins/$p/server_side" ]]; then
        rsync -a --relative "plugins/$p/server_side/" $DST
    fi
done
