import { ethers } from "hardhat";
import { Signer } from "ethers";
import { Empty__factory, Empty } from "../typechain-types";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    console.log("hero:", await hero.getAddress());

    const IToken: IToken = await ethers.getContractAt("");

    const Empty: Empty__factory = await ethers.getContractFactory("Empty");
    const empty: Empty = await Empty.deploy();

    const address: string = await empty.getAddress();
    console.log("address:", address);

    const tx = await empty.something.populateTransaction({
        value: 0,
    });
    const result: string = await hero.call(tx);
    console.log(result);

}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);
