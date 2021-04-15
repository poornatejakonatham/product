# product
A simple product catalog application.

## Create DB
```bash
sudo -u postgres psql -f schema.sql
```

## Bash
Set config in environment variables.
```bash
source ./env
```
If you get any error like this
```bash
bash: ./env-sample: Permission denied
```
Try following command.
```bash
chmod +x file_name
```

## Unit Testing
Simple testing.
```bash
go test -v -cover
```

Writing a "coverage profile". Results will be saved in a file.
```bash
go test -coverprofile=coverage.out
```

Function wise beakdown.
```bash
go tool cover -func=coverage.out
```

Get an HTML presentation of the source code decorated with coverage information.
```bash
go tool cover -html=coverage.out
```