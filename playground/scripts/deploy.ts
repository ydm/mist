import { ethers } from "hardhat";
import { Signer, TransactionResponse } from "ethers";
import { Lispiface } from "../typechain-types";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    const heroAddr: string = await hero.getAddress();
    console.log("HERO:", heroAddr, ethers.formatEther(await ethers.provider.getBalance(heroAddr)));

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

    const create: TransactionResponse = await hero.sendTransaction({
        type: 2,
        data: "0x61011c8061000d6000396000f3608060405260003560e01c806318160ddd141561001d576000610119565b806370a082311415610030576000610119565b8063a9059cbb1415610043576000610119565b8063dd62ed3e1415610056576000610119565b8063095ea7b31415610069576000610119565b806323b872dd141561007c576000610119565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001581602401527f756e7265636f676e697a65642066756e6374696f6e00000000000000000000008160440152606490fd5b9050",
    });
    const address: string = ethers.getCreateAddress(create);
    console.log("ADDR:", address);
    // console.log(await ethers.provider.getCode(address));

    const contract: Lispiface = await ethers.getContractAt("Lispiface", address)
    // const res = await contract.something({
    //     value: ethers.parseEther("1"),
    // });
    const tx = await contract.please.populateTransaction("asd", {
        value: 0, // ethers.parseEther("1"),
    });
    console.log("DATA:", tx.data);
    const result: string = await hero.call(tx);
    console.log("RESULT:", result);

    // const result: bigint = await contract.something.staticCall()
    // console.log("result:", result.toString(16));

    // await hero.sendTransaction({
    //     to: address,
    //     data: "0xaabbccddee00"
    // });

    // console.log(await ethers.provider.getCode(address));
    // console.log(
    //     await ethers.provider.getCode(address)
    // );
}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);
