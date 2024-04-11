// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Array {

    constructor() payable {}

    function something() external payable returns (uint256) {
        uint256[2] memory xs;
        xs[0] = 0x1000;
        xs[1] = 0x2000;
        return xs[0] + xs[1];
    }

    fallback() external payable {}

    receive() external payable {}
}
