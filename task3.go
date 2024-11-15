package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func writeXMLValue(key string, value interface{}, indent int) string {
	indentSpaces := strings.Repeat(" ", indent*2)
	valueStr := fmt.Sprintf("%v", value)

	return fmt.Sprintf("%s<%s>%s</%s>\n", indentSpaces, key, valueStr, key)
}

func processYAMLMap(data map[interface{}]interface{}, indent int) string {
	var xmlOutput string
	indentSpaces := strings.Repeat(" ", indent*2)
	var tagStack []string

	for key, value := range data {
		strKey := fmt.Sprintf("%v", key)

		switch v := value.(type) {
		case map[interface{}]interface{}:
			// Добавляем открывающий тег
			xmlOutput += fmt.Sprintf("%s<%s>\n", indentSpaces, strKey)
			tagStack = append(tagStack, strKey) // Добавляем тег в стек
			xmlOutput += processYAMLMap(v, indent+1)
			// Добавляем закрывающий тег
			xmlOutput += fmt.Sprintf("%s</%s>\n", indentSpaces, strKey)
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[interface{}]interface{}); ok {
					xmlOutput += fmt.Sprintf("%s<%s>\n", indentSpaces, strKey)
					tagStack = append(tagStack, strKey) // Добавляем тег в стек
					xmlOutput += processYAMLMap(itemMap, indent+1)
					xmlOutput += fmt.Sprintf("%s</%s>\n", indentSpaces, strKey)
				} else {
					xmlOutput += writeXMLValue("item", item, indent)
				}
			}
		default:
			xmlOutput += writeXMLValue(strKey, v, indent)
		}
	}
	return xmlOutput
}

func task3() {
	inputFile := "data/saturday.yaml" // Укажите путь к вашему YAML файлу
	outputFile := "output_3.xml"      // Укажите путь к выходному XML файлу

	yamlData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Ошибка чтения YAML файла: %v\n", err)
		return
	}

	var data map[interface{}]interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		fmt.Printf("Ошибка разбора YAML: %v\n", err)
		return
	}

	xmlOutput := processYAMLMap(data, 0)
	if err := ioutil.WriteFile(outputFile, []byte(xmlOutput), 0644); err != nil {
		fmt.Printf("Ошибка записи в XML файл: %v\n", err)
		return
	}

	fmt.Printf("Transformation successful! XML output written to %s\n", outputFile)
}
