# my-token
1. 20230418 同步
2. 根据个人需要，对[wallet-SDK](https://github.com/coming-chat/wallet-SDK.git) 做了精简和调整
3. 整体还在逐步开发中，后续框架结构可能会有较大变动。但已实现的钱包rpc可供参考使用，
4. 部分比较老，不再被维护的库，自行做了些调整后，手动添加到`core/lib`中了
5. 该项目宗旨是：尽可能减少外部依赖，独立单项目完成各钱包的集成

## 1. 当前进度
### 1.1 SUI
1. 导出keystore
2. 导入keystore
3. 助记词生成账户
4. 私钥生成账户
5. 铸造NFT（对接逻辑已完成，具体需要考虑合约进一步完善）
6. 获取NFT（同上）
7. 质押
8. 验证人节点信息查询、验证人节点数据等
9. Object基本操作，包括transfer
10. Token查询以及Token转账
11. 交易签名等
12. 预估交易手续费

### 1.2 Polkadot
1. keystore导入和导出
2. 密钥、助记词生成账户
3. 支持平行链转账
4. 账户余额查询
5. 交易生成和签名
6. 预估交易手续费


## 感谢
> 1. [wallet-SDK](https://github.com/coming-chat/wallet-SDK.git)


