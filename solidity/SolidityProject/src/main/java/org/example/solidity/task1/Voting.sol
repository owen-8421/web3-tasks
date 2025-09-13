// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

/**
 * @title Voting
 * @dev 一个简单的智能合约，用于给候选人投票。
 */

contract Voting {
    // 状态变量
    // 使用 mapping 来存储每位候选人（string）获得的票数（uint）。
    mapping(string => uint) public candidateVotes;

    // 为了能够重置所有投票，我们需要一个数组来追踪所有获得过投票的候选人。
    string[] private candidates;

    /**
     * @notice 为指定的候选人投一票。
     * @param _candidateName 候选人的名字。
     * 如果这是该候选人第一次得票，我们会将其记录到 candidates 数组中。
     */

    function vote(string memory _candidateName) public {
        // 在增加票数之前，检查这位候选人是否已经存在于我们的追踪数组中。
        // 如果票数是0，说明是第一次投票给他。
        if (candidateVotes[_candidateName] == 0) {
            candidates.push(_candidateName); // push
        }

        // 为该候选人的票数加 1
        candidateVotes[_candidateName]++;
    }

    /**
     * @notice 获取指定候选人的当前得票总数。
     * @param _candidateName 候选人的名字。
     * @return 候选人的总票数。
     */
    function getVotes(string memory _candidateName) public view returns (uint) {
        return candidateVotes[_candidateName];
    }

    /**
     * @notice 重置所有候选人的得票数。
     */
    function resetVotes() public {
        // 遍历我们追踪的所有候选人
        for (uint i = 0; i < candidates.length; i++) {
            string memory candidateName = candidates[i];
            // 使用 delete 关键字将该候选人的票数重置为 uint 的默认值 (0)
            delete candidateVotes[candidateName];
        }

        // 清空候选人追踪数组
        delete candidates;
    }
}