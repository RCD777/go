package main

import "flag"
import "fmt"
import "bufio"
import "io"
import "os"

var infile1 *string = flag.String("i1", "infile1", "input file1")
var infile2 *string = flag.String("i2", "infile1", "input file2")
var outfile *string = flag.String("o", "outfile", "output file")

func readValues(infile string) (values []string, err error) {
	file, err := os.Open(infile)
	if err != nil {
		fmt.Println("fail to open file ", infile)
		return
	}

	defer file.Close()

	br := bufio.NewReader(file)

	values = make([]string, 0)

	for {
		line, isPrefix, err1 := br.ReadLine()

		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		if isPrefix {
			fmt.Println("A too long line")
			return
		}

		str := string(line)

		values = append(values, str)
	}

	return
}

func writeValues(values []string, outfile string) error {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("fail to create the output file ", outfile)
		return err
	}

	defer file.Close()

	for _, value := range values {
		file.WriteString(value + "\n")
	}

	fmt.Println("write to the output file ", outfile)

	return nil
}

func main() {
	flag.Parse()

	if infile1 != nil {

		fmt.Println("infile1", *infile1)
	}

	values1, _ := readValues(*infile1)
	values2, _ := readValues(*infile2)

	outputValues := make([]string, 0)

	//core algorithm
	for _, value1 := range values1 {

		tobeSlim := false
		for _, value2 := range values2 {
			if value1 == value2 {
				tobeSlim = true
			}
		}

		if !tobeSlim {
			outputValues = append(outputValues, value1)
		}
	}

	for _, outputValue := range outputValues {
		fmt.Println(outputValue)
	}

	writeValues(outputValues, *outfile)

}
