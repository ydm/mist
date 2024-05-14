import { ethers } from "hardhat";
import { ContractTransaction, Signer, TransactionResponse } from "ethers";
import { IToken } from "../typechain-types";

const TWO: string = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";

async function main() {
    const hero: Signer = await ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    const heroAddr: string = await hero.getAddress();
    console.log("HERO:", heroAddr, ethers.formatEther(await ethers.provider.getBalance(heroAddr)));

    const create: TransactionResponse = await hero.sendTransaction({
        type: 2,
        data: "0x611ba08061000d6000396000f3608060405260003560e01c806318160ddd14156100265760206040516002548152f3611b9d565b806370a0823114156101c8576020604051600435806000811461004a5760006100e7565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff81168082146101a7576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd6101aa565b60005b905090505060006014528060005260406000205490508152f3611b9d565b8063a9059cbb14156104b2576020604051602435600435346101eb576000610288565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b508060008114610299576000610336565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff81168082146103f6576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd6103f9565b60005b90509050503461040a5760006104a7565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b9150508152f3611b9d565b8063dd62ed3e14156107d957602060405160243560043580600081146104d9576000610576565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214610636576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd610639565b60005b9050905050816000811461064e5760006106eb565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff81168082146107ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd6107ae565b60005b90509050506001601452806000526040600020601452816000526040600020549150508152f3611b9d565b8063095ea7b31415610a37576020604051602435600435346107fc576000610899565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b5080600081146108aa576000610947565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214610a07576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd610a0a565b60005b905090505081806001601452336000526040600020601452826000526040600020559150508152f3611b9d565b806323b872dd14156117c557602060405160443560243560043534610a5d576000610afa565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b508060008114610b0b576000610ba8565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214610c68576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd610c6b565b60005b90509050508160008114610c80576000610d1d565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214610ddd576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd610de0565b60005b905090505033818060008114610df7576000610e94565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214610f54576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd610f57565b60005b90509050508160008114610f6c576000611009565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff81168082146110c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd6110cc565b60005b9050905050600160145280600052604060002060145281600052604060002054915050831115611197576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001181602401527f544f444f3a206e6f7420616c6c6f7765640000000000000000000000000000008160440152606490fd61119a565b60005b508080600081146111ac576000611249565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff8116808214611309576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd61130c565b60005b905090505060006014528060005260406000205490508311156113ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001081602401527f544f444f3a206e6f7420656e6f756768000000000000000000000000000000008160440152606490fd6113cd565b60005b508282828180600081146113e257600061147f565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff811680821461153f576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd611542565b60005b90509050506000601452806000526040600020549050818060008114611569576000611606565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f7a65726f206164647265737300000000000000000000000000000000000000008160440152606490fd5b5073ffffffffffffffffffffffffffffffffffffffff81168082146116c6576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd6116c9565b60005b905090505060006014528060005260406000205490508481106116ed57600061178a565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001481602401527f696e73756666696369656e742062616c616e63650000000000000000000000008160440152606490fd5b5084810380600060145284600052604060002055508482018060006014528560005260406000205591505092505050925050508152f3611b9d565b806306fdde03141561187a5760206040516040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001081602401527f456d616373204c69737020546f6b656e000000000000000000000000000000008160440152606490fd8152f3611b9d565b806395d89b41141561192f5760206040516040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000381602401527f454c5400000000000000000000000000000000000000000000000000000000008160440152606490fd8152f3611b9d565b8063313ce567141561194a57602060405160128152f3611b9d565b80631249c58b1415611b0057602060405134611967576000611a04565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b50690358aeae441e1a0c000060025410611ab9576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001b81602401527f6d6178696d756d20746f6b656e206c696d6974207265616368656400000000008160440152606490fd611abc565b60005b50670de0b6b3a7640000600060145233600052604060002054018060006014523360005260406000205550670de0b6b3a764000060025401806002558152f3611b9d565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001581602401527f756e7265636f676e697a65642066756e6374696f6e00000000000000000000008160440152606490fd5b9050",
    });
    const address: string = ethers.getCreateAddress(create);
    console.log("ADDR:", address);

    // const address: string = "0x32BD158B9fbfC21FbB4E5d67A3aBc475b4343c59";
    const token: IToken = await ethers.getContractAt("IToken", address)

    // console.log(await token.connect(hero).mint());
    // const tx: ContractTransaction = await token.mint.populateTransaction({ value: 1 });
    // await hero.sendTransaction(tx);

    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());
    await token.connect(hero).mint();
    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());
    await token.connect(hero).mint();
    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());
    await token.connect(hero).mint();
    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());
    await token.connect(hero).mint();
    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());
    await token.connect(hero).transfer(TWO, ethers.parseEther("1"));
    console.log(await token.balanceOf(heroAddr), await token.balanceOf(TWO), await token.totalSupply());

    // await token.connect(hero).transfer("0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97", );

    // const SPENDER: string = "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97";
    // console.log("allowance before:", await token.allowance(heroAddr, SPENDER));
    // await token.connect(hero).approve(SPENDER, 0x100);
    // console.log("allowance after:", await token.allowance(heroAddr, SPENDER));

    

    // console.log("BEFORE:", await token.totalSupply());

    // console.log("DATA:", tx.data);
    // const result: TransactionResponse =
  
    // // const result = await token.something.send();
    // // const result: string = await hero.call(tx);
    // // console.log("RESULT:", result);
    // // console.log("RESULT.DATA:", result.data);

    // console.log("AFTER:", await token.totalSupply());

    // console.log(await ethers.provider.getStorage(address, 0));

    // console.log("BALANCE:", await token.balanceOf(heroAddr));
    // console.log("BALANCE:", await token.balanceOf.populateTransaction(heroAddr));
}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);
