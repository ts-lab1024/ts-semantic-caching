#include <iostream>
#include "message.pb.h"  // 生成的protobuf头文件

int main() {
    // 创建消息对象
    MyMessage msg;
    msg.set_id(123);
    msg.set_content("Hello, Protobuf!");

    // 打印消息内容
    std::cout << "ID: " << msg.id() << ", Content: " << msg.content() << std::endl;

    return 0;
}