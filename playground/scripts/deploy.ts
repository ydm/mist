import { ethers } from "hardhat";
import { ContractTransactionResponse, Signer, TransactionResponse } from "ethers";
import { Empty, Empty__factory } from "../typechain-types";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    
    // const Empty: Empty__factory = await ethers.getContractFactory("Empty");
    // const empty: Empty = await Empty.deploy();

    // const tx: null | ContractTransactionResponse = empty.deploymentTransaction()
    // if (tx == null) {
    //     // return;
    //     console.log("deploymentTransaction is null");
    // } else {
    //     console.log(tx);
    //     console.log("data:", tx.data);
    // }

    // const address: string = await empty.getAddress();
    // console.log("address:", address);
    // console.log("deployed:", await empty.getDeployedCode());

    // const query: TransactionResponse = await hero.sendTransaction({
    //     type: 0,
    // });

    // console.log(await hero.getAddress());
    // const create: TransactionResponse = await hero.sendTransaction({
    //     type: 2,
    //     data: "0x608060405260698060116000396000f3fe608060405260043610601c5760003560e01c8063a7a0d537146021575b600080fd5b60276029565b005b60a960008190555056fea264697066735822122027c0b72c8dd60a4929d903c99a82e73240fc9aac6147c1c13e0040829ce4285364736f6c63430008180033",
    // });
    // console.log(create);

    // const address: string = ethers.getCreateAddress(create);
    // console.log(address);

    // const address: string = "0xEDD45Aa00c15Eec3484d4a1F0C9895c5AdD95DC7";
    const address: string = "0xc37ed33276b2EDBd3D39bf721B7Dc4Dc3806Aa6C";
    console.log(await ethers.provider.getCode(address));
    console.log(
        await ethers.provider.getCode(address).
            then((code: string): number => code.length)
    );
}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);

/*
 * Constructor
 *

+----+----------+------------+-------------+----------------------------+
| pc | bytecode | assembly   | stack       | memory                     |
|----+----------+------------+-------------+----------------------------|
|  0 |     6080 | PUSH1 0x80 | 80          |                            |
|  2 |     6040 | PUSH1 0x40 | 40 80       |                            |
|  4 |       52 | MSTORE     |             | m[0x40:0x60] = 0x80        |
|  5 |       34 | CALLVALUE  | V           |                            |
|  6 |       80 | DUP1       | V V         |                            |
|  7 |       15 | ISZERO     | 1 V         |                            |
|  8 |     600f | PUSH1 0x0f | 0f 1 V      |                            |
|  a |       57 | JUMPI      | V           |                            |
|  b |     6000 | PUSH1 0x00 | 0 V         |                            |
|  d |       80 | DUP1       | 0 0 V       |                            |
|  e |       fd | REVERT     | V           |                            |
|  f |       5b | JUMPDEST   | V           |                            |
| 10 |       50 | POP        |             |                            |
| 11 |     6059 | PUSH1 0x60 | 59          |                            |
| 13 |       80 | DUP1       | 59 59       |                            |
| 14 |     6022 | PUSH1 0x22 | 22 59 59    |                            |
| 16 |     6000 | PUSH1 0x00 | 00 22 59 59 |                            |
| 18 |       39 | CODECOPY   | 59          | m[0:0x59] = Ib[0x22:+0x59] |
| 19 |     6069 | PUSH1 0x69 | 69 59       |                            |
| 1b |     6060 | PUSH1 0x60 | 60 69 59    |                            |
| 1d |       53 | MSTORE8    | 59          | m[0x60] = 0x69             |
| 1e |     6000 | PUSH1 0x00 | 00 59       |                            |
| 20 |       f3 | RETURN     |             |                            |
| 21 |       fe | INVALID    |             |                            |
+----+----------+------------+-------------+----------------------------+

*/
