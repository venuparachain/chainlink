package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type commonSetConfigArgs struct {
	onchainPubKeys         string
	offchainPubKeys        string
	configPubKeys          string
	peerIDs                string
	transmitters           string
	schedule               string
	f                      uint
	deltaProgress          time.Duration
	deltaResend            time.Duration
	deltaRound             time.Duration
	deltaGrace             time.Duration
	deltaStage             time.Duration
	maxRounds              uint8
	maxDurationQuery       time.Duration
	maxDurationObservation time.Duration
	maxDurationReport      time.Duration
	maxDurationAccept      time.Duration
	maxDurationTransmit    time.Duration
}

type dkgSetConfigArgs struct {
	commonSetConfigArgs
	dkgEncryptionPubKeys string
	dkgSigningPubKeys    string
	keyID                string
}

type vrfBeaconSetConfigArgs struct {
	commonSetConfigArgs
	confDelays string
}

func main() {
	e := helpers.SetupEnv(false)

	switch os.Args[1] {

	case "dkg-deploy":
		deployDKG(e)

	case "beacon-deploy":
		cmd := flag.NewFlagSet("beacon-deploy", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "dkg contract address")
		keyID := cmd.String("key-id", "", "key ID")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-deploy", "key-id")
		deployRecoveryBeacon(e, *dkgAddress, *keyID)

	case "dkg-add-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		addClientToDKG(e, *dkgAddress, *keyID, *clientAddress)

	case "dkg-remove-client":
		cmd := flag.NewFlagSet("dkg-add-client", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		clientAddress := cmd.String("client-address", "", "client address")
		helpers.ParseArgs(cmd, os.Args[2:], "dkg-address", "key-id", "client-address")
		removeClientFromDKG(e, *dkgAddress, *keyID, *clientAddress)

	case "dkg-set-config":
		cmd := flag.NewFlagSet("dkg-set-config", flag.ExitOnError)
		dkgAddress := cmd.String("dkg-address", "", "DKG contract address")
		keyID := cmd.String("key-id", "", "key ID")
		onchainPubKeys := cmd.String("onchain-pub-keys", "", "comma-separated list of OCR on-chain pubkeys")
		offchainPubKeys := cmd.String("offchain-pub-keys", "", "comma-separated list of OCR off-chain pubkeys")
		configPubKeys := cmd.String("config-pub-keys", "", "comma-separated list of OCR config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma-separated list of peer IDs")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		dkgEncryptionPubKeys := cmd.String("dkg-encryption-pub-keys", "", "comma-separated list of DKG encryption pubkeys")
		dkgSigningPubKeys := cmd.String("dkg-signing-pub-keys", "", "comma-separated list of DKG signing pubkeys")
		schedule := cmd.String("schedule", "", "comma-separted list of transmission schedule")
		f := cmd.Uint("f", 1, "number of faulty oracles")
		deltaProgress := cmd.Duration("delta-progress", 30*time.Second, "duration of delta progress")
		deltaResend := cmd.Duration("delta-resend", 10*time.Second, "duration of delta resend")
		deltaRound := cmd.Duration("delta-round", 10*time.Second, "duration of delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "duration of delta grace")
		deltaStage := cmd.Duration("delta-stage", 20*time.Second, "duration of delta stage")
		maxRounds := cmd.Uint("max-rounds", 3, "maximum number of rounds")
		maxDurationQuery := cmd.Duration("max-duration-query", 10*time.Millisecond, "maximum duration of query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 10*time.Second, "maximum duration of observation method")
		maxDurationReport := cmd.Duration("max-duration-report", 10*time.Second, "maximum duration of report method")
		maxDurationAccept := cmd.Duration("max-duration-accept", 10*time.Millisecond, "maximum duration of shouldAcceptFinalizedReport method")
		maxDurationTransmit := cmd.Duration("max-duration-transmit", 1*time.Second, "maximum duration of shouldTransmitAcceptedReport method")

		helpers.ParseArgs(cmd,
			os.Args[2:],
			"dkg-address",
			"key-id",
			"onchain-pub-keys",
			"offchain-pub-keys",
			"config-pub-keys",
			"peer-ids",
			"transmitters",
			"dkg-encryption-pub-keys",
			"dkg-signing-pub-keys",
			"schedule")

		commands := dkgSetConfigArgs{
			commonSetConfigArgs: commonSetConfigArgs{
				onchainPubKeys:         *onchainPubKeys,
				offchainPubKeys:        *offchainPubKeys,
				configPubKeys:          *configPubKeys,
				peerIDs:                *peerIDs,
				transmitters:           *transmitters,
				schedule:               *schedule,
				f:                      *f,
				deltaProgress:          *deltaProgress,
				deltaResend:            *deltaResend,
				deltaRound:             *deltaRound,
				deltaGrace:             *deltaGrace,
				deltaStage:             *deltaStage,
				maxRounds:              uint8(*maxRounds),
				maxDurationQuery:       *maxDurationQuery,
				maxDurationObservation: *maxDurationObservation,
				maxDurationReport:      *maxDurationReport,
				maxDurationAccept:      *maxDurationAccept,
				maxDurationTransmit:    *maxDurationTransmit,
			},
			dkgEncryptionPubKeys: *dkgEncryptionPubKeys,
			dkgSigningPubKeys:    *dkgSigningPubKeys,
			keyID:                *keyID,
		}

		setDKGConfig(e, *dkgAddress, commands)

	case "beacon-set-config":
		cmd := flag.NewFlagSet("beacon-set-config", flag.ExitOnError)
		beaconAddress := cmd.String("beacon-address", "", "VRF beacon contract address")
		confDelays := cmd.String("conf-delays", "1,2,3,4,5,6,7,8", "comma-separted list of 8 confirmation delays")
		onchainPubKeys := cmd.String("onchain-pub-keys", "", "comma-separated list of OCR on-chain pubkeys")
		offchainPubKeys := cmd.String("offchain-pub-keys", "", "comma-separated list of OCR off-chain pubkeys")
		configPubKeys := cmd.String("config-pub-keys", "", "comma-separated list of OCR config pubkeys")
		peerIDs := cmd.String("peer-ids", "", "comma-separated list of peer IDs")
		transmitters := cmd.String("transmitters", "", "comma-separated list transmitters")
		schedule := cmd.String("schedule", "", "comma-separted list of transmission schedule")
		f := cmd.Uint("f", 1, "number of faulty oracles")
		// TODO: Adjust default delta* and maxDuration* values below after benchmarking latency
		deltaProgress := cmd.Duration("delta-progress", 30*time.Second, "duration of delta progress")
		deltaResend := cmd.Duration("delta-resend", 10*time.Second, "duration of delta resend")
		deltaRound := cmd.Duration("delta-round", 10*time.Second, "duration of delta round")
		deltaGrace := cmd.Duration("delta-grace", 20*time.Second, "duration of delta grace")
		deltaStage := cmd.Duration("delta-stage", 20*time.Second, "duration of delta stage")
		maxRounds := cmd.Uint("max-rounds", 3, "maximum number of rounds")
		maxDurationQuery := cmd.Duration("max-duration-query", 10*time.Millisecond, "maximum duration of query")
		maxDurationObservation := cmd.Duration("max-duration-observation", 10*time.Second, "maximum duration of observation method")
		maxDurationReport := cmd.Duration("max-duration-report", 10*time.Second, "maximum duration of report method")
		maxDurationAccept := cmd.Duration("max-duration-accept", 10*time.Millisecond, "maximum duration of shouldAcceptFinalizedReport method")
		maxDurationTransmit := cmd.Duration("max-duration-transmit", 1*time.Second, "maximum duration of shouldTransmitAcceptedReport method")

		helpers.ParseArgs(cmd,
			os.Args[2:],
			"beacon-address",
			"onchain-pub-keys",
			"offchain-pub-keys",
			"config-pub-keys",
			"peer-ids",
			"transmitters",
			"schedule")

		commands := vrfBeaconSetConfigArgs{
			commonSetConfigArgs: commonSetConfigArgs{
				onchainPubKeys:         *onchainPubKeys,
				offchainPubKeys:        *offchainPubKeys,
				configPubKeys:          *configPubKeys,
				peerIDs:                *peerIDs,
				transmitters:           *transmitters,
				schedule:               *schedule,
				f:                      *f,
				deltaProgress:          *deltaProgress,
				deltaResend:            *deltaResend,
				deltaRound:             *deltaRound,
				deltaGrace:             *deltaGrace,
				deltaStage:             *deltaStage,
				maxRounds:              uint8(*maxRounds),
				maxDurationQuery:       *maxDurationQuery,
				maxDurationObservation: *maxDurationObservation,
				maxDurationReport:      *maxDurationReport,
				maxDurationAccept:      *maxDurationAccept,
				maxDurationTransmit:    *maxDurationTransmit,
			},
			confDelays: *confDelays,
		}

		setRecoveryBeaconConfig(e, *beaconAddress, commands)

	case "beacon-info":
		cmd := flag.NewFlagSet("beacon-info", flag.ExitOnError)
		beaconAddress := cmd.String("beacon-address", "", "Recovery beacon contract address")
		helpers.ParseArgs(cmd, os.Args[2:], "beacon-address")
		beacon := newRecoveryBeacon(common.HexToAddress(*beaconAddress), e.Ec)
		keyID, err := beacon.SKeyID(nil)
		helpers.PanicErr(err)
		fmt.Println("beacon key id:", hexutil.Encode(keyID[:]))
		keyHash, err := beacon.SProvingKeyHash(nil)
		helpers.PanicErr(err)
		fmt.Println("beacon proving key hash:", hexutil.Encode(keyHash[:]))

	case "dkg-setup":
		setupDKGNodes(e)
	case "ocr2recovery-setup":
		setupOCR2RecoveryNodes(e)
	default:
		panic("unrecognized subcommand: " + os.Args[1])
	}
}
