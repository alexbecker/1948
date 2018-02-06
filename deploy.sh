#! /bin/bash
set -e
[[ $DEPLOY_ADDR =~ localhost ]] && DST=$DIR || DST="$DEPLOY_ADDR:$DIR"
echo "Deploying to $DST"
cd local
rsync -a server $DST
rsync -a env.sh $DST
rsync -a server_side/ $DST
rsync -a --delete-after static $DST
exec ./deploy.sh
