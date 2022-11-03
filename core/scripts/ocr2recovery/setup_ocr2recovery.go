package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/core/cmd"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type jobType string

const (
	jobTypeDKG          jobType = "DKG"
	jobTypeOCR2Recovery jobType = "OCR2Recovery"
)

func setupOCR2RecoveryNodes(e helpers.Environment) {
	fs := flag.NewFlagSet("ocr2recovery-setup", flag.ExitOnError)

	keyID := fs.String("key-id", "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc", "key ID")
	useForwarder := fs.Bool("use-forwarder", false, "boolean to use the forwarder")
	confDelays := fs.String("conf-delays", "1,2,3,4,5,6,7,8", "8 confirmation delays")
	lookbackBlocks := fs.Int64("lookback-blocks", 1000, "lookback blocks")

	apiFile := fs.String("api", "../../../tools/secrets/apicredentials", "api credentials file")
	passwordFile := fs.String("password", "../../../tools/secrets/password.txt", "password file")
	databasePrefix := fs.String("database-prefix", "postgres://postgres:postgres@localhost:5432/ocr2recovery-test", "database prefix")
	databaseSuffixes := fs.String("database-suffixes", "sslmode=disable", "database parameters to be added")
	nodeCount := fs.Int("node-count", 6, "number of nodes")
	fundingAmount := fs.Int64("funding-amount", 5e17, "amount to fund nodes") // .1 ETH

	helpers.ParseArgs(fs, os.Args[2:])

	if *nodeCount < 6 {
		fmt.Println("Node count too low for OCR2Recovery job, need at least 6.")
		os.Exit(1)
	}

	delays := helpers.ParseIntSlice(*confDelays)
	if len(delays) != 8 {
		fmt.Println("confDelays must have a length of 8")
		os.Exit(1)
	}

	// Deploy DKG and Recovery contracts, and add the Recovery contract
	// as a consumer of DKG events.
	fmt.Println("Deploying DKG contract...")
	dkgAddress := deployDKG(e)

	fmt.Println("Deploying Recovery beacon...")
	recoveryBeaconAddress := deployRecoveryBeacon(e, dkgAddress.String(), *keyID)

	fmt.Println("Adding Recovery Beacon as DKG client...")
	addClientToDKG(e, dkgAddress.String(), *keyID, recoveryBeaconAddress.String())

	fmt.Println("Configuring nodes with OCR2Recovery jobs...")
	var (
		onChainPublicKeys  []string
		offChainPublicKeys []string
		configPublicKeys   []string
		peerIDs            []string
		transmitters       []string
		dkgEncrypters      []string
		dkgSigners         []string
		sendingKeys        [][]string
	)

	for i := 0; i < *nodeCount; i++ {
		flagSet := flag.NewFlagSet("run-ocr2recovery-job-creation", flag.ExitOnError)
		flagSet.String("api", *apiFile, "api file")
		flagSet.String("password", *passwordFile, "password file")
		flagSet.String("vrfpassword", *passwordFile, "vrfpassword file")
		flagSet.String("bootstrapPort", fmt.Sprintf("%d", 8000), "port of bootstrap")
		flagSet.Int64("chainID", e.ChainID, "the chain ID")

		flagSet.String("job-type", string(jobTypeOCR2Recovery), "the job type")

		// used by bootstrap template instantiation
		flagSet.String("contractID", dkgAddress.String(), "the contract to get peers from")

		// DKG args
		flagSet.String("keyID", *keyID, "")
		flagSet.String("dkg-address", dkgAddress.String(), "the contract address of the DKG")

		// Recovery args
		flagSet.String("recovery-beacon-address", recoveryBeaconAddress.String(), "the contract address of the Recovery Beacon")
		flagSet.Int64("lookback-blocks", *lookbackBlocks, "lookback blocks")
		flagSet.String("confirmation-delays", *confDelays, "confirmation delays")

		flagSet.Bool("dangerWillRobinson", true, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "peerID of first node")

		payload := setupNode(flagSet, i, *databasePrefix, *databaseSuffixes, *useForwarder)

		onChainPublicKeys = append(onChainPublicKeys, payload.OnChainPublicKey)
		offChainPublicKeys = append(offChainPublicKeys, payload.OffChainPublicKey)
		configPublicKeys = append(configPublicKeys, payload.ConfigPublicKey)
		peerIDs = append(peerIDs, payload.PeerID)
		transmitters = append(transmitters, payload.Transmitter)
		dkgEncrypters = append(dkgEncrypters, payload.DkgEncrypt)
		dkgSigners = append(dkgSigners, payload.DkgSign)
		sendingKeys = append(sendingKeys, payload.SendingKeys)
	}

	var nodesToFund []string
	for _, t := range transmitters[1:] {
		nodesToFund = append(nodesToFund, t)
	}

	var payees []common.Address
	var reportTransmitters []common.Address // all transmitters excluding bootstrap
	for _, t := range transmitters[1:] {
		payees = append(payees, e.Owner.From)
		reportTransmitters = append(reportTransmitters, common.HexToAddress(t))
	}

	fmt.Println("Funding transmitters...")
	helpers.FundNodes(e, nodesToFund, big.NewInt(*fundingAmount))

	fmt.Println("Generated dkg setConfig command:")
	dkgCommand := fmt.Sprintf(
		"go run . dkg-set-config -dkg-address %s -key-id %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -dkg-encryption-pub-keys %s -dkg-signing-pub-keys %s -schedule 1,1,1,1,1",
		dkgAddress.String(),
		*keyID,
		strings.Join(onChainPublicKeys[1:], ","),
		strings.Join(offChainPublicKeys[1:], ","),
		strings.Join(configPublicKeys[1:], ","),
		strings.Join(peerIDs[1:], ","),
		strings.Join(transmitters[1:], ","),
		strings.Join(dkgEncrypters[1:], ","),
		strings.Join(dkgSigners[1:], ","),
	)
	fmt.Println(dkgCommand)

	fmt.Println()
	fmt.Println("Generated recovery beacon setConfig command:")
	recoveryCommand := fmt.Sprintf(
		"go run . beacon-set-config -beacon-address %s -conf-delays %s -onchain-pub-keys %s -offchain-pub-keys %s -config-pub-keys %s -peer-ids %s -transmitters %s -schedule 1,1,1,1,1",
		recoveryBeaconAddress.String(),
		*confDelays,
		strings.Join(onChainPublicKeys[1:], ","),
		strings.Join(offChainPublicKeys[1:], ","),
		strings.Join(configPublicKeys[1:], ","),
		strings.Join(peerIDs[1:], ","),
		strings.Join(transmitters[1:], ","),
	)
	fmt.Println(recoveryCommand)
}

func setupNode(flagSet *flag.FlagSet, nodeIdx int, databasePrefix, databaseSuffixes string, useForwarder bool) *cmd.SetupOCR2RecoveryNodePayload {
	client := newSetupClient()
	app := cmd.NewApp(client)
	ctx := cli.NewContext(app, flagSet, nil)

	defer func() {
		err := app.After(ctx)
		helpers.PanicErr(err)
	}()

	err := app.Before(ctx)
	helpers.PanicErr(err)

	resetDatabase(client, ctx, nodeIdx, databasePrefix, databaseSuffixes)
	configureEnvironmentVariables((useForwarder) && (nodeIdx > 0))

	return setupOCR2RecoveryNodeFromClient(client, ctx)
}
