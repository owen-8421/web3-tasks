package task2

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。
*/

func PlusTen(input *int) *int {
	*input = *input + 10
	return input
}

// 测试
func Test() {
	input := 10
	res := PlusTen(&input)
	fmt.Println(*res)
}

/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/

func SliceMtl(input *[]int) {
	for i, v := range *input {
		(*input)[i] = v * 2
	}
	return
}

// 测试
func Test2() {
	input := []int{1, 2, 3, 4}
	input2 := []int{1, 2, 3, 4}
	SliceMtl(&input)
	for i, v := range input2 {
		fmt.Printf("before %d, after %d\n", v, input[i])
	}
}

/*
编写一个程序，使用 go 关键字启动两个协程。
一个协程打印从 1 到 10 的奇数，
另一个协程打印从 2 到 10 的偶数。
*/

func GoroutinePrint() {
	var wg sync.WaitGroup
	wg.Add(2) // 需要等待两个协程，所以add 2

	// 启动一个协程打印奇数

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Printf("奇数：%d\n", i)
			}
		}
	}()

	// 启动一个协程打印偶数
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Printf("偶数：%d\n", i)
			}
		}
	}()

	fmt.Println("主协程正在等待...")
	// Wait() 会阻塞当前协程（即 main 协程），
	// 直到 WaitGroup 的计数器变为零。
	wg.Wait()

	fmt.Println("所有协程执行完毕。")
}

/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），
并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

func GoroutineExp(tasks []func()) {
	// 创建一个带缓冲的通道，容量等于任务数，这样协程发送结果时不会被阻塞
	resultsChan := make(chan string, len(tasks))
	var wg sync.WaitGroup
	// 告诉 WaitGroup 我们需要等待多少个任务。
	wg.Add(len(tasks))
	fmt.Println("开始执行并发任务...")
	totalStart := time.Now()
	// 遍历所有任务，为每个任务启动一个协程。
	for i, task := range tasks {
		// 使用 go 关键字启动协程
		// 必须将 i 和 task 作为参数传递给匿名函数，以避免闭包问题！

		go func(taskID int, taskFn func()) {
			//defer 语句确保协程退出前一定会调用 Done()
			defer wg.Done()
			start := time.Now()
			// 执行任务函数
			taskFn()
			duration := time.Since(start) // 统计执行时间
			// 打印时间
			result := fmt.Sprintf("任务 %d 完成，耗时: %s", taskID+1, duration)
			// 将结果发送到通道
			resultsChan <- result
		}(i, task)
	}

	// 等待所有协程都调用 Done()，否则这里会一直阻塞。
	wg.Wait()
	fmt.Printf("所有任务都已完成。总耗时: %s\n", time.Since(totalStart))

	// 所有协程都已结束，可以安全地关闭通道
	close(resultsChan)

	// 从通道中读取所有结果并打印
	for result := range resultsChan {
		fmt.Println(result)
	}
}

// TestGoroutineExp 测试函数
func TestGoroutineExp() {
	tasks := []func(){
		func() {
			fmt.Println("正在执行任务 1: 下载文件...")
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
		},
		func() {
			fmt.Println("正在执行任务 2: 处理数据...")
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
		},
		func() {
			fmt.Println("正在执行任务 3: 保存到数据库...")
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
		},
		func() {
			fmt.Println("正在执行任务 4: 发送通知...")
			time.Sleep(time.Duration(rand.Intn(500)+200) * time.Millisecond)
		},
	}

	GoroutineExp(tasks)
}

/*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

const (
	pi = 3.1415926
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct { // 长方形
	Length float64
	Width  float64
}

func (r *Rectangle) Area() float64 {
	return r.Length * r.Width
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.Length + r.Width)
}

type Circle struct { // 圆
	Radius float64
}

func (r *Circle) Area() float64 {
	return pi * r.Radius * r.Radius
}

func (r *Circle) Perimeter() float64 {
	return float64(2) * pi * r.Radius
}

func TestObj() {
	r := Rectangle{Width: 4, Length: 5}
	c := Circle{Radius: 3}

	fmt.Printf("Rectangle's width: %.2f, length: %.2f, Area: %.2f, Perimeter: %.2f\n", r.Width, r.Length, r.Area(), r.Perimeter())
	fmt.Printf("Circle's Radius: %.2f, Area: %.2f, Perimeter: %.2f", c.Radius, c.Area(), c.Perimeter())
}
