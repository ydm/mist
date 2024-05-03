// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

interface Lispiface {
    function something() external payable returns (uint256);
    function pause() external payable returns (uint256);
    function please(string memory x) external payable;
}
