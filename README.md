# irishub-load

This simple tool is used to execute stress test for irishub.

## How to use

#### No.1  Create test accounts and send test-iris to these accounts

1) Copy config.json to $HOME (Set the parameters if necessary)
2) ./irishub-load init --config-dir=$HOME

#### No.2 Sign about tps * duration * 60 TXs, to avoid Sequence Conflict we use 4 different accounts (wenxi silei jiacheng haoyang)

./irishub-load signtx --config-dir=$HOME --tps=100 --duration=1 --account=wenxi

#### No.3 Broadcast tps * interval TXs for every interval seconds

./irishub-load broadcast --config-dir=$HOME --tps=50 --interval=5

## How to compile
1) dep ensure -v
2) go install
