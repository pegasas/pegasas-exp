# Python 内存模型

Python代码在运行是，Python解释器会向操作系统申请运行内存，将代码加载到内存中运行，如图所示：

Python 解释器为了利用好有限的内存空间，将内存进行了如图的划分：

![python-interpreter.png](./python-interpreter.png)

## 不可变类型 – 内存模型

不可变类型：数据在内存中一旦创建，就不能修改了。
Python 为了优化程序执行速度，将字符串、整数定义成了不可变类型，一旦声明出来，数据就不能修改了。

### 字符串操作

字符串是内存中使用特别的多的数据，所以 Python对字符串进行了优化，字符串是不可变数据类型，所以不能直接修改字符串内部的数据。
当我们通过变量修改数据时，内存中将变量指向了一个新的内存地址。原来的字符串数据依然存在，并没有修改。

![python-memory-architecture.png](./python-memory-architecture.png)

### 整数操作

整数和字符串一样，在程序中也是一个经常操作的数据。所以也对整数进行了优化，Python 解释器在加载的时候，将 -5~256 的整数直接在内存中创建好了开发人员要使用的时候直接使用即可，不需要创建对象。
整数也是不可变数据，如果需要修改变量中的整数数据时，就是将变量指向了一个新的内存地址，原来在内存中的数据不会收到影响。

![python-memory-architecture-1.png](./python-memory-architecture-1.png)

## 可变类型 – 内存模型

可变类型就是可以修改数据内部的数据，如列表
Python 中的列表可以存储多个数据，存储的多个数据可能要参与业务处理需要经常变化，所以列表中的数据在语法上被定义成了可以修改的数据。如图所示：

![python-memory-architecture-2.png](./python-memory-architecture-2.png)

## 参考

https://www.cnblogs.com/wsnan/p/15899945.html
