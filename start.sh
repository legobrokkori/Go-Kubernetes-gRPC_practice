#!/bin/sh

source app.env

echo "run db migtation"

echo $str
/app/migrate -path /app/migration -database "${DB_SOURCE}" -verbose up

echo "start the app"
exec "$@"
