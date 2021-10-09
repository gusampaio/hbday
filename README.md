### Usage

Run the applicattion: 
`$ docker compose up`

Create a user: 
`$ curl http://0.0.0.0:8080/hello/<username> --include --header "Content-Type: application/json"  --request "PUT" --data '{"date_of_birth": "YYYY-MM-DD"}'`

Check existent users: 
`$ curl http://0.0.0.0:8080/hello/all`

Check specific user: 
`$ curl http://0.0.0.0:8080/hello/<username>`