// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0 <0.9.0;

import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.2/contracts/token/ERC721/ERC721.sol";
import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.2/contracts/access/Ownable.sol";
import "https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.2/contracts/utils/ReentrancyGuard.sol";

/**
 * @title MyNFT
 * @dev 一个符合 ERC721 标准的 NFT 合约.
 * 继承自 OpenZeppelin 的 ERC721, Ownable 和 ReentrancyGuard.
 */

contract MyNFT is ERC721, Ownable, ReentrancyGuard {
    // tokenId 计数器
    uint256 private _nextTokenId;

    /**
     * @dev 构造函数，初始化 NFT 名称和符号.
     * 调用父合约 ERC721 和 Ownable 的构造函数.
     * @param initialOwner 合约的初始所有者地址.
     */

    constructor(address initialOwner)
    ERC721("My Awesome NFT", "MANFT")
    Ownable(initialOwner)
    {}

    /**
     * @dev 铸造一个新的 NFT.
     * 只有合约所有者可以调用.
     * @param recipient NFT 的接收者地址.
     * @param tokenURI NFT 的元数据链接 (例如 "ipfs://...").
     */

    function mintNFT(address recipient, string memory tokenURI)
    public
    onlyOwner
    nonReentrant
    {
        uint256 tokenId = _nextTokenId;
        _nextTokenId++;
        _safeMint(recipient, tokenId);
        _setTokenURI(tokenId, tokenURI);
    }
}