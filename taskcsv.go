package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func writeToCSV(day string, tagList []string) string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
		day, tagList[0], tagList[1], tagList[2],
		tagList[3], tagList[4], tagList[5], tagList[6])
}

func taskcsv() {
	content, err := ioutil.ReadFile("/media/egr/Data/code/infalab4/data/monday.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Называю столбцы
	var csvContent strings.Builder
	csvContent.WriteString("day,name,type,teacher,audience,building,start,end\n")

	lines := strings.Split(string(content), "\n")
	day := ""
	var tagList []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "day:") {
			day = strings.TrimSpace(line[strings.Index(line, ":")+1:])
		} else if strings.HasPrefix(line, "name:") {
			// Если tagList уже чет в себе имеет, то была пройдена строка, значит можно записываться в csv
			name := strings.TrimSpace(line[strings.Index(line, ":")+1:])
			tagList = []string{name} // обнуляем taglist для следующей строки, а затем обрабатываем каждый элемент yaml в лоб
		} else if strings.HasPrefix(line, "type:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
		} else if strings.HasPrefix(line, "teacher:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
		} else if strings.HasPrefix(line, "audience:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
		} else if strings.HasPrefix(line, "building:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
		} else if strings.HasPrefix(line, "start:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
		} else if strings.HasPrefix(line, "end:") {
			tagList = append(tagList, strings.TrimSpace(line[strings.Index(line, ":")+1:]))
			csvContent.WriteString(writeToCSV(day, tagList))
		}
	}

	if err := ioutil.WriteFile("schedule.csv", []byte(csvContent.String()), 0644); err != nil {
		fmt.Println("Error writing to CSV:", err)
	} else {
		fmt.Println("CSV file has been created successfully.")
	}
}
