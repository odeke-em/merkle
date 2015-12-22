package merkle

import (
	"bytes"
	"crypto/md5"
)

func HashSum(data []byte) []byte {
	fixedSizeSum := md5.Sum(data)
	return fixedSizeSum[:] // release the [Size]byte constraint => []byte
}

func concatenatePairs(data []byte, blockSize uint64) (result [][]byte) {
	i := uint64(0)
	n := uint64(len(data))

	for i < n {
		sects := [][]byte{}
		for j := 0; j < 2 && i < n; j++ {
			end := i + blockSize
			if end >= n {
				end = n
			}
			sects = append(sects, data[i:end])
			i = end
		}

		result = append(result, bytes.Join(sects, []byte{}))
	}

	return result
}

func TopLevelify(data []byte, blockSize uint64) (b []byte) {
	chunks := concatenatePairs(data, blockSize)

	for {
		n := len(chunks)

		topLevelSums := [][]byte{}
		for i := 0; i < n; {
			sects := [][]byte{}
			for j := 0; j < 2 && i < n; j++ {
				sects = append(sects, chunks[i])
				i++
			}

			joined := bytes.Join(sects, []byte{})
			topLevelSums = append(topLevelSums, HashSum(joined))
		}

		chunks = topLevelSums

		if n <= 1 {
			break
		}
	}

	if len(chunks) >= 1 {
		b = chunks[0]
	}

	return b
}
