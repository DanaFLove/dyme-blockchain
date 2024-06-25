{\rtf1\ansi\ansicpg1252\cocoartf2821
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 package app\
\
import (\
	"encoding/json"\
\
	"github.com/cosmos/cosmos-sdk/codec"\
)\
\
// The genesis state of the blockchain is represented here as a map of raw json\
// messages key'd by a identifier string.\
// The identifier is used to determine which module genesis information belongs\
// to so it may be appropriately routed during init chain.\
// Within this application default genesis information is retrieved from\
// the ModuleBasicManager which populates json from each BasicModule\
// object provided to it during init.\
type GenesisState map[string]json.RawMessage\
\
// NewDefaultGenesisState generates the default state for the application.\
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState \{\
	return ModuleBasics.DefaultGenesis(cdc)\
\}}