// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title SimpleERC20
 * @dev 一个简单的 ERC20 代币实现，参考 OpenZeppelin 的 IERC20.sol 接口.
 * 包含标准的 balanceOf, transfer, approve, transferFrom 功能.
 * 以及一个 owner 专属的 mint 功能.
 */
contract SimpleERC20 {
    // --- 状态变量 ---

    // 代币名称
    string public name;
    // 代币符号
    string public symbol;
    // 小数位数
    uint8 public decimals = 18;
    // 代币总供应量
    uint256 public totalSupply;

    // 存储每个账户的余额
    // mapping(地址 => 金额)
    mapping(address => uint256) public balanceOf;

    // 存储授权信息
    // mapping(所有者地址 => mapping(被授权地址 => 授权金额))
    mapping(address => mapping(address => uint256)) public allowance;

    // 合约所有者地址
    address public owner;

    // --- 事件 ---

    // 转账事件，记录代币的来源、去向和金额
    event Transfer(address indexed from, address indexed to, uint256 value);

    // 授权事件，记录授权人、被授权人和授权金额
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // --- 构造函数 ---

    /**
     * @dev 初始化合约，设置代币名称和符号.
     * 将合约部署者设置成所有者.
     */
    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
        owner = msg.sender;
    }

    // --- 核心功能 ---

    /**
     * @dev 转账功能.
     * 从调用者地址向 `_to` 地址转移 `_value` 数量的代币.
     * @param _to 接收方地址.
     * @param _value 转移的代币数量.
     * @return bool 是否成功.
     */
    function transfer(address _to, uint256 _value) public returns (bool) {
        require(_to != address(0), "ERC20: transfer to the zero address");
        require(balanceOf[msg.sender] >= _value, "ERC20: transfer amount exceeds balance");

        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;

        emit Transfer(msg.sender, _to, _value);
        return true;
    }

    /**
     * @dev 授权功能.
     * 授权 `_spender` 地址可以从调用者地址中提取最多 `_value` 数量的代币.
     * @param _spender 被授权的地址.
     * @param _value 授权的金额.
     * @return bool 是否成功.
     */
    function approve(address _spender, uint256 _value) public returns (bool) {
        require(_spender != address(0), "ERC20: approve to the zero address");

        allowance[msg.sender][_spender] = _value;

        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    /**
     * @dev 代扣转账功能.
     * 从 `_from` 地址向 `_to` 地址转移 `_value` 数量的代币.
     * 此操作需要 `_from` 地址预先通过 `approve` 函数授权.
     * @param _from 代币转出地址.
     * @param _to 代币接收地址.
     * @param _value 转移金额.
     * @return bool 是否成功.
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
        require(_from != address(0), "ERC20: transfer from the zero address");
        require(_to != address(0), "ERC20: transfer to the zero address");
        require(balanceOf[_from] >= _value, "ERC20: transfer amount exceeds balance");
        require(allowance[_from][msg.sender] >= _value, "ERC20: transfer amount exceeds allowance");

        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        allowance[_from][msg.sender] -= _value;

        emit Transfer(_from, _to, _value);
        return true;
    }

    // --- 管理员功能 ---

    /**
     * @dev 增发代币功能，仅限合约所有者调用.
     * @param _to 接收增发代币的地址.
     * @param _value 增发的数量.
     */
    function mint(address _to, uint256 _value) public {
        require(msg.sender == owner, "Only owner can mint tokens");
        require(_to != address(0), "ERC20: mint to the zero address");

        totalSupply += _value;
        balanceOf[_to] += _value;

        emit Transfer(address(0), _to, _value); // address(0) 代表增发
    }
}
