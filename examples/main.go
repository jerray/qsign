package main

import (
	"fmt"

	"github.com/jerray/qsign"
)

func main() {
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

		// The default filter
		// Filter: func(key, value string) {
		// 	return len(value) > 0
		// },

		// To use a hash.Hash other than md5
		// Hasher: func() hash.Hash {
		// 	return sha256.New()
		// },

		// To use a encoding other than hex
		// Encoder: func() qsign.Encoding {
		// 	return base64.StdEncoding
		// },
	})

	signature, _ := q.Sign(data)
	fmt.Printf("%s", string(signature))
	// 9a0a8659f005d6984697e2ca0a9cf3b7
}
