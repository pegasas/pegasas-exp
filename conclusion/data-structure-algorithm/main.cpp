#include<iostream>
#include<algorithm>
#include<cstdio>
#include<string.h>
#include<string>
#include<stack>
using namespace std;
typedef long long ll;
const int INF = 0x3f3f3f3f;
const int MAXN = 100000 + 100;
ll a[MAXN];
int L[MAXN], R[MAXN];
int main(){
    int n;
    while (scanf("%d", &n) != EOF && n){
        stack<int> s; s.push(0);
        ll ans = 0;
        for (int i = 1; i <= n + 1; i++){
            if (i <= n) scanf("%lld", &a[i]);
            else a[i] = 0;
            L[i] = R[i] = i;
            while (s.size() > 1 && a[s.top()] >= a[i]){//保证栈顶元素的下边没有大于自己的
                //如果a[栈顶元素] >= a[当前元素], 那么栈顶元素右边的第一个小于自己的就是这个元素
                R[s.top()] = i - 1;
                s.pop();
            }//循环保证有一个左边界
            L[i] = s.top();
            s.push(i);
        }

        for (int i = 0; i < n; i++) printf("%d ", L[i]); cout << endl;
        for (int i = 0; i < n; i++) printf("%d ", R[i]); cout << endl;

        for (int i = 1; i <= n; i++) ans = max(ans, (R[i] - L[i])*a[i]);
        printf("%lld\n", ans);
    }
    return 0;
}
//6 2 5 2 5 5 2