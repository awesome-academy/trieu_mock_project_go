#!/bin/bash

## This script tests the rate limiting middleware by sending concurrent requests to the /api/profile endpoint.

CONCURRENCY=10

do_curl() {
  local id=$1

  echo "==== [Request $id] START $(date '+%H:%M:%S.%3N') ===="

curl 'http://localhost:8080/api/profile' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'Accept-Language: vi,fr-FR;q=0.9,fr;q=0.8,en-US;q=0.7,en;q=0.6,ja;q=0.5' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InVzZXJAc3VuLWFzdGVyaXNrLmNvbSIsImV4cCI6MTc2ODMwMDQ0MywiaWF0IjoxNzY4MjE0MDQzfQ.zb2CueluJvFi7bRXakBQlZRjcyCkdIz8KJt0Q8iCVVU' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -b 'm=59b9:true' \
  -H 'Referer: http://localhost:8080/profile' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36' \
  -H 'X-Requested-With: XMLHttpRequest' \
  -H 'sec-ch-ua: "Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
    -w "\n[Request $id] HTTP:%{http_code} Time:%{time_total}s\n"

  echo "==== [Request $id] END ===="
  echo
}

export -f do_curl

seq 1 $CONCURRENCY | xargs -n1 -P$CONCURRENCY -I{} bash -c 'do_curl {}'
