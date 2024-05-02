nil (false) values:
  - the symbol nil
  - any empty list, quoted or not
  - the number 0
  
t (true) values:
  - the symbol t
  - any non-empty list, quoted or not
  - any number other than 0

Those functions do not result in an expression:
  - `(stop)`
  - `(return x)`
  - `(revert x)`

Limitations:
  - Code length can never exceed 2^16 bytes (64 megabytes).

Functions:

stop

address
origin
caller
call-value
calldata-load
calldata-size
code-size
gas-price
return-data-size
coinbase
timestamp
block-number
prev-randao
gas-limit
chain-id
self-balance
base-fee

* + - /
> < =
logand logior logxor

OPTIMIZER: (1) Delete PUSH[1-16] -> BYTES -> POP
