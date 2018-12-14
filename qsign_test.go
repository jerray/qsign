package qsign

import (
	"testing"
)

type weixinPayApp struct {
	AppID string `qsign:"appId"`
}

type weixinPayPackage struct {
	// used to test nested struct
	*weixinPayApp
	TimeStamp int64  `qsign:"timeStamp"`
	NonceStr  string `qsign:"-"`
	Package   string `qsign:"package"`
	SignType  string `qsign:"signType"`
	PaySign   string `qsign:"paySign"`
}

func TestQsignSetDelimiter(t *testing.T) {
	q := NewQsign(Options{})
	cases := []string{",", "|", "#", "&"}
	for _, c := range cases {
		q.SetDelimiter(c)
		actual := q.delimiter
		if actual != c {
			t.Errorf("expect delimiter is set to %s, actual is %s", c, actual)
		}
	}
}

func TestQsignSetConnector(t *testing.T) {
	q := NewQsign(Options{})
	cases := []string{":", ">", "-", "_"}
	for _, c := range cases {
		q.SetConnector(c)
		actual := q.connector
		if actual != c {
			t.Errorf("expect connector is set to %s, actual is %s", c, actual)
		}
	}
}

func TestQsignDigest(t *testing.T) {
	q := NewQsign(Options{})

	cases := []struct {
		input  interface{}
		expect string
	}{
		{
			input: struct {
			}{},
			expect: "",
		},
		{
			input: weixinPayPackage{
				weixinPayApp: &weixinPayApp{AppID: "wx6cfc34d48f33effe"},
				TimeStamp:    1503117550,
				NonceStr:     "9446",
				Package:      "prepay_id=wx20170819124333185b7b54140976921757",
				SignType:     "",
				PaySign:      "5C082C2524C0407B61053F82C584B527",
			},
			expect: "appId=wx6cfc34d48f33effe&package=prepay_id=wx20170819124333185b7b54140976921757&paySign=5C082C2524C0407B61053F82C584B527&timeStamp=1503117550",
		},
		{
			input: struct {
				Action          string
				Nonce           int
				Region          string
				SecretID        string `qsign:"SecretId"`
				SignatureMethod string
				Timestamp       string
				instanceIds0    string `qsign:"instanceIds_0"`
			}{
				Action:          "DescribeInstances",
				Nonce:           11886,
				Region:          "gz",
				SecretID:        "AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA",
				SignatureMethod: "HmacSHA256",
				Timestamp:       "1465185768",
				instanceIds0:    "ins-09dx96dg",
			},
			expect: "Action=DescribeInstances&Nonce=11886&Region=gz&SecretId=AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA&SignatureMethod=HmacSHA256&Timestamp=1465185768&instanceIds_0=ins-09dx96dg",
		},
	}

	for _, c := range cases {
		d, _ := q.Digest(c.input)
		actual := string(d)
		if actual != c.expect {
			t.Errorf("expect digest is %s, actual is %s", c.expect, actual)
		}
	}
}

func TestQsignSign(t *testing.T) {
	q := NewQsign(Options{})

	cases := []struct {
		input  interface{}
		expect string
	}{
		{
			input: struct {
			}{},
			expect: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			input: weixinPayPackage{
				weixinPayApp: &weixinPayApp{AppID: "wx6cfc34d48f33effe"},
				TimeStamp:    1503117550,
				NonceStr:     "9446",
				Package:      "prepay_id=wx20170819124333185b7b54140976921757",
				SignType:     "MD5",
				PaySign:      "5C082C2524C0407B61053F82C584B527",
			},
			expect: "3061a228d86f9fe42c3f260591ce4865",
		},
		{
			input: struct {
				Action          string
				Nonce           int
				Region          string
				SecretID        string `qsign:"SecretId"`
				SignatureMethod string
				Timestamp       string
				instanceIds0    string `qsign:"instanceIds_0"`
			}{
				Action:          "DescribeInstances",
				Nonce:           11886,
				Region:          "gz",
				SecretID:        "AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA",
				SignatureMethod: "HmacSHA256",
				Timestamp:       "1465185768",
				instanceIds0:    "ins-09dx96dg",
			},
			expect: "db10eadbfe3f5a84a344020d73577b52",
		},
	}

	for _, c := range cases {
		d, _ := q.Sign(c.input)
		actual := string(d)
		if actual != c.expect {
			t.Errorf("expect sign is %s, actual is %s", c.expect, actual)
		}
	}
}

func TestQsignPrefixSuffixGenerators(t *testing.T) {
	cases := []struct {
		input  interface{}
		prefix Generator
		suffix Generator
		expect string
	}{
		{
			input: struct {
			}{},
			prefix: func() string {
				return "p"
			},
			suffix: func() string {
				return "s"
			},
			expect: "ps",
		},
		{
			input: struct {
				AppID     string `qsign:"appId"`
				TimeStamp int64  `qsign:"timeStamp"`
				NonceStr  string `qsign:"-"`
			}{
				AppID:     "wx6cfc34d48f33effe",
				TimeStamp: 1503117550,
				NonceStr:  "9446",
			},
			suffix: func() string {
				return "&key=123456"
			},
			expect: "appId=wx6cfc34d48f33effe&timeStamp=1503117550&key=123456",
		},
		{
			input: struct {
				Action string
				Nonce  int
			}{
				Action: "DescribeInstances",
				Nonce:  11886,
			},
			prefix: func() string {
				return "https://github.com/"
			},
			expect: "https://github.com/Action=DescribeInstances&Nonce=11886",
		},
	}

	for _, c := range cases {
		q := NewQsign(Options{
			PrefixGenerator: c.prefix,
			SuffixGenerator: c.suffix,
		})
		d, _ := q.Digest(c.input)
		actual := string(d)
		if actual != c.expect {
			t.Errorf("expect digest is %s, actual is %s", c.expect, actual)
		}
	}
}
