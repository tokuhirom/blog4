#!/bin/bash

set -ex

# 引数の数を確認
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <access_token_id> <access_token_secret> <apprun_app_id> <image>"
    exit 1
fi

# 引数を変数に割り当て
access_token_id="$1"
access_token_secret="$2"
apprun_app_id="$3"
image="$4"

curl -s -u "$access_token_id:$access_token_secret" -o 'app.json' \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id"

# app.json から deploy.json を生成
# - container_registry.server と container_registry.username を削除(こうしないとエラーになる)
# - all_traffic_available を true に設定
jq --arg image "$image" '
  del(.components[0].deploy_source.container_registry.server, .components[0].deploy_source.container_registry.username)
  | .all_traffic_available = true
  | .components[0].deploy_source.container_registry.image = "\($image)"
' app.json > deploy.json

curl -s -u "$access_token_id:$access_token_secret" -X PATCH -d '@deploy.json' -o /dev/null \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id"

curl -s -u "$access_token_id:$access_token_secret" \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id/versions" | jq
