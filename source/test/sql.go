package main

import (
	"fmt"

	"../proxy/encode"
)

/*
func main() {
	mysqlDb := sqlinterface.DB{
		DbType:       "mysql",
		DatabaseName: "mydb",
		Table:        "students",
		Username:     "root",
		Password:     "123",
	}
	ra := sqlinterface.RowAccess{
		Column:  "ID",
		Indices: []int{1, 2, 3, 4, 5},
	}

	res := mysqlDb.GetRows(ra)
	row0 := res[0]
	fmt.Print(row0)
	fmt.Print(row0.NAME)
	//fmt.Print(indexToLetters(1))

}

func indexToLetters(index int) string {

	base := float64(26)

	// start at the base that is bigger and work your way down
	floatIndex := float64(index)
	leftOver := floatIndex

	columns := []int{}

	for leftOver > 0 {
		remainder := math.Mod(leftOver, base)

		if remainder == 0 {
			remainder = base
		}

		columns = append([]int{int(remainder)}, columns...)
		leftOver = (leftOver - remainder) / base
	}

	var buff bytes.Buffer

	for _, e := range columns {
		buff.WriteRune(toChar(e))
	}

	return buff.String()

}
func toChar(i int) rune {
	return rune('A' - 1 + i)
}

*/

func main() {

	ReadStree := encode.ReadFromFile("stree.txt")
	fmt.Printf("\nReadStree:\n")
	ReadStree.PrintTree()
	fmt.Printf("Order=%d\n", ReadStree.Order)

	a1 := "A1"
	leng := len(a1)
	numInStr := a1[1:leng]
	fmt.Printf("a1 = %s, leng = %d, numInStr = %s\n", a1, leng, numInStr)
}
