# Nova

## Flow list
* 3일에 한번씩 Withdraw


## Msg list (Tx)
1. oracle
    1. **MsgUpdateChainState**
2. inter-tx
    1. MsgRegisterZone
    2. MsgChangeRegisteredZoneInfo
    3. MsgDeleteRegisteredZone
    4. **MsgIcaDelegate**
    5. **MsgIcaUndelegate**
    6. **MsgIcaAutoStaking**
    7. **MsgIcaWithdraw**
3. gal
    1. MsgDeposit
    2. MsgUndelegate
    3. MsgUndelegateRecord
    4. MsgWithdrawRecord

## Event
* EmitTypedEvent으로 이벤트 구현 시
  * event type == proto package name + message name (e.g. nova.oracle.v1.ChainInfo)
  * event attribute key : proto field name
  * event attribute value : proto field's value. type is depend on proto field
* 
```sh
# query 시 {eventType}.{eventAttribute}={eventValue}
curl "localhost:26657/tx_search?query=\"message.sender='cosmos1...'\"&prove=true"
```

```json
{
   "jsonrpc": "2.0",
   "method": "subscribe",
   "id": "0",
   "params": {
      "query": "tm.event='eventCategory' AND eventType.eventAttribute='attributeValue'"
   }
}
```
