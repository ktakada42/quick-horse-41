#!/bin/bash

if [ ! -f .env ]; then
	echo "Error: .env file not found."
	exit 1
fi

export $(cat .env | xargs)

curl https://ipv4.mydns.jp/login.html -u $ID:$PASS
