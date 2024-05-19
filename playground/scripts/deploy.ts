import hre from "hardhat";
import type ethers from "ethers";
import { Signer, TransactionResponse } from "ethers";
// ContractTransaction
import { Lispiface } from "../typechain-types";

const TWO: string = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";

async function main() {
    const hero: Signer = await hre.ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    const heroAddress: string = await hero.getAddress();
    console.log("HERO:", heroAddress, hre.ethers.formatEther(await hre.ethers.provider.getBalance(heroAddress)));

    const createRequest: ethers.TransactionRequest = {
    };
    const create: TransactionResponse = await hero.sendTransaction({
        // type: 2,
        data: "0x600f8061000c6000396000f36080604052602060405160408152f3",
    });
    const contractAddress: string = hre.ethers.getCreateAddress(create);
    console.log("ADDR:", contractAddress);

    const contract: Lispiface = await hre.ethers.getContractAt("Lispiface", contractAddress);
    const tx = await contract.something.populateTransaction();
    const result: string = await hero.call(tx);
    console.log("RESULT:", result);

    // const contract: ethers.Contract = new hre.ethers.Contract(contractAddress, "[]");

    // const call: TransactionResponse = await hero.sendTransaction({
        
    // });
    // console.log(contract)
}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);
