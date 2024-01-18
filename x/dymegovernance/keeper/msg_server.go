package keeper

import (
	"bytes"
	"dymechain/x/dymegovernance/types"
	tokenmanager_types "dymechain/x/dymetokenmanager/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type msgServer struct {
	Keeper
}

type Balance struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type Pagination struct {
	NextKey string `json:"next_key"`
	Total   uint64 `json:"total"`
}

type BalanceResponse struct {
	Balances   []Balance  `json:"balances"`
	Pagination Pagination `json:"pagination"`
}

func getBalanceFromAddress(recipient string) (bool, string) {
	var balanceResponse BalanceResponse
	var client http.Client
	resp, err := client.Get("http://localhost:1317/cosmos/bank/v1beta1/balances/" + recipient)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}

	json.NewDecoder(resp.Body).Decode(&balanceResponse)
	var dymeBalance float64
	dymeBalance = 0
	if len(balanceResponse.Balances) > 0 {
		if balanceResponse.Balances[0].Denom == "udyme" {
			amountparsed, _ := strconv.Atoi(balanceResponse.Balances[0].Amount)
			dymeBalance = float64(float64(amountparsed) / float64(math.Pow(10, tokenmanager_types.DymeDecimals)))
		}
	}
	response := fmt.Sprintf(":information_source: Wallet address: %s\nBalance: %f udyme", recipient, dymeBalance)
	return true, response
}

func getBalanceFromWalletName(walletName string) (bool, string) {
	getAddressCommand := exec.Command("dymechaind", "keys", "show", walletName, "--output", "json", "--keyring-backend", types.KeyRingBackend)
	addressOutput, err := getAddressCommand.Output()
	if err != nil {
		return false, ":x: Wallet not found"
	}

	var addressOutputMap map[string]string
	json.Unmarshal([]byte(addressOutput), &addressOutputMap)

	address := addressOutputMap["address"]
	if address == "" {
		return false, ":x: Wallet not found"
	}

	balanceFetched, balanceInfo := getBalanceFromAddress(address)

	return balanceFetched, balanceInfo
}

func sendHttpRequest(requestUrl string, dataToSend map[string]string) string {
	jsonDataToSend, _ := json.Marshal(dataToSend)
	requestBody := bytes.NewBuffer(jsonDataToSend)

	resp, err := http.Post(requestUrl, "application/json", requestBody)
	if err != nil {
		fmt.Println("An Error Occured: ", err.Error())
		return "Error while making http request"
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	sb := string(body)
	return sb
}

func validateUser(walletName string) (bool, string) {
	walletCreationCommand := exec.Command("dymechaind", "query", "dymegovernance", "get-key-value", walletName, "--output", "json")
	stdout, err := walletCreationCommand.Output()
	if err != nil {
		return false, ":x: Error occured while authenticating user!"
	}

	var getKvResp types.QueryGetKeyValueResponse
	json.Unmarshal(stdout, &getKvResp)
	if getKvResp.Value == "" {
		return false, ":x: Unauthenticated User!"
	} else {
		validityTime, parseError := time.Parse("2006-01-02 15:04:05", getKvResp.Value)
		if parseError != nil {
			return false, ":x: Error while parsing date!"
		}

		nowTime, nowParseError := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
		if nowParseError != nil {
			return false, ":x: Error while parsing date!"
		}

		if nowTime.Unix() > validityTime.Unix() {
			return false, ":x: Access window expired. Please verify user again!"
		}
	}

	return true, ":white_check_mark: Valid User!"
}

func discordMessageListener(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			s.ChannelMessageSendReply(m.ChannelID, ":x: Proposal transacton error", m.Reference())
			fmt.Println("***********************discordMessageListener PANIC ERROR:", err)
		}
	}()
	botId := s.State.User.ID

	// Ignore all messages created by the bot itself (good practice)
	if m.Author.ID == botId {
		return
	}
	walletName := m.Author.ID
	fmt.Println("Wallet name" + walletName)
	// Handling account registration
	if strings.Contains(m.Content, "dymeregister") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}
		walletCreationCommand := exec.Command("dymechaind", "keys", "add", walletName, "--output", "json", "--keyring-backend", types.KeyRingBackend)
		stdout, err := walletCreationCommand.Output()
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, ":x: An error has occurred while creating wallet. Most likely the wallet already exists", m.Reference())
			return
		}

		var jsonMap map[string]string
		json.Unmarshal([]byte(stdout), &jsonMap)

		name := jsonMap["name"]
		address := jsonMap["address"]

		s.ChannelMessageSendReply(m.ChannelID, ":white_check_mark: Congrats! Wallet has been created successfully!\nName: "+string(name)+"\nAddress: "+string(address), m.Reference())
		return
	}

	// Fetching DYME balance of a Discord Dyme user
	if strings.Contains(m.Content, "getbalance") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		balanceFetched, balanceMessage := getBalanceFromWalletName(walletName)
		if !balanceFetched {
			s.ChannelMessageSendReply(m.ChannelID, balanceMessage, m.Reference())
			return
		}

		s.ChannelMessageSendReply(m.ChannelID, balanceMessage, m.Reference())
		return
	}

	// only for demo
	if strings.Contains(m.Content, "clearchat") {
		messages, _ := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		var messageIds []string
		for _, message := range messages {
			messageIds = append(messageIds, message.ID)
		}
		fmt.Println("messageIds", messageIds)
		s.ChannelMessagesBulkDelete(m.ChannelID, messageIds)
		return
	}

	// Sending OTP to the user for verification
	if strings.Contains(m.Content, "sendotp") {

		BACKEND_DOMAIN := os.Getenv("BACKEND_DOMAIN")

		dataToSend := map[string]string{
			"discordUserId": m.Author.ID,
		}
		response := sendHttpRequest(BACKEND_DOMAIN+"/discord/send-verification-otp", dataToSend)

		message := ""
		if !strings.Contains(response, "OTP sent successfully!") { // ignore backend API call for now (TESTING ONLY)
			message = ":x: Cannot send OTP. Error: " + response
		} else {
			message = ":white_check_mark: OTP sent successfully"
		}

		s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		return
	}

	// Verifiying OTP and opening transactional window for configured amount of time
	if strings.Contains(m.Content, "verifyotp") {

		BACKEND_DOMAIN := os.Getenv("BACKEND_DOMAIN")

		setkvParts := strings.Split(m.Content, " ")
		if len(setkvParts) != 2 {
			s.ChannelMessageSendReply(m.ChannelID, ":x: Please provide all arguments \n verifyotp <OTP>", m.Reference())
			return
		}
		dataToSend := map[string]string{
			"discordUserId": m.Author.ID,
			"otp":           setkvParts[1],
		}
		response := sendHttpRequest(BACKEND_DOMAIN+"/discord/verify-otp", dataToSend)
		fmt.Println(response)

		message := ""
		if response != "Verification successful!" {
			message = ":x: " + response
		} else {
			message = ":white_check_mark: " + response
			validityTime := time.Now().Add(time.Minute * 30).Format(types.TimeFormat)
			setOtpValidityCommand := exec.Command("dymechaind", "tx", "dymegovernance", "set-key-value", walletName, validityTime, "--from", "dymemaster1", "--output", "json", "--keyring-backend", types.KeyRingBackend, "--yes")
			_, err2 := setOtpValidityCommand.Output()

			if err2 != nil {
				s.ChannelMessageSendReply(m.ChannelID, ":x: Transaction failed for setting key value.", m.Reference())
				return
			}
		}

		s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		return
	}

	if strings.Contains(m.Content, "createproposal") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		createproposalCommandParts := strings.Split(m.Content, " ")
		response := ""
		if len(createproposalCommandParts) > 1 {
			proposalName := strings.ReplaceAll(m.Content, "createproposal ", "")
			proposalCreationCommand := exec.Command("dymechaind", "tx", "gov", "submit-legacy-proposal", "--title="+proposalName, "--description="+proposalName, "--type=Text", "--deposit=1000000DYME", "--from", walletName, "--output", "json", "--keyring-backend", types.KeyRingBackend, "--yes")
			proposalOutput, errP := proposalCreationCommand.Output()
			if errP != nil {
				s.ChannelMessageSendReply(m.ChannelID, ":x: Transaction failed for proposal submission. Error: "+errP.Error(), m.Reference())
				return
			}

			// Note:- Most likely this cannot be used because cosmos will not return logs (and hence the proposal ID) in broadcase_mode = sync
			// Alternative: use broadcase_mode = block and reduce block time in genesis (maybe this is not ideal but will work)
			print("proposal output ", proposalCreationCommand.String())
			response = ":white_check_mark: Proposal created successfully."
			proposalId := strings.Split(string(proposalOutput), "proposal_id")
			proposalId = strings.Split(proposalId[1], "proposal_messages")
			proposalIdParsed := strings.ReplaceAll(proposalId[0], "\"", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "value", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "key", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, ",", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "{", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "}", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "\\", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "[", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "]", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, ":", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "attributes", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "submit_proposal", "")
			proposalIdParsed = strings.ReplaceAll(proposalIdParsed, "type", "")

			s.ChannelMessageSend(types.ProposalChannelId, ":page_with_curl: New Proposal - "+proposalName+"\nProposal ID: "+proposalIdParsed)
		} else {
			response = ":x: Proposal needs a name"
		}

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	}

	if strings.Contains(m.Content, "viewproposal") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		viewwproposalCommandParts := strings.Split(m.Content, " ")
		response := ""
		if len(viewwproposalCommandParts) > 1 {
			proposalViewCommand := exec.Command("dymechaind", "query", "gov", "proposal", viewwproposalCommandParts[1], "--output", "json")
			proposalOutput, _ := proposalViewCommand.Output()
			if string(proposalOutput) == "" {
				s.ChannelMessageSendReply(m.ChannelID, ":x: Proposal not found", m.Reference())
				return
			}
			var jsonObjGeneric map[string]interface{}
			json.Unmarshal([]byte(proposalOutput), &jsonObjGeneric)
			indentedJSON, _ := json.MarshalIndent(jsonObjGeneric, "", "    ")

			response = ":white_check_mark: Proposal created successfully \n Output:" + string(indentedJSON)
		} else {
			response = ":x: Enter a valid Proposal ID"
		}

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	}

	// stake 1 DYME (exact) in order to be eligible for voting
	if strings.Contains(m.Content, "stakedyme") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		response := ""
		stakeDymecommand := exec.Command("dymechaind", "tx", "dymegovernance", "stakedyme", "--from", walletName, "--keyring-backend", types.KeyRingBackend, "--output", "json", "--yes")
		stakeDymecommandOutput, _ := stakeDymecommand.Output()
		if string(stakeDymecommandOutput) == "" || strings.Contains(string(stakeDymecommandOutput), "failed to execute") {
			s.ChannelMessageSendReply(m.ChannelID, ":x:Stake dyme transaction failed [ERR1]: "+string(stakeDymecommandOutput), m.Reference())
			return
		}

		// TODO: make this dynamic by on-chain query
		validatorAddress := "cosmosvaloper1hrzkun2xgxu4gl2gsq3y6juzucwk446d2v3j8m"
		stakeCommand := exec.Command("dymechaind", "tx", "staking", "delegate", validatorAddress, "1000000udyme", "--from", walletName, "--keyring-backend", types.KeyRingBackend, "--output", "json", "--yes")
		stakeCommandOutput, _ := stakeCommand.Output()
		if string(stakeCommandOutput) == "" || strings.Contains(string(stakeCommandOutput), "failed to execute") {
			s.ChannelMessageSendReply(m.ChannelID, ":x:Stake dyme transaction failed [ERR2]: "+string(stakeCommandOutput), m.Reference())
			return
		}

		response = ":white_check_mark: 1 DYME staked successfully"

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	}

	if strings.Contains(m.Content, "voteonproposal") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		response := ""
		voteCommandParts := strings.Split(m.Content, " ")
		if len(voteCommandParts) > 1 {
			if voteCommandParts[2] != "yes" && voteCommandParts[2] != "no" && voteCommandParts[2] != "no_with_veto" && voteCommandParts[2] != "abstain" {
				response = ":x: Vote value can be yes, no, no_with_veto or abstain"
			} else {
				voteCommand := exec.Command("dymechaind", "tx", "gov", "vote", voteCommandParts[1], voteCommandParts[2], "--from", walletName, "--keyring-backend", types.KeyRingBackend, "--output", "json", "--yes")
				voteCommandOutput, _ := voteCommand.Output()
				if string(voteCommandOutput) == "" || strings.Contains(string(voteCommandOutput), "failed to execute") {
					s.ChannelMessageSendReply(m.ChannelID, ":x: Vote transaction failed", m.Reference())
					return
				}

				response = ":white_check_mark: Voted successfully"
			}

		} else {
			response = ":x: Enter a valid Proposal ID"
		}

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	}

	if strings.Contains(m.Content, "viewvotes") {
		isValidUser, validationMessage := validateUser(walletName)
		if !isValidUser {
			s.ChannelMessageSendReply(m.ChannelID, validationMessage, m.Reference())
			return
		}

		viewwproposalCommandParts := strings.Split(m.Content, " ")
		response := ""
		if len(viewwproposalCommandParts) > 1 {
			proposalViewCommand := exec.Command("dymechaind", "query", "gov", "votes", viewwproposalCommandParts[1], "--output", "json")
			proposalOutput, _ := proposalViewCommand.Output()
			if string(proposalOutput) == "" {
				s.ChannelMessageSendReply(m.ChannelID, ":x: Proposal not found", m.Reference())
				return
			}
			var jsonObjGeneric map[string]interface{}
			json.Unmarshal([]byte(proposalOutput), &jsonObjGeneric)
			indentedJSON, _ := json.MarshalIndent(jsonObjGeneric, "", "    ")

			response = ":white_check_mark: Proposal created successfully \n Output:" + string(indentedJSON)
		} else {
			response = ":x: Enter a valid Proposal ID"
		}

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	}
}

func initdymechaindiscordBot() {
	fmt.Println("------------------ STARTING DYME DISCORD BOT ON-CHAIN -----------------------")
	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")

	dg, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		panic("initdymechaindiscordBot ERROR: " + err.Error())
	}
	dg.AddHandler(discordMessageListener)
	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("******************** initdymechaindiscordBot error opening connection ******************")
		fmt.Println(err)
		return
	}

	fmt.Println("--------------STARTED DYME Discord Bot successfully-----------------------")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {

	go initdymechaindiscordBot()
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
