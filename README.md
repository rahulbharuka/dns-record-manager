# DNS Record Manager
This Go app automates DNS records management for a single domain (specified by env `$DOMAIN`).

### Features
- View list of available servers.
- View list of all published DNS records.
- Add/Remove a server to/from cluster (subdomain).
- All endpoints are idempotent.

### APIs
- `GET /servers` - Returns list of all servers with DNS status and allowed actions.
- `GET /dns-records` - Returns list of all DNS A records published for the domain.
- `POST /servers/:id/add` - Adds the server to rotation (by adding DNS A record).
- `POST /servers/:id/remove` - Removes the server from rotation (by removing DNS A record).

### External Dependencies
- Uses `ClearDB MySQL` as a persistent data storage layer.
- Uses `Gin` web framework for routing.
- Uses `AWS Route53` service for hosting DNS records.

### Assumptions
- A server (IP) can be part of at most one cluster (subdomain).

### Building binary
For Linux:
    
    $ env GOOS=linux GOARCH=amd64 go build -o bin/dns-record-manager -v .

For MacOS:
    
    $ go build -o bin/dns-record-manager -v .

### Running app locally
    $ heroku local web

Please make sure that following environment variable are set.

    AWS_ACCESS_KEY_ID
    AWS_SECRET_ACCESS_KEY
    CLEARDB_DATABASE_URL
    DATABASE_URL
    DB_HOST
    DB_NAME
    DB_PASSWORD
    DB_USER
    DOMAIN
    HOSTED_ZONE_ID
    PORT

### Testing
- Added basic tests.
