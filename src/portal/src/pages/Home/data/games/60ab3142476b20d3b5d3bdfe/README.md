# [pafka] 快速体验多级存储内核在消息队列中的优化

## 你会获得什么？
本次任务提供了一个基于多级存储内核优化的技术组件案例(消息队列)，你将会在5分钟的时间内，快速体验并了解消息队列组件在数据收集和处理的使用。

## 背景

随着现代计算机体系架构技术的飞速发展，出现了很多先进的存储技术，比如 SSD、非易失性存储等。这使得现代存储呈现出了比以往更为复杂丰富的多级存储架构，同时赋予了存储器件从未有过的新特性（比如具有数据持久性的内存）。MemArk（[memark.io](https://memark.io)） 技术社区依托现代存储硬件架构，研发针对多级存储的优化内核，为大数据和人工智能应用提供构存储优化方案，成功案例包括 Kafka, Redis, RocksDB, Elasticsearch 等。

在本次实验中，我们将评估消息队列系统在多级存储内核上的优化方案。消息队列在人工智能和大数据领域拥有非常广阔的落地场景，例如，世界上财富 100 强公司中的 80% 公司内部都有使用到消息队列系统 Kafka 做数据搜集和处理工作。但是， Kafka 由于其性能瓶颈存在于落盘 IO 上，导致企业必须配备众多的节点来满足整个系统的吞吐需求，从而导致了巨额的硬件成本和运营费用。比如根据美团公开的[技术报告](https://tech.meituan.com/2021/01/14/kafka-ssd.html)，其搭建了超过 6,000台 的 Kafka 集群来满足吞吐需求。

Pafka （https://github.com/4paradigm/pafka） 由 MemArk 社区主导开发的，基于多级存储内核的 Kafka 优化方案。特别的，Pafka 是基于具有傲腾持久内存的多级存储内核的 Kafka 优化方案。持久内存是目前工业界内存架构上的具有革命性的突破技术，其具有大容量、低成本、高性能等特点（[持久内存（PMem）科普与入门](https://memark.io/index.php/2021/04/27/pmem_intro-2/)）。Pafka 利用持久内存的高速持久化的特点，大幅提高单节点的吞吐和降低延迟，其在单节点的吞吐和延迟性能均能相比较普通 Kafka 提升 10 倍以上，因此可以带来集群总投入的 10 倍节省。

## 实验步骤

> 我们已经将 Pafka 以及实验运行环境包装在 Python 脚本下，如果想参照完整的 Pafka 运行方式，请参照我们的 Github：https://github.com/4paradigm/pafka 。本实验将在 Jupyter Notebook 里进行，在算力平台上启动环境以后，按照 Notebook 里设置好的步骤和默认参数逐步运行即可。

本次评估一共将会跑四个实验：基于 Kafka 的生产者和消费者性能实验，以及基于 Pafka 的生产者和消费者实验。<mark>请将最后一个单元格输出的内容截图保存，作为比赛完成依据（参照步骤 8），截图存档并发送至 opensource@4paradigm.com 。</mark>

**注意：由于已经封装在 Notebook，所以以下每一个步骤（对应每一个单元格）只需要点击如下界面上的执行的小箭头即可。请务必按照顺序从第一个单元格开始依次执行。个别单元格可能执行时间稍长（几分钟），请等待执行完成。**

![img](https://ftp.bmp.ovh/imgs/2021/05/db37fcc322ae50ea.png)

正式实验步骤解释如下。

1. 环境初始化，启动 Kafka 

   ![image-20210518141327758](https://ftp.bmp.ovh/imgs/2021/05/f789b88da2a07ba4.png)

   如果启动成功，预期会最后出现 `kafka started` 的输出，如下。
   ![image-20210506141005856](https://ftp.bmp.ovh/imgs/2021/05/6bf4e456cf0e44d7.png)

2. 进行 Kafka 生产者性能实验，命令和预期输出如下，注意观察最后高亮的性能数字，包括吞吐和延迟，稍后和 Pafka 实验进行比较。
 ![image-20210506141005856](https://ftp.bmp.ovh/imgs/2021/05/7cc27e7438e91b06.png)

   最后的预期输出
   ![img](https://ftp.bmp.ovh/imgs/2021/05/d18cad2f8baf4247.png)

3. 进行 Kafka 消费者性能实验，命令和预期输出如下，注意观察最后高亮的性能数字，包括吞吐和延迟，作为和稍后 Pafka 实验进行比较。
   ![img](https://z3.ax1x.com/2021/05/24/gvPfkn.png)

   ![kafka_consumer_res](https://ftp.bmp.ovh/imgs/2021/05/889d2c97767e5692.png)

4. 停止 Kafka，并且启动 Pafka（启动 Pafka 时间可能稍长，请耐心等待），命令和预期输出如下
   ![image-20210506150635234](https://ftp.bmp.ovh/imgs/2021/05/ac266866ad4be8b9.png)

5. 进行 Pafka 生产者性能实验，命令和预期输出如下，注意观察最后高亮的性能数字，包括吞吐和延迟，和之前 Kafka 实验性能进行比较。
   ![img](https://ftp.bmp.ovh/imgs/2021/05/451443261edbca4c.png)

   预期输出如下
   ![img](https://ftp.bmp.ovh/imgs/2021/05/c36c96a1ceb58be6.png)

6. 进行 Pafka 消费者性能实验，命令和预期输出如下，注意观察最后高亮的性能数字，包括吞吐和延迟，和之前 Kafka 实验性能进行比较。
 ![img](https://ftp.bmp.ovh/imgs/2021/05/3b4b9b3e280a7bf0.png)
   ![img](https://ftp.bmp.ovh/imgs/2021/05/d3a3d4fe2f311249.png)

7. 停止 Pafka 实验并且退出
    ![img](https://ftp.bmp.ovh/imgs/2021/05/0e9910c17991855b.png)
8. 运行结果总结和打印，将会输出上面实验 Kafka 和 Pafka 的吞吐和延迟性能对比，预期参考输出如下（注意，因为实验运行时环境的区别，输出的性能数字会略有出入，如果出现和以下参考结果相差很大的数字，请和我们开发人员联系）。<mark>请保存如下面的输出结果，截图存档，作为本次实验完成的依据，截图存档并发送至 opensource@4paradigm.com 。</mark>

   ![img](https://ftp.bmp.ovh/imgs/2021/05/2b1cad1ccdae6049.png)


## 实验结果解读

如上的实验输出结果可以看到，使用基于持久内存的 Pafka ，不论在吞吐还是延迟上，均比 Kafka 有显著的优势。**但注意，由于实验环境的限制，本次实验的性能数字仅具有参考性，但是并不能代表真实生产环境下的数据，主要由以下限制引起。**

1. 仅使用单个生产者和消费者，workload 上并没有完全打满持久内存的带宽
2. 生产者和消费者在同一个物理节点，忽略了网络带宽的影响
3. 同一个物理机上可能存在多个虚拟机，互相之间存在物理资源竞争（比如 CPU 计算资源、持久内存的带宽），影响性能

作为参考，下图（Figure 1）给出在一个标准的环境内（配备 100 Gb 网络的分布式环境），当 workload 可以使得性能逼近硬件瓶颈的时候，在不同存储介质上 Kafka 和 Pafka 的性能优势表现。可以看到，相比较于数据中心常用的 SATA SSD 的硬件配置，在单节点上 Pafka 在吞吐和延迟上均能达到 20x 的性能优势。

![image-20210507095519804](https://ftp.bmp.ovh/imgs/2021/05/1ce7bb0e96c2a2c1.png) *Figure 1. Kafka 和 Pafka 吞吐和延迟的性能对比*

如果我们考虑硬件拥有成本。假设为了整个集群可以达到 20 GB/sec 的消息吞吐，以下图片（Figure 2）给出了两种方案所需要的机器数量和成本比较。可以看到，如果使用基于持久内存的 Pafka，其成本可以下降 10 倍左右。

![image-20210507101848517](https://z3.ax1x.com/2021/05/24/gjVF4s.png)
*Figure 2. Kafka 和 Pafka 系统的硬件成本对比*

如果你对于以上实际环境中 Pafka 能达到的优势感兴趣，可以参考我们的 Github Repo 关于更多的描述（https://github.com/4paradigm/pafka），也欢迎随时联系我们的开发人员（卢冕或者张浩）进行深入的讨论！

## 开发团队和支持

> MemArk ([memark.io](https://memark.io)) 是由第四范式主导的，并由 Intel 等赞助的，推动现代存储架构在企业中落地和价值实现的开放技术社区。MemArk 技术社区专注于研发多级存储架构内核，以及基于内核的优化应用方案，助力企业实现先进存储架构演进，最大化利用现代存储技术价值。

Pafka 是由 MemArk 技术社区研发，如果你对Pafka运行和配置有任何的疑问或者反馈，可以在以下渠道获得更多资料和支持：

-    Github repo: https://github.com/4paradigm/pafka
-    MemArk技术论坛：MemArk 技术论坛下的 “Pafka 使用和开发”专区 https://discuss.memark.io/
-    [MemArk Slack Channels](https://join.slack.com/t/memarkworkspace/shared_invite/zt-o1wa5wqt-euKxFgyrUUrQCqJ4rE0oPw)
-    开发人员邮件：[zhanghao@4paradigm.com](mailto:zhanghao@4paradigm.com); lumian@4paradigm.com
-    MemArk 邮件：contact@memark.io%     


