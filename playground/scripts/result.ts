import { ethers } from "hardhat";
import { Signer, TransactionResponse } from "ethers";
import { Lispiface } from "../typechain-types";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    const heroAddr: string = await hero.getAddress();
    console.log("HERO:", heroAddr, ethers.formatEther(await ethers.provider.getBalance(heroAddr)));

    const create: TransactionResponse = await hero.sendTransaction({
        type: 2,
        data: "0x60198061000c6000396000f3608060405261032161012360206040518383018152f3915050",
    });
    const address: string = ethers.getCreateAddress(create);
    console.log("ADDR:", address);

    const contract: Lispiface = await ethers.getContractAt("Lispiface", address)
    const tx = await contract.something.populateTransaction({
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
