#include<bits/stdc++.h>
#include<cstdio>
#include<cstring>
#include<vector>
#include<queue>
using namespace std;
#define INF 1000000
const int maxn=1000+10;
struct Edges{
    int from,to,d,p;
    Edges(int from=0,int to=0,int d=0,int p=0):from(from),to(to),d(d),p(p){}
};
vector<Edges>edges;
vector<int>g[maxn];
int n,m;
int vis[maxn];
int d[maxn];
int last;

struct node{
    int u,d,p; //节点u，起始节点到当前节点的距离、以及开销
    node(int u=0,int d=0,int p=0):u(u),d(d),p(p){}
    bool operator < (const node& nt)const {
        if(d>nt.d)return true;
        if(d==nt.d && p>nt.p)return true;
        return false ;
    }
};

void dijkstra(int first){
    memset(vis,0,sizeof(vis));
    priority_queue<node>q;
    q.push(node(first,0,0));
    for(int i=1;i<=n;i++)d[i]=INF;
    d[first]=0;

    while(!q.empty()){
        node v=q.top();q.pop();
        int u=v.u;
        if(u==last){
            printf("%d %d\n",v.d,v.p);
            return ;
        }

        if(vis[u])continue;
        vis[u]=1;


        for(int i=0;i<g[u].size();i++){
            Edges e=edges[g[u][i]];
            if(d[e.to]>=d[e.from]+e.d ){   //这不要忘了=，因为两点之间可能有多条路径，长度相同选择消费最小的，这么没注意到WA了一次
                d[e.to]=d[e.from]+e.d;
                q.push(node(e.to,d[e.to],v.p+e.p));
            }
        }
    }
}
int main(){
    while(scanf("%d%d",&n,&m)==2 && n){
        int a,b,d,p;
        for(int i=0;i<m;i++)
        {
            scanf("%d%d%d%d",&a,&b,&d,&p);
            edges.push_back(Edges(a,b,d,p));
            int l=edges.size();
            edges.push_back(Edges(b,a,d,p));
            g[a].push_back(l-1);
            g[b].push_back(l);
        }

        int first;
        scanf("%d%d",&first,&last);
        dijkstra(first);

        for(int i=0;i<n;i++)g[i].clear(); //清空容器
        edges.clear();
    }
    return 0;
}
