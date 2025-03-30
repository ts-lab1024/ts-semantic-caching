#include <iostream>
#include <string>
#include <cstring>
#include <unistd.h>
#include <arpa/inet.h>
#include "message.pb.h"

const int PORT = 8080;
const char* SERVER_IP = "127.0.0.1";

void sendMessage(int socket, const MyMessage& message) {
    std::string serialized = message.SerializeAsString();
    uint32_t size = htonl(serialized.size());
    send(socket, &size, sizeof(size), 0);
    send(socket, serialized.data(), serialized.size(), 0);
}

MyMessage receiveMessage(int socket) {
    uint32_t size;
    recv(socket, &size, sizeof(size), 0);
    size = ntohl(size);

    char* buffer = new char[size];
    recv(socket, buffer, size, 0);

    MyMessage message;
    message.ParseFromArray(buffer, size);
    delete[] buffer;

    return message;
}

int main() {
    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) {
        std::cerr << "Socket creation failed" << std::endl;
        return 1;
    }

    sockaddr_in server_addr{};
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    inet_pton(AF_INET, SERVER_IP, &server_addr.sin_addr);

    if (connect(sock, (sockaddr*)&server_addr, sizeof(server_addr))) {
        std::cerr << "Connection failed" << std::endl;
        return 1;
    }

    // 发送消息
    MyMessage send_msg;
    send_msg.set_id(1);
    send_msg.set_content("Hello from sender!");
    sendMessage(sock, send_msg);

    // 接收消息
    MyMessage recv_msg = receiveMessage(sock);
    std::cout << "Received message: ID=" << recv_msg.id() << ", Content=" << recv_msg.content() << std::endl;

    close(sock);
    return 0;
}