## 运行
```
https://play.golang.com/
```

## 提取
```
package main

import (
	"fmt"
)

func main() {
	T := "Bs6oSyqZYt8i0xDCWaLd2RV9JheGO7TvfXu1rQzMw3gHb5FKEckNInAUjp4mPl"
	TP := []int{50, 0, 18, 9, 4, 3, 0, 36, 19, 46, 32, 44, 42, 32, 23, 55, 35, 2, 39, 53, 6, 32, 46, 7, 49, 31, 56, 11, 14, 50, 2, 47, 54, 27, 32, 23, 43, 11}
	tmpR := []rune{}
	for i := 0; i < len(TP); i++ {
		tmpR = append(tmpR, rune(T[TP[i]]))
	}
	fmt.Println(string(tmpR))
}
```
## 解码
```
package main

import (
	"bytes"
	"math/big"
)

// Base58字符集
var base58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Decoding函数解码Base58编码的字符串
func Decoding(str string) string {
	strByte := []byte(str)
	ret := big.NewInt(0)
	for _, byteElem := range strByte {
		index := bytes.IndexByte(base58, byteElem)
		ret.Mul(ret, big.NewInt(58))
		ret.Add(ret, big.NewInt(int64(index)))
	}

	return string(ret.Bytes())
}

func main() {
	// 待解码的Base58编码字符串
	encodedStr := "kBLtSoBrdFfbgf9U16MnqfFZcvjiDk6KAGf9Hi"

	// 调用Decoding函数进行解码
	decodedStr := Decoding(encodedStr)

	// 输出解码结果
	println("解码后的字符串:", decodedStr)
}
```

## 结果
```
解码后的字符串: Harry-zklcdc/go-proxy-bingai
```
