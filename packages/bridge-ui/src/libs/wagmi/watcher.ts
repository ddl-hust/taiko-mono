import { watchAccount } from '@wagmi/core';

import { chains, isSupportedChain } from '$libs/chain';
import { refreshUserBalance } from '$libs/util/balance';
import { checkForPausedContracts } from '$libs/util/checkForPausedContracts';
import { isSmartContract } from '$libs/util/isSmartContract';
import { getLogger } from '$libs/util/logger';
import { account, connectedSmartContractWallet } from '$stores/account';
import { switchChainModal } from '$stores/modal';
import { connectedSourceChain } from '$stores/network';

import { config } from './client';

const log = getLogger('wagmi:watcher');

let isWatching = false;
let unWatchAccount: () => void;

export async function startWatching() {
  checkForPausedContracts();

  if (!isWatching) {
    unWatchAccount = watchAccount(config, {
      async onChange(data) {
        await checkForPausedContracts();
        log('Account changed', data);
        account.set(data);

        refreshUserBalance();
        const { chainId, address } = data;

        if (chainId && address) {
          let smartWallet = false;
          try {
            smartWallet = (await isSmartContract(address, Number(chainId))) || false;
          } catch (error) {
            console.error('Error checking for smart contract wallet', error);
          } finally {
            connectedSmartContractWallet.set(smartWallet);
          }
        }

        // We need to check if the chain is supported, and if not
        // we present the user with a modal to switch networks.
        if (chainId && !isSupportedChain(Number(chainId))) {
          log('Unsupported chain', chainId);
          switchChainModal.set(true);
          return;
        } else if (chainId) {
          // When we switch networks, we are actually selecting
          // the source chain.
          const srcChain = chains.find((c) => c.id === Number(chainId));
          if (srcChain) connectedSourceChain.set(srcChain);
          refreshUserBalance();
        }
      },
    });

    isWatching = true;
  }
}

export function stopWatching() {
  unWatchAccount();
  isWatching = false;
}
