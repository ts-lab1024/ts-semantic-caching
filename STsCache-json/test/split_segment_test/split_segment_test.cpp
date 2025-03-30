#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string>
#include <regex>
#include <sstream>
#include <vector>

int CountFieldNumInSegment(const std::string& input) {
  int num = 0;
  std::vector<std::string> tokens;
  std::istringstream stream(input);
  std::string token;

  while (std::getline(stream, token, '#')) {
    tokens.push_back(token);
  }

  std::string fields = tokens[1];
  num = std::count(fields.begin(), fields.end(), ',')+1;

  return num;
}

int main() {
  std::string segment = "{(cpu.hostname=host_1)}#{usage_system[int64],usage_idle[int64],usage_nice[int64]}#{empty}#{mean,5m}";

  int field_num = CountFieldNumInSegment(segment);
  std::cout<<field_num<<std::endl;

  return 0;
}

