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
        data: "0x604e8061000d6000396000f3fe60806040523461001057600061001b565b602060405160018152fd5b506000610029576000610034565b6020604051607b8152f35b506020604051600061004757600361004a565b60025b8152f3",
    });
    const address: string = ethers.getCreateAddress(create);
    console.log("addr:", address);
    // console.log(await ethers.provider.getCode(address));

    const contract: Lispiface = await ethers.getContractAt("Lispiface", address)
    // const res = await contract.something({
    //     value: ethers.parseEther("1"),
    // });
    const tx = await contract.something.populateTransaction({
        value: 0, // ethers.parseEther("1"),
    });
    const result: string = await hero.call(tx)
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
