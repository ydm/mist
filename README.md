# Mist

[sad]: ## "that's a joke"
The [much-needed][sad] Lisp language for the Ethereum EVM.

### Running Mist

```
$ make
$ ./mist < source.mist
```

### Quickstart

Want to waste some resources? Say no more:

```Lisp
(defun fib (i)
  (if (< i 2)
      i
    (+ (fib (- i 2))
       (fib (- i 1)))))

(return (fib 0xBAD))
```

Public contracts most likely would include a `dispatcher`:

```Lisp
(defvar *balances*    (mapping address uint256))
(defvar *allowances*  (mapping address (mapping address uint256)))
(defvar *totalSupply* uint256)

(defun totalSupply () *totalSupply*)
(defun balanceOf (address) (gethash *balances* address))
(defun allowance (owner spender) (gethash *allowances* owner spender))

(defun transfer (to value) ...)
(defun transferFrom (from to value) ...)

(defun approve (spender value)
  (puthash *allowances* value (caller) spender)
  (emit3 "Approval(address,address,uint256)" owner spender value)
  t)

(dispatch
 ("totalSupply()"                          totalSupply)
 ("balanceOf(address)"                       balanceOf)
 ("allowance(address,address)"               allowance)
 ("transfer(address,uint256)"                 transfer)
 ("transferFrom(address,address,uint256)" transferFrom)
 ("approve(address,uint256)"                   approve))
```

For an example implementation of an ERC-20 token, please see [charm.mist](examples/charm.mist).

# Standard library

#### `nil` (`false`) values:
  - the number `0`
  - the symbol `nil`
  - any empty list, quoted or not

#### `t` (`true`) values:
  - any number other than `0`
  - the symbol `t`
  - any non-empty list, quoted or not

#### Functions mapped to native instructions
  - `(stop)`
  - `(-)` e.g. `(- 20 5)`
  - `(/)` e.g. `(/ 10 2)`
  - `(%)` e.g. `(% 10 3)`
  - `(+%)` e.g. `(+% a b c)`, shorthand for `(% (+ a b) c)`
  - `(*%)` e.g. `(*% a b c)`, shorthand for `(% (* a b) c)`
  - `(**)` e.g. `(** x y)`, x to the power of y
  - `(expt)`, an alias for `(**)`
  - `(<)` e.g. `(< a b)` results in `t` if `a<b`, `nil` otherwise
  - `(>)`
  - `(=)` e.g. `(= a b)` results in `t` if `a==b`, `nil` otherwise
  - `(not)` e.g. `(not t)` results in `nil` and `(not nil)` results in `t`
  - `(zerop)`, an alias for `(not)`
  - `(~)`, e.g. `(~ a)` returns the bitwise complement of `a`
  - `(lognot)`, an alias for `(~)`
  - `(byte)`, e.g. `(byte index word)`, see opcode `BYTE`
  - `(<<)`, e.g. `(<< a b)` returns integer `a` with its bits shifted by `b` bit positions, `(<< 1 4)` results in `16`
  - `(>>)`
  - `(current-address)`, see opcode `ADDRESS`
  - `(balance)`, see opcode `BALANCE`
  - `(origin)`, see opcode `ORIGIN`
  - `(caller)`, a.k.a. `msg.sender`
  - `(call-value)`, value sent with the transaction
  - `(calldata-load)`, low-level access to input transaction data
  - `(calldata-size)`, low-level access to the size of the input data
  - `(code-size)`, low-level access to the size of the currently running code
  - `(gas-price)`, a.k.a. effective gas price
  - `(coinbase)`, current block's beneficiary address
  - `(timestamp)`, current block's timestamp
  - `(block-number)`, current block's number
  - `(prev-randao)`, latest RANDAO mix of the post beacon state of the previous block
  - `(gas-limit)`, current block's gas limit
  - `(chain-id)`
  - `(self-balance)`, balance of currently executing account
  - `(base-fee)`, current block's base fee
  - `(available-gas)`, amount of available gas (after paying for this instruction)

#### Variadic:
  - `(+)`, e.g. `(+ 1 2 3 4 5)`
  - `(*)`
  - `(&)` and its alias `(logand)`
  - `(|)` and its alias `(logior)`
  - `(&)` and its alias `(logxor)`

#### Builtins:
  - `(case)`, standard Lisp `(case)`, see `examples/case*.mist` for examples
  - `(defconst)`, give a name to a constant expression
  - `(defun)`, e.g. `(defun NAME ARGLIST BODY...)`, define NAME as function
  - `(defvar)`, e.g. `(defvar totalSupply uint256)`, create a *storage* variable
  - `(emit3)`, e.g. `(emit3 "Transfer(address,address,uint256)" from to value)`, emit a Log with 3 topics
  - `(ether)`, e.g. `(ether "1")` results in `1e18`
  - `(gethash TABLE KEYS...)`, access values in a mapping, e.g. `(gethash balances owner)` or `(gethash allowances owner spender)`
  - `(if COND A B)` results in `A` if `COND` holds and `B` otherwise
  - `(progn BODY...)` executes all BODY expressions in a sequence and yields the result of the last one
  - `(puthash TABLE VALUE KEYS...)`, analogous to `(gethash)`, e.g. `(puthash balances value owner)` or `(puthash allowances value owner spender)`
  - `(return VALUE-OR-STRING)`
  - `(revert VALUE-OR-STRING)`
  - `(selector STRING)`
  - `(setq SYMBOL VALUE)` assigns `VALUE` to the *storage* variable named `SYMBOL`

#### Macros:
  - `(<=)`
  - `(>=)`
  - `(apply FUNCTION ARGUMENTS...)`
  - `(dispatch)`, see `examples/token.mist`
  - `(let VARLIST BODY...)`
  - `(unless COND BODY...)` if `COND` yields `nil`, do `BODY`, else return nil
  - `(when COND BODY...)` if `COND` yields `t`, do `BODY`, else return nil

#### Notes:

These are the only functions do not result in an expression:
  - `(stop)`
  - `(return x)`
  - `(revert x)`

### Limitations:
  - Code length can't exceed 2^16 bytes (64 kilobytes).

### TODO:
  - Implement real macros.
  - Proper error reporting, no `panic()`s.
  - `Segment` should be an interface instead of a stateful mess.
  - Solidity offers a rich assortment of opcode optimizations; perhaps reuse?
  - Clean up the public/private mess.
  - Proper documentation, more examples, and more tests.
  - Deal with linters.
