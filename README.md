# PostgreSQL Driver Benchmarks: pgx vs sqlx

This project benchmarks the performance of two popular PostgreSQL drivers for Go: `pgx` and `sqlx`. It compares their performance in various scenarios, including simple queries, concurrent operations, and mixed workloads.

## Prerequisites

- Go 1.20 or later
- Docker
- PostgreSQL

## Setup

1. Start a PostgreSQL container:

```
docker run --name postgres-container -e POSTGRES_PASSWORD=mysecretpassword -d postgres  
```

2. Run the benchmarks:

```
go test -bench=. -benchmem -benchtime=10s > benchmark_results.txt 2>&1
```

## Results

The results will be saved in `benchmark_results.txt`.

