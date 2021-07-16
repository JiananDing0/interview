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
  
##### Unix和windows下的进程与线程
* Unix/Linux系统中没有明确区分线程与进程的定义，两者都被task_struct所指代。在Linux2.4版本后才逐渐出现了线程的定义，但它本质上依然使用了数据结构task_struct。创建进程时，linux使用的底层系统调用(system call)为```fork()```方法。
* Windows系统中有明确区分的线程与进程。创建进程时，windows系统的系统调用为```creatprocess()```方法。


### 计算机网络
##### TCP/IP
* 三次握手四次挥手
  - TIME_WAIT状态：主动关闭方在收到被动关闭方的FIN包后并返回ACK后，会进入TIME_WAIT状态，期间端口不会被复用，这个状态的持续时间为```2 * MSL```。
    * 一般情况下，MSL的规定时间为2分钟，而一个web服务器最大的端口数为65535个，如果这个服务器作为客户端不停的和服务端不停的创建短连接，就会导致有大量的TCP进入TIME_WAIT状态。
    * ```SO_REUSEADDR```设置为1在TIME_WAIT时允许套接字端口复用；设置为0TIME_WAIT时不允许允许套接字端口复用。可以解决以上问题

### Golang
##### 为什么Golang中可以返回局部变量(local variable)地址？
* 引用官方文档的回答如下，从文中可以得知，如果编译器(compiler)无法证明局部变量在返回后未被引用，则将局部变量在堆(heap)中初始化。因此我们可以引用局部变量的地址。
> How do I know whether a variable is allocated on the heap or the stack?
>
>From a correctness standpoint, you don't need to know. Each variable in Go exists as long as there are references to it. The storage location chosen by the implementation is irrelevant to the semantics of the language.
>
>The storage location does have an effect on writing efficient programs. When possible, the Go compilers will allocate variables that are local to a function in that function's stack frame. **However, if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors.** Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.
>
>In the current compilers, if a variable has its address taken, that variable is a candidate for allocation on the heap. However, a basic escape analysis recognizes some cases when such variables will not live past the return from the function and can reside on the stack.
