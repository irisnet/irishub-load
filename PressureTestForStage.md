# Pressure Test For Stage

本次压力测试旨在 在stage环境下 使用多个账户互相转账，此过程分为两步：

1. 读取一个config.json文件， 生成本节点下的签名交易，存到一个以时间戳命名的文件里面 例如res_signed_tx_1547623712
2. 读取交易数据，进行压力测试

## Init signed data

**Command**

```bash
mock gen-signed-tx-separately --chain-id shilei-qa --home /Users/zjb/.iriscli/ --tps 200 --duration 10 --bots 4 --account-index 0

```

**Parameters**

- `chain-id`：block chain id
- `home`：config.json所在目录
- `tps`：整体的期望tps
- `duration`：测试时间（单位 min） 一开始可以先短一点 比如10min
- `account-index`：account id in the config.json
- `bots`：使用本软件的本次测试的节点数 一般为 4个 或 5个

该指令把测试数据保存在了$HOME/output 文件夹下

## broadcast signed tx data

**Command**

```bash
mock broadcast-signed-tx-separately --output {output-dir} --node {node-url} --tps={max broadcast speed} --duration={duration} --bots={num of test node} --commit={block commit time in config}

mock broadcast-signed-tx-separately --output /Users/zjb/output/res_signed_tx_1547620358 --node http://localhost:1317 --tps 200 --bots 4 --commit 5
```

**Parameters**

- `output`：测试数据的文件路径
- `node`：lcd ip port addr
- `tps`：整体的期望tps
- `duration`：测试时间（单位 min） 一开始可以先短一点 比如10min
- `bots`：使用本软件的本次测试的节点数 一般为 4个 或 5个
- `commit`：是指区块的期望打包时间 一般为5s
