package prometheus

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"fmt"
)

func Listtargets(){

}

func Checktargets(){

}

func Updatetargets(){

}

func Deletetargets(){

}

func Getfileinfo(filename string, filter string) (mem_usage_bytes int64) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file failed,err:", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Text()
		//filter ="VmRSS:(.*)kB"
		mem_usage_string, _ := regexp.Compile(filter)
		match := mem_usage_string.FindString(value)
		if match != "" {
			digital, _ := regexp.Compile("[0-9]+")
			number_str := digital.FindString(match)
			number, _ := strconv.Atoi(number_str)
			fmt.Println(number)
			break
		}
	}
	fmt.Println( " mem_usage_bytes ", mem_usage_bytes)
	return mem_usage_bytes
}
