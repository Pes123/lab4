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

	//в теории чтобы потом модно
	tagMap := map[string]int{
		"name:":     0,
		"type:":     1,
		"teacher:":  2,
		"audience:": 3,
		"building:": 4,
		"start:":    5,
		"end:":      6,
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "day:") {
			day = strings.TrimSpace(line[strings.Index(line, ":")+1:])
		} else {
			for prefix, index := range tagMap {
				if strings.HasPrefix(line, prefix) {
					if prefix == "name:" {
						tagList = make([]string, len(tagMap))
					}
					tagList[index] = strings.TrimSpace(line[strings.Index(line, ":")+1:])
					if prefix == "end:" {
						csvContent.WriteString(writeToCSV(day, tagList))
					}
					break
				}
			}
		}
	}

	if err := ioutil.WriteFile("schedule.csv", []byte(csvContent.String()), 0644); err != nil {
		fmt.Println("Error writing to CSV:", err)
	} else {
		fmt.Println("CSV file has been created successfully.")
	}
}
