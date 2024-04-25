# Folder configuration

```
lgtmeme/
├── .github/         # GitHub Actions configurations
├── .vscode/         # Visual Studio Code settings files
├── cmd/             # Application's entry point
│   └── lgtmeme/
│       └── main.go  # main function
├── config/          # Configuration files
├── db/              # Migration files and seed data
├── docker/          # Docker files
├── document/        # Documents
├── internal/        # Application's source code
│   ├── dto/         # Data Transfer Object structures
│   ├── handler/     # HTTP handlers (controllers)
│   ├── middleware/  # Middleware
│   ├── model/       # Data models
│   ├── repository/  # Data access layer (DB, Redis)
│   ├── service/     # Internal and external API requests and business logic
│   ├── setup/       # setup Echo
│   └── util/        # Utility functions and wrappers
├── script/          # Scripts files
├── test/            # Test
│   ├── endpoint/    # Endpoint tests
│   ├── mock/        # mock files
│   └── testutil/    # Utility functions fot test
└── view/            # Next.js
    ├── out/         # Static files exported by SSG
    └── src/         # Next.js source code
```
