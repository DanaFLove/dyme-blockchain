rm -rf ~/.dymechaind

./dymechaind init $1

./dymechaind keys add masterwallet --keyring-backend test 

MY_VALIDATOR_ADDRESS=$(./dymechaind keys show masterwallet -a --keyring-backend test)

./dymechaind add-genesis-account $MY_VALIDATOR_ADDRESS 100000000000stake

./dymechaind gentx masterwallet 100000000stake --chain-id dymechain --keyring-backend test
./dymechaind collect-gentxs

./dymechaind start --rpc.laddr="tcp://0.0.0.0:26657" --api.enable --p2p.external-address=$2:26656