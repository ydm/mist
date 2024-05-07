package mist

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func EncodeWithSignature(signature string, args ...any) string {
	opening := strings.Index(signature, "(")
	closing := strings.Index(signature, ")")

	name := signature[:opening]

	inputTypes := strings.Split(
		signature[opening+1:closing],
		",",
	)
	inputs := make([]string, len(inputTypes))
	for i, t := range inputTypes {
		inputs[i] = fmt.Sprintf(`{"type": "%s"}`, t)
	}

	def := fmt.Sprintf(`[{
        "type": "function",
        "name": "%s",
        "inputs": [%s],
        "outputs": [%s],
        "stateMutability": "payable"
    }]`, name, strings.Join(inputs, ","), "")

	encoder, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		panic(err)
	}

	encoded, err := encoder.Pack(name, args...)
	if err != nil {
		panic(err)
	}

	var hex strings.Builder
	for _, e := range encoded {
		fmt.Fprintf(&hex, "%02x", e)
	}

	ans := hex.String()
	if (len(ans)-8)%64 != 0 {
		panic("TODO")
	}

	return ans
}
