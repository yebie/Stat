#!/bin/bash
sleep 5 # for db start
PaddleStat -dbconnect=${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_DB}?sslmode=disable \
	-key=${COOKIE_KEY} \
	-version=${PADDLE_VERSION}
