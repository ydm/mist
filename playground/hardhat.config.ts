import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

function getenv(name: string): string {
    return "" + process.env[name]
}

const MAINNET: string = getenv("MAINNET")
const HOLESKY: string = "https://ethereum-holesky-rpc.publicnode.com";

function accounts(): string[] {
    const PRIVATE_KEY = getenv("PRIVATE_KEY");

    if (PRIVATE_KEY && !/[0-9a-fA-F]{64}/.test(PRIVATE_KEY)) {
        throw new Error();
    }

    const ans: string[] = !!PRIVATE_KEY ? [PRIVATE_KEY] : [];
    return ans;
}

const config: HardhatUserConfig = {
    networks: {
        mainnet: {
            chainId: 1,
            url: MAINNET,
            accounts: accounts(),
        },
        holesky: {
            chainId: 17000,
            url: HOLESKY,
            accounts: accounts(),
        },
        hardhat: {
            forking: {
                enabled: true,
                url: HOLESKY,
            },
        },
        node: {
            url: "http://127.0.0.1:8545",
        },
    },
    solidity: {
        version: "0.8.24",
        settings: {
            optimizer: {
                enabled: true,
                runs: 5000,
            },
        },
    },
};

export default config;
