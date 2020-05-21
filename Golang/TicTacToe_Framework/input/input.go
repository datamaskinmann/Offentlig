package input

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadString() string {
	value, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Split(value, "\n")[0]
}

func ReadInt() (int, error) {
	value, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	intValue, conversionError := strconv.Atoi(strings.Split(value, "\n")[0])

	if conversionError != nil {
		return -1, conversionError
	}

	return intValue, nil
}

func ReadIntSlice() ([]int, error) {
	value, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	var args []string = strings.Split(value, " ")

	var retSlice []int

	for _, a := range args {
		intValue, conversionError := strconv.Atoi(strings.Split(a, "\n")[0])
		if conversionError != nil {
			return []int{}, conversionError
		}
		retSlice = append(retSlice, intValue)
	}

	return retSlice, nil
}
