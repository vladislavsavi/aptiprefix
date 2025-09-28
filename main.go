package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 1. Определение флага командной строки
	var prefix string
	// Устанавливаем пустое значение по умолчанию ("")
	flag.StringVar(&prefix, "prefix", "", "Префикс, который нужно удалить из имен файлов в текущей директории. ОБЯЗАТЕЛЬНЫЙ.")

	// Переопределяем функцию вывода справки, чтобы она была более информативной при ошибке
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: %s --prefix <значение>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Удаляет заданный префикс из имен файлов в текущей директории.")
		fmt.Fprintln(os.Stderr, "Пример: antiprefix --prefix ef3e_beg_")
		fmt.Fprintln(os.Stderr, "\nФлаги:")
		flag.PrintDefaults()
	}

	// Разбор аргументов
	flag.Parse()

	// 2. Обязательная проверка наличия префикса
	if prefix == "" {
		fmt.Fprintln(os.Stderr, "\nОшибка: Префикс не может быть пустым. Используйте флаг --prefix.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Поиск и переименование файлов с префиксом '%s'...\n", prefix)

	// 3. Чтение текущей директории
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при чтении директории: %v\n", err)
		os.Exit(1)
	}

	// 4. Итерация, фильтрация и переименование
	var renamedCount int
	for _, entry := range entries {
		oldName := entry.Name()

		// Проверяем, что это обычный файл
		if !entry.Type().IsRegular() {
			continue
		}

		// Проверяем, начинается ли имя файла с заданного префикса
		if strings.HasPrefix(oldName, prefix) {
			// Удаляем префикс
			newName := strings.TrimPrefix(oldName, prefix)

			// Проверка: новое имя не должно быть пустым
			if newName == "" {
				fmt.Printf("Пропуск %s: новое имя было бы пустым (невозможно).\n", oldName)
				continue
			}

			// Переименование
			err := os.Rename(oldName, newName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при переименовании файла %s в %s: %v\n", oldName, newName, err)
			} else {
				fmt.Printf("Переименовано: %s -> %s\n", oldName, newName)
				renamedCount++
			}
		}
	}

	fmt.Printf("Завершено. Всего переименовано файлов: %d\n", renamedCount)
}
