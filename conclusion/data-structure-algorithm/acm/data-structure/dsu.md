author: HeRaNO, JuicyMio, Xeonacid, sailordiary, ouuan

![](images/disjoint-set.svg)

## 引入

并查集是一种用于管理元素所属集合的数据结构，实现为一个森林，其中每棵树表示一个集合，树中的节点表示对应集合中的元素。

顾名思义，并查集支持两种操作：

-   合并（Union）：合并两个元素所属集合（合并对应的树）
-   查询（Find）：查询某个元素所属集合（查询对应的树的根节点），这可以用于判断两个元素是否属于同一集合

并查集在经过修改后可以支持单个元素的删除、移动；使用动态开点线段树还可以实现可持久化并查集。

??? warning
    并查集无法以较低复杂度实现集合的分离。

## 初始化

初始时，每个元素都位于一个单独的集合，表示为一棵只有根节点的树。方便起见，我们将根节点的父亲设为自己。

???+ note "实现"
    === "C++"
        ```cpp
        struct dsu {
          vector<size_t> pa;
        
          explicit dsu(size_t size) : pa(size) { iota(pa.begin(), pa.end(), 0); }
        };
        ```
    
    === "Python"
        ```python
        class Dsu:
            def __init__(self, size):
                self.pa = list(range(size))
        ```

## 查询

我们需要沿着树向上移动，直至找到根节点。

![](images/disjoint-set-find.svg)

???+ note "实现"
    === "C++"
        ```cpp
        size_t dsu::find(size_t x) { return pa[x] == x ? x : find(pa[x]); }
        ```
    
    === "Python"
        ```python
        def find(self, x):
            return x if self.pa[x] == x else self.find(self.pa[x])
        ```

### 路径压缩

查询过程中经过的每个元素都属于该集合，我们可以将其直接连到根节点以加快后续查询。

![](images/disjoint-set-compress.svg)

???+ note "实现"
    === "C++"
        ```cpp
        size_t dsu::find(size_t x) { return pa[x] == x ? x : pa[x] = find(pa[x]); }
        ```
    
    === "Python"
        ```python
        def find(self, x):
            if self.pa[x] != x:
                self.pa[x] = self.find(self.pa[x])
            return self.pa[x]
        ```

## 合并

要合并两棵树，我们只需要将一棵树的根节点连到另一棵树的根节点。

![](images/disjoint-set-merge.svg)

???+ note "实现"
    === "C++"
        ```cpp
        void dsu::unite(size_t x, size_t y) { pa[find(x)] = find(y); }
        ```
    
    === "Python"
        ```python
        def union(self, x, y):
            self.pa[self.find(x)] = self.find(y)
        ```

### 启发式合并

合并时，选择哪棵树的根节点作为新树的根节点会影响未来操作的复杂度。我们可以将节点较少或深度较小的树连到另一棵，以免发生退化。

??? note "具体复杂度讨论"
    由于需要我们支持的只有集合的合并、查询操作，当我们需要将两个集合合二为一时，无论将哪一个集合连接到另一个集合的下面，都能得到正确的结果。但不同的连接方法存在时间复杂度的差异。具体来说，如果我们将一棵点数与深度都较小的集合树连接到一棵更大的集合树下，显然相比于另一种连接方案，接下来执行查找操作的用时更小（也会带来更优的最坏时间复杂度）。
    
    当然，我们不总能遇到恰好如上所述的集合——点数与深度都更小。鉴于点数与深度这两个特征都很容易维护，我们常常从中择一，作为估价函数。而无论选择哪一个，时间复杂度都为 $O (m\alpha(m,n))$，具体的证明可参见 References 中引用的论文。
    
    在算法竞赛的实际代码中，即便不使用启发式合并，代码也往往能够在规定时间内完成任务。在 Tarjan 的论文[^tarjan1984worst]中，证明了不使用启发式合并、只使用路径压缩的最坏时间复杂度是 $O (m \log n)$。在姚期智的论文[^yao1985expected]中，证明了不使用启发式合并、只使用路径压缩，在平均情况下，时间复杂度依然是 $O (m\alpha(m,n))$。
    
    如果只使用启发式合并，而不使用路径压缩，时间复杂度为 $O(m\log n)$。由于路径压缩单次合并可能造成大量修改，有时路径压缩并不适合使用。例如，在可持久化并查集、线段树分治 + 并查集中，一般使用只启发式合并的并查集。

按节点数合并的参考实现：

???+ note "实现"
    === "C++"
        ```cpp
        struct dsu {
          vector<size_t> pa, size;
        
          explicit dsu(size_t size_) : pa(size_), size(size_, 1) {
            iota(pa.begin(), pa.end(), 0);
          }
        
          void unite(size_t x, size_t y) {
            x = find(x), y = find(y);
            if (x == y) return;
            if (size[x] < size[y]) swap(x, y);
            pa[y] = x;
            size[x] += size[y];
          }
        };
        ```
    
    === "Python"
        ```python
        class Dsu:
            def __init__(self, size):
                self.pa = list(range(size))
                self.size = [1] * size
        
            def union(self, x, y):
                x, y = self.find(x), self.find(y)
                if x == y:
                    return
                if self.size[x] < self.size[y]:
                    x, y = y, x
                self.pa[y] = x
                self.size[x] += self.size[y]
        ```

## 删除

要删除一个叶子节点，我们可以将其父亲设为自己。为了保证要删除的元素都是叶子，我们可以预先为每个节点制作副本，并将其副本作为父亲。

???+ note "实现"
    === "C++"
        ```cpp
        struct dsu {
          vector<size_t> pa, size;
        
          explicit dsu(size_t size_) : pa(size_ * 2), size(size_ * 2, 1) {
            iota(pa.begin(), pa.begin() + size_, size_);
            iota(pa.begin() + size_, pa.end(), size_);
          }
        
          void erase(size_t x) {
            --size[find(x)];
            pa[x] = x;
          }
        };
        ```
    
    === "Python"
        ```python
        class Dsu:
            def __init__(self, size):
                self.pa = list(range(size, size * 2)) * 2
                self.size = [1] * size * 2
        
            def erase(self, x):
                self.size[self.find(x)] -= 1
                self.pa[x] = x
        ```

## 移动

与删除类似，通过以副本作为父亲，保证要移动的元素都是叶子。

???+ note "实现"
    === "C++"
        ```cpp
        void dsu::move(size_t x, size_t y) {
          auto fx = find(x), fy = find(y);
          if (fx == fy) return;
          pa[x] = fy;
          --size[fx], ++size[fy];
        }
        ```
    
    === "Python"
        ```python
        def move(self, x, y):
            fx, fy = self.find(x), self.find(y)
            if fx == fy:
                return
            self.pa[x] = fy
            self.size[fx] -= 1
            self.size[fy] += 1
        ```

## 复杂度

### 时间复杂度

同时使用路径压缩和启发式合并之后，并查集的每个操作平均时间仅为 $O(\alpha(n))$，其中 $\alpha$ 为阿克曼函数的反函数，其增长极其缓慢，也就是说其单次操作的平均运行时间可以认为是一个很小的常数。

[Ackermann 函数](https://en.wikipedia.org/wiki/Ackermann_function)  $A(m, n)$ 的定义是这样的：

$A(m, n) = \begin{cases}n+1&\text{if }m=0\\A(m-1,1)&\text{if }m>0\text{ and }n=0\\A(m-1,A(m,n-1))&\text{otherwise}\end{cases}$

而反 Ackermann 函数 $\alpha(n)$ 的定义是阿克曼函数的反函数，即为最大的整数 $m$ 使得 $A(m, m) \leqslant n$。

时间复杂度的证明 [在这个页面中](./dsu-complexity.md)。

### 空间复杂度

显然为 $O(n)$。

## 带权并查集

我们还可以在并查集的边上定义某种权值、以及这种权值在路径压缩时产生的运算，从而解决更多的问题。比如对于经典的「NOI2001」食物链，我们可以在边权上维护模 3 意义下的加法群。

## 例题

???+ note "[UVa11987 Almost Union-Find](https://onlinejudge.org/index.php?option=com_onlinejudge&Itemid=8&category=229&page=show_problem&problem=3138)"
    实现类似并查集的数据结构，支持以下操作：
    
    1.  合并两个元素所属集合
    2.  移动单个元素
    3.  查询某个元素所属集合的大小及元素和
    
    ??? note "参考代码"
        === "C++"
            ```cpp
            --8<-- "docs/ds/code/dsu/dsu_1.cpp"
            ```
        
        === "Python"
            ```python
            --8<-- "docs/ds/code/dsu/dsu_1.py"
            ```

## 习题

[「NOI2015」程序自动分析](https://uoj.ac/problem/127)

[「JSOI2008」星球大战](https://www.luogu.com.cn/problem/P1197)

[「NOI2001」食物链](https://www.luogu.com.cn/problem/P2024)

[「NOI2002」银河英雄传说](https://www.luogu.com.cn/problem/P1196)

## 其他应用

[最小生成树算法](../graph/mst.md) 中的 Kruskal 和 [最近公共祖先](../graph/lca.md) 中的 Tarjan 算法是基于并查集的算法。

相关专题见 [并查集应用](../topic/dsu-app.md)。

## 参考资料与拓展阅读

1.  [知乎回答：是否在并查集中真的有二分路径压缩优化？](https://www.zhihu.com/question/28410263/answer/40966441)
2.  Gabow, H. N., & Tarjan, R. E. (1985). A Linear-Time Algorithm for a Special Case of Disjoint Set Union. JOURNAL OF COMPUTER AND SYSTEM SCIENCES, 30, 209-221.[PDF](https://dl.acm.org/doi/pdf/10.1145/800061.808753)

[^tarjan1984worst]: Tarjan, R. E., & Van Leeuwen, J. (1984). Worst-case analysis of set union algorithms. Journal of the ACM (JACM), 31(2), 245-281.[ResearchGate PDF](https://www.researchgate.net/profile/Jan_Van_Leeuwen2/publication/220430653_Worst-case_Analysis_of_Set_Union_Algorithms/links/0a85e53cd28bfdf5eb000000/Worst-case-Analysis-of-Set-Union-Algorithms.pdf)

[^yao1985expected]: Yao, A. C. (1985). On the expected performance of path compression algorithms.[SIAM Journal on Computing, 14(1), 129-133.](https://epubs.siam.org/doi/abs/10.1137/0214010?journalCode=smjcat)

author: orzAtalod

本部分内容转载并修改自 [时间复杂度 - 势能分析浅谈](https://www.luogu.com.cn/blog/Atalod/shi-jian-fu-za-du-shi-neng-fen-xi-qian-tan)，已取得原作者授权同意。

## 定义

### 阿克曼函数

这里，先给出 $\alpha(n)$ 的定义。为了给出这个定义，先给出 $A_k(j)$ 的定义。

定义 $A_k(j)$ 为：

$$
A_k(j)=\left\{
\begin{aligned}
&j+1& &k=0&\\
&A_{k-1}^{(j+1)}(j)& &k\geq1&
\end{aligned}
\right.
$$

即阿克曼函数。

这里，$f^i(x)$ 表示将 $f$ 连续应用在 $x$ 上 $i$ 次，即 $f^0(x)=x$，$f^i(x)=f(f^{i-1}(x))$。

再定义 $\alpha(n)$ 为使得 $A_{\alpha(n)}(1)\geq n$ 的最小整数值。注意，我们之前将它描述为 $A_{\alpha(n)}(\alpha(n))\geq n$，反正他们的增长速度都很慢，值都不超过 4。

### 基础定义

每个节点都有一个 rank。这里的 rank 不是节点个数，而是深度。节点的初始 rank 为 0，在合并的时候，如果两个节点的 rank 不同，则将 rank 小的节点合并到 rank 大的节点上，并且不更新大节点的 rank 值。否则，随机将某个节点合并到另外一个节点上，将根节点的 rank 值 +1。这里根节点的 rank 给出了该树的高度。记 x 的 rank 为 $rnk(x)$，类似的，记 x 的父节点为 $fa(x)$。我们总有 $rnk(x)+1\leq rnk(fa(x))$。

为了定义势函数，需要预先定义一个辅助函数 $level(x)$。其中，$level(x)=\max(k:rnk(fa(x))\geq A_k(rnk(x)))$。当 $rnk(x)\geq1$ 的时候，再定义一个辅助函数 $iter(x)=\max(i:rnk(fa(x))\geq A_{level(x)}^i(rnk(x))$。这些函数定义的 $x$ 都满足 $rnk(x)>0$ 且 $x$ 不是某个树的根。

上面那些定义可能让你有点头晕。再理一下，对于一个 $x$ 和 $fa(x)$，如果 $rnk(x)>0$，总是可以找到一对 $i,k$ 令 $rnk(fa(x))\geq A_k^i(rnk(x))$，而 $level(x)=\max(k)$，在这个前提下，$iter(x)=\max(i)$。$level$ 描述了 $A$ 的最大迭代级数，而 $iter$ 描述了在最大迭代级数时的最大迭代次数。

对于这两个函数，$level(x)$ 总是随着操作的进行而增加或不变，如果 $level(x)$ 不增加，$iter(x)$ 也只会增加或不变。并且，它们总是满足以下两个不等式：

$$
0\leq level(x)<\alpha(n)
$$

$$
1\leq iter(x)\leq rnk(x)
$$

考虑 $level(x)$、$iter(x)$ 和 $A_k^j$ 的定义，这些很容易被证明出来，就留给读者用于熟悉定义了。

定义势能函数 $\Phi(S)=\sum\limits_{x\in S}\Phi(x)$，其中 $S$ 表示一整个并查集，而 $x$ 为并查集中的一个节点。定义 $\Phi(x)$ 为：

$$
\Phi(x)=
\begin{cases}
\alpha(n)\times \mathit{rnk}(x)& \mathit{rnk}(x)=0\ \text{或}\ x\ \text{为某棵树的根节点}\\
(\alpha(n)-\mathit{level}(x))\times \mathit{rnk}(x)-iter(x)& \text{otherwise}
\end{cases}
$$

然后就是通过操作引起的势能变化来证明摊还时间复杂度为 $\Theta(\alpha(n))$ 啦。注意，这里我们讨论的 $union(x,y)$ 操作保证了 $x$ 和 $y$ 都是某个树的根，因此不需要额外执行 $find(x)$ 和 $find(y)$。

可以发现，势能总是个非负数。另，在开始的时候，并查集的势能为 $0$。

## 证明

### union(x,y) 操作

其花费的时间为 $\Theta(1)$，因此我们考虑其引起的势能的变化。

这里，我们假设 $rnk(x)\leq rnk(y)$，即 $x$ 被接到 $y$ 上。这样，势能增加的节点仅有 $x$（从树根变成非树根），$y$（秩可能增加）和操作前 $y$ 的子节点（父节点的秩可能增加）。我们先证明操作前 $y$ 的子节点 $c$ 的势能不可能增加，并且如果减少了，至少减少 $1$。

设操作前 $c$ 的势能为 $\Phi(c)$，操作后为 $\Phi(c')$，这里 $c$ 可以是任意一个 $rnk(c)>0$ 的非根节点，操作可以是任意操作，包括下面的 find 操作。我们分三种情况讨论。

1.  $iter(c)$ 和 $level(c)$ 并未增加。显然有 $\Phi(c)=\Phi(c')$。
2.  $iter(c)$ 增加了，$level(c)$ 并未增加。这里 $iter(c)$ 至少增加一，即 $\Phi(c')\leq \Phi(c)-1$，势能函数减少了，并且至少减少 1。
3.  $level(c)$ 增加了，$iter(c)$ 可能减少。但是由于 $0<iter(c)\leq rnk(c)$，$iter(c)$ 最多减少 $rnk(c)-1$，而 $level(c)$ 至少增加 $1$。由定义 $\Phi(c)=(\alpha(n)-level(c))\times rnk(c)-iter(c)$，可得 $\Phi(c')\leq\Phi(c)-1$。
4.  其他情况。由于 $rnk(c)$ 不变，$rnk(fa(c))$ 不减，所以不存在。

所以，势能增加的节点仅可能是 $x$ 或 $y$。而 $x$ 从树根变成了非树根，如果 $rnk(x)=0$，则一直有 $\Phi(x)=\Phi(x')=0$。否则，一定有 $\alpha(x)\times rnk(x)\geq(\alpha(n)-level(x))\times rnk(x)-iter(x)$。即，$\Phi(x')\leq \Phi(x)$。

因此，唯一势能可能增加的点就是 $y$。而 $y$ 的势能最多增加 $\alpha(n)$。因此，可得 $union$ 操作均摊后的时间复杂度为 $\Theta(\alpha(n))$。

### find(a) 操作

如果查找路径包含 $\Theta(s)$ 个节点，显然其查找的时间复杂度是 $\Theta(s)$。如果由于查找操作，没有节点的势能增加，且至少有 $s-\alpha(n)$ 个节点的势能至少减少 $1$，就可以证明 $find(a)$ 操作的时间复杂度为 $\Theta(\alpha(n))$。为了避免混淆，这里用 $a$ 作为参数，而出现的 $x$ 都是泛指某一个并查集内的结点。

首先证明没有节点的势能增加。很显然，我们在上面证明过所有非根节点的势能不增，而根节点的 $rnk$ 没有改变，所以没有节点的势能增加。

接下来证明至少有 $s-\alpha(n)$ 个节点的势能至少减少 $1$。我们上面证明过了，如果 $level(x)$ 或者 $iter(x)$ 有改变的话，它们的势能至少减少 $1$。所以，只需要证明至少有 $s-\alpha(n)$ 个节点的 $level(x)$ 或者 $iter(x)$ 有改变即可。

回忆一下非根节点势能的定义，$\Phi(x)=(\alpha(n)-level(x))\times rnk(x)-iter(x)$，而 $level(x)$ 和 $iter(x)$ 是使 $rnk(fa(x))\geq A_{level(x)}^{iter(x)}(rnk(x))$ 的最大数。

所以，如果 $root_x$ 代表 $x$ 所处的树的根节点，只需要证明 $rnk(root_x)\geq A_{level(x)}^{iter(x)+1}(rnk(x))$ 就好了。根据 $A_k^i$ 的定义，$A_{level(x)}^{iter(x)+1}(rnk(x))=A_{level(x)}(A_{level(x)}^{iter(x)}(rnk(x)))$。

注意，我们可能会用 $k(x)$ 代表 $level(x)$，$i(x)$ 代表 $iter(x)$ 以避免式子过于冗长。这里，就是 $rnk(root_x)\geq A_{k(x)}(A_{k(x)}^{i(x)}(x))$。

当你看到这的时候，可能会有一种「这啥玩意」的感觉。这意味着你可能需要多看几遍，或者跳过一些内容以后再看。

这里，我们需要一个外接的 $A_{k(x)}$，意味着我们可能需要再找一个点 $y$。令 $y$ 是搜索路径上在 $x$ 之后的满足 $k(y)=k(x)$ 的点，这里「搜索路径之后」相当于「是 $x$ 的祖先」。显然，不是每一个 $x$ 都有这样一个 $y$。很容易证明，没有这样的 $y$ 的 $x$ 不超过 $\alpha(n)-2$ 个。因为只有每个 $k$ 的最后一个 $x$ 和 $a$ 以及 $root_a$ 没有这样的 $y$。

我们再强调一遍 $fa(x)$ 指的是路径压缩 **之前**  $x$ 的父节点，路径压缩 **之后**  $x$ 的父节点一律用 $root_x$ 表示。对于每个存在 $y$ 的 $x$，总是有 $rnk(y)\geq rnk(fa(x))$。同时，我们有 $rnk(fa(x))\geq A_{k(x)}^{i(x)}(rnk(x))$。由于 $k(x)=k(y)$，我们用 $k$ 来统称，即，$rnk(fa(x))\geq A_k^{i(x)}(rnk(x))$。我们需要造一个 $A_k$ 出来，所以我们可以不关注 $iter(y)$ 的值，直接使用弱化版的 $rnk(fa(y))\geq A_k(rnk(y))$。

如果我们将不等式组合起来，神奇的事情就发生了。我们发现，$rnk(fa(y))\geq A_k^{i(x)+1}(rnk(x))$。也就是说，为了从 $rnk(x)$ 迭代到 $rnk(fa(y))$，至少可以迭代 $A_k$ 不少于 $i(x)+1$ 次而不超过 $rnk(fa(y))$。

显然，有 $rnk(root_y)\geq rnk(fa(y))$，且 $rnk(x)$ 在路径压缩时不变。因此，我们可以得到 $rnk(root_x)\geq A_k^{i(x)+1}(rnk(x))$，也就是说 $iter(x)$ 的值至少增加 1，如果 $rnk(x)$ 没有增加，一定是 $level(x)$ 增加了。

所以，$\Phi(x)$ 至少减少了 1。由于这样的 $x$ 节点至少有 $s-\alpha(n)-2$ 个，所以最后 $\Phi(S)$ 至少减少了 $s-\alpha(n)-2$，均摊后的时间复杂度即为 $\Theta(\alpha(n)+2)=\Theta(\alpha(n))$。

## 为何并查集会被卡

这个问题也就是问，如果我们不按秩合并，会有哪些性质被破坏，导致并查集的时间复杂度不能保证为 $\Theta(m\alpha(n))$。

如果我们在合并的时候，$rnk$ 较大的合并到了 $rnk$ 较小的节点上面，我们就将那个 $rnk$ 较小的节点的 $rnk$ 值设为另一个节点的 $rnk$ 值加一。这样，我们就能保证 $rnk(fa(x))\geq rnk(x)+1$，从而不会出现类似于满地 compile error 一样的性质不符合。

显然，如果这样子的话，我们破坏的就是 $union(x,y)$ 函数「y 的势能最多增加 $\alpha(n)$」这一句。

存在一个能使路径压缩并查集时间复杂度降至 $\Omega(m\log_{1+\frac{m}{n}}n)$ 的结构，定义如下：

二项树（实际上和一般的二项树不太一样），其中 j 是常数，$T_k$ 为一个 $T_{k-1}$ 加上一个 $T_{k-j}$ 作为根节点的儿子。

![我们的二项树](./images/dsu_complexity1.png)

边界条件，$T_1$ 到 $T_j$ 都是一个单独的点。

令 $rnk(T_k)=r_k$，这里我们有 $r_k=(k-1)/j$（证明略）。每轮操作，我们将它接到一个单节点上，然后查询底部的 $j$ 个节点。也就是说，我们接到单节点上的时候，单节点的势能提高了 $(k-1)/j+1$。在 $j=\lfloor\frac{m}{n}\rfloor$，$i=\lfloor\log_{j+1}\frac{n}{2}\rfloor$，$k=ij$ 的时候，势能增加量为：

$$
\alpha(n)\times((ij-1)/j+1)=\alpha(n)\times((\lfloor\log_{\lfloor\frac{m}{n}\rfloor+1}\frac{n}{2}\rfloor\times \lfloor\frac{m}{n}\rfloor-1)/\lfloor\frac{m}{n}\rfloor+1)
$$

变换一下，去掉所有的取整符号，就可以得出，势能增加量 $\geq \alpha(n)\times(\log_{1+\frac{m}{n}}n-\frac{n}{m})$，m 次操作就是 $\Omega(m\log_{1+\frac{m}{n}}n-n)=\Omega(m\log_{1+\frac{m}{n}}n)$。

## 关于启发式合并

由于按秩合并比启发式合并难写，所以很多 dalao 会选择使用启发式合并来写并查集。具体来说，则是对每个根都维护一个 $size(x)$，每次将 $size$ 小的合并到大的上面。

所以，启发式合并会不会被卡？

首先，可以从秩参与证明的性质来说明。如果 $size$ 可以代替 $rnk$ 的地位，则可以使用启发式合并。快速总结一下，秩参与证明的性质有以下三条：

1.  每次合并，最多有一个节点的秩上升，而且最多上升 1。
2.  总有 $rnk(fa(x))\geq rnk(x)+1$。
3.  节点的秩不减。

关于第二条和第三条，$siz$ 显然满足，然而第一条不满足，如果将 $x$ 合并到 $y$ 上面，则 $siz(y)$ 会增大 $siz(x)$ 那么多。

所以，可以考虑使用 $\log_2 siz(x)$ 代替 $rnk(x)$。

关于第一条性质，由于节点的 $siz$ 最多翻倍，所以 $\log_2 siz(x)$ 最多上升 1。关于第二三条性质，结论较为显然，这里略去证明。

所以说，如果不想写按秩合并，就写启发式合并好了，时间复杂度仍旧是 $\Theta(m\alpha(n))$。
