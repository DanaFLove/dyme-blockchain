{\rtf1\ansi\ansicpg1252\cocoartf2821
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 package app\
\
import (\
	"github.com/cosmos/cosmos-sdk/codec"\
	"github.com/cosmos/cosmos-sdk/codec/types"\
	"github.com/cosmos/cosmos-sdk/std"\
	"github.com/cosmos/cosmos-sdk/x/auth/tx"\
\
	"dymechain/app/params"\
)\
\
// makeEncodingConfig creates an EncodingConfig for an amino based test configuration.\
func makeEncodingConfig() params.EncodingConfig \{\
	amino := codec.NewLegacyAmino()\
	interfaceRegistry := types.NewInterfaceRegistry()\
	marshaler := codec.NewProtoCodec(interfaceRegistry)\
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)\
\
	return params.EncodingConfig\{\
		InterfaceRegistry: interfaceRegistry,\
		Marshaler:         marshaler,\
		TxConfig:          txCfg,\
		Amino:             amino,\
	\}\
\}\
\
// MakeEncodingConfig creates an EncodingConfig for testing\
func MakeEncodingConfig() params.EncodingConfig \{\
	encodingConfig := makeEncodingConfig()\
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)\
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)\
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)\
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)\
	return encodingConfig\
\}}