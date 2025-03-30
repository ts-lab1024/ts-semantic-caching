#include <iostream>
#include <string>
#include <cstring>
#include <unistd.h>
#include <arpa/inet.h>
#include "message.pb.h"

const int PORT = 8080;

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
    int server_sock = socket(AF_INET, SOCK_STREAM, 0);
    if (server_sock < 0) {
        std::cerr << "Socket creation failed" << std::endl;
        return 1;
    }

    sockaddr_in server_addr{};
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    server_addr.sin_addr.s_addr = INADDR_ANY;

    if (bind(server_sock, (sockaddr*)&server_addr, sizeof(server_addr)) < 0) {
        std::cerr << "Bind failed" << std::endl;
        return 1;
    }

    if (listen(server_sock, 1) < 0) {
        std::cerr << "Listen failed" << std::endl;
        return 1;
    }

    std::cout << "Waiting for connection..." << std::endl;
    sockaddr_in client_addr{};
    socklen_t client_len = sizeof(client_addr);
    int client_sock = accept(server_sock, (sockaddr*)&client_addr, &client_len);
    if (client_sock < 0) {
        std::cerr << "Accept failed" << std::endl;
        return 1;
    }

    // 接收消息
    MyMessage recv_msg = receiveMessage(client_sock);
    std::cout << "Received message: ID=" << recv_msg.id() << ", Content=" << recv_msg.content() << std::endl;

    // 发送消息
    MyMessage send_msg;
    send_msg.set_id(2);
    send_msg.set_content("Hello from receiver!");
    sendMessage(client_sock, send_msg);

    close(client_sock);
    close(server_sock);
    return 0;
}