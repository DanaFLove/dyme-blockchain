{\rtf1\ansi\ansicpg1252\cocoartf2821
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 package app_test\
\
import (\
	"os"\
	"testing"\
	"time"\
\
	"github.com/cosmos/cosmos-sdk/simapp"\
	simulationtypes "github.com/cosmos/cosmos-sdk/types/simulation"\
	"github.com/cosmos/cosmos-sdk/x/simulation"\
	"github.com/stretchr/testify/require"\
	abci "github.com/tendermint/tendermint/abci/types"\
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"\
	tmtypes "github.com/tendermint/tendermint/types"\
\
	"dymechain/app"\
)\
\
func init() \{\
	simapp.GetSimulatorFlags()\
\}\
\
var defaultConsensusParams = &abci.ConsensusParams\{\
	Block: &abci.BlockParams\{\
		MaxBytes: 200000,\
		MaxGas:   2000000,\
	\},\
	Evidence: &tmproto.EvidenceParams\{\
		MaxAgeNumBlocks: 302400,\
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration\
		MaxBytes:        10000,\
	\},\
	Validator: &tmproto.ValidatorParams\{\
		PubKeyTypes: []string\{\
			tmtypes.ABCIPubKeyTypeEd25519,\
		\},\
	\},\
\}\
\
// BenchmarkSimulation run the chain simulation\
// Running using starport command:\
// `starport chain simulate -v --numBlocks 200 --blockSize 50`\
// Running as go benchmark test:\
// `go test -benchmem -run=^$ -bench ^BenchmarkSimulation ./app -NumBlocks=200 -BlockSize 50 -Commit=true -Verbose=true -Enabled=true`\
func BenchmarkSimulation(b *testing.B) \{\
	simapp.FlagEnabledValue = true\
	simapp.FlagCommitValue = true\
\
	config, db, dir, logger, _, err := simapp.SetupSimulation("goleveldb-app-sim", "Simulation")\
	require.NoError(b, err, "simulation setup failed")\
\
	b.Cleanup(func() \{\
		db.Close()\
		err = os.RemoveAll(dir)\
		require.NoError(b, err)\
	\})\
\
	encoding := app.MakeEncodingConfig()\
\
	app := app.New(\
		logger,\
		db,\
		nil,\
		true,\
		map[int64]bool\{\},\
		app.DefaultNodeHome,\
		0,\
		encoding,\
		simapp.EmptyAppOptions\{\},\
	)\
\
	// Run randomized simulations\
	_, simParams, simErr := simulation.SimulateFromSeed(\
		b,\
		os.Stdout,\
		app.BaseApp,\
		simapp.AppStateFn(app.AppCodec(), app.SimulationManager()),\
		simulationtypes.RandomAccounts,\
		simapp.SimulationOperations(app, app.AppCodec(), config),\
		app.ModuleAccountAddrs(),\
		config,\
		app.AppCodec(),\
	)\
\
	// export state and simParams before the simulation error is checked\
	err = simapp.CheckExportSimulation(app, config, simParams)\
	require.NoError(b, err)\
	require.NoError(b, simErr)\
\
	if config.Commit \{\
		simapp.PrintStats(db)\
	\}\
\}}