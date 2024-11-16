package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/Guesstrain/EthBankok/models"
	"github.com/gin-gonic/gin"
	"github.com/onflow/go-ethereum"
	"github.com/onflow/go-ethereum/accounts/abi"
	"github.com/onflow/go-ethereum/common"
	"github.com/onflow/go-ethereum/ethclient"
)

const (
	infuraURL          = "https://sepolia.infura.io/v3/d68e6d7c2e5c42fbb30fe563ada8f432"
	contractAddressHex = "0xAF7785F8dDFC9629949eDdb07Ba14d53Fc853C14"
	privateKeyHex      = "83476b334581d65de37aee25349b36c5b5e917c7acddb3fc9ff50ef5feb87eaa" // Replace with your wallet's private key
	walletAddressHex   = "0xDb90007b986c2711d3814b1760EDC3b2DfB71e76"
	contractABI        = `[
{
"inputs": [
{
"internalType": "string",
"name": "name_",
"type": "string"
},
{
"internalType": "string",
"name": "symbol_",
"type": "string"
}
],
"stateMutability": "nonpayable",
"type": "constructor"
},
{
"inputs": [
{
"internalType": "address",
"name": "spender",
"type": "address"
},
{
"internalType": "uint256",
"name": "allowance",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "needed",
"type": "uint256"
}
],
"name": "ERC20InsufficientAllowance",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "sender",
"type": "address"
},
{
"internalType": "uint256",
"name": "balance",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "needed",
"type": "uint256"
}
],
"name": "ERC20InsufficientBalance",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "approver",
"type": "address"
}
],
"name": "ERC20InvalidApprover",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "receiver",
"type": "address"
}
],
"name": "ERC20InvalidReceiver",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "sender",
"type": "address"
}
],
"name": "ERC20InvalidSender",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "spender",
"type": "address"
}
],
"name": "ERC20InvalidSpender",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "owner",
"type": "address"
}
],
"name": "OwnableInvalidOwner",
"type": "error"
},
{
"inputs": [
{
"internalType": "address",
"name": "account",
"type": "address"
}
],
"name": "OwnableUnauthorizedAccount",
"type": "error"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "owner",
"type": "address"
},
{
"indexed": true,
"internalType": "address",
"name": "spender",
"type": "address"
},
{
"indexed": false,
"internalType": "uint256",
"name": "value",
"type": "uint256"
}
],
"name": "Approval",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "user",
"type": "address"
},
{
"indexed": true,
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"indexed": false,
"internalType": "uint256",
"name": "amount",
"type": "uint256"
},
{
"indexed": false,
"internalType": "uint256",
"name": "dueDate",
"type": "uint256"
},
{
"indexed": false,
"internalType": "uint256",
"name": "loanId",
"type": "uint256"
}
],
"name": "LoanGiven",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"indexed": true,
"internalType": "address",
"name": "user",
"type": "address"
},
{
"indexed": false,
"internalType": "uint256",
"name": "amount",
"type": "uint256"
},
{
"indexed": false,
"internalType": "uint256",
"name": "loanId",
"type": "uint256"
}
],
"name": "LoanRepaid",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"indexed": false,
"internalType": "address",
"name": "spender",
"type": "address"
},
{
"indexed": false,
"internalType": "uint256",
"name": "amountApproved",
"type": "uint256"
}
],
"name": "MerchantInitialized",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "previousOwner",
"type": "address"
},
{
"indexed": true,
"internalType": "address",
"name": "newOwner",
"type": "address"
}
],
"name": "OwnershipTransferred",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "newRepaymentManager",
"type": "address"
}
],
"name": "RepaymentManagerSet",
"type": "event"
},
{
"anonymous": false,
"inputs": [
{
"indexed": true,
"internalType": "address",
"name": "from",
"type": "address"
},
{
"indexed": true,
"internalType": "address",
"name": "to",
"type": "address"
},
{
"indexed": false,
"internalType": "uint256",
"name": "value",
"type": "uint256"
}
],
"name": "Transfer",
"type": "event"
},
{
"inputs": [
{
"internalType": "address",
"name": "owner",
"type": "address"
},
{
"internalType": "address",
"name": "spender",
"type": "address"
}
],
"name": "allowance",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "spender",
"type": "address"
},
{
"internalType": "uint256",
"name": "value",
"type": "uint256"
}
],
"name": "approve",
"outputs": [
{
"internalType": "bool",
"name": "",
"type": "bool"
}
],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "account",
"type": "address"
}
],
"name": "balanceOf",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "checkUpkeep",
"outputs": [
{
"internalType": "bool",
"name": "upkeepNeeded",
"type": "bool"
},
{
"internalType": "uint256[]",
"name": "loanIdsToRepay",
"type": "uint256[]"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "decimals",
"outputs": [
{
"internalType": "uint8",
"name": "",
"type": "uint8"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "user",
"type": "address"
},
{
"internalType": "address",
"name": "merchant",
"type": "address"
}
],
"name": "getLoansBetween",
"outputs": [
{
"internalType": "string",
"name": "",
"type": "string"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "merchant",
"type": "address"
}
],
"name": "getTransactions",
"outputs": [
{
"internalType": "string",
"name": "",
"type": "string"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "uint256",
"name": "amount",
"type": "uint256"
}
],
"name": "initializeMerchantApproval",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"internalType": "uint256",
"name": "amount",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "daysUntilDue",
"type": "uint256"
}
],
"name": "lendToMerchant",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "",
"type": "address"
},
{
"internalType": "address",
"name": "",
"type": "address"
}
],
"name": "loanCount",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "",
"type": "address"
},
{
"internalType": "address",
"name": "",
"type": "address"
},
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"name": "loanIds",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"name": "loans",
"outputs": [
{
"internalType": "uint256",
"name": "loanId",
"type": "uint256"
},
{
"internalType": "address",
"name": "buyer",
"type": "address"
},
{
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"internalType": "uint256",
"name": "amount",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "dueDate",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "repaidAmount",
"type": "uint256"
},
{
"internalType": "bool",
"name": "isRepaid",
"type": "bool"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "",
"type": "address"
}
],
"name": "merchantCurrentLoans",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "",
"type": "address"
}
],
"name": "merchantMaxLoanLimits",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "name",
"outputs": [
{
"internalType": "string",
"name": "",
"type": "string"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "nextLoanId",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "owner",
"outputs": [
{
"internalType": "address",
"name": "",
"type": "address"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "uint256[]",
"name": "loanIdsToRepay",
"type": "uint256[]"
}
],
"name": "performUpkeep",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [],
"name": "renounceOwnership",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "user",
"type": "address"
},
{
"internalType": "uint256",
"name": "loanId",
"type": "uint256"
},
{
"internalType": "uint256",
"name": "amount",
"type": "uint256"
}
],
"name": "repayLoan",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [],
"name": "repaymentManager",
"outputs": [
{
"internalType": "address",
"name": "",
"type": "address"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "merchant",
"type": "address"
},
{
"internalType": "uint256",
"name": "maxLimit",
"type": "uint256"
}
],
"name": "setMerchantMaxLoanLimit",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "_repaymentManager",
"type": "address"
}
],
"name": "setRepaymentManager",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [],
"name": "symbol",
"outputs": [
{
"internalType": "string",
"name": "",
"type": "string"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [],
"name": "totalSupply",
"outputs": [
{
"internalType": "uint256",
"name": "",
"type": "uint256"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "to",
"type": "address"
},
{
"internalType": "uint256",
"name": "value",
"type": "uint256"
}
],
"name": "transfer",
"outputs": [
{
"internalType": "bool",
"name": "",
"type": "bool"
}
],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "from",
"type": "address"
},
{
"internalType": "address",
"name": "to",
"type": "address"
},
{
"internalType": "uint256",
"name": "value",
"type": "uint256"
}
],
"name": "transferFrom",
"outputs": [
{
"internalType": "bool",
"name": "",
"type": "bool"
}
],
"stateMutability": "nonpayable",
"type": "function"
},
{
"inputs": [
{
"internalType": "address",
"name": "newOwner",
"type": "address"
}
],
"name": "transferOwnership",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
}
]

自动还款合约[
{
"inputs": [],
"name": "checkUpkeep",
"outputs": [
{
"internalType": "bool",
"name": "upkeepNeeded",
"type": "bool"
},
{
"internalType": "uint256[]",
"name": "loanIdsToRepay",
"type": "uint256[]"
}
],
"stateMutability": "view",
"type": "function"
},
{
"inputs": [
{
"internalType": "uint256[]",
"name": "loanIdsToRepay",
"type": "uint256[]"
}
],
"name": "performUpkeep",
"outputs": [],
"stateMutability": "nonpayable",
"type": "function"
}
]
`
)

func GetTransactionHistory(c *gin.Context) {
	TransactionList := make([]models.Transaction, 0)
	TargetParam := c.Query("target")
	targetAddress := common.HexToAddress(TargetParam)

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/d68e6d7c2e5c42fbb30fe563ada8f432")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the Ethereum client"})
		return
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ABI"})
		return
	}

	lendEvent := parsedABI.Events["LoanGiven"]

	lendQuery := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{lendEvent.ID},
			{common.BytesToHash(targetAddress.Bytes())},
		},
	}

	Lendlogs, err := client.FilterLogs(context.Background(), lendQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter logs"})
		return
	}
	for _, lLog := range Lendlogs {
		event, err := parsedABI.Unpack("Transfer", lLog.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unpack event data"})
			return
		}

		from := lLog.Topics[1].Hex()
		to := lLog.Topics[2].Hex()
		value := event[0].(*big.Int)
		TransactionList = append(TransactionList, models.Transaction{from, to, value})

		fmt.Printf("lend from: %s, to: %s, value: %s\n", from, to, value.String())
	}

	c.JSON(http.StatusOK, TransactionList)
}

func GetTransactionHistoryMerchants(c *gin.Context) {
	TransactionList := make([]models.Transaction, 0)
	TargetParam := c.Query("target")
	targetAddress := common.HexToAddress(TargetParam)

	client, err := ethclient.Dial("https://sepolia.infura.io/v3/d68e6d7c2e5c42fbb30fe563ada8f432")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the Ethereum client"})
		return
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ABI"})
		return
	}

	lendEvent := parsedABI.Events["LoanGiven"]

	lendQuery := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{lendEvent.ID},
			nil,
			{common.BytesToHash(targetAddress.Bytes())},
		},
	}

	Lendlogs, err := client.FilterLogs(context.Background(), lendQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter logs"})
		return
	}
	for _, lLog := range Lendlogs {
		event := struct {
			Amount  *big.Int
			DueDate *big.Int
			LoanId  *big.Int
		}{}
		err = parsedABI.UnpackIntoInterface(&event, "LoanGiven", lLog.Data)

		from := lLog.Topics[1].Hex()
		to := lLog.Topics[2].Hex()
		fmt.Println("DueDate: ", event.DueDate)
		fmt.Println("LoanId: ", event.LoanId)
		value := event.Amount
		TransactionList = append(TransactionList, models.Transaction{from, to, value})

		fmt.Printf("lend from: %s, to: %s, value: %s\n", from, to, value.String())
	}

	c.JSON(http.StatusOK, TransactionList)
}

func GetAllLoans(c *gin.Context) error {
	TargetParam := c.Query("target")

	// Connect to the Ethereum network
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Create an instance of the contract
	contractAddr := common.HexToAddress(contractAddressHex)
	merchantAddr := common.HexToAddress(TargetParam)

	// Use the ABI encoding to generate the appropriate call data for `getTransactions()`
	// Here, we assume you've already parsed the contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Pack the arguments for the `getTransactions` function
	data, err := contractAbi.Pack("getTransactions", merchantAddr)
	if err != nil {
		log.Fatalf("Failed to pack arguments: %v", err)
	}

	// Prepare call message to invoke the contract
	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	// Execute the call
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Error while calling the getTransactions function: %v", err)
	}
	var loans []models.Loan
	var jsonString string
	err = contractAbi.UnpackIntoInterface(&jsonString, "getTransactions", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}
	cleanedString := strings.TrimSpace(jsonString)
	cleanedString2 := strings.TrimPrefix(cleanedString, "~")
	err = json.Unmarshal([]byte(cleanedString2), &loans)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	c.JSON(http.StatusOK, loans)
	return nil
}
