# Web 3 Storage Proxy Server

## APIS

### Upload File

#### Request

```json
POST /api/v1/files

{
    "data": "asjkhdkjashdas", // File content
    "fileName": "test.txt",
    "path": "/", // Default is /
}
```

#### Response

```json
201 Created

{
    "fileName": "test.txt",
    "fileSize": 12345,
    "path": "/",
    "dataHash": "hash",
    "dataId": "00000000-0000-0000-0000-000000000001", // Internal Data Id, not the CID
    "downloadUrl": "https://ipfs.tech/12345",
    "dateCreated": "2020-01-01T00:00:00.000Z",
    "dateModified": "2020-01-01T00:00:00.000Z"
}
```

### List File

#### Request
```json
POST /api/v1/files

{
    "path": "/"
}
```

#### Response

```json
200 OK

{
 "files": [
        {
            "fileName": "test.txt",
            "fileSize": 12345,
            "path": "/",
            "dataHash": "hash",
            "dataId": "00000000-0000-0000-0000-000000000001", // Internal Data Id, not the CID
            "downloadUrl": "https://ipfs.tech/12345",
            "dateCreated": "2020-01-01T00:00:00.000Z",
            "dateModified": "2020-01-01T00:00:00.000Z"
        }
        ...
    ]
}
```

### Get File

#### Request
```json
POST /api/v1/files

{
    "path": "/test.txt"
}
```

#### Response

```json
200 OK

{
    "fileName": "test.txt",
    "fileSize": 12345,
    "path": "/",
    "dataHash": "hash",
    "dataId": "00000000-0000-0000-0000-000000000001", // Internal Data Id, not the CID
    "downloadUrl": "https://ipfs.tech/12345",
    "dateCreated": "2020-01-01T00:00:00.000Z",
    "dateModified": "2020-01-01T00:00:00.000Z"
}
```

### Update File

#### Request

```json
POST /api/v1/files

{
    "data": "hellothere", // File content
    "fileName": "test.txt",
    "path": "/", // Default is /
}
```

#### Response

```json
201 Created

}
    "fileName": "test.txt",
    "fileSize": 123,
    "path": "/",
    "dataHash": "hash",
    "dataId": "00000000-0000-0000-0000-000000000001", // Internal Data Id, not the CID
    "downloadUrl": "https://ipfs.tech/12345",
    "dateCreated": "2020-01-01T00:00:00.000Z",
    "dateModified": "2020-01-02T00:00:00.000Z
}
```

### Delete File

#### Request
```json
POST /api/v1/files

{
    "path": "/test.txt"
}
```

#### Response

```json
200 OK
```


## Database Schema

```sql
CREATE TABLE files IF NOT EXISTS


```

a/a/b/a/a
a/a/c/a/a
d/a/c/a/a

select * from files where path like 'a/%'

## Improvements

* Add auth, split file system by user