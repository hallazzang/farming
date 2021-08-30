#!/bin/sh

set -o xtrace

MNEMONIC="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
GENESIS_COINS=10000000000stake,10000000000uatom,10000000000uiris

killall farmingd

rm -r localnet

mkdir -p localnet

farmingd init test --home localnet --chain-id localnet

read -p "Are you sure? " -n 1 -r

echo $MNEMONIC | farmingd keys add validator --home localnet --recover --keyring-backend test
farmingd add-genesis-account $(farmingd keys show validator -a --home localnet --keyring-backend test) $GENESIS_COINS --home localnet

farmingd gentx validator 1000000000stake --home localnet --chain-id localnet --keyring-backend test
farmingd collect-gentxs --home localnet

sed -i '' 's/timeout_commit = "5s"/timeout_commit = "3s"/g' localnet/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "2s"/g' localnet/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' localnet/config/config.toml
sed -i '' 's/enable = false/enable = true/g' localnet/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' localnet/config/app.toml

farmingd start --home localnet --inv-check-period 1 --x-crisis-skip-assert-invariants false
