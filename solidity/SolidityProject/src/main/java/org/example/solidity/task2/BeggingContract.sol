// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * @title BeggingContract
 * @dev 一个允许用户捐赠以太币，并允许所有者提款的合约.
 */
contract BeggingContract {
    // --- 状态变量 ---

    // 合约所有者
    address public owner;

    // 记录每个地址的总捐赠金额
    // mapping(捐赠者地址 => 捐赠总金额)
    mapping(address => uint256) public donations;

    // --- 事件 ---

    // 记录每次捐赠的事件
    event DonationReceived(address indexed from, uint256 amount);

    // --- 修改器 ---

    /**
     * @dev 限制函数只能被合约所有者调用.
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Ownable: caller is not the owner");
        _;
    }

    // --- 构造函数 ---

    /**
     * @dev 设置合约部署者为所有者.
     */
    constructor() {
        owner = msg.sender;
    }

    // --- 核心功能 ---

    /**
     * @dev 接收捐赠的函数.
     * `payable` 关键字允许此函数接收以太币.
     */
    function donate() public payable {
        require(msg.value > 0, "Donation amount must be greater than zero");

        // 累加该地址的捐赠金额
        donations[msg.sender] += msg.value;

        // 触发捐赠事件
        emit DonationReceived(msg.sender, msg.value);
    }

    /**
     * @dev 提款函数，将合约中的所有余额提取到所有者账户.
     * 只有所有者可以调用.
     */
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        // 使用 call 方法转账，更安全
        (bool success, ) = owner.call{value: balance}("");
        require(success, "Transfer failed.");
    }

    /**
     * @dev 查询指定地址的捐赠总额.
     * @param _donor 查询的地址.
     * @return uint256 该地址的总捐赠金额.
     */
    function getDonation(address _donor) public view returns (uint256) {
        return donations[_donor];
    }
}
