cmake_minimum_required(VERSION 3.22)
project(cuda_exp CUDA)

set(CMAKE_CUDA_STANDARD 17)

add_executable(cuda_exp main.cu)

set_target_properties(cuda_exp PROPERTIES
        CUDA_SEPARABLE_COMPILATION ON)


