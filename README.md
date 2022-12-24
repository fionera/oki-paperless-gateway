# Oki Paperless Gateway

A simple tool to forward scans from a OKI Printer to Paperless. Just setup a profile with http as protocol and enter the path to this tool. 

**You have to prepend // to the url**
´//gw.example.com´

## Config
Set the following env vars:
```
PAPERLESS_URL="https://example.com"
PAPERLESS_USER="user"
PAPERLESS_PASS="pass"
```