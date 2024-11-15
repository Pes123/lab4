package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func transformRe(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var xml string
	xml += "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<config>\n"
	previousIndent := 0
	var tagStack []string

	re := regexp.MustCompile(`^([\w:-]+):\s*(.*)`)

	for _, line := range lines {
		// Вычисление текущего уровня отступа
		thisIndent := len(line) - len(strings.TrimSpace(line))
		thisIndent /= 2 //  Так как отступы в два пробела, это сделано чтобы типо инкремент работал норм xd
		line = strings.TrimSpace(line)
		// Игнорирую разделительные линии и пустые строки
		if line == "---" || line == "..." || line == "" {
			continue
		}

		// Закрываю тег в случае уменьшения уровня отступа
		for previousIndent > thisIndent {
			xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"</%s>\n", tagStack[len(tagStack)-1])
			tagStack = tagStack[:len(tagStack)-1]
			previousIndent--
		}
		// По заданию надо использовать регулярку, ну я и использую просто выбирая потом группы..
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			tag := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])
			// Если значение отсутствует просто открываю обычный тег и добавляю его в стек
			if value == "" {
				xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"<%s>\n", tag)
				tagStack = append(tagStack, tag)
				previousIndent++
			} else {
				xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"<%s>%s</%s>\n", tag, value, tag)
			}
		}
	}
	// Закрываю оставшиеся открытые теги
	for len(tagStack) > 0 {
		// В каком то моменте появляляс лишний отступ поэтому -1 (костыль лол)
		xml += fmt.Sprintf(strings.Repeat(" ", (len(tagStack)-1)*2)+"</%s>\n", tagStack[len(tagStack)-1])
		tagStack = tagStack[:len(tagStack)-1]
	}
	return xml, nil
}
func reTask() {

	inputFile := "data/monday.yaml"
	outputFile := "output_re.xml"

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	xml, err := transformRe(string(data))
	if err != nil {
		fmt.Printf("Error transforming data: %v\n", err)
		return
	}

	err = ioutil.WriteFile(outputFile, []byte(xml), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Transformation successful! XML output written to %s\n", outputFile)
}
