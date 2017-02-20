#! /bin/bash
[[ $HOST =~ localhost ]] && DST=$DIR || DST="$HOST:$DIR"
rsync -a go/bin/server $DST
rsync -a static $DST
rsync -a local/env.sh $DST
rsync -a local/server_side/ $DST
/bin/bash local/deploy.sh
