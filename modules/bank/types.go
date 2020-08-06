package bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/types"
	json2 "github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	maxMsgLen             = 5
	ModuleName            = "bank"
)

var (
	_ types.Msg = MsgSend{}
	_ types.Msg = MsgMultiSend{}

	cdc = types.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

//type MsgSend struct {
//	Inputs  []Input  `json:"inputs"`
//	Outputs []Output `json:"outputs"`
//}
//
//type MsgSend struct {
//	Inputs  []Input  `json:"inputs"`
//	Outputs []Output `json:"outputs"`
//}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

func (msg MsgMultiSend) Route() string { return ModuleName }

// Implements Msg.
func (msg MsgMultiSend) Type() string { return "send" }

// Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return errors.New("invalid input coins")
	}
	if len(msg.Outputs) == 0 {
		return errors.New("invalid output coins")
	}
	// make sure all inputs and outputs are individually valid
	var totalIn, totalOut types.Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Add(in.Coins...)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Add(out.Coins...)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return errors.New("inputs and outputs don't match")
	}
	return nil
}

// Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := cdc.MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return json2.MustSort(b)
}

// Implements Msg.
func (msg MsgMultiSend) GetSigners() []types.AccAddress {
	addrs := make([]types.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}
// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr types.AccAddress, amount types.Coins) MsgSend {
	return MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amount}
}

func (msg MsgSend) Route() string {
	return ModuleName
}

func (msg MsgSend) Type() string {
	return "send"
}

func (msg MsgSend) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return errors.New("missing sender address")
	}
	if msg.ToAddress.Empty() {
		return errors.New("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return errors.New("invalid coins")
	}
	if !msg.Amount.IsAllPositive() {
		return errors.New("invalid coins")
	}
	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	bz, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return types.MustSortJSON(bz)
}

func (msg MsgSend) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.FromAddress}
}

////----------------------------------------
//// Input
//
//// Transaction Input
//type Input struct {
//	Address types.AccAddress `json:"address"`
//	Coins   types.Coins      `json:"coins"`
//}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := cdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(bin)
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", in.Address.String())))
	}
	if in.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !in.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", in.Coins))
	}
	return nil
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr types.AccAddress, coins types.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

////----------------------------------------
//// Output
//
//// Transaction Output
//type Output struct {
//	Address types.AccAddress `json:"address"`
//	Coins   types.Coins      `json:"coins"`
//}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := cdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(bin)
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", out.Address.String())))
	}
	if out.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !out.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", out.Coins))
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr types.AccAddress, coins types.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}


type tokenStats struct {
	LooseTokens  types.Coins `json:"loose_tokens"`
	BondedTokens types.Coins `json:"bonded_tokens"`
	BurnedTokens types.Coins `json:"burned_tokens"`
	TotalSupply  types.Coins `json:"total_supply"`
}

func (ts tokenStats) Convert() interface{} {
	return rpc.TokenStats{
		LooseTokens:  ts.LooseTokens,
		BondedTokens: ts.BondedTokens,
		BurnedTokens: ts.BurnedTokens,
		TotalSupply:  ts.TotalSupply,
	}
}

func registerCodec(cdc types.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "irishub/bank/Send")
	cdc.RegisterConcrete(MsgMultiSend{}, "irishub/bank/MultiSend")

	//cdc.RegisterConcrete(&Params{}, "irishub/Auth/Params")
}
