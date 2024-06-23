### Scala语法

Option语法:

Scala鼓励在变量和函数返回值可能不会引用任何值的时候使用Option类型。Option有两个子类别，Some和None。在没有值的时候，使用None。如果有值可以引用，就使用Some来包含这个值。
当程序回传Some的时候，代表这个函式成功地返回了一个String，此时就可以透过get()函数拿到那个String，如果程序返回的是None，则代表没有字符串可以返回。在返回None时，也就是没有String
的时候，如果还是要调用get()来取得String的话，Scala一样是会抛出一个NoSuchElementException异常。此时可以选用另外一个方法，getOrElse。这个方法在这个Option是Some的实例时返回对
应的值，而在是None的实例时返回传入的参数。换句话说，传入getOrElse的参数实际上是默认返回值。


