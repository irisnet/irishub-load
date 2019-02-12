# Mock Test

本项目旨在一次性生成多个已签名的数据，此生成过程可大致分为两步：

1. 初始化一个大账户并使用该账户生成多个子账户，用于后续转账
2. 批量生成账户并生成已签名的交易数据

## Init mock faucet account

**Command**

```bash
mock faucet-init --faucet-name {faucet-name} --seed="recycle light kid ..." --sub-faucet-num {sub-faucet-num} --home {config-home} --chain-id {chain-id} --node {node}
```

**Parameters**

- `faucet-name`：faucet name
- `seed`：faucet seed
- `sub-faucet-num`：num of sub faucet account
- `home`：home for save config 
- `chain-id`：chain id
- `node`：lcd addr

## Gen signed tx data

**Command**

```bash
mock gen-signed-tx --num {num} --receiver {receiver-address} --home {config-home} --chain-id {chain-id} --node {node-url}
```

**Parameters**

- `num`：num of signed tx which need generate
- `receiver`：receiver address
- `home`：home of config file
- `chain-id`：chain id
- `node`：lcd addr