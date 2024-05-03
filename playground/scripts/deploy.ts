import { ethers } from "hardhat";
import { Signer, TransactionResponse } from "ethers";
import { Lispiface } from "../typechain-types";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    console.log("hero:", await hero.getAddress());

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
        data: "0x60fa8061000c6000396000f360806040527f08c379a0000000000000000000000000000000000000000000000000000000006080526020608452600360a4527f617364000000000000000000000000000000000000000000000000000000000060c45260646080fd5063a7a0d53760003560e01c1461007357600061007e565b602060405160698152f35b50638456cb5960003560e01c146100965760006100a1565b602060405160798152f35b507f08c379a0000000000000000000000000000000000000000000000000000000006080526020608452600360a4527f617364000000000000000000000000000000000000000000000000000000000060c45260646080fd",
    });
    const address: string = ethers.getCreateAddress(create);
    console.log("addr:", address);
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
    console.log(result);

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
