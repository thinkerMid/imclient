该目录下都是IFuture的实现

IFuture用于异步处理消息内容的一种形式，类似于某Action发送了某个消息出去，把消息ID给了IFuture，让它监听并处理返回的内容

使用场景有 
    1.某种Action执行场景下导致Action无法处理Receive内容 
    2.Action自身没有Receive内容
    3.Action的Receive处理逻辑唯一化