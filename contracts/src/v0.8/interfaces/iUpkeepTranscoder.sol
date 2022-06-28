// SPDX-License-Identifier: MIT

import "../UpkeepFormat.sol";

pragma solidity ^0.8.0;

interface iUpkeepTranscoder {
  function transcodeUpkeeps(
    UpkeepFormat fromVersion,
    UpkeepFormat toVersion,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}