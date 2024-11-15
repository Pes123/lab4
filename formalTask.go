package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func formal() {
	inputFile := "data/saturday.yaml"
	outputFile := "output_formal.xml"
	yamlFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Ошибка чтения YAML файла: %s\n", err)
	}

	xmlOutput, err := TransformWithFormal(string(yamlFile))
	if err != nil {
		log.Fatalf("Ошибка парсинга YAML файла: %s\n", err)
	}

	err = ioutil.WriteFile(outputFile, []byte(xmlOutput), 0644)
	if err != nil {
		log.Fatalf("Ошибка записи XML файла: %s\n", err)
	}

	fmt.Printf("Transformation successful! XML output written to %s\n", outputFile)
}

func TransformWithFormal(data string) (string, error) {
	lines := strings.Split(data, "\n")
	var xmlResult string
	xmlResult += "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<config>\n" // Начало XML
	var tagStack []string                                                 // Стек для хранения текущих тегов
	var previousIndent int                                                // Хранит предыдущий уровень отступа
	var listTagName string                                                // Временная переменная для имени тега списка

	// Регулярные выражения для обработки различных форматов строк
	re := regexp.MustCompile(`^([\w:-]+):\s*(.*)`)
	listRe := regexp.MustCompile(`^\s*-\s+(.*)`)
	commentRe := regexp.MustCompile(`^\s*#(.*)`)
	inlineListRe := regexp.MustCompile(`^([\w:-]+):\s*\[(.*)\]`)

	// Проход по каждой строке
	for i, line := range lines {
		thisIndent := len(line) - len(strings.TrimSpace(line)) // Вычисление текущего уровня отступа
		// Так как отступы в два пробела, это сделано чтобы типо инкремент работал норм xd
		thisIndent /= 2
		line = strings.TrimSpace(line)

		// Игнорирование маркеров и пустых строк
		if line == "---" || line == "..." || line == "" {
			continue
		}

		// Закрытие тегов при уменьшении уровня отступа
		for previousIndent > thisIndent {
			xmlResult += fmt.Sprintf("%s</%s>\n", strings.Repeat("  ", previousIndent-1), tagStack[len(tagStack)-1])
			tagStack = tagStack[:len(tagStack)-1] // Удаление последнего тега из стека
			previousIndent--
		}

		// Обработка комментариев
		if handleComment(line, commentRe, thisIndent, &xmlResult) {
			continue
		}

		// Обработка строчных списков
		if handleInlineList(line, inlineListRe, thisIndent, &xmlResult) {
			continue
		}

		// Обработка тегов с значениями
		if handleTagValue(line, re, &xmlResult, &tagStack, &previousIndent, &listTagName, lines, i) {
			continue
		}

		// Обработка списков
		handleList(line, listRe, thisIndent, &xmlResult, listTagName)
	}

	// Закрытие оставшихся открытых тегов
	for len(tagStack) > 0 {
		previousIndent--
		xmlResult += fmt.Sprintf("%s</%s>\n", strings.Repeat("  ", previousIndent), tagStack[len(tagStack)-1])
		tagStack = tagStack[:len(tagStack)-1]
	}

	xmlResult += "</config>\n"
	return xmlResult, nil
}

// escapeXML заменяет специальные символы на соответствующие XML сущности
func escapeXML(value string) string {
	replacements := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"'":  "&apos;",
		"\"": "&quot;",
		"$":  "&#36;",
	}
	for old, new := range replacements {
		value = strings.ReplaceAll(value, old, new) // Замена в строке
	}
	return value
}

// handleComment обрабатывает строки комментариев и добавляет их в XML
func handleComment(line string, commentRe *regexp.Regexp, indent int, xmlResult *string) bool {
	commentMatches := commentRe.FindStringSubmatch(line)
	if commentMatches != nil {
		lastComment := strings.TrimSpace(commentMatches[1])
		if lastComment != "" {
			*xmlResult += fmt.Sprintf("%s<!-- %s -->\n", strings.Repeat("  ", indent), escapeXML(lastComment)) // Добавление комментария в формате XML
		}
		return true
	}
	return false
}

// handleInlineList обрабатывает "строчные" списки и добавляет их в XML
func handleInlineList(line string, inlineListRe *regexp.Regexp, indent int, xmlResult *string) bool {
	inlineListMatches := inlineListRe.FindStringSubmatch(line)
	if inlineListMatches != nil {
		tag := inlineListMatches[1]
		values := strings.Split(strings.TrimSpace(inlineListMatches[2]), ",") // Разделение элементов списка по запятой
		for _, value := range values {
			*xmlResult += fmt.Sprintf("%s<%s>%s</%s>\n", strings.Repeat("  ", indent), tag, escapeXML(strings.TrimSpace(value)), tag) // Добавление каждого элемента в XML
		}
		return true
	}
	return false
}

// handleTagValue обрабатывает строки с тегом и значением
func handleTagValue(line string, re *regexp.Regexp, xmlResult *string, tagStack *[]string, previousIndent *int, listTagName *string, lines []string, index int) bool {
	matches := re.FindStringSubmatch(line)
	if matches != nil {
		tag := matches[1]
		value := strings.TrimSpace(matches[2])

		isNextLineListItem := false
		if index+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[index+1])
			isNextLineListItem = regexp.MustCompile(`^\s*-\s+(.*)`).MatchString(nextLine)
			if isNextLineListItem {
				*listTagName = tag // Запоминаем тег списка, чтобы потом все теги в "списке" были с его именем
			}
		}

		if value == "" {
			if !isNextLineListItem {
				*xmlResult += fmt.Sprintf("%s<%s>\n", strings.Repeat("  ", *previousIndent), tag)
				*tagStack = append(*tagStack, tag)
				*previousIndent++
			}
		} else {
			*xmlResult += fmt.Sprintf("%s<%s>%s</%s>\n", strings.Repeat("  ", *previousIndent), tag, escapeXML(value), tag)
		}
		return true
	}
	return false
}

// handleList обрабатывает строки со списками и добавляет их в XML
func handleList(line string, listRe *regexp.Regexp, indent int, xmlResult *string, listTagName string) {
	listMatches := listRe.FindStringSubmatch(line)
	if listMatches != nil {
		listValue := strings.TrimSpace(listMatches[1])
		if listTagName != "" {
			*xmlResult += fmt.Sprintf("%s<%s>%s</%s>\n", strings.Repeat("  ", indent-1), listTagName, escapeXML(listValue), listTagName)
		}
	}
}
