package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	tests := []struct {
		description string

		// Test input
		route  string
		method string
		body   string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "404 route",
			route:         "/i-dont-exist",
			method:        "GET",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "404 page not found\n",
		},
		{
			description:   "C Chain RPC route",
			route:         "/ext/bc/C/rpc",
			method:        "GET",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "",
		},
		{
			description:   "Get some block 0x1000",
			route:         "/ext/bc/C/rpc",
			method:        "POST",
			body:          `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1000",false],"id":1}`,
			expectedError: false,
			expectedCode:  200,
			expectedBody: `{"jsonrpc":"2.0","id":1,"result":{"blockExtraData":"0x","difficulty":"0x1","extDataHash":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","extraData":"0x","gasLimit":"0x7a1200","gasUsed":"0x520c","hash":"0x71bf35e949fd8a4ab11ad141def5caa704d4aa9072adb236bc4e2c53cdb28dfc","logsBloom":"0x0012000000001000000000000000000000000000000000000000000000020000000000000000000000000000000000080000000080000000000000040000100000040000000000000040000000200001000000400000000040008000000000000000000000000000000000080000000000000000000000000100000000000000000008400000000000000000000000000000040010400000000080040000000008008000200000020000000000080000000000018000000100000000001000000000000000000000000000000000010001010000000000900000000200000000000a000000000000000000000000000000000040000000280000000000000000","miner":"0x0100000000000000000000000000000000000000","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","number":"0x1000","parentHash":"0x73a42f8c66433ba964f309ef67ef5395fad641f6a7f4c0f8355e62cedf86d594","receiptsRoot":"0x970f0a58616549cbad3d7f1260c05fae5da0157435e775529d1353b94f50604c","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x28c","stateRoot":"0x39cef27c57d5da7357c76a9397a44242b461948b8136afaf68a3f60a69cc10c6","timestamp":"0x6146bb69","totalDifficulty":"0x1000","transactions":["0x3c0d5d3d878ae4dd19f73f5d8bf59938e1d7a57221548ca41a0194aafdc527dc"],"transactionsRoot":"0xe7e18b0f402df710a3d5b24037bb6093073900a7259df862f56d245cfbb0a28c","uncles":[]}}
`,
		},
		{
			description: "Valid eth_sendRawTransaction",
			route:       "/ext/bc/C/rpc",
			method:      "POST",
			body: `{"id":2,"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf8680685051f4d5c00833d090094ac2d0226ade52e6b4ccc97359359f01e34d50352834c4b40802aa004f4410acfc636f5d65e75ca37448fd0b1e90632661ad3ea071594b646e7621fa0186d78ddbeef43bc372fcc650ad4e77e0629206426cd183d751e9ddcc8d5e777"]}
`,
			expectedError: false,
			expectedCode:  200,
			expectedBody: `{"jsonrpc":"2.0","id":2,"error":{"code":-32000,"message":"invalid sender"}}
`,
		},
		{
			description: "Blocked eth_sendRawTransaction",
			route:       "/ext/bc/C/rpc",
			method:      "POST",
			body: `{"id":2,"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf9040380808094100000000000000000000000000000000000000380b903a4c5adc539000000000000000000000000000000000000000000000000000000000003907e00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000700000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000000ce95a91f50cd86421f18822d2996de377f6e56efbab59a4d164155ada292605bd21c008f7e7726c98b3eed9ac7e880417405eade31ddf8e9bd5286609f50f0d1c1ed1c05db54d096b150eef65d4e31aa963529b7156b4ba60607db8729c5202f829da4d7da26c6423bd67a2b4a1bc77c08f8a9e9bd7600ee05379652051408ca54ac8ca92b4a214025f75277dc0e3a1846fc3929764752ccdc44dfb2b5cfe9e8b088b93d65fe866ec6f6738a37c678b7da9a681696c5148a005524f97e9b3492fb5754e5293502a58942f639aa5778ec98e7d097419146ee179be36d77ff5fba85a2dbd1ac6dbf8a77f1a6d1fd7f094a768fe6b28897f7569a6bde8ce9521ba57f20332695cd9b0a4bde0a379933fbf2c98403034a27714187d02bf1d4579ea07c5e1b56b6b855d82203d37ba5fab0e301159a66d721e40e62bf3843a9d5d750732671a2e0e0089199a5aabde98aeb3d4abb3344e673d22ee9c35cb2eba0a936f0bf94cab2b6908ba75cc8e59070453e0b4e182aec3616e97745a78c2fac172c11ca0f0e923a0f564d03a731e269154fb07bd9142b9115848a32364c519af1ee75035a04fbdd8f3d5f435aa091423ae0eaa4eea0b4041544cc6520cf42ec6d0f887028e"]}
`,
			expectedError: false,
			expectedCode:  400,
			expectedBody:  "Bad Request",
		},
		{
			description: "Valid eth_sendTransaction",
			route:       "/ext/bc/C/rpc",
			method:      "POST",
			body: `{"id":2,"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","gas":"0x76c0","gasPrice":"0x9184e72a000","value":"0x9184e72a","data":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}]}
`,
			expectedError: false,
			expectedCode:  200,
			expectedBody: `{"jsonrpc":"2.0","id":2,"error":{"code":-32000,"message":"unknown account"}}
`,
		},
		{
			description: "Blocked eth_sendTransaction",
			route:       "/ext/bc/C/rpc",
			method:      "POST",
			body: `{"id":2,"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","to":"0x1000000000000000000000000000000000000003","gas":"0x76c0","gasPrice":"0x9184e72a000","value":"0x9184e72a","data":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}]}
`,
			expectedError: false,
			expectedCode:  400,
			expectedBody:  "Bad Request",
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		var req *http.Request

		if test.method == "GET" {
			req, _ = http.NewRequest(
				test.method,
				test.route,
				nil,
			)
		} else {
			body := []byte(test.body)
			req, _ = http.NewRequest(
				test.method,
				test.route,
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		resBody, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(resBody), test.description)
	}
}
