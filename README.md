# pProxy Server

This is the pProxy server, use with the pProxy client

To get started run the below commands on linux (You will need to adjust for windows)

```bash
git clone https://github.com/PyroChiliarch/pProxy.git
cd pProxy
go get github.com/google/uuid
go run .
```

By default, the proxy listens on port 8080.
This is currently hardcoded, you can change it in main.go line 36

The pProxy command line tool and library are available here:
https://www.lexaloffle.com/bbs/?tid=141188
