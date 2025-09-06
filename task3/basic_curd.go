package task3

import (
	"leetcode/database"
	"log"
)

/*
题目1：基本CRUD操作
假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
要求 ：
编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
*/

/*
INSERT INTO students (name, age, grade) VALUES ('张三', 20, '三年级');
SELECT * FROM students WHERE age > 18;
UPDATE students SET grade = '四年级' WHERE name = '张三';
DELETE FROM students WHERE age < 15;
*/
const TableName = "students"

// Students 学生表数据结构
type Students struct {
	Id    uint64 `gorm:"column:id" json:"id"`       // 主键
	Name  string `gorm:"column:name" json:"name"`   // 学生姓名
	Age   int    `gorm:"column:age" json:"age"`     // 年龄
	Grade string `gorm:"column:grade" json:"grade"` // 学生年级
}

// StudentService 支持学生表的基础CRUD
type StudentService struct {
}

// TblName 获取表名
func (s *StudentService) TblName() string {
	return TableName
}

// Upsert 实现对students表的 upsert功能
// INSERT INTO students (name, age, grade) VALUES ('张三', 20, '三年级');
// UPDATE students SET grade = '四年级' WHERE name = '张三';
func (s *StudentService) Upsert(student *Students, name string) error {
	myDB, _ := database.GetDB("local")
	var tmp Students
	err := myDB.Where("name = ?", name).
		Find(&tmp).Error
	if err != nil {
		log.Printf("select failed: %s", err.Error())
		return err
	}
	if tmp.Id > 0 { // 数据库中已经存在，则更新数据
		err = myDB.Where("name = ?", name).Updates(student).Error
		return err
	} else { // 做insert操作
		err = myDB.Save(&student).Error
		return err
	}
}

func (s *StudentService) Select(minAge int) ([]*Students, error) {
	// SELECT * FROM students WHERE age > 18;
	myDB, err := database.GetDB("local")
	var res []*Students
	err = myDB.Where("age > ?", minAge).Select("*").Find(&res).Error
	return res, err
}

func (s *StudentService) Delete(maxAge int) error {
	// DELETE FROM students WHERE age < 15;
	myDB, _ := database.GetDB("local")
	err := myDB.Where("age > ?", maxAge).Delete("*").Error
	return err
}
