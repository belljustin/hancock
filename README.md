# Hancock

A cryptographic signing service.

# API

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

A cryptographic signature algorithm.

## Hash
string

The cryptographic hash used to create the digest.
