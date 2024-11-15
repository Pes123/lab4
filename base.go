package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func transform_base(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var xml string
	xml += "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<config>\n"
	previousIndent := 0   // Хранит уровень отступа предыдущей строки
	var tagStack []string // Стек для хранения открытых тегов XML

	for _, line := range lines {
		// Вычисление текущего уровня отступа
		thisIndent := len(line) - len(strings.TrimSpace(line))
		thisIndent /= 2 // Так как отступы в два пробела, это сделано чтобы типо инкремент работал норм xd
		line = strings.TrimSpace(line)

		// Игнорируем разделительные линии и пустые строки
		if line == "---" || line == "..." || line == "" {
			continue
		}

		// Закрываем тег в случае уменьшения уровня отступа относительно предыдущей строки
		for previousIndent > thisIndent {
			xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"</%s>\n", tagStack[len(tagStack)-1])
			tagStack = tagStack[:len(tagStack)-1]
			previousIndent--
		}

		// Разделяем строку на тег и значение
		parts := strings.SplitN(line, ":", 2)
		// Если значение отсутствует, открываем тег
		if strings.TrimSpace(parts[1]) == "" { 
			xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"<%s>\n", parts[0])
			tagStack = append(tagStack, parts[0]) // Добавляем тег в стек
			previousIndent++
			continue
		}
		// Если тег
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			tag := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			xml += fmt.Sprintf(strings.Repeat(" ", thisIndent*2)+"<%s>%s</%s>\n", tag, value, tag)
		}
	}

	// Закрываем оставшиеся открытые теги
	for len(tagStack) > 0 {
		xml += fmt.Sprintf(strings.Repeat(" ", (len(tagStack)-1)*2)+"</%s>\n", tagStack[len(tagStack)-1])
		tagStack = tagStack[:len(tagStack)-1]
	}
	return xml, nil
}

func base() {
	inputFile := "data/monday.yaml"
	outputFile := "output_base.xml"

	// Читаем данные из входного файла
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Преобразуем данные в формат XML
	xml, err := transform_base(string(data))
	if err != nil {
		fmt.Printf("Error transforming data: %v\n", err)
		return
	}

	// Записываем получившийся XML в выходной файл
	err = ioutil.WriteFile(outputFile, []byte(xml), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Transformation successful! XML output written to %s\n", outputFile) // Успешное завершение
}
