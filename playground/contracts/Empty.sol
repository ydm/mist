// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Empty {
    uint256 one;

    function something() external payable {
        one = 169;
    }

    constructor() payable {}

    /*
    fallback() external payable {
    }

    receive() external payable {
    }
    */
}
