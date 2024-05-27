// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Empty {

    event Transfer(address indexed from, address indexed to, uint256 value);

    uint256 one;

    function something() public pure returns (uint256) {
        return 69;
    }

    function name() public payable {
        emit Transfer(msg.sender, msg.sender, 123);
    }
}
