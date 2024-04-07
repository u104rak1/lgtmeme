# LGTMeme
LGTMeme is a simple LGTM (Looks Good To Me) image creation service. You can quickly insert LGTM characters into your favorite image and copy markdown to the clipboard. Brighten up the monotonous code review and approval process with humorous LGTM image.





## URL





## Endpoints

<details>
<summary><b>Auth server</b></summary>

| Endpoint                   | Path                        |
| :------------------------- | :-------------------------- |
| Authorization Endpoint     | `/auth-api/authorize`       |
| Health Check               | `/auth-api/health`    |
| Token Endpoint             | `/connect/token`            |
| Login Form                 | `/login`                    |
| Access Token JWKS Endpoint | `/internal/v1/connect/jwks` |
| ID Token JWKS Endpoint     | `/connect/jwks`             |
</details>





## Supported Authentication/Authorization flow
<details>
<summary><b>Client Credentials Flow (with GET images)</b></summary>

```mermaid
sequenceDiagram
  participant browser as Browser
  participant client_server as Client Server
  participant auth_server as Auth Server
  participant resource_server as Resource Server
  participant db as DB
  participant Redis as Redis

  browser->>+client_server: Request home view (GET / )
  client_server->>+auth_server: Request Access Token (GET /auth-api/token)
  auth_server->>+db: Validate credentials
  db-->>-auth_server: Credentials valid
  auth_server->>auth_server: Generate Access Token
  auth_server-->>-client_server: Access Token
  client_server->>+Redis: Store Access Token
  Redis-->>-client_server: Cached Access Token
  client_server-->>-browser: HTML, CSS and JS
  browser->>+client_server: Request images (GET /client-api/images)
  client_server->>+Redis: Load Access Token
  Redis-->>-client_server: Access Token
  client_server->>+resource_server: Request images (GET /resource-api/images)
  resource_server->>+auth_server: Request public key (GET /auth-api/jwks)
  auth_server-->>-resource_server: public key
  resource_server->>+Redis: Store public key for one day
  Redis->>-resource_server: Cached public key
  resource_server->>resource_server: Validate the scope included in the access token
  resource_server->>+db: Get images
  db-->>-resource_server: images
  resource_server-->>-client_server: images
  client_server-->>-browser: images
```
</details>





<details>
<summary><b>Authorization Code Flow (Transition to admin screen)</b></summary>

```mermaid
sequenceDiagram
  participant browser as Browser
  participant client_server as Client Server
  participant auth_server as Auth Server
  participant db as DB
  participant Redis as Redis

  browser->>+client_server: Request admin view (GET /admin )
  client_server->>+Redis: Load Access Token or Refresh Token
  Redis-->>-client_server: Not found Access Token and Refresh Token
  client_server->>+Redis: Store state and nonce
  Redis-->>-client_server: Cached state and nonce
  client_server-->>-browser: 302 Redirect
  browser->>+auth_server: Request authorize (GET /auth-api/authorize)
  auth_server->>+db: Validate credentials
  db-->>-auth_server: Credentials valid
  auth_server->>+Redis: Load login session
  Redis-->>-auth_server: Not logged in
  auth_server->>+Redis: Store query parmas (pre authentication)
  Redis-->>-auth_server: Cached query params
  auth_server-->>-browser: 302 Redirect
  browser->>+auth_server: Request login view (GET /login )
  auth_server-->>-browser: HTML, CSS and JS
  browser->>+auth_server: Send username and password (POST /auth-api/login )
  auth_server->>+db: Validate username and password
  db-->>-auth_server: username and password valid
  auth_server->>+Redis: Store login session
  Redis-->>-auth_server: Cached login session
  auth_server->>+Redis: Load query params
  Redis-->>-auth_server: Create redirect URL using query params
  auth_server-->>-browser: 200 redirect URL
  browser->>+auth_server: Re Request authorize (GET /auth-api/authorize)
  auth_server->>+Redis: Store query parmas (authorization context)
  Redis-->>-auth_server: Cached query params
  auth_server-->>-client_server: 302 redirect (/client-api/admin/callback)
  client_server->>+auth_server: Request Access Token (GET /auth-api/token)
  auth_server->>+db: Validate credentials
  db-->>-auth_server: Credentials valid
  auth_server->>auth_server: Generate Access Token and Refresh Token and ID Token
  auth_server->>+db: Upsert Refresh Token
  db-->>-auth_server: Save Refresh Token
  auth_server-->>-client_server: Access Token and more
  client_server->>+Redis: Load public key (If not, GET /auth-api/jwks)
  Redis-->>-client_server: Validate ID Token with public key
  client_server->>+Redis: Validate state and nonce
  Redis-->>-client_server: state and nonce valid
  client_server->>+Redis: Store Access Token and Refresh Token
  Redis-->>-client_server: Cached Access Token and Refresh Token
  client_server-->>browser: 302 Redirect
  browser->>+client_server: Request admin view (GET /admin)
  client_server-->>-browser: HTML, CSS and JS
```
</details>





## Diagrams
<details>
<summary><b>Folder configuration diagram</b></summary>

```
lgtmeme/
├── .github/         # GitHub Actions configurations
├── .vscode/         # Visual Studio Code settings files
├── cmd/             # Application's entry point
│   └── lgtmeme/
│       └── main.go  # main function
├── config/          # Configuration files (DB, logger, constants, etc.)
├── db/              # Migration files and seed data for development
├── docker/          # Docker files
├── internal/        # Application's source code
│   ├── dto/         # Data Transfer Object structures
│   ├── handler/     # HTTP handlers (controllers)
│   ├── middleware/  # Middleware
│   ├── model/       # Data models
│   ├── repository/  # Data access layer (DB, Redis)
│   ├── service/     # Internal and external API requests and business logic
│   └── util/        # Utility functions and wrappers
├── script/          # Scripts files
├── test/            # Endpoint tests
└── view/            # Next.js
     ├── out/        # Static files exported by SSG
     └── src/        # Next.js source code
```
</details>

<details>
<summary><b>ER diagram</b></summary>

```mermaid
erDiagram
    health_checks {
        string key PK "PRIMARY KEY"
        string value "NOT NULL"
    }
    users {
        uuid id PK "PRIMARY KEY"
        string name "UNIQUE, NOT NULL"
        text password "NOT NULL"
        string role "NOT NULL"
    }
    oauth_clients {
        uuid id PK "PRIMARY KEY"
        string name "NOT NULL"
        uuid client_id "UNIQUE, NOT NULL"
        string client_secret "UNIQUE, NOT NULL"
        text redirect_uri "NOT NULL"
        text application_url "NOT NULL"
        string client_type "NOT NULL"
    }
    master_scopes {
        string code PK "PRIMARY KEY"
        text description "NOT NULL"
    }
    oauth_clients_scopes {
        uuid client_id PK "FOREIGN KEY"
        string scope_code PK "FOREIGN KEY"
    }
    master_application_types {
        string type PK "PRIMARY KEY"
    }
    oauth_clients_application_types {
        uuid client_id PK "FOREIGN KEY"
        string application_type PK "FOREIGN KEY"
    }
    refresh_tokens {
        string token PK "PRIMARY KEY"
        uuid user_id FK "FOREIGN KEY"
        uuid client_id FK "FOREIGN KEY"
        text scopes "NOT NULL"
    }
    images {
        uuid id PK "PRIMARY KEY"
        text url "NOT NULL"
        string keyword "NOT NULL"
        integer used_count "NOT NULL"
        boolean reported "NOT NULL"
        boolean confirmed "NOT NULL"
        timestamp created_at "NOT NULL"
    }

    users ||--|| refresh_tokens : ""
    oauth_clients ||--|| oauth_clients_scopes : ""
    master_scopes ||--|| oauth_clients_scopes : ""
    oauth_clients ||--|| oauth_clients_application_types : ""
    master_application_types ||--|| oauth_clients_application_types : ""
    oauth_clients ||--o{ refresh_tokens : ""
```
</details>





<details>
<summary><b>Architecture diagram</b></summary>

```mermaid

```
</details>





## Developer knowledge

## Redis
```
docker exec -it lgtmeme_redis redis-cli
keys *
get [key]
```

## Grant execution permission to the script
```
chmod +x ./script/*.sh
```

## Generate RSA key pair
```
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
cat private_key.pem | base64 | tr -d '\n'
cat public_key.pem | base64 | tr -d '\n'