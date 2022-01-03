package integration_tests

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"time"
)

func (s *IntegrationTestSuite) TestHappyPath() {
	s.Run("Bring up chain, and test the happy path", func() {
		err := s.approveERC20()
		s.Require().NoError(err, "error approving spending balance for the gravity contract on behalf of the first validator")

		allowance, err := s.getERC20AllowanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), gravityContract)
		s.Require().NoError(err, "error getting allowance of gravity contract spending on behalf of first validator")
		s.Require().Equal(UInt256Max().BigInt(), allowance.BigInt(), "spending allowance not set correctly, got: %s", allowance.String())

		balance, err := s.getEthBalanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address))
		s.Require().NoError(err, "error getting first validator balance")
		s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())

		err = s.sendToCosmos(s.chain.validators[0].keyInfo.GetAddress(), sdk.NewInt(200))
		s.Require().NoError(err, "error sending test denom to cosmos")

		sendToEthereumMsg := types.NewMsgSendToEthereum(
			s.chain.validators[0].keyInfo.GetAddress(),
			s.chain.validators[1].ethereumKey.address,
			sdk.Coin{Denom: "DDS", Amount: sdk.NewInt(100)},
			sdk.Coin{Denom: "DDS", Amount: sdk.NewInt(1)},
		)

		s.Require().Eventuallyf(func() bool {
			val := s.chain.validators[0]
			keyring, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			response, err := s.chain.sendMsgs(*clientCtx, sendToEthereumMsg)
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			if response.Code != 0 {
				if response.Code != 32 {
					s.T().Log(response)
				}
				return false
			}
			return true
		}, 10*time.Second, time.Millisecond, "unable to send to ethereum")

		// Send some test token to cosmos

		// Validate updated balances

		// send some test token across the bridge
	})
}