# Cruds

Cruds is a 

## Commands

### Testing

#### GET data
`curl http://localhost:8080/cards`
`curl http://localhost:8080/cards/1`

#### POST data
`curl -X POST -H "Content-Type: application/json" -d @testdata/card-0001.json http://localhost:8080/cards`

### Performance

`hey -z 1m http://localhost:8080/cards/1`
