run-postgres:
    docker run -d --name hbday_db -e POSTGRES_USER=gusampaio -e POSTGRES_PASSWORD=gusampaio_pass -e POSTGRES_DB=hbday_db -v db:/var/lib/postgresql/data postgres

docker run -v migrations:migrations migrate/migrate -path=/migrations/ -database postgres://localhost:5432/hbday_db up 2