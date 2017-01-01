#! /bin/bash
source local/env.sh

[[ $HOST =~ localhost ]] && DST=$DIR || DST="$HOST:$DIR"
rsync -a go/bin/server $DST
rsync -a static $DST
rsync -a local/env.sh $DST
rsync -a init_server.sh $DST
rsync -a schema.sql $DST
