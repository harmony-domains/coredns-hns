package onens

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestONENSParse(t *testing.T) {
	tests := []struct {
		key                string
		inputFileRules     string
		err                string
		connection         string
		ethlinknameservers []string
	}{
		{ // 0
			".",
			`onens {
			}`,
			"Testfile:2 - Error during parsing: no connection",
			"",
			nil,
		},
		{ // 1
			".",
			`onens {
			   connection
			}`,
			"Testfile:2 - Error during parsing: invalid connection; no value",
			"",
			nil,
		},
		{ // 2
			".eth.link",
			`onens {
			  connection /home/test/.ethereum/geth.ipc
			  ethlinknameservers ns1.ethdns.xyz
			}`,
			"",
			"/home/test/.ethereum/geth.ipc",
			[]string{"ns1.ethdns.xyz."},
		},
		{ // 3
			".",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"http://localhost:8545/",
			[]string{"ns1.ethdns.xyz.", "ns2.ethdns.xyz."},
		},
		{ // 4
			".",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"",
			nil,
		},
		{ // 5
			".",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"",
			nil,
		},
		{ // 6
			"tls://.:8053",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"",
			nil,
		},
		{ // 7
			".:8053",
			`onens {
			  connection http://localhost:8545/ bad
			}`,
			"Testfile:2 - Error during parsing: invalid connection; multiple values",
			"",
			nil,
		},
		{ // 8
			".:8053",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"",
			nil,
		},
		{ // 9
			".:8053",
			`onens {
			  connection http://localhost:8545/
			  ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz
			}`,
			"",
			"",
			nil,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("onens", test.inputFileRules)
		c.Key = test.key
		connection, ethlinknameservers, err := onensParse(c)

		if test.err != "" {
			if err == nil {
				t.Fatalf("Failed to obtain expected error at test %d", i)
			}
			if err.Error() != test.err {
				t.Fatalf("Unexpected error \"%s\" at test %d", err.Error(), i)
			}
		} else {
			if err != nil {
				t.Fatalf("Unexpected error \"%s\" at test %d", err.Error(), i)
			} else {
				if test.connection != "" && connection != test.connection {
					t.Fatalf("Test %d connection expected %v, got %v", i, test.connection, connection)
				}
				if test.ethlinknameservers != nil {
					if len(ethlinknameservers) != len(test.ethlinknameservers) {
						t.Fatalf("Test %d ethlinknameservers expected %v entries, got %v", i, len(test.ethlinknameservers), len(ethlinknameservers))
					}
					for j := range test.ethlinknameservers {
						if ethlinknameservers[j] != test.ethlinknameservers[j] {
							t.Fatalf("Test %d ethlinknameservers expected %v, got %v", i, test.ethlinknameservers[j], ethlinknameservers[j])
						}
					}
				}
			}
		}
	}
}
