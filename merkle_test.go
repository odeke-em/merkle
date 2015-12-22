package merkle

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func callerPath() string {
	_, p, _, _ := runtime.Caller(0)
	return p
}

func readFile(p string) string {
	f, err := os.Open(p)
	if err != nil {
		return ""
	}

	defer f.Close()
	data, _ := ioutil.ReadAll(f)
	return string(data)
}

func ownContent() string {
	return readFile(callerPath())
}

func TestChunkifyDataLessThanBlockSize(t *testing.T) {
	cases := []string{
		ownContent(),
		"abcdefghijklmnopqrstuvwxyz",
		"0",
		"aE??l?LiٱD[ʝ?VUk?Ԏw??5???4?h?'?0???W?Ңu?׋=?J??-ݓL?)5nKz??W??egV????U?????-",
		"11110000000      120202022020afeareallll183755^2qsssa",
	}

	for _, tc := range cases {
		data := []byte(tc)
		expected := HashSum(data)
		oneExtraLen := uint64(len(data) + 1)
		topLevel := TopLevelify(data, oneExtraLen)

		if !bytes.Equal(topLevel, expected) {
			t.Errorf("given %q, expected %x got %x", tc, expected, topLevel)
		}
	}
}

func TestChunkifyOnUnitLevel(t *testing.T) {
	cases := []struct {
		c1, c2 string
	}{
		{c1: "abcdefghijklmno", c2: "123456801*^&~~"},
		{c1: "82821B--+^&!&!i+wwlwla", c2: "5??r?????x%L??W?N????x?c????"},
	}

	for _, tc := range cases {
		concat2 := tc.c1 + tc.c2
		bconcat2 := []byte(concat2)
		expected := HashSum(bconcat2)

		ln := len(concat2)
		if ln%2 == 1 {
			ln += 1
		}

		halfLen := uint64(ln / 2)

		topLevel := TopLevelify(bconcat2, halfLen)

		if !bytes.Equal(topLevel, expected) {
			t.Errorf("given l: %q, r: %q: expected %x got %x", tc.c1, tc.c2, expected, topLevel)
		}
	}
}
