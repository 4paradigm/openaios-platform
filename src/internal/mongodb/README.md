# MongoDB 工具包

因为之前比较好用的mgo库已经无人维护，所以我们使用了mongodb的官方sdk：MongoDB Go Driver（https://github.com/mongodb/mongo-go-driver）
但是因为这个SDK中的函数都较为底层不方便使用，我们自己在它的基础上又封装了一层。

当然，目前来看这个工具包只实现了一些简单的功能，如果有额外需求可以自行添加，下面简单介绍一下每个函数

- GetMongodbClient：创建一个mongodb的client
- KillMongodbClient：回收一个mongodb的client
- CreateIndex：为collection创建一个index（非unique、目前只支持单个key）
- CreateUniqueIndex：为collection创建一个 unique index（可以是多个key组成的对）
- CountDocuments：在collection下对制定条件的document计数
- InsertOneDocument：插入一个document
- DeleteOneDocument：删除一个document
- UpdateOneDocument：更新一个document，支持多个operator，但是每个operator只能出现一次
- FindOneDocuemnt：找到一个document
- FindDoucments：找到满足条件的多个document，目前只支持对单个key的比较
- FindDocumentsByMultiKey：通过多个key找到满足条件的多个documents，目前只支持等于
- UpdateOrInsertOneDocument：更新或者插入一个document，支持多个operator，但是每个operator只能出现一次

我们MongoDB的存储格式可以参考doc/developer的下的database-structure.md文件。