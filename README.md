目录结构：
client：客户端，每隔10ms向node1和node2发送交易数据
initPackage：初始包，用于创建区块链和创世块
node1：节点1，接收客户端发来的交易数据、初始包发来的区块链数据以及其他节点发来的区块数据。同时负责挖矿、验证、打包区块等
node2：痛node1功能一致
someCompose：实现以上包中各个功能所需的函数，例如：验证区块合法、添加区块等。
以上四个包中的main文件都需要运行。
