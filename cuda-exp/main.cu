#include <stdio.h>
#include <iostream>
__global__ void add(int a, int b){
    int c = a + b;
    int d = c + a;
    printf("c: %d\n", c);
}

void test_add(int a, int b){
    add<<<1, 1>>>(a, b);
    cudaDeviceReset(); // 同步设备函数的结果
}

int main() {
    int a = 1, b = 2;
    test_add(a, b);
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
