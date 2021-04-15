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
