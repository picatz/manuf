# manuf

Go package and CLI tool for listing OUIs.

## Install

```console
$ go install github.com/picatz/manuf@latest
...
```

## Usage

The `manuf` CLI tool can be used with tools like `grep` and `jq` to filter results. Records are fetched over HTTPS on first use
from [`manuf.csv`]("https://raw.githubusercontent.com/picatz/manuf/main/manuf.csv"), which is then cached in a local directroy,
`/Users/$USER/Library/Caches/manuf.csv` on macOS. After 30 days, the records are fetched again to refresh the cache.

```go
$ manuf | grep "Apple, Inc."
{"Registry":"MA-L","Assignment":"608B0E","OrganizationName":"Apple, Inc.","OrganizationAddress":"1 Infinite Loop Cupertino CA US 95014"}
{"Registry":"MA-L","Assignment":"88B291","OrganizationName":"Apple, Inc.","OrganizationAddress":"1 Infinite Loop Cupertino CA US 95014"}
{"Registry":"MA-L","Assignment":"C42AD0","OrganizationName":"Apple, Inc.","OrganizationAddress":"1 Infinite Loop Cupertino CA US 95014"}
{"Registry":"MA-L","Assignment":"CCD281","OrganizationName":"Apple, Inc.","OrganizationAddress":"1 Infinite Loop Cupertino CA US 95014"}
...
$ manuf | jq 'select(.OrganizationName == "Apple, Inc.")'
...
$ manuf | grep "Apple, Inc." | wc -l
    973
$ manuf | jq -r .OrganizationName | sort -n | uniq -c | sort -rn | head -n 15
1013 Cisco Systems, Inc
 973 Apple, Inc.
 906 HUAWEI TECHNOLOGIES CO.,LTD
 687 Samsung Electronics Co.,Ltd
 490 Intel Corporate
 380 Huawei Device Co., Ltd.
 343 ARRIS Group, Inc.
 270 IEEE Registration Authority
 267 zte corporation
 257 Texas Instruments
 229 Private
 154 TP-LINK TECHNOLOGIES CO.,LTD.
 150 Hewlett Packard
 148 Dell Inc.
 139 Juniper Networks
$ manuf | jq -r .Registry | sort | uniq
CID
IAB
MA-L
MA-M
MA-S
$ manuf | jq 'select(.Registry == "MA-L")'
...
```
