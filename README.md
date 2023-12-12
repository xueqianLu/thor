# Vechain 攻击测试
测试代码共7个分支
- test-base : 基准测试，不包含任何攻击代码，增加了交易发送，部署脚本等内容
- test-case1: 交易丢弃，方式为不向外广播收到的交易
- test-case2: 交易泛洪，方式为广播交易时每个交易对每个peer广播5000次
- test-case3: 区块丢弃，同步到新区块后不向外提供区块
- test-case4: 区块时间戳修改，方式为修改区块的时间戳为当前时间+15秒
- test-case5: 区块重复，方式为在同一个块号同时产生两个不同的区块广播出去
- test-case6: 区块泛洪，方式为广播区块时每个区块对每个peer广播5000次
- test-case7: 区块延迟，方式为生产区块后，延迟10秒再向外广播

每个测试分支的最后一条提交记录为修改的攻击代码.

# 测试
## 0. 测试环境
Ubuntu 或其他linux操作系统，最好是 8核以上CPU，16G以上内存，硬盘空间 40G 以上

## 1. 安装依赖
安装 `git,docker,docker compose` 等工具，具体安装方法请自行搜索.

## 2. 测试方法
在 `test-base` 分支下，执行 `runtest.sh` 即可.

## 3. 测试结果
测试脚本运行结束后，在 `deploy/` 目录下产生每个测试项目的数据目录，其中`data`目录为最后一次测试的数据，其余目录为每次测试的数据。
```
./deploy/
├── Makefile
├── config
├── scripts
├── test.sh
├── testdata_base
├── testdata_case1
├── testdata_case2
├── testdata_case3
├── testdata_case4
├── testdata_case5
├── testdata_case6
└── testdata_case7
```
每个测试数据的目录结构如下, `node.log` 为节点运行日志，`sender.log` 为交易发送程序的运行日志, `node6/report.csv`是本次测试各个收益节点在每个区块的余额情况，也就是收益总额：
```shell
./deploy/testdata_base/
├── bootnode
│   └── instance-28a75abe331c89eb-v3
├── node0
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node1
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node2
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node3
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node4
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node5
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   └── sender.log
├── node6
│   ├── instance-28a75abe331c89eb-v3
│   ├── node.log
│   ├── report.csv
│   └── sender.log
└── query.log
```