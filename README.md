# interview

### 操作系统
##### Linux文件权限
* 第一个字母代表文件类型，如字母“d”代表目录(directory)，字母“p“代表管道(pipeline)，字母“s”代表套接字(socket)，字符‘-’代表普通文件，如二进制文件等。

##### 进程通信
* 信号(signal)
* 管道(pipeline)：分为有名管道和无名管道。
  - 有名管道：可以用Linux系统命令```mknod```或命令```mkfifo```创建， 也可以在c里面使用同名的系统函数(function)创建。
    * 有名管道在无人读取或无人写入时会造成阻塞。
  - 无名管道：最常见的无名管道为shell中常用的符号``` | ```， 在父进程运行完毕后，使用```fork()```将父进程的地址拷贝到子进程。
    * 无名管道中‘-’符号的应用：```tar -cvf - /home | tar -xvf -``` 命令表示了将```/home```目录打包再解包的过程，中间用无名管道和‘-’符号帮忙。
* 共享内存(shared memory)
* 套接字(socket)：
  - 套接字可以用于互联网中不同电脑之间的交流，同样也可以用于同一台电脑上不同进程之间的通信。
  - 套接字读取三个变量(variable)，分别是domain，type 和 protocol。

##### 线程通信
* 互斥锁(Mutex Lock)
* 信号量(Semaphore), 个人理解本质上信号量包含了互斥锁

##### 线程和进程的区别
* 线程看eip，
* 进程看cr3：x86构造中的control register 3控制了虚拟内存，进程间的切换本质上来说也就是让

### 计算机网络
##### TCP/IP
* 三次握手四次挥手
  - TIME_WAIT状态：主动关闭方在收到被动关闭方的FIN包后并返回ACK后，会进入TIME_WAIT状态，期间端口不会被复用，这个状态的持续时间为```2 * MSL```。
    * 一般情况下，MSL的规定时间为2分钟，而一个web服务器最大的端口数为65535个，如果这个服务器作为客户端不停的和服务端不停的创建短连接，就会导致有大量的TCP进入TIME_WAIT状态。
    * ```SO_REUSEADDR```设置为1在TIME_WAIT时允许套接字端口复用；设置为0TIME_WAIT时不允许允许套接字端口复用。可以解决以上问题

### Golang
##### 为什么Golang中可以返回局部变量(local variable)地址？
* 引用官方文档的回答如下，从文中可以得知，如果编译器(compiler)无法证明局部变量在返回后未被引用，则将局部变量在堆(heap)中初始化。
> How do I know whether a variable is allocated on the heap or the stack?
>
>From a correctness standpoint, you don't need to know. Each variable in Go exists as long as there are references to it. The storage location chosen by the implementation is irrelevant to the semantics of the language.
>
>The storage location does have an effect on writing efficient programs. When possible, the Go compilers will allocate variables that are local to a function in that function's stack frame. **However, if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors.** Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.
>
>In the current compilers, if a variable has its address taken, that variable is a candidate for allocation on the heap. However, a basic escape analysis recognizes some cases when such variables will not live past the return from the function and can reside on the stack.
