## hyperledger国密改造的方法?

参考同济区块链研究院与twgc国米改造方案

```
* bccsp/bccsp.go 定义了一系列接口:
* bccsp.Key : 对称加密的密钥或签名/非对称加密的公私钥
* bccsp.KeyGenOpts : bccsp.Key生成器配置
* bccsp.KeyDerivOpts : bccsp.Key驱动配置
* bccsp.KeyImportOpts : bccsp.Key导入配置
* bccsp.HashOpts : bccsp摘要算法配置
* bccsp.SignerOpts : bccsp签名器配置
* bccsp.EncrypterOpts : bccsp对称加密器配置
* bccsp.DecrypterOpts : bccsp对称解密器配置
* bccsp.BCCSP : bccsp接口，提供完整的密钥生成、对称加解密、签名/验签、摘要等功能
```