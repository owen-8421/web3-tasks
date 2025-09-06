package task3

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

/*
Sqlx入门
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

*/

// Employee 结构体用于映射 employees 表的数据
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

// findEmployeesByDepartment 使用 sqlx.Select 查询多行数据
func findEmployeesByDepartment(db *sqlx.DB, department string) ([]Employee, error) {
	var employees []Employee
	query := "SELECT id, name, department, salary FROM employees WHERE department = ?"

	// Select 会自动查询、映射并填充到 employees 切片中
	err := db.Select(&employees, query, department)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

// findHighestPaidEmployee 使用 sqlx.Get 查询薪资最高的员工
func findHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee
	query := "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1"

	// Get 会自动查询、映射并将结果填充到 employee 结构体中
	// 如果查询结果为空，它会返回 sql.ErrNoRows
	err := db.Get(&employee, query)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func test() {
	// 连接数据库
	dsn := "username:password@tcp(localhost:3306)/dbname"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}
	TechEmployees, err := findEmployeesByDepartment(db, "技术部")
	fmt.Println("技术部员工:", TechEmployees)

	// 查询工资最高的员工
	MaxSalaryEmployee, err := findHighestPaidEmployee(db)
	fmt.Println("薪资最高的员工:", MaxSalaryEmployee)
}

/*题目2：实现类型安全映射
假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
要求 ：
定义一个 Book 结构体，包含与 books 表对应的字段。
编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
*/

// Book 结构体用于映射 books 表的数据。
// `db` 标签告诉 sqlx 如何将数据库列名（如 'title'）映射到 Go 结构体的字段（如 'Title'）。
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// findExpensiveBooks 查询价格高于指定阈值的书籍。
// 它接收一个 sqlx.DB 连接实例和最低价格，返回一个 Book 切片或一个错误。
func findExpensiveBooks(db *sqlx.DB, minPrice float64) ([]Book, error) {
	// 声明一个 Book 切片用于存放查询结果。
	var books []Book

	// 定义我们的 SQL 查询语句，使用 '?'作为占位符。
	query := "SELECT id, title, author, price FROM books WHERE price > ?"

	// 使用 db.Select 来执行查询。
	// sqlx 会安全地处理参数，并自动将结果映射到 books 切片中。
	err := db.Select(&books, query, minPrice)
	if err != nil {
		// 如果查询出错（例如，数据库连接问题），则返回错误。
		return nil, err
	}

	// 返回填充好的切片和 nil 错误。
	return books, nil
}

func test2() {
	// --- 1. 连接数据库 ---
	// 请将 DSN 替换为您的数据库连接信息
	dsn := "root:your_password@tcp(127.0.0.1:3306)/sqlx_books_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	priceThreshold := 50.0
	fmt.Printf("--- 查询价格大于 %.2f 元的书籍 ---\n", priceThreshold)

	expensiveBooks, err := findExpensiveBooks(db, priceThreshold)
	if err != nil {
		log.Fatalf("查询书籍失败: %v", err)
	}

	// --- 4. 打印结果 ---
	if len(expensiveBooks) == 0 {
		fmt.Println("未找到价格高于该值的书籍。")
	} else {
		for _, book := range expensiveBooks {
			fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %.2f\n", book.ID, book.Title, book.Author, book.Price)
		}
	}
}
