import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const { MAINNET }: { MAINNET: string } = process.env;
const HOLESKY: string = "https://ethereum-holesky-rpc.publicnode.com";

function accounts(): string[] {
    const { PRIVATE_KEY }: { PRIVATE_KEY: string } = process.env;

    if (PRIVATE_KEY && !/[0-9a-fA-F]{64}/.test(PRIVATE_KEY)) {
        throw new Error();
    }

    const ans: string[] = !!PRIVATE_KEY ? [PRIVATE_KEY] : [];
    return ans;
}

const config: HardhatUserConfig = {
    networks: {
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
