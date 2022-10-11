# Tester Bot Types

* **pusher** : block reward를 auth-claim하고 host -> controller로 자동으로 돈 넘겨주는 테스트용 봇
* **minter** : 넘어온 자산 balance 확인하고 deposit. 그다음 특정 시점 이후 snatom을 lazy minting하는 테스트 용 봇
* **taker** : 특정 시점 이후 무작위 물량을 withdraw해서 snAsset을 wAsset으로 교환하는 테스트 용 봇
  * taker를 minter와 통합할 수도 있음

# Cmd
## 1) Setup
```shell
# novad도 동일하게
gaiad keys add <keyname> --recover
gaiad keys list
gaiad config node <host>:26657 # tcp://localhost:26657
gaiad config trust-node true
gaiad config chain-id <zone> # default chain-id가 gaia / nova 라 별도 설정 x

nbot=nova1lds58drg8lvnaprcue2sqgfvjnz5ljlkq9lsyf
nclient=nova1ma2378jplsnj3chlcgwqfx8p8g92z7dzxdvxmd
gclient=cosmos1ma2378jplsnj3chlcgwqfx8p8g92z7dztqzacq
ghost=cosmos1mn52fx79ckzam60qa4ctgt5meyrk09xc6y6j8adqz3wynww7mczqhj6wur

gaiad query account $gclient # 제대로 연결된 지 확인
```

## 2) Core logics
```shell
# 1. Send uatom to nova
gaiad tx ibc-transfer transfer transfer channel-233 $nclient 987uatom --from $gclient --chain-id gaia

# 2. Check ibc asset on Nova
novad query bank balances $nclient
novad query ibc-transfer denom-trace <ibc-denom>

# 3. Deposit ibc asset on Nova (enter the ibc delegation queue)
novad tx gal deposit gaia $nclient $nclient 80000ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
novad query gal deposit-records gaia $nclient

# 4. Claim snAsset after oracle update
novad query gal claimable gaia $nclient
novad tx gal claim gaia $nclient

# 5. Delegate ibc asset on host chain by bot

# 6. Undelegate asset from nova's validator on specific zone (enter the ibc undelegation queue)
novad tx gal pending-undelegate gaia $nclient $nclient 10000000000000snuatom

# 7. Undelegate ibc asset on host chain by bot
# 8. Withdraw ibc asset on host chain by bot (ibc-transfer from host chain to Nova)

# 9. Withdraw wAsset after waiting for a set delegation period
novad tx gal withdraw gaia $nclient
```