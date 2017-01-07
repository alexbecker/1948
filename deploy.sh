#! /bin/bash
[[ $HOST =~ localhost ]] && DST=$DIR || DST="$HOST:$DIR"
rsync -a go/bin/server $DST
rsync -a static $DST
rsync -a local/env.sh $DST
rsync -a local/server_side/ $DST
for p in $(ls plugins); do
    if [[ -e "plugins/$p/server_side" ]]; then
        rsync -a --relative "plugins/$p/server_side/" $DST
    fi
done
rsync -a install_server.sh $DST
