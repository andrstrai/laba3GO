package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type CalcFunc func(x, e float64) (float64, uint, error)

type Entry struct {
	Name  string
	Fn    CalcFunc
	Exact func(float64) float64
}

// maxIter — предел итераций, чтобы медленный ряд (Лейбница) не «зависал» при очень маленьком ε
const maxIter = 100_000_000

var allFuncs = []Entry{
	{"sin(x)", sinTaylor, math.Sin},
	{"cos(x)", cosTaylor, math.Cos},
	{"exp(x)", expTaylor, math.Exp},
	{"π (Лейбниц)", piLeibniz, func(_ float64) float64 { return math.Pi }},
	{"π (Нилаканта)", piNilakantha, func(_ float64) float64 { return math.Pi }},
	{"e", eSeries, func(_ float64) float64 { return math.E }},
}

// функция для очистки экрана (работает и на Windows, и на macOS/Linux)
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func sinTaylor(x float64, e float64) (float64, uint, error) {
	var cnt uint = 0
	x = math.Mod(x, 2*math.Pi)
	var term float64 = 1 * (math.Pow(x, 2*0+1))
	var sum float64 = term
	var k float64 = 1
	for math.Abs(term) >= e {
		term = -term * x * x / ((2 * k) * (2*k + 1))
		sum += term
		k++
		cnt++
	}
	return sum, cnt, nil

}

func cosTaylor(x float64, e float64) (float64, uint, error) {
	var cnt uint = 0
	x = math.Mod(x, 2*math.Pi)
	var term float64 = 1
	sum := term
	n := 1.0

	for math.Abs(term) >= e {
		term = -term * x * x / (2*n - 1) / (2 * n)
		sum += term
		n++
		cnt++
	}
	return sum, cnt, nil
}

func expTaylor(x float64, e float64) (float64, uint, error) {
	if math.Abs(x) > 700 {
		return 0.0, 0, errors.New("Ошибка: значение модуля аргумента слишком велико, выйдите в меню" +
			" и установите его в пределах от -700 до 700")
	}

	if x < 0 {
		res, cnt, err := expTaylor(-x, e)
		if err != nil {
			return 0.0, 0, err
		}
		return 1 / res, cnt, nil
	}

	var term float64 = 1
	sum := term
	n := 1.0
	var cnt uint = 0

	for math.Abs(term) >= e {
		term = term * x / n
		sum += term
		n++
		cnt++
	}
	return sum, cnt, nil
}

// piLeibniz вычисляет π по ряду Лейбница: π/4 = 1 - 1/3 + 1/5 - 1/7 + ...
func piLeibniz(_ float64, e float64) (float64, uint, error) {
	sum := 0.0       // частичная сумма ряда для π/4
	sign := 1.0      // знак очередного члена: +1, -1, +1, ...
	var cnt uint = 0 // число сложенных членов (итераций)

	for {
		term := sign / float64(2*cnt+1) // член ряда: ±1/(2n+1), n = cnt
		// в конце сумма умножается на 4, значит и ошибка тоже: сравниваем с ε член × 4
		if 4*math.Abs(term) < e {
			break // член меньше точности — дальше не считаем
		}
		if cnt >= maxIter {
			return 0, cnt, errors.New("ряд Лейбница: слишком много итераций, увеличьте ε")
		}
		sum += term
		sign = -sign
		cnt++
	}
	return 4 * sum, cnt, nil
}

// piNilakantha вычисляет π по ряду Нилаканты: π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - ...
func piNilakantha(_ float64, e float64) (float64, uint, error) {
	sum := 3.0       // ряд начинается с 3
	sign := 1.0      // первый член после тройки идёт с плюсом
	k := 1.0         // номер члена
	var cnt uint = 0 // число сложенных членов

	term := 4 / ((2 * k) * (2*k + 1) * (2*k + 2)) // первый член 4/(2·3·4)
	for term >= e {
		sum += sign * term
		sign = -sign
		k++
		cnt++
		term = 4 / ((2 * k) * (2*k + 1) * (2*k + 2)) // следующий член сразу по формуле
	}
	return sum, cnt, nil
}

// eSeries вычисляет число e по ряду 1/0! + 1/1! + 1/2! + 1/3! + ... (это exp(1)).
func eSeries(_ float64, e float64) (float64, uint, error) {
	term := 1.0      // 1/0! = 1
	sum := term      // частичная сумма
	n := 1.0         // номер следующего члена
	var cnt uint = 0 // число итераций

	for {
		term /= n // 1/n! = (1/(n-1)!) / n — следующий член из предыдущего
		if term < e {
			break
		}
		sum += term
		n++
		cnt++
	}
	return sum, cnt, nil
}

// printResult печатает результат одного вычисления: приближённое и точное значения,
func printResult(name string, approx, exact float64, iters uint, e float64) {
	fmt.Println(name)
	fmt.Printf("Точность ε:          %v\n", e)
	fmt.Printf("Приближённое:        %.15f\n", approx)
	fmt.Printf("Точное (пакет math): %.15f\n", exact)
	fmt.Printf("Погрешность:         %.2e\n", math.Abs(approx-exact))
	fmt.Printf("Итераций:            %d\n", iters)
}

// calcAndPrint считает одну запись из allFuncs и печатает результат или ошибку
func calcAndPrint(en Entry, x, e float64) {
	approx, iters, err := en.Fn(x, e)
	if err != nil {
		fmt.Println(err)
		return
	}
	printResult(en.Name, approx, en.Exact(x), iters, e)
}

func printTable(x float64, e float64) {
	fmt.Printf("\n\t\t\t   Сравнительная таблица функций при точности %v\n", e)
	fmt.Println(strings.Repeat("-", 115))
	fmt.Printf("%-20s | %-30s | %-30s | %-15s | %-10s\n",
		"Функция", "Приближённое", "Точное", "Погрешность", "Итераций")
	fmt.Println(strings.Repeat("-", 115))

	for _, en := range allFuncs {
		approx, iters, err := en.Fn(x, e)
		if err != nil {
			fmt.Printf("%-20s | %-30s | %-30s | %-15s | %-10s\n",
				en.Name, "ошибка", "ошибка", "—", "—")
			continue
		}
		exact := en.Exact(x)
		eps := math.Abs(approx - exact)
		fmt.Printf("%-20s | %-30v | %-30v | %-15.2e | %d\n",
			en.Name, approx, exact, eps, iters)
	}
	fmt.Println(strings.Repeat("-", 115))

}

// основная функция
func main() {
	reader := bufio.NewReader(os.Stdin)
A:
	for {
		fmt.Println("Для начала работы пожалуйста, нажмите Enter, для завершения введите 0.....")
		inputStart, _ := reader.ReadString('\n')
		inputStart = strings.TrimSpace(inputStart)
		if inputStart == "0" {
			return
		} else {
			clearScreen()
			fmt.Println("Укажите точность вычислений в формате ( 0 < epsilon < 1): ")
			e, _ := reader.ReadString('\n')
			e = strings.TrimSpace(e)
			fmt.Println("Укажите значение аргумента функции: ")
			x, _ := reader.ReadString('\n')
			x = strings.TrimSpace(x)
			if x == "" || e == "" {
				fmt.Println("Некорректный ввод! Пожалуйста, введите корректные значения.")
				continue
			}
			eVal, err1 := strconv.ParseFloat(e, 64)
			xVal, err2 := strconv.ParseFloat(x, 64)
			if err1 != nil || err2 != nil || eVal <= 0 || eVal >= 1 {
				fmt.Println("Некорректный ввод! Пожалуйста, введите корректные значения.")
				continue
			}
			for {
				fmt.Println("Меню команд:")
				fmt.Println("1 - Расчёт функции sin(x) по формуле Тейлора")
				fmt.Println("2 - Расчёт функции cos(x) по формуле Тейлора")
				fmt.Println("3 - Расчёт функции exp(x) по формуле Тейлора")
				fmt.Println("4 - Расчёт числа π с помощью ряда Лейбница")
				fmt.Println("5 - Расчёт числа π с помощью ряда Нилаканты")
				fmt.Println("6 - Расчёт числа e с помощью ряда Тейлора")
				fmt.Println("7 - Вывести таблицу со всей статистикой")
				fmt.Println("0 - Выйти из меню команд")
				fmt.Println("Выберите пункт из меню управления: ")

				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)
				a, err := strconv.Atoi(input)
				if err != nil {
					fmt.Println("Некорректный ввод! Пожалуйста, введите число.")
					continue
				}

				switch a {
				case 0:
					continue A
				case 1, 2, 3, 4, 5, 6:
					calcAndPrint(allFuncs[a-1], xVal, eVal)
				case 7:
					printTable(xVal, eVal)
				default:
					fmt.Println("Нет такого пункта меню!")
				}
				fmt.Println("Нажмите Enter для продолжения работы......")
				reader.ReadString('\n')
				clearScreen()
				fmt.Println("Заданные значения:")
				fmt.Println("epsilon = ", eVal)
				fmt.Println("x = ", xVal)
			}
		}

	}

}
