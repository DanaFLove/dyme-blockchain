#Note: - this fully reset older blockchain data

rm -rf ~/.dymechaind
rm -rf ~/.dymechain/

mkdir ~/.dymechain/
mkdir ~/.dymechain/config

cp ./chain-configs/* ~/.dymechain/config/

dymechaind init $1

dymechaind keys add dymemaster1 --keyring-backend test 

MY_VALIDATOR_ADDRESS=$(dymechaind keys show dymemaster1 -a --keyring-backend test)

# total balance of validator
dymechaind add-genesis-account $MY_VALIDATOR_ADDRESS 10000000000000udyme

# bonded tokens of validator
dymechaind gentx dymemaster1 1000000udyme --chain-id dymechain --keyring-backend test
dymechaind collect-gentxs

dymechaind keys add dymemaster2 --keyring-backend test
dymechaind keys add dymemaster3 --keyring-backend test
dymechaind keys add dymemaster4 --keyring-backend test