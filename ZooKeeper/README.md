## [ZooKeeper](https://pdos.csail.mit.edu/6.824/papers/zookeeper.pdf)阅读笔记
#### 1.Introduction
##### Features
* **Wait-free**: no blocking primitives
* **FIFO client ordering**: enables clients to submit operations asynchronously
* **Linearizable writes**: enable efficient implementation of the service
##### Support an API that allow developers to implement their own primitives
> When designing our coordination service, we moved away from implementing specific primitives on the server side, and instead we opted for exposing an API that enables application developers to implement their own primitives.

#### 2.The ZooKeeper Service
##### 2.1 Service Overview
* *Znode* is similar to tree nodes. Zookeeper has a set of znodes to store data.
  - Zookeeper use standard UNIX notation to refer a znode.
  - Every znode can have a child, every znode must have a parent.
  - Znode has 2 types: regular and ephemeral(暂时的). Ephemeral znodes can be removed manually or automatically after the session create it terminates.
  - Sequential flag: if a znode is created with sequential flag set, then its sequence number must be no smaller than any other nodes that are already created under the same parent.
  - Watch flag: if a operation is created with watch flag set, clients who create the operation will receive timely notifications of changes without requiring polling
##### 2.2 Client API
* All client API has a synchronized version and a unsynchronized version. The unsynchronized version uses a callback.
* All update operation expected a version number, and version number not match will result in failure. If version number is -1, there is no version number check.
##### 2.3 Zookeeper Guarantees
* **Asynchronous Linearizable writes**: all requests that update the state of ZooKeeper are serializable and respect precedence;
* **FIFO client order**: all requests from a given client are executed in the order that they were sent by the client
* Sync API is similar to flush, sync causes a server to apply all pending write requests before processing the read without the overhead of a full write.
##### 2.4 Examples of Primitives
* ZooKeeper’s **ordering guarantees** allow efficient reasoning about system state, and **watches(the watch flag)** allow for efficient waiting.
* Configuration Management（配置管理）：创建名为z<sub>c</sub>的znode，每个从z<sub>c</sub>中获取配置信息的进程都会将watch flag设置为true。因此任何的对配置的更新都会被所有使用该配置的进程发现。每次更新导致原有的watch flag被删除以后，重新将watch flag设置为true。
* Rendezvous（）：用户创建一个名为z<sub>r</sub>的Rendezvous Znode，并在其中存储希望给别人共享的信息，保留完整的z<sub>r</sub>路径名称传递给别人即可实现消息的传递。
* Group Membership（）：用户创建一个名为z<sub>g</sub>的ephemeral znode代表组。利用ephemaral可以获取创建它的session的当前状态这个特性，每当有新的成员加入该组时，在z<sub>g</sub>下创建一个新的ephemeral znode，并储存该成员的相关信息。可以遍历这个组中的所有成员获取成员名单，也可以设置一个watch flag获取任何有关于组的更新。
* Simple Locks（）：一种无阻塞的分布式锁的实现方式，创建一个znode即为加锁，其他成员利用**exist** API获取是否已经上锁，利用watch flag来第一时间获取解锁通知。有多种写法来实现如读写锁之类的性质。

TODO：根据源码写出ephemeral node如何获取创建它的session的当前状态。

#### 4.Zookeeper Implementation
* Read Requests: serviced from the local replica of each server database
* Write Requests: processed by anagreement protocol.
* Request Processor: calculate the state(也许是版本号) after the write request is applied, and generate write transaction if the version number matches. Otherwise generate error transaction.
* Atomic Broadcast: use an atomic broadcast protocol Zab. Zab guarantees that changes broadcast by a leader are delivered in the order they were sent and all changes from previous leaders are delivered to an established leader before it broadcasts its own changes.
* Replicated Database: replicated database is an in-memory database (on each server) containing the entire data tree. For recoverability, we have a replay log (a write-ahead log, in our case) of committed operations and generate periodic snapshots of the in-memory database.
  - We do not lock the ZooKeeper state to take the snapshot; instead, we do a depth first scan of the tree atomically reading each znode’s data and meta-data and writing them to disk.
* Client-server interactions:
  - Global gaurantee of ```sync``` API:
  - > In order for this to work, the follower must be sure that the leader is still the leader. If there are pending transactions that commit, then the server does not suspect the leader. If the pending queue is empty, the leader needs to issue a null transaction to commit and orders the sync after that transaction. This has the nice property that when the leader is under load, no extra broadcast traffic is generated. In our implementation, timeouts are set such that leaders realize they are not leaders before followers abandon them, so we do not issue the null transaction.

TODO: 查明zookeeper如何实现idempotent：版本号控制，当replay的transaction的版本号比snapshot的版本号小时，不做任何操作
