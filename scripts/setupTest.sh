#!/bin/bash

export POSTGRES_PASSWORD=gusampaio_pass
export POSTGRES_USER=gusampaio
export POSTGRES_DB=hbday_db
export DATABASE_HOST=localhost
export DATABASE_PORT=5432

echo "Setting up database"
container_id=$(docker run -d --name hbday_db -p 5432:5432 -e POSTGRES_USER=gusampaio -e POSTGRES_PASSWORD=gusampaio_pass -e POSTGRES_DB=hbday_db -v db:/var/lib/postgresql/data postgres)

echo "Running Tests"
go test

echo "Tearing down db"
docker stop $container_id
docker rm $container_id