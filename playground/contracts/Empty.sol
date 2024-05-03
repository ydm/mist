// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Empty {
    uint256 one;

    function something() external payable returns (uint256) {
        revert("asd");
        return 69;
    }

    constructor() payable {}

    /*
    fallback() external payable {
    }

    receive() external payable {
    }
    */
}
