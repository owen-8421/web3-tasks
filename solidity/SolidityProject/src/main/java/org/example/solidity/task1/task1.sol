// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract task1 {

    /*
     * @notice 反转一个字符串。
     * @param _input 需要被反转的原始字符串。
     * @return 反转后的新字符串。
     * 这个函数是 "pure" 类型，因为它既不读取也不修改合约的状态。它仅仅是根据输入计算并返回一个输出。
     * 解题思路： 1. 将 string 转换为 bytes 数组，这样我们才能访问到每一个字节。
     *          2. 获取字节数组的长度，新建输出结果
     *          3. 遍历，进行反转操作
     */

    function reverseString(string memory _input) public pure returns (string memory) {
        // 1. 将 string 转换为 bytes 数组，这样我们才能访问到每一个字节。
        bytes memory inputBytes = bytes(_input);
        // 2. 获取字节数组的长度。
        uint length = inputBytes.length;
        bytes memory reversedBytes = new bytes(length);
        // 遍历，进行反转操作。
        for (uint i = 0; i < length; i++) {
            reversedBytes[i] = inputBytes[length - 1 - i];
        }
        return string(reversedBytes);
    }

    /*
   * @notice 整数转罗马数字
     * @param _input 需要被转化的整数
     * @return 转化后的罗马数字
     * 这个函数是 "pure" 类型，因为它既不读取也不修改合约的状态。它仅仅是根据输入计算并返回一个输出。
     */

    function intToRoman(uint256 _num) public pure returns (string memory) {
        require(_num > 0 && _num < 4000, "Input must be between 1 and 3999");

        uint256[] memory values = new uint256[](13);
        values[0] = 1000;
        values[1] = 900;
        values[2] = 500;
        values[3] = 400;
        values[4] = 100;
        values[5] = 90;
        values[6] = 50;
        values[7] = 40;
        values[8] = 10;
        values[9] = 9;
        values[10] = 5;
        values[11] = 4;
        values[12] = 1;

        string[] memory symbols = new string[](13);
        symbols[0] = "M";
        symbols[1] = "CM";
        symbols[2] = "D";
        symbols[3] = "CD";
        symbols[4] = "C";
        symbols[5] = "XC";
        symbols[6] = "L";
        symbols[7] = "XL";
        symbols[8] = "X";
        symbols[9] = "IX";
        symbols[10] = "V";
        symbols[11] = "IV";
        symbols[12] = "I";

        string memory result = "";
        uint256 num = _num;

        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = string.concat(result, symbols[i]);
            }
        }

        return result;
    }



    /**
     * @notice 一个内部辅助函数，返回单个罗马字符对应的整数值。
     * @param _char 单个罗马字符的字节表示。
     * @return 对应的整数值。如果字符无效，则返回0。
     */
    function getValue(bytes1 _char) private pure returns (uint) {
        if (_char == 'I') return 1;
        if (_char == 'V') return 5;
        if (_char == 'X') return 10;
        if (_char == 'L') return 50;
        if (_char == 'C') return 100;
        if (_char == 'D') return 500;
        if (_char == 'M') return 1000;
        return 0; // 对于无效字符返回0
    }

    /*
     * @notice 罗马数字转整数。
     * @param _input 需要被转化的罗马数字
     * @return 转化后的整数
     * 这个函数是 "pure" 类型，因为它既不读取也不修改合约的状态。它仅仅是根据输入计算并返回一个输出。
     */
    function romanToInt (string memory _input) public pure returns (int memory) {
        bytes memory s = bytes(_input);
        uint length = s.length;

        // 输入不能为空
        if (length == 0) {
            return 0;
        }

        uint result = 0;

        for (uint i = 0; i < length; i++) {
            uint currentValue = getValue(s[i]);
            // 检查是否有下一个字符，并且当前值是否小于下一个值
            if (i + 1 < length) {
                uint nextValue = getValue(s[i+1]);
                if (currentValue < nextValue) {
                    // 应用减法规则
                    result -= currentValue;
                } else {
                    // 应用加法规则
                    result += currentValue;
                }
            } else {
                // 这是最后一个字符，直接相加
                result += currentValue;
            }
        }

        return result;
    }

    /*
     * @notice 合并两个有序数组
     * @param  需要被合并有序数组：_input1， _input2
     * @return 合并后的数组 output
     */

    function mergeArray (uint256[] memory _a, uint256[] memory _b) public pure returns (uint256[] memory) {
        uint256[] memory result = new uint256[](_a.length + _b.length);
        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;

        while (i < _a.length && j < _b.length ) {
            if (_a[i] < _b[j]) {
                result[k] = _a[i];
                i++;
            } else {
                result[k] = _b[j];
                j++;
            }
            k++;
        }

        while (i < _a.length) {
            result[k] = _a[i];
            i++;
            k++;
        }

        while (j < _b.length) {
            result[k] = _b[j];
            j++;
            k++;
        }

        return result;

    }

    /*
    * @notice 二分查找: 在一个有序数组中查找目标值。
     * @param  有序数组：arr， 需要查找的目标target
     * @return 查到的下标？还是说是否查到？
     */

    function mergeArray (uint256[] memory arr, uint256 memory target) public pure returns (bool memory) {
        uint256 left = 0;
        uint256 right = arr.length - 1;

        while (left <= right) {
            uint256 mid = left + (right - left) / 2;
            if (arr[mid] == target) {
                return true;
            }
            if (arr[mid] < target) {
                left = mid + 1
            } else {
                if (mid == 0) {
                    return false;
                }
                right = mid - 1;
            }
        }
        return false;
    }

}
