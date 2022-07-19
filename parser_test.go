package halvector

import (
	"bytes"
	"testing"
)

type stage struct {
	hal    string
	expect string
	err    error
}

// var stages = []stage{
// 	{
// 		hal:    "en-GB;q=0.8",
// 		expect: `[{"code":"en","region":"GB","quality":0.8}]`,
// 	},
// 	{
// 		hal:    "en-GB",
// 		expect: `[{"code":"en","region":"GB","quality":1.0}]`,
// 	},
// 	{
// 		hal:    "en;q=0.8",
// 		expect: `[{"code":"en","quality":0.8}]`,
// 	},
// 	{
// 		hal:    "az-AZ",
// 		expect: `[{"code":"az","region":"AZ","quality":1.0}]`,
// 	},
// 	{
// 		hal:    "fr-CA,fr;q=0.8",
// 		expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"fr","quality":0.8}]`,
// 	},
// 	{
// 		hal:    "fr-CA,*;q=0.8",
// 		expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"*","quality":0.8}]`,
// 	},
// 	{
// 		hal:    "fr-150",
// 		expect: `[{"code":"fr","region":"150","quality":1.0}]`,
// 	},
// 	{
// 		hal:    "fr-CA,fr;q=0.8,en-US;q=0.6,en;q=0.4,*;q=0.1",
// 		expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"fr","quality":0.8},{"code":"en","region":"US","quality":0.6},{"code":"en","quality":0.4},{"code":"*","quality":0.1}]`,
// 	},
// 	{
// 		hal:    "fr-CA, fr;q=0.8,  en-US;q=0.6,en;q=0.4,    *;q=0.1",
// 		expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"fr","quality":0.8},{"code":"en","region":"US","quality":0.6},{"code":"en","quality":0.4},{"code":"*","quality":0.1}]`,
// 	},
// 	{
// 		hal:    "fr-CA,fr;q=0.2,en-US;q=0.6,en;q=0.4,*;q=0.5",
// 		expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"en","region":"US","quality":0.6},{"code":"*","quality":0.5},{"code":"en","quality":0.4},{"code":"fr","quality":0.2}]`,
// 	},
// 	{
// 		hal:    "zh-Hant-cn",
// 		expect: `[{"code":"zh", script: 'Hant',"region":"cn","quality":1.0}]`,
// 	},
// 	{
// 		hal:    "zh-Hant-cn;q=1, zh-cn;q=0.6, zh;q=0.4",
// 		expect: `[{"code":"zh", script: 'Hant',"region":"cn","quality":1.0},{"code":"zh","region":"cn","quality":0.6},{"code":"zh","quality":0.4}]`,
// 	},
// 	{
// 		hal:    "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
// 		expect: `[{"code":"ru","region":"RU","quality":0.9},{"code":"en","region":"US","quality":0.8},{"code":"en","quality":0.7}]`,
// 	},
// }

func TestParser(t *testing.T) {
	var stages = []stage{
		{
			hal: "en",
			expect: `[
	{
		"code": "en",
		"quality": 1.0
	}
]`,
		},
		{
			hal: "en-GB",
			expect: `[
	{
		"code": "en",
		"region": "GB",
		"quality": 1.0
	}
]`,
		},
		{
			hal: "en-Latin-GB",
			expect: `[
	{
		"code": "en",
		"script": "Latin",
		"region": "GB",
		"quality": 1.0
	}
]`,
		},
		{
			hal: "en-GB;q=0.8",
			expect: `[
	{
		"code": "en",
		"region": "GB",
		"quality": 0.8
	}
]`,
		},
		{
			hal: "en;q=0.8",
			expect: `[
	{
		"code": "en",
		"quality": 0.8
	}
]`,
		},
		{
			hal: "az-AZ",
			expect: `[
	{
		"code": "az",
		"region": "AZ",
		"quality": 1.0
	}
]`,
		},
		{
			hal: "fr-CA,fr;q=0.8",
			expect: `[
	{
		"code": "fr",
		"region": "CA",
		"quality": 1.0
	},
	{
		"code": "fr",
		"quality": 0.8
	}
]`,
		},
		{
			hal: "fr-CA,*;q=0.8",
			expect: `[
	{
		"code": "fr",
		"region": "CA",
		"quality": 1.0
	},
	{
		"code": "*",
		"quality": 0.8
	}
]`,
		},
		// {
		// 	hal:    "fr-150",
		// 	expect: `[{"code":"fr","region":"150","quality":1.0}]`,
		// },
		// {
		// 	hal:    "fr-CA,fr;q=0.8,en-US;q=0.6,en;q=0.4,*;q=0.1",
		// 	expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"fr","quality":0.8},{"code":"en","region":"US","quality":0.6},{"code":"en","quality":0.4},{"code":"*","quality":0.1}]`,
		// },
		// {
		// 	hal:    "fr-CA, fr;q=0.8,  en-US;q=0.6,en;q=0.4,    *;q=0.1",
		// 	expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"fr","quality":0.8},{"code":"en","region":"US","quality":0.6},{"code":"en","quality":0.4},{"code":"*","quality":0.1}]`,
		// },
		// {
		// 	hal:    "fr-CA,fr;q=0.2,en-US;q=0.6,en;q=0.4,*;q=0.5",
		// 	expect: `[{"code":"fr","region":"CA","quality":1.0},{"code":"en","region":"US","quality":0.6},{"code":"*","quality":0.5},{"code":"en","quality":0.4},{"code":"fr","quality":0.2}]`,
		// },
		// {
		// 	hal:    "zh-Hant-cn",
		// 	expect: `[{"code":"zh", script: 'Hant',"region":"cn","quality":1.0}]`,
		// },
		// {
		// 	hal:    "zh-Hant-cn;q=1, zh-cn;q=0.6, zh;q=0.4",
		// 	expect: `[{"code":"zh", script: 'Hant',"region":"cn","quality":1.0},{"code":"zh","region":"cn","quality":0.6},{"code":"zh","quality":0.4}]`,
		// },
		// {
		// 	hal:    "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
		// 	expect: `[{"code":"ru","region":"RU","quality":0.9},{"code":"en","region":"US","quality":0.8},{"code":"en","quality":0.7}]`,
		// },
	}
	for _, stg := range stages {
		t.Run(stg.hal, func(t *testing.T) {
			var buf bytes.Buffer
			vec := Acquire()
			if err := vec.ParseStr(stg.hal); err != nil {
				t.Error(err)
				return
			}
			_ = vec.Beautify(&buf)
			if stg.expect != buf.String() {
				t.Errorf("expect: %s\ngot: %s", stg.expect, buf.String())
			}
		})
	}
}
