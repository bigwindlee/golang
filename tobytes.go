/*  ==== Knowledge Points ====
1. io.Reader interface
2. built-in copy()
3. slice
4. idiom of for...range loop, related to interface slice
*/

package main

import (
	"fmt"
	"io"
)

type StringPair struct{ first, second string }

func (pair *StringPair) Exchange() {
	pair.first, pair.second = pair.second, pair.first
}

func (pair StringPair) String() string {
	return fmt.Sprintf("%q+%q", pair.first, pair.second)
}

func (pair *StringPair) Read(data []byte) (n int, err error) {
	if pair.first == "" && pair.second == "" {
		return 0, io.EOF
	}
	if pair.first != "" {
		n = copy(data, pair.first)
		pair.first = pair.first[n:]
	}
	if n < len(data) && pair.second != "" {
		m := copy(data[n:], pair.second)
		pair.second = pair.second[m:]
		n += m
	}
	return n, nil
}

func ToBytes(reader io.Reader, size int) ([]byte, error) {
	data := make([]byte, size)
	n, err := reader.Read(data)
	if err != nil {
		return data, err
	}
	// the data slice is resliced to reduce its length to the number of bytes
	// actually read
	return data[:n], nil
}

func main() {
	const size = 16
	robert := &StringPair{"Robert L.", "Stevenson"}
	david := StringPair{"David", "Balfour"}

	// idiom of for...range loop, related to interface slice
	for _, pair := range []fmt.Stringer{robert, &david} {
		fmt.Println(pair)
	}
	fmt.Println()

	for _, reader := range []io.Reader{robert, &david} {
		raw, err := ToBytes(reader, size)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%q\n", raw)
	}

	fmt.Println()
	for _, pair := range []fmt.Stringer{robert, &david} {
		fmt.Println(pair)
	}
}

/*
[Output]
"Robert L."+"Stevenson"
"David"+"Balfour"

"Robert L.Stevens"
"DavidBalfour"

""+"on"
""+""

*/
