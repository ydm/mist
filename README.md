nil (false) values:
  - the symbol nil
  - any empty list, quoted or not
  - the number 0

t (true) values:
  - the symbol t
  - any non-empty list, quoted or not
  - any number other than 0

Mapped to instructions
  - (stop)
  - (-)
  - (/)
  - (%)
  - (+%)
  - (*%)
  - (**)
  - (expt)
  - (<)
  - (>)
  - (=)
  - (not)
  - (zerop)
  - (~)
  - (lognot)
  - (byte)
  - (<<)
  - (>>)
  - (address)
  - (balance)
  - (origin)
  - (caller)
  - (call-value)
  - (calldata-load)
  - (calldata-size)
  - (code-size)
  - (gas-price)
  - (coinbase)
  - (timestamp)
  - (block-number)
  - (prev-randao)
  - (gas-limit)
  - (chain-id)
  - (self-balance)
  - (base-fee)
  - (available-gas)

Variadic:
  - (+ )
  - (* )
  - (&) (logand)
  - (|) (logior)
  - (&) (logxor)

Builtins:
  - (defconst)
  - (defun)
  - (if)
  - (progn)
  - (return)
  - (revert)
  - (unless)
  - (when)

These functions do not result in an expression:
  - `(stop)`
  - `(return x)`
  - `(revert x)`

Limitations:
  - Code length can never exceed 2^16 bytes (64 megabytes).
