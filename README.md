# Scheduled Jobs

### tree

```
.
├── config
│   └── config.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── models
│   │   ├── job.go
│   │   └── job_history.go
│   ├── repository
│   │   ├── job.go
│   │   └── job_history.go
│   ├── scheduler
│   │   ├── ports.go
│   │   └── scheduler.go
│   └── task
│       ├── echo
│       │   └── echo.go
│       └── task.go
├── main.go
├── pkg
│   ├── logger
│   │   └── logger.go
│   ├── postgres
│   │   ├── config.go
│   │   └── postgres.go
│   └── uuid
│       └── uuid.go
├── README.md
└── schema.sql

11 directories, 19 files
```