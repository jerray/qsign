# Qsign

[![Build Status](https://travis-ci.org/jerray/qsign.svg?branch=master)](https://travis-ci.org/jerray/qsign)

Generate signature for Golang struct value.

## Requirements

* Go version >= 1.8

## Signing Method

It uses the signing method widely used by tencent and wechat APIs.

For example, we have such data to be signed:

```
appid: wxd930ea5d5a258f4f
mch_id: 10000100
device_info: 1000
body: test
nonce_str: ibuaiVcKdpRxkhJA
```

First step, we make a query string using the data. Key is sorted by ASCII order. Then we have string A.

```go
A := "appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA"
```

Next, we prepend or append some secret data to A. Then we have our digest B. Here, we append a
secret key.

```go
B := A + "&key=192006250b4c09247ec02edce69f6a2d"
// appid=wxd930ea5d5a258f4f&body=test&device_info=1000&mch_id=10000100&nonce_str=ibuaiVcKdpRxkhJA&key=192006250b4c09247ec02edce69f6a2d
```

Finaly, calculate the checksum of the digest using some hash method. Then we get the signature.

```go
md5(B)
```

## Usage

Qsign implements the signing method mentioned before. It uses reflection to get structs' fields, determines which
fields will appear in the digest. By default field name is used as key. You can asign a "qsign" tag to that field
to change key string or ignore that field. Qsign also support "json", "yaml" and "xml" tags optionally.

Tag "qsign" has the highest priority. A field with tag `qsign:"-"` will be ignored.

```go
data := struct {
	AppId      string `qsign:"appid"`
	MchId      int    `qsign:"mch_id"`
	DeviceInfo string `qsign:"device_info"`
	Body       string `qsign:"body"`
	NonceStr   string `qsign:"nonce_str"`
	IgnoreMe   string `qsign:"-"`
}{
	AppId:      "wxd930ea5d5a258f4f",
	MchId:      10000100,
	DeviceInfo: "1000",
	Body:       "test",
	NonceStr:   "ibuaiVcKdpRxkhJA",
	IgnoreMe:   "won't be used to generate digest",
}

q := qsign.NewQsign(&qsign.Options{
	SuffixGenerator: func() string {
		return "&key=192006250b4c09247ec02edce69f6a2d"
	},
})

signature, _ := q.Sign(data)
fmt.Printf("%s", string(signature))
// 9a0a8659f005d6984697e2ca0a9cf3b7
```

## Limitations

Array and Slice types of field are not supported.
But if a field's type has a method `String() string`, Qsign will treat it as string field.

## License

MIT
