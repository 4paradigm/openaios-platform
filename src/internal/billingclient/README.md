# Billing Server Client

我们通过openapi生成了billing server client端的代码，但为了更方便的使用，我们在它的基础上包装了一层，主要函数的功能如下：

- GetBillingClient：创建billing server的client
- InitUserBillingAccount：初始化用户的账户，目前初始金额是hardcode到代码里的，之后可以移到数据库
- GetUserBalance：获取用户的余额
- GetOneComputeUnit：获取一个算力规格的具体内容
- GetComputeUnitListByUserID：获取一个用户可以使用的所有算力规格
- GetComputeUnitListByGroupName：获取指定group下所有的算力规格
- GetComputeUnitPrice：获取指定算力规格的价格 
- AddComputeunitGroupToUser：给用户添加算力规格的group

目前这个工具包会被webserver以及webhook使用。