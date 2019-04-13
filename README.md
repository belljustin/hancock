# Hancock

A cryptographic signing service.

# Server

A json REST server can be started with the command:

```
./hancock server
```

The following endpoints are available at `http://127.0.0.1:8000` by default.

## List Keys

- Method: `GET`
- Path:   `/keys` 
- Query Params:
    - `alg` Algorithm: filter the returned keys based on the algorithm. 

Response:
```
{
    "keys": [ Key ]
}
```

Get a list of keys.

## Get Key

- Method: `GET`
- Path:   `/keys/{id}`

Response:
```
{
    "key": Key
}
```

Get a single key.

## Create Key

- Method: `POST`
- Path:   `/keys`

Request:
```
{
    "alg": Algorithm
}
```

Response:
```
{
    "key": Key
}
```

Create a key with the provided Algorithm.

## Sign

- Method: `POST`
- Path:   `/keys/{id}/signature`

Request:
```
{
    "digest": string,
    "hash": Hash
}
```

Response:
```
{
    "signature": string
}
```

Sign the provided digest with a key. The `digest` is hashed with the function specified by `hash`.

# Models

## Key
object

```
{
    "id": UUID,
    "alg": Algorithm,
    "owner": string,
    "pub": string
}
```

- id: unique identifier of the key.
- alg: indicates the algorithm used to generate the key.
- owner: the owner of the key.
- pub: the public component of the key. 

## UUID
string

16 byte Universal Unique IDentifier as defined in [RFC 4122](https://tools.ietf.org/html/rfc4122).

## Algorithm
string
- rsa

A cryptographic signature algorithm.

## Hash
string
- sha256

The cryptographic hash used to create the digest.

# Client

The hancock binary also provides a cli for interacting with the server.

## New Key

Generate a new key with the specified `$ALG`.

```
$ ./hancock keys new -alg=$ALG

id: 8b40edd6-926d-40a6-a825-7e7db2a6aae9
alg: rsa
owner: belljust.in/justin
pub: MIIBCgKCAQEAztwZg/HaDqy1Iu37foCg+Ew3WA4YKISefKuTIK0t5OShyX1IjgR3BOSX8syN5TTfXITA6KfL/kDdUC1qWsM6zz08v57V888ICU7P9fhmARCPJl4L54XnO4BUZWjVI79V/M0T8dN+PQhanLJXIlF+01PJponvjr+LNWgYW4habxzl3MWECtYy5oKjRjyLfyltrEpchfBmefgdL353XWb7ftI+XwGQwLJLif9zTIvs88cAr1XXHcxlZW52i3pHYX1XWA3FIB8FB5ubaWFv0BvcHtfnADhqwWdNpbbkqzBXKaaKamAWtafwy3Zm61/i7xfkjwW5rPu5xjPpGzqJ7iqiHQIDAQAB   
```

flags:
- alg, required: the algorithm to use in key generation
