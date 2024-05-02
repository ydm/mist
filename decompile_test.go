package mist_test

import (
	"fmt"
	"testing"

	"github.com/ydm/mist"
)

func TestDecompile(t *testing.T) {
	t.Parallel()
	// fmt.Println(mist.Decompile("0x60178061000c6000396000f340526000506020604051602060108181019150508152f3"))
	fmt.Println(mist.Decompile("0x3373fffffffffffffffffffffffffffffffffffffffe14604d57602036146024575f5ffd5b5f35801560495762001fff810690815414603c575f5ffd5b62001fff01545f5260205ff35b5f5ffd5b62001fff42064281555f359062001fff015500"))
}
