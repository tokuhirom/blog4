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

# curl -f はステータスコードが 400 以上の場合にエラーとなる

curl -f --no-progress-meter -u "$access_token_id:$access_token_secret" -o 'app.json' \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id"

# app.json から deploy.json を生成
# - container_registry の認証系フィールド (server / username / password) を削除する。
#   ghcr.io の public パッケージなので pull に認証は要らない。server だけを送ると
#   AppRun API が Username と Password も必須と判定して 400 を返す:
#     "Authentication to Container Registry requires Server, Username, and Password."
#   どこから引くかは image のフルパス (ghcr.io/...) で決まるので server は不要。
# - all_traffic_available を true に設定
jq --arg image "$image" '
  del(.components[0].deploy_source.container_registry.server,
      .components[0].deploy_source.container_registry.username,
      .components[0].deploy_source.container_registry.password)
  | .all_traffic_available = true
  | .components[0].deploy_source.container_registry.image = "\($image)"
' app.json > deploy.json

# PATCH が 400 等で失敗したときにレスポンスボディ(エラー原因)を確認できるよう、
# --fail-with-body を使い、ボディは /dev/null に捨てずに標準出力へ出す
curl --fail-with-body --no-progress-meter -u "$access_token_id:$access_token_secret" -X PATCH -d '@deploy.json' \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id"

curl -f --no-progress-meter -u "$access_token_id:$access_token_secret" \
  "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api/applications/$apprun_app_id/versions" | jq
