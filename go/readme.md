# go面经

## 变量逃逸

什么情况西会发生变量逃逸?

1. 变量大小不确定
2. 变量类型不确定例如interface
3. 对象内存超过8kb
4. 函数返回局部变量指针
5. 闭包

## 字符串

go语言字符串实现:

```go
type stringStruct struct {
	str unsafe.Pointer		//字符串首地址，指向底层字节数组的指针
	len int					//字符串长度
}
```

字符串拼接:

```go
func main(){
	a := []string{"a", "b", "c"}
	//方式1：+
	ret := a[0] + a[1] + a[2]
	//方式2：fmt.Sprintf
	ret := fmt.Sprintf("%s%s%s", a[0],a[1],a[2])
	//方式3：strings.Builder
	var sb strings.Builder
	sb.WriteString(a[0])
	sb.WriteString(a[1])
	sb.WriteString(a[2])
	ret := sb.String()
	//方式4：bytes.Buffer
	buf := new(bytes.Buffer)
	buf.Write(a[0])
	buf.Write(a[1])
	buf.Write(a[2])
	ret := buf.String()
	//方式5：strings.Join
	ret := strings.Join(a,"")
}
// 拼接性能,strings.Join ≈ strings.Builder > bytes.Buffer > "+" > fmt.Sprintf
```

## interface

```go
type eface struct {
    _type *_type // 接口的动态类型元数据
    data  unsafe.Pointer // 数据指针
}


type iface struct {
    tab   *itab
    data  unsafe.Pointer
}


type itab struct {
    inter  *interfacetype
    _type  *_type // 接口的动态类型元数据
    hash   uint32 // 类型hash值从动态类型元数据copy
    _      [4]byte
    fun    [1]uintptr // 动态类型实现接口的方法地址
}


type interfacetype struct {
    typ      _type
    pkgpath  name // 包路径
    mhdr     []imethod // 接口方法
}


var rw io.ReadWriter


f, _ := os.Open("eggo.txt")
rw = f
```

interface比较的原理,判断_type和data是否相等,两个interface都等于nil时也相同.

2个nil可能不相等，**两个nil只有在类型相同时才相等**。

 **itab缓存**

Go语言会把用到的itab结构体缓存起来，并且以**<接口类型, 动态类型>**组合为key，以***itab**为value，构造一个哈希表，用于存储与查询itab信息。

## go slice扩容机制

Go <= 1.17

如果当前容量小于1024，则判断所需容量是否大于原来容量2倍，如果大于，当前容量加上所需容量；否则当前容量乘2。

如果当前容量大于1024，则每次按照1.25倍速度递增容量，也就是每次加上cap/4。

## map的底层实现

```go
type hmap struct{
    count int //元素个数，len的返回值
    flags uint2 // 
    B uint8  // buckets数组长度的对数  2^B为buckets个数
    noverflow uint8 // overflow的bucket的近似数
    hash0 uint32   //hash函数
    buckets  unsafe.Pointer //指向buckets数组 如果元素个数为0 即为nil
    oldbuckets unsafe.Pointer // 扩容时候会是old 的两倍
    nevacuate uinptr // 扩容进度 小于此地址的buckerts完成了迁移
    extra *mapextra  //扩充项
    
}

// 一个bucket存储8个键值对,键hash值高8位存储在tophash
type bmap struct {
	tophash [bucketCnt]uint8
}

// bmap在运行时会进行替换runtime.bmap
type bmap struct {
	topbits  [8]uint8
	keys     [8]keytype
	values   [8]valuetype
	pad      uintptr
	overflow uintptr
}

// 采用&运算寻找桶
// golang用链地址法解决hash冲突问题
// golang map扩容机制
// 等量扩容: 当B>=15B大于15时，以15计算，如果溢出桶 >= 2^15次方,触发等量扩容,当B小于15时，以B计算，如果溢出桶 >= 大于2^B次方，触发等量扩容
// 增量扩容: loadfactor>6.5 count/2^B
// golang的扩容使用渐进式扩容
```



## channel实现

```go
type hchan struct {
    qcount   uint           // 数组长度，即已有元素个数
    dataqsiz uint           // 数组容量，即可容纳元素个数
    buf      unsafe.Pointer // 数组地址
    elemsize uint16         // 元素大小
    closed   uint32
    elemtype *_type // 元素类型
    sendx    uint   // 下一次写下标位置
    recvx    uint   // 下一次读下标位置
    recvq    waitq  // 读等待队列
    sendq    waitq  // 写等待队列
    lock     mutex
}

//channel支持交替的读写（称send为写，recv为读，更简洁），有缓冲channel内的缓冲数组会被作为一个“环型”来使用，当下标超过数组容量后会回到第一个位置，所以要分别记录读写下标的位置

// 等待队列是runtime.sudog类型的双向链表，sudog中会记录是哪个协程在等待，等待哪一个channel等等。
// 特别是sudog.elem这个成员，对于recvq中的sudog而言，它代表recv到数据以后要存储到哪里；对于sendq中的sudog而言，它代表要send的数据在哪里
type sudog struct {
    g runtime.g
    elem 
    c hchan
}
```

## GC工作原理

三色标记法加混合写屏障

白色：不确定对象（默认色）；黑色：存活对象。灰色：存活对象，子对象待处理。

一次完整的GC分为四个阶段：

1. 准备标记（需要STW），开启写屏障。
2. 开始标记
3. 标记结束（STW），关闭写屏障
4. 清理（并发）

基于插入写屏障和删除写屏障在结束时需要STW来重新扫描栈，带来性能瓶颈。**混合写屏障**分为以下四步：

1. GC开始时，将栈上的全部对象标记为黑色（不需要二次扫描，无需STW）；
2. GC期间，任何栈上创建的新对象均为黑色
3. 被删除引用的对象标记为灰色
4. 被添加引用的对象标记为灰色

## go内存管理

tcmalloc模型

一些基本概念：
页Page：一块8K大小的内存空间。Go向操作系统申请和释放内存都是以页为单位的。
span : 内存块，一个或多个连续的 page 组成一个 span 。如果把 page 比喻成工人， span 可看成是小队，工人被分成若干个队伍，不同的队伍干不同的活。
sizeclass : 空间规格，每个 span 都带有一个 sizeclass ，标记着该 span 中的 page 应该如何使用。使用上面的比喻，就是 sizeclass 标志着 span 是一个什么样的队伍。
object : 对象，用来存储一个变量数据内存空间，一个 span 在初始化时，会被切割成一堆等大的 object 。假设 object 的大小是 16B ， span 大小是 8K ，那么就会把 span 中的 page 就会被初始化 8K / 16B = 512 个 object 。所谓内存分配，就是分配一个 object 出去。

1.mheap

一开始go从操作系统索取一大块内存作为内存池，并放在一个叫mheap的内存池进行管理，mheap将一整块内存切割为不同的区域，并将一部分内存切割为合适的大小。

mheap.spans ：用来存储 page 和 span 信息，比如一个 span 的起始地址是多少，有几个 page，已使用了多大等等。

2.mcentral

用途相同的span会以链表的形式组织在一起存放在mcentral中。这里用途用**sizeclass**来表示，就是该span存储哪种大小的对象。

找到合适的 span 后，会从中取一个 object 返回给上层使用。

3.**mcache**

为了提高内存并发申请效率，加入缓存层mcache。每一个mcache和处理器P对应。Go申请内存首先从P的mcache中分配，如果没有可用的span再从mcentral中获取。

### Go 可以限制运行时操作系统线程的数量吗？ 常见的goroutine操作函数有哪些？

可以，使用runtime.GOMAXPROCS(num int)可以设置线程数目。该值默认为CPU逻辑核数，如果设的太大，会引起频繁的线程切换，降低性能。

runtime.Gosched()，用于让出CPU时间片，让出当前goroutine的执行权限，调度器安排其它等待的任务运行，并在下次某个时候从该位置恢复执行。
runtime.Goexit()，调用此函数会立即使当前的goroutine的运行终止（终止协程），而其它的goroutine并不会受此影响。runtime.Goexit在终止当前goroutine前会先执行此goroutine的还未执行的defer语句。请注意千万别在主函数调用runtime.Goexit，因为会引发panic。

### 如何控制协程数目

`GOMAXPROCS` 限制的是同时执行用户态 Go 代码的操作系统线程的数量，但是对于被系统调用阻塞的线程数量是没有限制的。`GOMAXPROCS` 的默认值等于 CPU 的逻辑核数，同一时间，一个核只能绑定一个线程，然后运行被调度的协程。因此对于 CPU 密集型的任务，若该值过大，例如设置为 CPU 逻辑核数的 2 倍，会增加线程切换的开销，降低性能。对于 I/O 密集型应用，适当地调大该值，可以提高 I/O 吞吐率

### mutex有几种模式？

```go
// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
	state int32
	sema  uint32
}
// Mutex 结构体有 state 和 sema 两个字段组成，共8字节。其中 state 表示锁的状态，sema 是信号量，用于管理等待队列。

// 使用 Mutex 无需格外初始化，state 默认值为 0，表示处于解锁状态。
```



1. 正常模式

   在正常模式下，等待者 waiter 会进入到一个 FIFO(先进先出) 队列，在获取锁时 waiter 会按照先进先出的顺序获取。

   当唤醒一个 waiter 时它不会立即获取锁，而是要与新来的 goroutine 进行竞争。这种情况下新来的 goroutine 具有优势，因为它已经运行在 CPU 上，而且这种新来的 goroutine 数量可能不止一个，所以唤醒的 waiter 大概率获取不到锁。如果唤醒的 waiter 依旧获取不到锁的情况，那么它会被添加到队列的前面。

   如果 waiter 获取不到锁的时间超出了1 毫秒，当前状态将被切换为饥饿模式。

   新来一个goroutine 时会尝试一次获取锁，如果获取不到我们就视其为watier，并将其添加到FIFO队列里。

2. 在正常模式下，每次新来的 goroutine 都可能会抢走锁，就这会导致等待队列中的 waiter 可能永远也获取不到锁，从而产生饥饿问题。所以，为了应对高并发抢锁场景下的公平性，官方引入了饥饿模式。

   在饥饿模式下，锁会直接交给队列最前面的 waiter。新来的 goroutine 即使在锁未被持有情况下，也不会参与竞争锁，同时也不会进行自旋等待，而直接将添加到队列的尾部。

   如果拥有锁的 waiter 发现有以下两种情况，它将切换回正常模式：

   - 它是等待队列里的最后一个 waiter，再也没有其它 waiter
   - 它等待的时间小于1毫秒

## gmp调度模型

go进行调度过程：

- 某个线程尝试创建一个新的G，那么这个G就会被安排到这个线程的G本地队列LRQ中，如果LRQ满了，就会分配到全局队列GRQ中；
- 尝试获取当前线程的M，如果无法获取，就会从空闲的M列表中找一个，如果空闲列表也没有，那么就创建一个M，然后绑定G与P运行。
- 进入调度循环：
  - 找到一个合适的G
  - 执行G，完成以后退出

### Go什么时候发生阻塞？阻塞时，调度器会怎么做。

- **channel阻塞**：当goroutine读写channel发生阻塞时，会调用gopark函数，该G脱离当前的M和P，调度器将新的G放入当前M。
- **系统调用**：当某个G由于系统调用陷入内核态，该P就会脱离当前M，此时P会更新自己的状态为Psyscall，M与G相互绑定，进行系统调用。结束以后，若该P状态还是Psyscall，则直接关联该M和G，否则使用闲置的处理器处理该G。
- **系统监控**：当某个G在P上运行的时间超过10ms时候，或者P处于Psyscall状态过长等情况就会调用retake函数，触发新的调度。
- **主动让出**：由于是协作式调度，该G会主动让出当前的P（通过GoSched），更新状态为Grunnable，该P会调度队列中的G运行。

### Go中GMP有哪些状态？

G的状态：

**_Gidle**：刚刚被分配并且还没有被初始化，值为0，为创建goroutine后的默认值

**_Grunnable**： 没有执行代码，没有栈的所有权，存储在运行队列中，可能在某个P的本地队列或全局队列中(如上图)。

**_Grunning**： 正在执行代码的goroutine，拥有栈的所有权(如上图)。

**_Gsyscall**：正在执行系统调用，拥有栈的所有权，与P脱离，但是与某个M绑定，会在调用结束后被分配到运行队列(如上图)。

**_Gwaiting**：被阻塞的goroutine，阻塞在某个channel的发送或者接收队列(如上图)。

**_Gdead**： 当前goroutine未被使用，没有执行代码，可能有分配的栈，分布在空闲列表gFree，可能是一个刚刚初始化的goroutine，也可能是执行了goexit退出的goroutine(如上图)。

**_Gcopystac**：栈正在被拷贝，没有执行代码，不在运行队列上，执行权在

**_Gscan** ： GC 正在扫描栈空间，没有执行代码，可以与其他状态同时存在。

P的状态：

**_Pidle** ：处理器没有运行用户代码或者调度器，被空闲队列或者改变其状态的结构持有，运行队列为空

**_Prunning** ：被线程 M 持有，并且正在执行用户代码或者调度器(如上图)

**_Psyscall**：没有执行用户代码，当前线程陷入系统调用(如上图)

**_Pgcstop** ：被线程 M 持有，当前处理器由于垃圾回收被停止

**_Pdead** ：当前处理器已经不被使用

M的状态：

**自旋线程**：处于运行状态但是没有可执行goroutine的线程，数量最多为GOMAXPROC，若是数量大于GOMAXPROC就会进入休眠。

**非自旋线程**：处于运行状态有可执行goroutine的线程。



### 如果有一个G一直占用资源怎么办？什么是work stealing算法？

如果有个goroutine一直占用资源，那么GMP模型会**从正常模式转变为饥饿模式**（类似于mutex），允许其它goroutine使用work stealing抢占（禁用自旋锁）。

work stealing算法指，一个线程如果处于空闲状态，则帮其它正在忙的线程分担压力，从全局队列取一个G任务来执行，可以极大提高执行效率。

## atomic的原理

atomic源码位于`sync\atomic`。通过阅读源码可知，atomic采用**CAS**（CompareAndSwap）的方式实现的。所谓CAS就是使用了CPU中的原子性操作。在操作共享变量的时候，CAS不需要对其进行加锁，而是通过类似于乐观锁的方式进行检测，总是假设被操作的值未曾改变（即与旧值相等），并一旦确认这个假设的真实性就立即进行值替换。本质上是**不断占用CPU资源来避免加锁的开销**。

## select的实现原理

select源码位于`src\runtime\select.go`，最重要的`scase` 数据结构为：

```go
type scase struct {
	c    *hchan         // chan
	elem unsafe.Pointer // data element
}
```

scase.c为当前case语句所操作的channel指针，这也说明了一个case语句只能操作一个channel。

scase.elem表示缓冲区地址：

- caseRecv ： scase.elem表示读出channel的数据存放地址；
- caseSend ： scase.elem表示将要写入channel的数据存放地址；

select的主要实现位于：`select.go`函数：其主要功能如下：

\1. 锁定scase语句中所有的channel

\2. 按照随机顺序检测scase中的channel是否ready

2.1 如果case可读，则读取channel中数据，解锁所有的channel，然后返回(case index, true)

2.2 如果case可写，则将数据写入channel，解锁所有的channel，然后返回(case index, false)

2.3 所有case都未ready，则解锁所有的channel，然后返回（default index, false）

\3. 所有case都未ready，且没有default语句

3.1 将当前协程加入到所有channel的等待队列

3.2 当将协程转入阻塞，等待被唤醒

\4. 唤醒后返回channel对应的case index

4.1 如果是读操作，解锁所有的channel，然后返回(case index, true)

4.2 如果是写操作，解锁所有的channel，然后返回(case index, false)

## context包的作用

`context`可以用来在`goroutine`之间传递上下文信息，相同的`context`可以传递给运行在不同`goroutine`中的函数，上下文对于多个`goroutine`同时使用是安全的，`context`包定义了上下文类型，可以使用`background`、`TODO`创建一个上下文，在函数调用链之间传播`context`，也可以使用`WithDeadline`、`WithTimeout`、`WithCancel` 或 `WithValue` 创建的修改副本替换它，听起来有点绕，其实总结起就是一句话：**`context`的作用就是在不同的`goroutine`之间同步请求特定的数据、取消信号以及处理请求的截止日期**。