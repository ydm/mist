// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

contract Array {

    mapping(uint256 => mapping(uint256 => uint256)) m;

    constructor() payable {}

    /*
    function something() external payable returns (uint256) {
        uint256[2] memory xs;
        xs[0] = 0x1000;
        xs[1] = 0x2000;
        return xs[0] + xs[1];
    }
    */

    function somethingElse(uint256 a, uint256 b) external payable returns (uint256) {
        return m[a][b];
    }

    fallback() external payable {}

    receive() external payable {}
}
