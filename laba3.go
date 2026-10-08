package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

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

func sinTaylor(x float64, e float64) float64 {
	var term float64 = 1 * (math.Pow(x, 2*0+1))
	var d float64 = 0
	var sum float64 = term
	var k float64 = 1
	for math.Abs(term-d) >= e {
		d = term
		term = -term * x * x / ((2 * k) * (2*k + 1))
		sum += term
		k++
	}
	return sum

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
			fmt.Println("Укажите значение аругмента функции: ")
			x, _ := reader.ReadString('\n')
			x = strings.TrimSpace(x)
			if x == "" || e == "" {
				fmt.Println("Некорректный ввод! Пожалуйста, введите корректные значения.")
				continue
			}
			parts := strings.Fields(e + " " + x)
			eVal, err1 := strconv.ParseFloat(parts[0], 64)
			xVal, err2 := strconv.ParseFloat(parts[1], 64)
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
				fmt.Println("5 - Расчёт числа π с помощью ряда Нилакана")
				fmt.Println("6 - Расчёт числа е с помощью ряда Тейлора")
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
				case 1:
					fmt.Println("Функция sin(x) по формуле Тейлора: ")
					fmt.Println("Точность: ", eVal)
					fmt.Println("sin(", xVal, ") = ", sinTaylor(xVal, eVal))
				case 2:
					//функция Расчёт функции cos(x) по формуле Тейлора
				case 3:
					//функция Расчёт функции exp(x) по формуле Тейлора
				case 4:
					//функция Расчёт числа π с помощью ряда Лейбница
				case 5:
					//функция Расчёт числа π с помощью ряда Нилакана
				case 6:
					//функция Расчёт числа е с помощью ряда Тейлора
				default:
					fmt.Println("Нет такого пункта меню!")
				}
				fmt.Println("Нажмите Enter для продолжения работы......")
				inputSides, _ := reader.ReadString('\n')
				inputSides = strings.TrimSpace(inputSides)
				clearScreen()
				fmt.Println("Заданные значения:")
				fmt.Println("epsilon = ", eVal)
				fmt.Println("x = ", xVal)
			}
		}

	}

}
