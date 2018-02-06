#! /bin/bash
[[ $DEPLOY_ADDR =~ localhost ]] && DST=$DIR || DST="$DEPLOY_ADDR:$DIR"
for p in $(ls plugins); do
    if [[ -e "plugins/$p/server_side" ]]; then
        rsync -a --relative "plugins/$p/server_side/" $DST
    fi
done
