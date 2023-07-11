drcom-go
===
A jlu drcom client implemented in Go programming language. Improved from go-jlu-drcom-client.

## Quick start
```bash
go build .
./drcom_go --username zhangs2121 --password Pa$$w0rb --mac 12:34:56:78:9a:bc
```
```bash
docker build -t drcom_go .
docker run --name drcom_go drcom_go
```

## Acknowledgement
- [nweclient.py](https://github.com/drcoms/jlu-drcom-client/blob/master/newclient.py)
- [jlu-drcom-java](https://github.com/drcoms/jlu-drcom-client/tree/master/jlu-drcom-java)
- [go-jlu-drcom-client](https://github.com/wucongyou/go-jlu-drcom-client/tree/6b4d6b4c839742279a26c133eab934692926fb9b)