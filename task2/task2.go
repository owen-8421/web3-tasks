package task2

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
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

/*
题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/

type Person struct {
	Name string
	Age  int32
}

type Employee struct {
	EmployeeID string
	Person
}

func (e *Employee) PrintInfo() string {
	return fmt.Sprintf("EmployeeID: %s, Name: %s, Age: %d", e.EmployeeID, e.Name, e.Age)
}

func TestObj2() {
	e := Employee{EmployeeID: "1001", Person: Person{Age: 26, Name: "张三"}}
	res := e.PrintInfo()
	fmt.Println(res)
}

/*
编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从 1 到 10 的整数，
并将这些整数发送到通道中；另一个协程从通道中接收这些整数并打印出来。
*/

func ChannelPrint() {
	// 创建一个用于传输整数的通道
	ch := make(chan int)

	// 创建一个WaitGroup，用于等待协程完成
	var wg sync.WaitGroup

	wg.Add(2) // 我们需要等待两个协程（一个生产者，一个消费者）

	// 1 启动生产者
	go func() {
		// 任务完成时，通知wait group
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			fmt.Printf("生产者：发送 -> %d\n", i)
			ch <- i
		}
		// 数据发送完毕后，关闭通道
		// 这是通知消费方没有更多数据的信号
		close(ch)
	}()

	// 2. 启动“消费者”协程
	go func() {
		// 任务完成时，通知 WaitGroup
		defer wg.Done()

		// 使用 for range 循环从通道接收数据，直到通道被关闭
		// 接收到的信息打印
		for number := range ch {
			fmt.Printf("消费者：接收 <- %d\n", number)
		}
	}()

	// 3. 主协程等待
	fmt.Println("主协程：等待所有子协程执行完毕...")
	// Wait() 会阻塞，直到 WaitGroup 的计数器归零
	wg.Wait()

	fmt.Println("主协程：所有任务完成，程序即将退出。")
}

/*
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

/*
缓冲机制：这意味着生产者可以立即向通道发送 10 个整数，
而无需等待消费者接收。只有当缓冲区满了（即通道中已经有10个未被读取的元素），
生产者再次尝试发送时才会被阻塞。同样，如果缓冲区为空，消费者尝试接收时会被阻塞。

缓冲通道解耦生产者和消费者速度的能力，允许它们在一定程度上独立运行，从而提高整体程序的效率
*/

func ChannelPrintBuffer() {
	// 创建一个用于传输整数的通道,容量设置为 10。
	// 这意味着生产者可以连续发送 10 个整数而不会被阻塞。
	ch := make(chan int, 10)

	// 创建一个WaitGroup，用于等待协程完成
	var wg sync.WaitGroup

	wg.Add(2) // 我们需要等待两个协程（一个生产者，一个消费者）

	// 1 启动生产者
	go func() {
		// 任务完成时，通知wait group
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			fmt.Printf("生产者：发送 -> %d\n", i)
			ch <- i
		}
		// 数据发送完毕后，关闭通道
		// 这是通知消费方没有更多数据的信号
		close(ch)
	}()

	// 2. 启动“消费者”协程
	go func() {
		// 任务完成时，通知 WaitGroup
		defer wg.Done()

		// 使用 for range 循环从通道接收数据，直到通道被关闭
		// 接收到的信息打印
		for number := range ch {
			fmt.Printf("消费者：接收 <- %d\n", number)
		}
	}()

	// 3. 主协程等待
	fmt.Println("主协程：等待所有子协程执行完毕...")
	// Wait() 会阻塞，直到 WaitGroup 的计数器归零
	wg.Wait()

	fmt.Println("主协程：所有任务完成，程序即将退出。")
}

/*
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。

当多个协程并发执行时，可能会发生以下情况：协程A读取了 counter 的值为5，
此时协程B也读取了值为5。然后协程A计算后将6写回，协程B计算后也将6写回。
结果是两个协程都执行了递增，但 counter 的值只增加了1。这就是数据竞争，
它会导致最终的计数值不可预知，且几乎肯定会小于10000。

*/

// SafeCounter 定义一个共享的计数器结构体
type SafeCounter struct {
	mu      sync.Mutex // 互斥锁，用于保护下面的计数值
	counter int        // 共享的计数值
}

// Inc 方法用于增加计数器的值
// 它会先获取锁，执行递增操作，然后释放锁
func (c *SafeCounter) Inc() {
	c.mu.Lock() // 加锁
	// 在 Lock 和 Unlock 之间的代码是临界区，
	// 有锁的保护，我们可以确保 c.counter++ 这个递增操作是原子性的，不会被其他协程中断。
	c.counter++
	c.mu.Unlock() // 解锁
}

// GetVal 方法用于获取计数器的当前值
// 它同样需要加锁来保证读取的是一个一致性的值
func (c *SafeCounter) GetVal() int {
	c.mu.Lock()
	// defer 语句确保在函数返回前一定会执行 Unlock 操作
	// 保证在读取当前值时，counter不被同时修改
	defer c.mu.Unlock()
	// 使用 defer 来调用 Unlock 它能确保无论函数从哪个路径返回
	// 例如，即使中间发生了 panic，解锁操作都一定会被执行，从而避免死锁。
	return c.counter
}

func LockExp() {
	// 创建 SafeCounter 实例
	sc := SafeCounter{}
	// 使用 WaitGroup 来等待所有协程完成
	var wg sync.WaitGroup

	numGoroutines := 10   // 协程数量
	numIncrements := 1000 // 每个协程的递增次数

	// 1. 启动 10 个协程
	for i := 0; i < numGoroutines; i++ {
		// 告诉 WaitGroup 需要等待一个协程
		wg.Add(1)
		go func() {
			// 在协程退出时，通知 WaitGroup 该协程已完成
			defer wg.Done()
			// 2. 每个协程对计数器递增 1000 次
			for j := 0; j < numIncrements; j++ {
				sc.Inc()
			}
		}()
	}

	// 3. 等待所有协程执行完毕
	wg.Wait()

	// 4. 所有协程完成后，打印最终的计数值
	// 期望的结果是 10 * 1000 = 10000
	fmt.Printf("最终计数器的值: %d\n", sc.GetVal())
}

/*
使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。

sync/atomic 包的基本使用。
原子操作的概念及其与锁的区别。
如何在无锁的情况下保证并发数据安全。
*/

// AtomicCounter 定义一个无锁的原子计数器
type AtomicCounter struct {
	counter int64 // 使用 int64/uint64 以便使用 atomic 包中的函数
}

// Inc 方法使用原子操作来递增计数器
func (c *AtomicCounter) Inc() {
	// atomic.AddInt64 会原子性地给 c.counter 加上 1，并返回新值。
	// 这是一个单一的、不可中断的操作。
	atomic.AddInt64(&c.counter, 1)
}

// GetVal 方法使用原子操作来读取计数器的值
func (c *AtomicCounter) GetVal() int64 {
	// atomic.LoadInt64 会原子性地读取 c.counter 的值。
	// 这可以防止在读取过程中，值被其他协程修改，从而读取到不一致的数据。
	return atomic.LoadInt64(&c.counter)
}

func AtomicExp() {
	// 创建 AtomicCounter 实例
	ac := AtomicCounter{}

	// 使用 WaitGroup 来等待所有协程完成
	var wg sync.WaitGroup

	numGoroutines := 10   // 协程数量
	numIncrements := 1000 // 每个协程的递增次数

	// 1. 启动 10 个协程
	for i := 0; i < numGoroutines; i++ {
		// 告诉 WaitGroup 需要等待一个协程
		wg.Add(1)
		go func() {
			// 在协程退出时，通知 WaitGroup 该协程已完成
			defer wg.Done()
			// 2. 每个协程对计数器递增 1000 次
			for j := 0; j < numIncrements; j++ {
				ac.Inc()
			}
		}()
	}

	// 3. 等待所有协程执行完毕
	wg.Wait()

	// 4. 所有协程完成后，打印最终的计数值
	// 期望的结果是 10 * 1000 = 10000
	fmt.Printf("最终计数器的值: %d\n", ac.GetVal())
}
