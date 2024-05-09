// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.24;

import {IERC20Metadata} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";

interface IToken is IERC20Metadata {
    function mint() external;
}
