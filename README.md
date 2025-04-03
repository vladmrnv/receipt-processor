# Receipt Processor

## Requirements

- Docker (or run with Go 3.23.4)

## Dependencies
```go
require github.com/google/uuid v1.6.0
```
Used for generating UUID values for IDs

## Running with Docker
1. Clone this repository
2. Build this Docker Image
```zsh
docker build -t receipt-processor .
```
3. Run this Docker container
```zsh
docker run -p 8080:8080 receipt-processor
```
4. This process is now running on http://localhost:8080

## Running the Tests
In order to run full test suite

```zsh
go test ./...
```

In order to run test coverage
```zsh
go test ./... -cover
```

Run specific test suites
```zsh
go test ./handlers
go test ./store
etc...
```

## Testing API Endpoints

This was all tested on Postman

### Process Receipt
```go
POST /receipts/process
```

Example request:
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```

Expected Response:
```json
{
    "id": "ef8ee7f4-ecc2-410e-9c80-1bbb1aee28fe"
}
```

### Get Point Values

```go
GET /receipts/{id}/points
```

Expected Response:
```json
{
    "points": 28
}
```
