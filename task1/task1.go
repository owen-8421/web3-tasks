package task1

import "slices"

/**
136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，
其余每个元素均出现两次。找出那个只出现了一次的元素。可以使用 for 循环遍历数组，
结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，
然后再遍历 map 找到出现次数为1的元素。
*/

/*
只出现一次的数字，利用异或的特性，相同的数字异或后为0，只剩下只出现一次的数字
*/

func singleNumber(nums []int) int {
	res := 0
	for _, v := range nums {
		res ^= v
	}
	return res
}

/*
9. 回文数
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
*/

/*
题解：如果是负数，回文数不成立，返回false
如果是正数，将数字回文操作一次，如果和原数相同，返回true，否则false
*/
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	// 回文数 将数字回文操作一次，如果和x相等，将返回true
	newX := 0
	num := x
	for num != 0 {
		newX = newX * 10
		tp := num % 10 // 获取个位数
		newX += tp     // 新的数
		num = num / 10
	}
	return newX == x
}

/*
20. 有效的括号
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。

*/

/*
使用堆栈，判断是左括号的是否，将右括号入栈，
遇到右括号的时候，判断栈顶是否一致，一致就出栈，不一致就返回false

遍历完成后，检查堆栈是否为空
*/
func isValid(s string) bool {
	// 使用堆栈，看括号是否能够match
	if len(s)%2 != 0 {
		return false
	}
	st := []rune{}
	for _, v := range s {
		switch v {
		case '(':
			st = append(st, ')') // 入栈对应的右括号
		case '[':
			st = append(st, ']')
		case '{':
			st = append(st, '}')
		default: // c 是右括号
			if len(st) == 0 || st[len(st)-1] != v {
				return false // 没有左括号，或者左括号类型不对
			}
			st = st[:len(st)-1] // 出栈
		}
	}
	return len(st) == 0
}

/*
14. 最长公共前缀

编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。

示例 1：

输入：strs = ["flower","flow","flight"]
输出："fl"

简单题：使用暴力解法，按列比较，直到找到不相等的列
返回下标
*/
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	//取每一个单词的同一位置的字母，看是否相同。
	s0 := strs[0]
	for j, c := range s0 {
		for _, s := range strs {
			if j == len(s) || s[j] != byte(c) { // 找到不相等的列，或者缺失
				return s0[:j]
			}
		}
	}
	return s0
}

/*
66. 加一
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。
这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
将大整数加 1，并返回结果的数字数组。
*/

/*
数字操作，注意最高位进位
*/

func plusOne(digits []int) []int {
	// 数字按从左到右，从最高位到最低位排列，所以需要从右到左遍历。计算数组 + 1的值

	length := len(digits)
	plusX := 0
	for i := length - 1; i >= 0; i-- {
		//判断数字加1是否大于10
		if i == length-1 { // 最小的数加一
			plusX = digits[i] + 1
			digits[i] = plusX % 10 // 最新的值
			plusX = plusX / 10
		} else {
			if plusX > 0 {
				plusX = digits[i] + plusX // 进位和当前位相加
				digits[i] = plusX % 10
				plusX = plusX / 10
			}
		}
	}
	if plusX <= 0 {
		return digits
	} else {
		res := make([]int, 0)
		res = append(res, plusX)
		res = append(res, digits...)
		return res
	}
}

/*
26. 删除有序数组中的重复项
给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，
使每个元素 只出现一次 ，返回删除后数组的新长度。元素的 相对顺序 应该保持
一致 。然后返回 nums 中唯一元素的个数。
考虑 nums 的唯一元素的数量为 k ，你需要做以下事情确保你的题解可以被通过：
更改数组 nums ，使 nums 的前 k 个元素包含唯一元素，并按照它们最初在
nums 中出现的顺序排列。nums 的其余元素与 nums 的大小不重要。
返回 k 。
*/

/*
思路：使用快慢指针，将非重复数字移动到最左边，返回下标 slow + 1
*/

func removeDuplicates(nums []int) int {
	// 使用双指针，一快一慢，将非重复数移动到最左边
	slow := 0
	fast := 0
	for fast < len(nums) {
		if nums[slow] < nums[fast] {
			slow++
			nums[slow] = nums[fast]
		}
		fast++
	}
	return slow + 1
}

/*
56. 合并区间
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回 一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间 。
*/

func merge(intervals [][]int) (ans [][]int) {
	// 按照左端点从小到大排序
	slices.SortFunc(intervals, func(p, q []int) int { return p[0] - q[0] })
	for _, p := range intervals {
		m := len(ans)
		if m > 0 && p[0] <= ans[m-1][1] { // 可以合并
			ans[m-1][1] = max(ans[m-1][1], p[1]) // 更新右端点最大值
		} else { // 不相交，无法合并
			ans = append(ans, p) // 新的合并区间
		}
	}
	return
}

/*
1. 两数之和
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target
的那 两个 整数，并返回它们的数组下标。
你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
你可以按任意顺序返回答案。
*/

/*
获取两数之和等于target的两个数的下标，
使用哈希表记录已经遍历过的数： key = 数值, value = 下标，遍历到一个新的数cur_val，
检查hash map中是否有target - cur_val
如果有，返回 对应下标，如果没有，将当前值写入hash map
*/

func twoSum(nums []int, target int) []int {
	index := make(map[int]int, 0)
	for i, v := range nums {
		if k, ok := index[target-v]; !ok {
			index[v] = i
		} else {
			return []int{i, k}
		}
	}
	return []int{}
}
