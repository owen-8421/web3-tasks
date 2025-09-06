package task3

import "time"

/*
题目2：事务语句
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表
（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

/*
-- 创建 accounts 表
CREATE TABLE accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    balance DECIMAL(10, 2) NOT NULL CHECK (balance >= 0)
);

-- 创建 transactions 表
CREATE TABLE transactions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    from_account_id INT,
    to_account_id INT,
    amount DECIMAL(10, 2) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    FOREIGN KEY (to_account_id) REFERENCES accounts(id)
);

-- 插入示例数据
-- 账户A (id=1) 余额为 500 元
-- 账户B (id=2) 余额为 1000 元
INSERT INTO accounts (id, balance) VALUES (1, 500.00);
INSERT INTO accounts (id, balance) VALUES (2, 1000.00);

-- 声明变量以便灵活修改
SET @from_account = 1; -- 转出账户ID (账户A)
SET @to_account = 2;   -- 转入账户ID (账户B)
SET @transfer_amount = 100; -- 转账金额

-- 开始事务
START TRANSACTION;

-- 1. 查询并锁定转出账户的余额，防止其他事务干扰
SELECT balance INTO @current_balance FROM accounts WHERE id = @from_account FOR UPDATE;

-- 2. 检查余额是否足够
IF @current_balance >= @transfer_amount THEN
    -- 3. 如果余额足够，执行更新操作
    -- 3.1 从账户A扣款
    UPDATE accounts SET balance = balance - @transfer_amount WHERE id = @from_account;

    -- 3.2 向账户B存款
    UPDATE accounts SET balance = balance + @transfer_amount WHERE id = @to_account;

    -- 3.3 在 transactions 表中记录日志
    INSERT INTO transactions (from_account_id, to_account_id, amount)
    VALUES (@from_account, @to_account, @transfer_amount);

    -- 提交事务，使所有更改永久生效
    COMMIT;
    SELECT '转账成功！' AS '结果';

ELSE
    -- 4. 如果余额不足，回滚事务
    ROLLBACK;
    SELECT '余额不足，转账失败，事务已回滚。' AS '结果';

END IF;
*/

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Account 模型对应 accounts 表
// 在实际金融应用中，建议使用支持高精度的 decimal 类型代替 float64
type Account struct {
	ID        uint    `gorm:"primaryKey"`
	Balance   float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Transaction 模型对应 transactions 表
type Transaction struct {
	ID            uint    `gorm:"primaryKey"`
	FromAccountID uint    `gorm:"not null"`
	ToAccountID   uint    `gorm:"not null"`
	Amount        float64 `gorm:"not null"`
	CreatedAt     time.Time
}

// Transfer 函数执行转账操作
// 它接收一个 *gorm.DB 实例，以及转账所需的参数
func Transfer(db *gorm.DB, fromID, toID uint, amount float64) error {
	// GORM 的核心事务逻辑
	// 所有在此函数内的数据库操作都属于同一个事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 0. 检查参数有效性
		if amount <= 0 {
			return errors.New("转账金额必须为正数")
		}
		if fromID == toID {
			return errors.New("不能给自己转账")
		}

		// 1. 锁定转出账户，防止并发问题 (相当于 SQL 的 SELECT ... FOR UPDATE)
		var fromAccount Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fromAccount, fromID).Error; err != nil {
			return fmt.Errorf("找不到转出账户: %w", err)
		}

		// 2. 检查余额是否足够
		if fromAccount.Balance < amount {
			return errors.New("账户余额不足")
		}

		// 3. 从转出账户扣款
		// 使用 gorm.Expr 可以安全地执行数据库层面的计算
		result := tx.Model(&fromAccount).Update("balance", gorm.Expr("balance - ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("扣款失败，未找到账户")
		}

		// 4. 向转入账户加款
		result = tx.Model(&Account{}).Where("id = ?", toID).Update("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}
		// 检查是否真的更新了收款方，如果没有找到该账户，则回滚
		if result.RowsAffected == 0 {
			return errors.New("加款失败，未找到收款账户")
		}

		// 5. 创建交易记录
		transaction := Transaction{
			FromAccountID: fromID,
			ToAccountID:   toID,
			Amount:        amount,
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		// 返回 nil，GORM 将自动提交事务
		return nil
	})
	return err
}
