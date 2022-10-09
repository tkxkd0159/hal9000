# Tester Bot Types

* **pusher** : block reward를 auth-claim하고 host -> controller로 자동으로 돈 넘겨주는 테스트용 봇
* **minter** : 넘어온 자산 balance 확인하고 deposit. 그다음 특정 시점 이후 snatom을 lazy minting하는 테스트 용 봇
* **taker** : 특정 시점 이후 무작위 물량을 withdraw해서 snAsset을 wAsset으로 교환하는 테스트 용 봇
  * taker를 minter와 통합할 수도 있음

# Cmd
```shell
# novad도 동일하게
gaiad keys add <keyname> --recover
gaiad config node <host>:26657 # tcp://localhost:26657
gaiad config trust-node true
gaiad config chain-id <zone> # default chain-id가 gaia / nova 라 별도 설정 x

gaiad query account <account_cosmos> # 제대로 연결된 지 확인
```