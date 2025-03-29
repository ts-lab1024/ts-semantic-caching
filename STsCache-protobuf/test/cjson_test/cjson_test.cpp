#include <stdio.h>
#include <stdlib.h>
#include "cjson/cJSON.h"

int main() {
  // 创建 JSON 数据
  cJSON *root = cJSON_CreateObject();
  cJSON_AddStringToObject(root, "name", "John");
  cJSON_AddNumberToObject(root, "age", 30);

  // 打印 JSON 数据
  char *json_string = cJSON_Print(root);
  printf("JSON Data:\n%s\n", json_string);

  // 解析 JSON 数据
  cJSON *parsed_json = cJSON_Parse(json_string);
  if (parsed_json) {
    cJSON *name = cJSON_GetObjectItem(parsed_json, "name");
    if (name && name->type == cJSON_String) {
      printf("Name: %s\n", name->valuestring);
    }
    cJSON_Delete(parsed_json);
  } else {
    printf("Failed to parse JSON data: %s\n", cJSON_GetErrorPtr());
  }

  // 清理
  cJSON_Delete(root);
  free(json_string);

  return 0;
}