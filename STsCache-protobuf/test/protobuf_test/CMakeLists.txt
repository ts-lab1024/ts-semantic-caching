find_package(Protobuf REQUIRED)

# 设置proto文件路径
set(PROTO_FILE ${CMAKE_SOURCE_DIR}/test/protobuf_test/message.proto)

# 生成protobuf的C++源文件
protobuf_generate_cpp(PROTO_SRCS PROTO_HDRS ${PROTO_FILE})

# 添加可执行文件
add_executable(protobuf_test protobuf_test.cpp ${PROTO_SRCS} ${PROTO_HDRS})
add_executable(sender sender.cpp ${PROTO_SRCS} ${PROTO_HDRS})
add_executable(receiver receiver.cpp ${PROTO_SRCS} ${PROTO_HDRS})

# 链接protobuf库
target_link_libraries(protobuf_test PRIVATE ${Protobuf_LIBRARIES})
target_link_libraries(sender PRIVATE ${Protobuf_LIBRARIES})
target_link_libraries(receiver PRIVATE ${Protobuf_LIBRARIES})

# 包含生成的protobuf头文件路径
target_include_directories(protobuf_test PRIVATE ${CMAKE_CURRENT_BINARY_DIR})
target_include_directories(sender PRIVATE ${CMAKE_CURRENT_BINARY_DIR})
target_include_directories(receiver PRIVATE ${CMAKE_CURRENT_BINARY_DIR})