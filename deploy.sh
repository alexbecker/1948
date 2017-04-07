#! /bin/bash
[[ $DEPLOY_ADDR =~ localhost ]] && DST=$DIR || DST="$DEPLOY_ADDR:$DIR"
rsync -a go/bin/server $DST
rsync -a static $DST
rsync -a local/env.sh $DST
rsync -a local/server_side/ $DST
/bin/bash local/deploy.sh
