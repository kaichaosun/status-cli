package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	gethbridge "github.com/status-im/status-go/eth-node/bridge/geth"
	"github.com/status-im/status-go/protocol"
	"github.com/status-im/status-go/wakuv2"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Test sending a message",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("added task: ", cCtx.Args().First())

					privateKey, err := crypto.GenerateKey()
					if err != nil {
						fmt.Println(err)
						return err
					}

					enrBootstrap := "enrtree://AL65EKLJAUXKKPG43HVTML5EFFWEZ7L4LOKTLZCLJASG4DSESQZEC@prod.status.nodes.status.im"
					config := &wakuv2.Config{}
					config.EnableDiscV5 = true
					config.DiscV5BootstrapNodes = []string{enrBootstrap}
					config.DiscoveryLimit = 20
					node, err := wakuv2.New("", "", config, nil, nil, nil, nil, nil)
					if err != nil {
						fmt.Println(err)
						return err
					}

					messenger, err := protocol.NewMessenger(
						"testnode",
						privateKey,
						gethbridge.NewNodeBridge(nil, nil, node),
						uuid.New().String(),
						nil,
						// options...,
					)
					fmt.Println(messenger)
					if err != nil {
						fmt.Println(err)
						return err
					}

					err = messenger.Init()
					if err != nil {
						fmt.Println(err)
						return err
					}

					messenger.Start()

					recipientKey, err := crypto.GenerateKey()
					pkString := hex.EncodeToString(crypto.FromECDSAPub(&recipientKey.PublicKey))
					chat := protocol.CreateOneToOneChat(pkString, &recipientKey.PublicKey, messenger.)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
