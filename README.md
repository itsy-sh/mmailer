

```text
SELECT_STRATEGY=RoundRobin \
SERVICES="generic:smtp://user:pass@smtp.server.com:25 mailjet:pubkeyXXXX:secretkeyYYYY" \
go run mmailerd.go 
Services:
 - Generic: posthooks are not implmented, adding smtp://user:pass@smtp.server.com:25
 - Mailjet: add the following posthook url  example.com/path/to/mmailer/posthook?key=&service=mailjet
Select Strategy: RoundRobin
Retry Strategy:  None

> Send mail by HTTP POST example.com/path/to/mmailer/send?key=

Starting server, :8080
```