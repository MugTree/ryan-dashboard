# Sensors project

## moving parts

- go + templ + datastar
- daisy ui and tailwind

Got a basic POC working where data (fake cpu usage data) comes in from an external api and gets fed through and is and rendered.

## To run for dev

```cmd
go run cmd/sensor/main.go
```

open another terminal

```cmd
make start-dev
```
