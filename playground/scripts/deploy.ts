import hre, { network } from "hardhat";
import type ethers from "ethers";
import { Signer, TransactionResponse, TransactionReceipt, ContractTransactionResponse, ContractTransactionReceipt } from "ethers";
// ContractTransaction
import { Empty, IToken, Lispiface } from "../typechain-types";

const TWO: string = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8";
const THREE: string = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045";

async function impersonate<T>(
    who: string,
    f: (signer: ethers.Signer) => Promise<T>
): Promise<T | undefined> {
    if (hre.network.name === "hardhat") {
        const norm: string = hre.ethers.getAddress(who);
        await network.provider.request({
            method: "hardhat_impersonateAccount",
            params: [norm],
        });
        const signer = await hre.ethers.getSigner(who);
        try {
            return await f(signer);
        } finally {
            network.provider.request({
                method: "hardhat_stopImpersonatingAccount",
                params: [norm],
            });
        }
    }
}

async function printLogs(tx: ContractTransactionResponse, iface: ethers.Interface) {
    const receipt: ContractTransactionReceipt | null = await tx.wait();
    if (receipt != null) {
        receipt.logs.forEach((log: ethers.Log): void => {
            const topics = log.topics.slice();
            const data = log.data;
            const desc: ethers.LogDescription | null = iface.parseLog({
                topics,
                data,
            });
            if (desc != null) {
                console.log(desc.name, desc.args.toString());
            }
        });
    }
}

async function main() {
    const hero: Signer = await hre.ethers.getSigners().then(
        (signers: Signer[]): Signer => signers[0]
    );
    const heroAddress: string = await hero.getAddress();
    console.log("HERO:", heroAddress, hre.ethers.formatEther(await hre.ethers.provider.getBalance(heroAddress)));

    const create: TransactionResponse = await hero.sendTransaction({
        // type: 2,
        data: "0x6109a18061000d6000396000f3608060405260003560e01c806306fdde03141561011257602060405161010a5b6100cf5b3461002f5760006100cc565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001781602401527f66756e6374696f6e206973206e6f742070617961626c650000000000000000008160440152606490fd5b90565b50606060405160208152601181602001527f4c75636b7920436861726d20546f6b656e0000000000000000000000000000008160400152f390565b8152f361099e565b806395d89b41141561017257602060405161016a5b61012f610023565b50606060405160208152600581602001527f434841524d0000000000000000000000000000000000000000000000000000008160400152f390565b8152f361099e565b8063313ce567141561019d5760206040516101955b61018f610023565b50600090565b8152f361099e565b806318160ddd14156101c05760206040516101b85b60025490565b8152f361099e565b806370a0823114156102b45760206040516102ac6004355b6101e0610023565b50610297815b8060a01c6101f5576000610292565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000f81602401527f696e76616c6964206164647265737300000000000000000000000000000000008160440152606490fd5b905090565b50600060145280600052604060002054905090565b8152f361099e565b8063a9059cbb14156106505760206040516106486024356004355b6102d7610023565b506102e1816101e6565b5061063f8282335b8061038f576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000e81602401527f696e76616c69642073656e6465720000000000000000000000000000000000008160440152606490fd610392565b60005b5081610439576040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001081602401527f696e76616c6964207265636569766572000000000000000000000000000000008160440152606490fd61043c565b60005b506106388383835b8061050357613dbb836002540180600255116104615760006104fe565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000000c81602401527f65786365656465642063617000000000000000000000000000000000000000008160440152606490fd5b6105d2565b6105d161050f826101d8565b5b84811061052f57848103806000601452846000526040600020556105cc565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001481602401527f696e73756666696369656e742062616c616e63650000000000000000000000008160440152606490fd5b905090565b5b50816105e657826002540380600255610602565b826105f0836101d8565b01806000601452836000526040600020555b5081817fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6020604051878152a360009250505090565b9250505090565b50600191505090565b8152f361099e565b8063dd62ed3e14156106a85760206040516106a06024356004355b610674816101e6565b5061067e826101e6565b5060016014528060005260406000206014528160005260406000205491505090565b8152f361099e565b8063095ea7b3141561074757602060405161073f6024356004355b6106cb610023565b5061073660018383335b828060016014528260005260406000206014528360005260406000205550836106ff57600061072e565b81817f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9256020604051878152a360005b935050505090565b50600191505090565b8152f361099e565b806323b872dd14156108a85760206040516108a06044356024356004355b61076d610023565b50610777816101e6565b50610781826101e6565b5061088a8333835b610883610796838361066b565b5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811461087b578481106107d9576107d4600086830386866106d5565b610876565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001681602401527f696e73756666696369656e7420616c6c6f77616e6365000000000000000000008160440152606490fd5b61087e565b60005b905090565b9250505090565b506108968383836102e9565b5060019250505090565b8152f361099e565b8063355274ea14156108cb5760206040516108c35b613dbb90565b8152f361099e565b80631249c58b14156109015760206040516108f95b6108e8610023565b506108f66001336000610444565b90565b8152f361099e565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081527f000000000000000000000000000000000000000000000000000000000000002081600401527f000000000000000000000000000000000000000000000000000000000000001581602401527f756e7265636f676e697a65642066756e6374696f6e00000000000000000000008160440152606490fd5b9050",
    });
    const contractAddress: string = hre.ethers.getCreateAddress(create);
    console.log("ADDR:", contractAddress);

    const charm: IToken = await hre.ethers.getContractAt("IToken", contractAddress);
    // const empty: Empty = await (await hre.ethers.getContractFactory("Empty")).deploy()

    console.log("NAME:", await charm.name());
    console.log("SYMBOL:", await charm.symbol());
    console.log("DECIMALS:", await charm.decimals());
    console.log("CAP:", await charm.cap());

    console.log("[1] BALANCE 0:", await charm.balanceOf(hero), await charm.totalSupply());
    printLogs(await charm.mint(), charm.interface);
    console.log("[1] BALANCE 1:", await charm.balanceOf(hero), await charm.totalSupply());
    printLogs(await charm.mint(), charm.interface);
    console.log("[1] BALANCE 2:", await charm.balanceOf(hero), await charm.totalSupply());
    printLogs(await charm.mint(), charm.interface);
    console.log("[1] BALANCE 3:", await charm.balanceOf(hero), await charm.totalSupply());
    printLogs(await charm.mint(), charm.interface);
    console.log("[1] BALANCE 4:", await charm.balanceOf(hero), await charm.totalSupply());
    // await charm.transfer(hre.ethers.ZeroAddress, 1);
    console.log("[1] BALANCE 5:", await charm.balanceOf(hero), await charm.totalSupply());

    console.log("[2] BALANCE 0:", await charm.balanceOf(TWO), await charm.totalSupply());
    printLogs(await charm.connect(hero).transfer(TWO, 2), charm.interface);
    console.log("[2] BALANCE 1:", await charm.balanceOf(TWO), await charm.totalSupply());

    console.log("[1] BALANCE 6:", await charm.balanceOf(hero), await charm.totalSupply());

    await impersonate(THREE, async (three: Signer): Promise<void> => {
        console.log("[3] BALANCE 0:", await charm.balanceOf(three), await charm.totalSupply());
        printLogs(await charm.connect(three).mint(), charm.interface);
        printLogs(await charm.connect(three).mint(), charm.interface);
        printLogs(await charm.connect(three).mint(), charm.interface);
        console.log("[3] BALANCE 1:", await charm.balanceOf(three), await charm.totalSupply());

        console.log("[3] ALLOWANCE 0:", await charm.allowance(three, hero));
        const tx: ContractTransactionResponse = await charm.connect(three).approve(hero, 3);
        printLogs(tx, charm.interface);

        // await charm.connect(three).approve(hero, hre.ethers.MaxUint256);
        console.log("[3] ALLOWANCE 1:", await charm.allowance(three, hero));
    });

    console.log("[1] BALANCE 7:", await charm.balanceOf(hero), await charm.totalSupply());
    console.log("[2] BALANCE 2:", await charm.balanceOf(TWO), await charm.totalSupply());
    console.log("[3] BALANCE 2:", await charm.balanceOf(THREE), await charm.totalSupply());

    console.log("[3] ALLOWANCE 2:", await charm.allowance(THREE, hero));
    printLogs(await charm.connect(hero).transferFrom(THREE, TWO, 1), charm.interface);
    console.log("[3] ALLOWANCE 3:", await charm.allowance(THREE, hero));

    console.log("[1] BALANCE 8:", await charm.balanceOf(hero), await charm.totalSupply());
    console.log("[2] BALANCE 3:", await charm.balanceOf(TWO), await charm.totalSupply());
    console.log("[3] BALANCE 3:", await charm.balanceOf(THREE), await charm.totalSupply());

    // await charm.connect(hero).transferFrom(THREE, TWO, 1);
    // await charm.connect(hero).transferFrom(THREE, TWO, 1);
    // // await charm.connect(hero).transferFrom(THREE, TWO, 1);
    // console.log("[3] ALLOWANCE 2:", await charm.allowance(THREE, hero));

    // const before: bigint = await hre.ethers.provider.getBalance(hero);
    // const interact: TransactionResponse = await hero.sendTransaction({
    //     to: contractAddress,
    //     // data: "0x1249c58b",
    //     data: "0x313ce567",
    //     value: hre.ethers.parseEther("1"),
    // });

    // const receipt: TransactionReceipt | null = await interact.wait();
    // if (receipt != null) {
    //     console.log("RESULT:", await receipt.getResult());
    // }

    // const after: bigint = await hre.ethers.provider.getBalance(hero);
    // console.log("DIFFERENCE:", hre.ethers.formatEther(after-before));
    // const result: string = await hero.call(interact);
    // console.log("RESULT:", result, "(", parseInt(result, 16), ")");

    // const tx = await charm.name.populateTransaction();
    // const result: string = await hero.call(tx);
    // console.log("RESULT:", result, "(", parseInt(result, 16), ")");
}

main().then(
    () => process.exit(0),
    (reason: any) => {
        console.error(reason);
        process.exit(1);
    }
);
