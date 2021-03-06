# interview

### 操作系统
#### 一些专有名词
* POSIX表示可移植操作系统接口(Portable Operating System Interface of UNIX，缩写为 POSIX)。POSIX标准定义了操作系统应该为应用程序提供的接口标准，是IEEE为要在各种UNIX操作系统上运行的软件而定义的一系列API标准的总称。
* XSI代表的是Unix System V，是unix最早的商业化版本，由AT&T开发。

#### Linux文件权限
* 第一个字母代表文件类型，如字母“d”代表目录(directory)，字母“p“代表管道(pipeline)，字母“s”代表套接字(socket)，字符‘-’代表普通文件，如二进制文件等。

#### 线程和进程
##### 线程(Thread)：
* 管理：理论上来说，每个线程由线程控制块(Thread Control Block)统一管理。每个线程有属于自己的程序计数器(Program Counter)，堆(heap)，栈(stack)和指令集(instruction set)。
* 交流：由于线程共享内存空间，他们的交流可以用最简单的修改内存中的变量来完成。在某些特殊情况下，如需要加锁时，采用信号量(semaphore)等方法来完成。
  - 互斥锁(mutex lock)从本质上来看是一种特殊情况下的信号量，因此可以看作是同一种东西。
* 切换：线程的切换不涉及虚拟地址空间(virtual memory address)的切换，因为同一进程下的不同线程**共享**同一片内存空间。
  
##### 进程(Process)：
* 管理：理论上来说，每个进程由进程控制块(Process Control Block)统一管理。进程控制块是一种数据结构，包含进程状态(Running，Waiting等状态)，程序计数器，CPU寄存器等数据。
* 交流：进程间的交流有多种方式，其中有代表性的有XSI-IPC(InterProcess Communication)，其中包含消息队列(message queue)，信号量(semaphore)与共享内存(shared memory)。进程中同样也有其他通讯方式也有其他方式如信号(signal)，管道(pipeline)和套接字(socket)。以下展开写一些需要注意的通信方式：
  - XSI-IPC:消息队列：父子进程利用同一个键值(shared memory key)来获取标识符(identifier,即ID)，使用的C标准库为```sys/msg.h```。
  - XSI-IPC:信号量：与POSIX定义的信号量类似，两者都可用于线程之间的通信。使用的C标准库为```sys/sem.h```。
  - XSI-IPC:共享内存：父子进程利用同一个键值(shared memory key)来获取标识符(identifier,即ID)，使用的C标准库为```sys/shm.h```。
  - 管道：分为有名管道和无名管道。
    * 有名管道：可以用Linux系统命令```mknod```或命令```mkfifo```创建， 也可以在c里面使用同名的系统函数(function)创建。有名管道在无人读取或无人写入时会造成阻塞。
    * 无名管道：最常见的无名管道为shell中常用的符号``` | ```， 在父进程运行完毕后，使用```fork()```将父进程的地址拷贝到子进程。需要注意的是无名管道中‘-’符号的应用，例子```tar -cvf - /home | tar -xvf -``` 命令表示了将```/home```目录打包再解包的过程，在这个过程中‘-’符号代表输出和输入文件。
  - 套接字：可以用于互联网中不同电脑之间的交流，同样也可以用于同一台电脑上不同进程之间的通信。它读取三个变量(variable)，分别是domain，type 和 protocol。
* 切换：在x86架构下，进程之间的切换需要参考寄存器cr3(control register 3)中储存的虚拟地址与实际地址的映射表。切换时会将当前进程的状态储存到当前进程的进程控制块中，并读取下一个进程的控制块。
* TLB：理论上来说，进程的切换需要清空cpu中原有的TLB，因为不同进程的虚拟地址可能相同，但对应的物理地址却不同。为了提高效率，可以不清除kernel space的TLB，因为即使是不同进程的kernel space也是不冲突的，同时还有一种叫做ASID(address space ID)的技术给不同的进程一个不同的ID存在TLB中，这样就不需要任何的清除操作，因为即使是不同进程的同一个虚拟地址也会有区别，不会影响寻址。
  
##### 为什么要设计进程与线程两种概念？
> 操作系统模型中，进程有两个功能： 
> 1. 进程是任务的调度执行基本单位
> 2. 进程代表了某一种资源的所有权
> 
> 线程的出现就是将这两个功能分离开来了：线程负责执行任务进程的任务，是任务的调度和执行基本单位；进程决定了某一些资源所有权，这样的好处是：
> 
> 操作系统中有两个重要概念：并发和隔离。在分离线程和进程后，线程和并发有关系，进程和隔离有关系.
> * 并发：提高硬件利用率，进程的上下文切换比线程的上下文切换效率低，所以线程可以提高并发的效率
> * 隔离：计算机的资源是共享的，当程序发生奔溃时，需要保证这些资源要被回收，进程的资源是独立的，奔溃时不会影响其他程 序的进行，线程资源是共享的，奔溃时整个进程也会奔溃
  
##### Unix和windows下的进程与线程：
* Unix/Linux系统中没有明确区分线程与进程的定义，两者都被task_struct所指代。在Linux2.4版本后才逐渐出现了线程的定义，但它本质上依然使用了数据结构task_struct。创建进程时，linux使用的底层系统调用(system call)为```fork()```方法。
* Windows系统中有明确区分的线程与进程。创建进程时，windows系统的系统调用为```creatprocess()```方法。

##### 线程切换和进程切换的区别：[文章链接](https://www.cnblogs.com/lfri/p/12597297.html)
##### 为什么协程切换代价比进程切换低：[知乎链接](https://www.zhihu.com/question/308641794)

##### I/O 模式
* [各类型I/O的简单介绍](https://www.jianshu.com/p/786d351e06d8)
* [Linux IO模式及 select、poll、epoll详解](https://segmentfault.com/a/1190000003063859)
* Epoll的LT和ET分别是什么：

##### 协程(co-routine)的定义及用法：[文章链接](https://medium.com/@alex_random/%E5%AF%B9%E4%BA%8E%E5%8D%8F%E7%A8%8B-coroutine-%E4%BD%A0%E5%BF%85%E9%A1%BB%E7%9F%A5%E9%81%93%E4%BA%8B%E6%83%85%E9%83%BD%E5%9C%A8%E8%BF%99%E9%87%8C%E4%BA%86-%E5%86%85%E9%99%84%E4%BB%A3%E7%A0%81-125d07739758)


### 计算机网络
#### TCP/IP
##### 三次握手四次挥手
* TIME_WAIT状态：主动关闭方在收到被动关闭方的FIN包后并返回ACK后，会进入TIME_WAIT状态，期间端口不会被复用，这个状态的持续时间为```2 * MSL```。
  - 一般情况下，MSL的规定时间为2分钟，而一个web服务器最大的端口数为65535个，如果这个服务器作为客户端不停的和服务端不停的创建短连接，就会导致有大量的TCP进入TIME_WAIT状态。
  - ```SO_REUSEADDR```设置为1在TIME_WAIT时允许套接字端口复用；设置为0TIME_WAIT时不允许允许套接字端口复用。可以解决以上问题
* [TCP和HTTP的Keep-alive代表了什么](https://zhuanlan.zhihu.com/p/224595048)
##### 快重传和快恢复
* 快重传指的是不等待timeout，而是在收到三次丢失某个包裹的ack后，直接进行重传
* 快恢复指的是在出现超时重传，或者收到第三个重复的 ack 时（快重传），TCP 会把慢启动门限 ssthresh 的值设置为 cwnd 值的一半，同时 cwnd = ssthresh。
* 一直超时怎么办？在超时重传的情况下，TCP最多会重试7次，每次时间是之前的两倍，如果第一次时间为1，则7次以后需要共127秒才会断开连接
##### TCP连接如何确定使用哪一个文件描述符(file descriptor)
> 1.对于第二个描述符connfd它有没有绑定这个IP与port，是绑定当前端口还是其他端口？这个可以直接通过netstat命令来看。假如最开始bind时的port是7070.使用的是多进程服务器。$netstat -a | grep 7070你会发现，有一个7070端口在监听，有一个7070已经连接建立了。所以connfd就是使用的与listenfd相同的port。
> 
> 2.如果有一个消息来，那么首先由IP地址选定了这台主机（假设IP地址+域名信息，能够确定这台主机），根据port又确定了这个进程，那么server怎么知道这个消息是应该给描述符listenfd还是给connfd的呢？这个问题应该说就有问题，消息来了不是说传给哪个套接字描述符。而是本身信息就是靠TCP连接来传递的。而TCP连接又是通过本地IP+本地port+外地IP+外地port来唯一确定的。那么当我们读取，或者写入某个套接字描述符的时候，它是描述的哪一个TCP连接，那么信息就会在那个连接之间传递。
> 

### 数据库
#### MySQL

##### 数据库中锁的类型
* 共享锁(shared lock)：也叫读锁，简称S锁。
* 排他锁(exclusive lock)：也叫写锁，简称X锁。
> | 兼容性      | S     | X     |
> | ---------- | :-----------:  | :-----------: |
> | S     | T     | F     |
> | X     | F     | F     |

##### MyISAM引擎
* 事务：MyISAM不支持事务，因此没有ACID(Atomicity, Consistency, Isolation, Durability)，隔离级别等概念。
* 锁：只有表级锁，分为读锁和写锁，因此导致MyISAM是一款读性能较好，并发写性能较差的存储引擎。但是MyISAM不会产生死锁，因为MyISAM总是一次性获得所需的全部锁，要么全部满足，要么全部等待。

##### InnoDB引擎
* 事务：InnoDB引擎支持事务，因此也具有隔离级别这个概念，会在下面一一介绍。
* InnoDB中事务的隔离级别：
  - 读未提交(READ UNCOMMITTED)：在读未提交隔离级别下，事务A可以读取到事务B修改过但未提交的数据。可能出现脏读(dirty read)，不可重复读(non-repeatable read)以及幻读(phantom)问题。实现方式为：事务对当前被读取的数据不加锁；但事务在更新某数据的瞬间（就是发生更新的瞬间），必须先对其加行级共享锁，直到事务结束才释放。
  - 读已提交(READ COMMITTED)：在读已提交隔离级别下，事务A只能在事务B修改过并且已提交后才能读取到事务B修改的数据。可能出现不可重复读(non-repeatable read)以及幻读(phantom)问题。实现方式为：事务对当前被读取的数据加行级共享锁（当读到时才加锁），一旦读完该行，立即释放该行级共享锁；但事务在更新某数据的瞬间（就是发生更新的瞬间），必须先对其加行级排他锁，直到事务结束才释放。
  - 可重复读(REPEATABLE READ)：在可重复读隔离级别下，事务A只能在事务B修改过数据并提交后，且自己也提交事务后，才能读取到事务B修改的数据。可能出现幻读(phantom)问题。实现方式为：事务在读取某数据的瞬间（就是开始读取的瞬间），必须先对其加行级共享锁，直到事务结束才释放；但事务在更新某数据的瞬间（就是发生更新的瞬间），必须先对其加行级排他锁，直到事务结束才释放。
  - 可串行化(SERIALIZABLE)：在这种隔离级别下，如果事务A和事务B同时对一张表进行读写/写读/写写操作，则一方会阻塞直到另一方全部完成。不会出现以上问题。实现方式为：事务在读取数据时，必须先对其加表级共享锁，直到事务结束才释放；但事务在更新数据时，必须先对其加表级排他锁，直到事务结束才释放。
  - **注意**：不可重复读的重点是修改，幻读的重点在于新增或者删除。
* InnoDB中的表级锁：
  - 共享锁(shared lock)：表级别的共享锁，也叫读锁，简称S锁。
  - 排他锁(exclusive lock)：表级别的排他锁，也叫写锁，简称X锁。
  - 意向共享锁(intention shared lock)：表级别的意向共享锁，简称IS锁。
  - 意向排他锁(intention exclusive lock)：表级别的意向排他锁，简称IX锁。
  > IS、IX锁是表级锁，它们的提出仅仅为了在之后加表级别的S锁和X锁时可以快速判断表中的记录是否被上锁，以避免用遍历的方式来查看表中有没有上锁的记录，也就是说其实IS锁和IX锁是兼容的，IX锁和IX锁是兼容的。下表体现了表级别中各种锁的兼容性：
  > | 兼容性      | S     | X     | IS     | IX     |
  > | ---------- | :-----------:  | :-----------: | :-----------: | :-----------: |
  > | S     | T     | F     | T     | F     |
  > | X     | F     | F     | F     | F     |
  > | IS     | T     | F     | T     | T     |
  > | IX     | F     | F     | T     | T     |
* InnoDB中的行级锁：InnoDB通过给索引项加锁来实现行锁，如果没有索引，则通过隐藏的聚簇索引(clustered index)来对记录加锁。此时加上的锁同样为读锁或写锁，行为与表级别的锁类似。
  - 记录锁(Record Lock)：对单个记录的索引加对应的锁，当别的记录想获取该锁时，遵从兼容性表格给出阻塞或不阻塞的结果。
  - 间隙锁(Gap Lock)：锁定一个范围，但不包括记录本身
  - 组合锁(Next-key Lock)：记录锁和间隙锁的组合，锁定记录本身以及其范围。
  > Next-key Lock的使用场景是什么？
  > 
  > 默认隔离级别REPEATABLE-READ下，InnoDB中行锁默认使用算法Next-Key Lock，只有当查询的索引是唯一索引或主键时，InnoDB会对Next-Key Lock进行优化，将其降级为Record Lock，即仅锁住索引本身，而不是范围。当查询的索引为辅助索引时，InnoDB则会使用Next-Key Lock进行加锁。InnoDB对于辅助索引有特殊的处理，不仅会锁住辅助索引值所在的范围，还会将其下一键值加上Gap LOCK。
  > 
  > 假如我们有表格：
  > | A(primary key) | B(key) |
  > | :---: | :---:|
  > | 1 | 1 |
  > | 3 | 1 |
  > | 5 | 3 |
  > | 7 | 6 |
  > | 10 | 8 |
  > 
  > 然后开启一个会话执行`SELECT * FROM e4 WHERE b=3 FOR UPDATE;`
  > 
  > 因为通过索引b来进行查询，所以InnoDB会使用Next-Key Lock进行加锁，并且索引b是非主键索引，所以还会对主键索引a进行加锁。对于主键索引a，仅仅对值为5的索引加上Record Lock（因为之前的规则）。而对于索引b，需要加上Next-Key Lock索引，锁定的范围是(1,3]。除此之外，还会对其下一个键值加上Gap Lock，即还有一个范围为(3,6)的锁。 大家可以再新开一个会话，执行下面的SQL语句，会发现都会被阻塞。
  > 
  > * 此时如果执行`SELECT * FROM e4 WHERE a = 5 FOR UPDATE;`，会产生阻塞，因为主键a被锁
  > * 此时如果执行`INSERT INTO e4 SELECT 4,2;`，会产生阻塞，因为插入行b的值为2，在锁定的(1,3]范围内
  > * 此时如果执行`INSERT INTO e4 SELECT 6,5;`，会产生阻塞，因为插入行b的值为5，在锁定的(3,6)范围内
  > 
  > InnoDB引擎采用Next-Key Lock来解决幻读问题。因为Next-Key Lock是锁住一个范围，所以就不会产生幻读问题。但是需要注意的是，InnoDB只在Repeatable Read隔离级别下使用该机制。

* MVCC在MySQL的InnoDB中的实现：
> 在InnoDB中，会在每行数据后添加两个额外的隐藏的值来实现MVCC，这两个值一个记录这行数据何时被创建，另外一个记录这行数据何时过期（或者被删除）。 在实际操作中，存储的并不是时间，而是事务的版本号，每开启一个新事务，事务的版本号就会递增。 在可重读Repeatable reads事务隔离级别下：
> 
> * SELECT时，读取创建版本号<=当前事务版本号，删除版本号为空或>当前事务版本号。
> * INSERT时，保存当前事务版本号为行的创建版本号
> * DELETE时，保存当前事务版本号为行的删除版本号
> * UPDATE时，插入一条新纪录，保存当前事务版本号为行创建版本号，同时保存当前事务版本号到原来删除的行
> 
> 通过MVCC，虽然每行记录都需要额外的存储空间，更多的行检查工作以及一些额外的维护工作，但可以减少锁的使用，大多数读操作都不用加锁，读数据操作很简单，性能很好，并且也能保证只会读取到符合标准的行，也只锁住必要行。
  
* InnoDB死锁及解决：死锁原理十分简单，据一些资料显示死锁只会在并发删除时才会出现，有待考证。
* InnoDB是如何利用日志(log)保证ACID特性的：[文章链接](https://zhuanlan.zhihu.com/p/267143540)

##### 聚簇索引（InnoDB使用）和非聚簇索引（MyISAM使用）
> InnoDB使用的是聚簇索引，将主键组织到一棵B+树中，而行数据就储存在叶子节点上，若查找主键，则按照B+树的检索算法即可查找到对应的叶节点，之后获得行数据。若进行条件搜索（非主键，假设辅助索引存在）则需要两个步骤：第一步在辅助索引B+树中检索对应的条件，到达其叶子节点获取对应的主键。第二步使用主键在主索引B+树种再执行一次B+树检索操作，最终到达叶子节点即可获取整行数据。
> 
> 相对的，MyISAM使用的是非聚簇索引，同样使用B+树的结构，MyISAM所有的索引都直接指向叶子结点，因此查找时可以直接找到结果一步到位。
> 
> 看上去聚簇索引的效率明显要低于非聚簇索引，因为每次使用辅助索引检索都要经过两次B+树查找，这不是多此一举吗？聚簇索引的优势在哪？
> 
> 由于行数据和叶子节点存储在一起，同一页中会有多条行数据，访问同一数据页不同行记录时，已经把页加载到了Buffer中，再次访问的时候，会在内存中完成访问，不必访问磁盘。这样主键和行数据是一起被载入内存的，找到叶子节点就可以立刻将行数据返回了，如果按照主键 ID 来组织数据，获得数据更快。
> 
> 辅助索引使用主键作为"指针"而不是使用地址值作为指针的好处是，减少了当出现行移动或者数据页分裂时辅助索引的维护工作，使用主键值当作指针会让辅助索引占用更多的空间，换来的好处是InnoDB在移动行时无须更新辅助索引中的这个"指针"。也就是说行的位置（实现中通过16K的Page来定位）会随着数据库里数据的修改而发生变化（前面的B+树节点分裂以及Page的分裂），使用聚簇索引就可以保证不管这个主键B+树的节点如何变化，辅助索引树都不受影响。

##### MySQL慢查询及解决
详见收藏夹文章。

##### MySQL8.0对于GIS的支持对比PostgreSQL
[腾讯云文章](https://cloud.tencent.com/developer/article/1739424)

#### Redis
##### Redis的持久化储存
* RDB存储当前数据库的snapshot，包括所有的key-value pairs
* AOF存储所有的写入日志，包括delete，update，insert的各种数据
* 在重启redis的时候，优先查看AOF，如果没有再利用RDB还原redis数据库，因为AOF如果采用onesecond更新模式，最多只会损失一秒的记录，而RDB因为异步写入等原因，很可能和原来数据有较大差异

#####

### 数据结构
#### 跳表(skip list): [知乎链接](https://zhuanlan.zhihu.com/p/54869087)


### 编程中常用关键字
#### static的作用：
* [C语言中的作用](https://www.cnblogs.com/stonejin/archive/2011/09/21/2183313.html)
* [JAVA中的作用](https://blog.csdn.net/mawei7510/article/details/83412304)
* [为什么静态方法只能访问其他静态的东西](https://blog.csdn.net/x6696/article/details/80798471?utm_term=%E9%9D%99%E6%80%81%E6%96%B9%E6%B3%95%E4%B8%BA%E4%BB%80%E4%B9%88%E5%8F%AA%E8%83%BD%E8%AE%BF%E9%97%AE%E9%9D%99%E6%80%81%E6%95%B0%E6%8D%AE&utm_medium=distribute.pc_aggpage_search_result.none-task-blog-2~all~sobaiduweb~default-5-80798471&spm=3001.4430)
#### volatile的作用：
* [维基百科的解释](https://zh.wikipedia.org/wiki/Volatile%E5%8F%98%E9%87%8F)
* [JAVA中的作用](https://juejin.cn/post/6844903502418804743)


### Java
#### [Java的内存泄漏](https://blog.csdn.net/m0_38110132/article/details/81986334)
* > 内存泄露，是从操作系统的角度上来阐述的，形象的比喻就是“操作系统可提供给所有进程的存储空间(虚拟内存空间)正在被某个进程榨干”，导致的原因就是程序在运行的时候，会不断地动态开辟的存储空间，这些存储空间在在运行结束之后后并没有被及时释放掉。应用程序在分配了某段内存之后，由于设计的错误，会导致程序失去了对该段内存的控制，造成了内存空间的浪费。
#### [Java的Minor GC会导致STW现象吗](https://www.zhihu.com/question/29114369)
#### [Java的HashMap，HashTable以及ConcurrentHashMap](https://blog.csdn.net/xiaobaixiongxiong/article/details/88574343)


### Golang
#### 为什么Golang中可以返回局部变量(local variable)地址？
* 引用官方文档的回答如下，从文中可以得知，如果编译器(compiler)无法证明局部变量在返回后未被引用，则将局部变量在堆(heap)中初始化。因此我们可以引用局部变量的地址。
> How do I know whether a variable is allocated on the heap or the stack?
>
>From a correctness standpoint, you don't need to know. Each variable in Go exists as long as there are references to it. The storage location chosen by the implementation is irrelevant to the semantics of the language.
>
>The storage location does have an effect on writing efficient programs. When possible, the Go compilers will allocate variables that are local to a function in that function's stack frame. **However, if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors.** Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.
>
>In the current compilers, if a variable has its address taken, that variable is a candidate for allocation on the heap. However, a basic escape analysis recognizes some cases when such variables will not live past the return from the function and can reside on the stack.

#### Go Routine底层实现
* [Golang调度器以及GMP模型](https://learnku.com/articles/41728)

#### Go 垃圾回收总结
* Go的垃圾回收官方形容为 非分代 非紧缩 写屏障 并发标记清理。非分代是golang GC区别于JVM GC分代模型的特点；非紧缩意味着在回收垃圾的过程中，不需要像复制算法那样移动内存中的对象，这样避免STW过长；标记清理算法的字面解释，就是将可达的内存块进行标记mark，最后没有标记的不可达内存块将进行清理sweep；Golang中实现标记功能的算法就是三色标记法，Golang里面三色标记法会造成错标问题，使用写屏障来解决这种问题，而JVM里面的CMS和G1解决错标或者漏标问题的算法分别是Increment Update和SATB。

### C++
* 虚函数表
