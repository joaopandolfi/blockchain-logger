package migrations

import (
	"fmt"
	"logger/remotes/blockchain"
	"logger/remotes/postgres"

	"github.com/google/uuid"
	"github.com/joaopandolfi/blackwhale/utils"
)

func Migrate() {
	if postgres.Driver() == nil {
		utils.CriticalError("[Migrations] - Database is not connected")
		return
	}
	postgres.Driver().AutoMigrate(
		&blockchain.Block{},
	)
}

func Terraform() error {
	if postgres.Driver() == nil {
		utils.CriticalError("[Terraform] - Database is not connected")
		return fmt.Errorf("[Terraform] - Database is not connected")
	}

	var genesisBlock blockchain.Block

	postgres.Driver().First(&genesisBlock)
	if genesisBlock.ID == uuid.Nil {

		err := blockchain.Get().GenerateGenesis()
		if err != nil {
			msg := "[Terraform] - Creating genesis block"
			utils.CriticalError(msg, err.Error())
			return fmt.Errorf("%s : %w", msg, err)
		}

		genesisBlock = blockchain.Get().GenesisBlock

		tx := postgres.Driver().Create(&genesisBlock)

		if tx.Error != nil {
			return fmt.Errorf("terraforming genesis block: %w", tx.Error)
		}
	}

	return nil
}
