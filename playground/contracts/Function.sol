// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Function {

    constructor() payable {}

    function something() external payable returns (uint256) {
        return f(1, 2, 3);
    }

    function f(uint256 a, uint256 b, uint256 c) private pure returns (uint256) {
        return a + b + c;
    }

    fallback() external payable {}

    receive() external payable {}
}
