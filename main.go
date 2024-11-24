package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Log_file_name string `json:"log_file_name"`
	URL           string
	Nums          []float64 `json:"nums"`
}

func main() {
	var config Config

	lf := createLogFile("log.txt")
	defer lf.Close()
	stdin, err := parseArgs(os.Args)
	if err != nil {
		log.Println("Для считывания массива из файла используйте --file, а из stdin --stdin", err)
		os.Exit(1)
	}

	if err = decoding(&config); err != nil {
		log.Println("Ошибка при считывании из json файла:", err)
		os.Exit(1)
	}

	if config.Log_file_name != "" {
		os.Rename("log.txt", config.Log_file_name)
	}

	if stdin {
		if config.Nums, err = input(); err != nil {
			log.Println("Не удалось считать массив чисел из стандартного ввода", err)
		} else {
			log.Println("Массив успешно считан из stdin")
			sum(config.Nums)
		}
	} else {
		if err = checkConfig(&config); err != nil {
			log.Println("Ошибка массива", err)
		} else {
			log.Println("Массив успешно считан из json файла")
			log.Println("Массив:", config.Nums)
			sum(config.Nums)
		}
	}

	if err := responceStatus(config.URL); err != nil {
		log.Println("Ошибка при выполнении Get запроса:", err)
	}
}

// Функция ниже считывает аргументы командной строки
func parseArgs(args []string) (bool, error) {
	if len(args) != 2 {
		err := fmt.Errorf("ошибка: недопустимое значение конфигурации: %s", args[1:])
		return false, err
	}
	var flag bool
	var err error
	switch args[1] {
	case "--file":
		flag = false
	case "--stdin":
		flag = true
	default:
		flag = false
		err = fmt.Errorf("ошибка: недопустимое значение конфигурации: %s", args[1:])
	}
	return flag, err
}

// Функция ниже проверяет json файл на корректный ввод массива
func checkConfig(config *Config) error {
	var err error
	if len(config.Nums) == 0 {
		err = errors.New("ошибка: в файле congig.json задан пустой массив")
	}
	return err
}

// Функция ниже считывает и обрабатывает массив из stdin
func input() ([]float64, error) {
	var nums []float64
	fmt.Println("Введите числа массива с клавиатуры. Все числа должны разделяться пробелом.\nПример ввода пяти чисел:\n1 2 3 4 5 ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nums, err
	}
	text = strings.TrimSpace(text)
	numbers := strings.Split(text, " ")
	nums = make([]float64, len(numbers))
	for i, val := range numbers {
		nums[i], err = strconv.ParseFloat(val, 64)
		if err != nil {
			break
		}
	}
	return nums, err
}

func createLogFile(fname string) *os.File {
	lf, err := os.Create(fname)
	if err != nil {
		log.Fatal("Ошибка при создании", fname, " файла.", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(multiWriter)
	return lf
}

// Функция ниже десериализирует json файл в структуру типа Config
func decoding(config *Config) error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(data), config); err != nil {
		return err
	}
	return nil
}
func sum(nums []float64) {
	sum := 0.0
	for _, val := range nums {
		sum += val
	}
	log.Println("Посчитанная сумма всех чисел в массиве:", sum)
}

// Функция ниже делает Get запрос на указанный URL в json файле и выводит статус ответа
func responceStatus(URL string) error {
	var err error
	if URL != "" {
		resp, err := http.Get(URL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			log.Println("Статус ответа: 200;", "URL:", URL)
		} else {
			log.Println("Неожиданный статус ответа:", resp.StatusCode, "Ожидаемый статус ответа: 200.", "URL:", URL)
		}
	} else {
		log.Println("Получена пустая строка вместо URL ссылки")
	}
	return err
}
