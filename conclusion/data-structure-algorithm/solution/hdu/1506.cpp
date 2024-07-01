#include<bits/stdc++.h>
#include<iostream>
#include<algorithm>
#include<cstdio>
#include<string.h>
#include<string>
#include<stack>
#include<vector>
#include<queue>
using namespace std;
typedef long long ll;
const int INF = 0x3f3f3f3f;
const int MAXN = 100000 + 100;
ll a[MAXN];
ll L[MAXN], R[MAXN];

int main() {
    int n;
    while (scanf("%d", &n) != EOF && n) {
        stack<int> s;
        s.push(-1);
        for (int i = 0; i < n; i++) {
            scanf("%lld", &a[i]);
            L[i] = R[i] = i;
            while(s.size() > 1 && a[s.top()] >= a[i]) {
                R[s.top()] = i - 1;
                s.pop();
            }
            L[i] = s.top() + 1;
            s.push(i);
        }

        ll ans = 0;
        for (int i = 0; i < n; i++) ans = max(ans, (R[i] - L[i] + 1)*a[i]);
        printf("%lld\n", ans);
    }
    return 0;
}
