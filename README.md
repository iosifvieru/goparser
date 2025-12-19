# Go nginx log parser

This is a basic Nginx log parser written in Go.

After running this program for now you will receive a distribution of HTTP status codes from the access file.

This is the output for 3.5GB of log files available at [Harvard Dataverse - Online shopping store - Web server Logs](https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/3QBYB5)

```yml
Total requests: 10365077
HTTP 404: 105011 requests (1.01 %)
HTTP 304: 340228 requests (3.28 %)
HTTP 400: 529 requests (0.01 %)
HTTP 502: 798 requests (0.01 %)
HTTP 408: 112 requests (0.00 %)
HTTP 200: 9579824 requests (92.42 %)
HTTP 301: 67553 requests (0.65 %)
HTTP 500: 14266 requests (0.14 %)
HTTP 403: 5634 requests (0.05 %)
HTTP 206: 3 requests (0.00 %)
HTTP 302: 199835 requests (1.93 %)
HTTP 504: 103 requests (0.00 %)
HTTP 499: 50852 requests (0.49 %)
HTTP 401: 323 requests (0.00 %)
HTTP 405: 6 requests (0.00 %)
```

## Usage

```
    go build -o main.out main.go
    ./main.out -file=PATH_TO_FILE
```
