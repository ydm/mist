// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Log {

    event Approval(address owner, address spender, uint256 value);
    event Transfer(address indexed from, address indexed to, uint256 value);

    function something() external payable {
        emit Approval(address(0x10), address(0x20), f(0x1234));
    }

    function f(uint256 x) internal pure returns(uint256) {
        return x + x;
    }
}
