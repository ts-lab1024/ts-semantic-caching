#include <cjson/cJSON.h>
#include <iostream>
#include <string>
#include <vector>
#include "../include/semantic_meta_value_json/semantic_meta_value_json.h"

namespace semantic_json {
std::string serialize(const SemanticMetaValue& metaValue) {
  // 创建根对象
  cJSON* root = cJSON_CreateObject();

  // 添加 semantic_meta
  cJSON_AddStringToObject(root, "semantic_meta", metaValue.semantic_meta.c_str());

  // 添加 series_array
  cJSON* seriesArray = cJSON_CreateArray();
  for (const auto& series : metaValue.series_array) {
    cJSON* seriesItem = cJSON_CreateObject();
    cJSON_AddStringToObject(seriesItem, "series_segment", series->series_segment.c_str());

    // 添加 values
    cJSON* valuesArray = cJSON_CreateArray();
    for (const auto& sample : series->values) {
      cJSON* sampleItem = cJSON_CreateObject();
      cJSON_AddNumberToObject(sampleItem, "timestamp", sample->timestamp);

      // 添加 value
      cJSON* valueArray = cJSON_CreateArray();
      for (const auto& val : sample->value) {
        cJSON_AddItemToArray(valueArray, cJSON_CreateNumber(val));
      }
      cJSON_AddItemToObject(sampleItem, "value", valueArray);

      cJSON_AddItemToArray(valuesArray, sampleItem);
    }
    cJSON_AddItemToObject(seriesItem, "values", valuesArray);

    cJSON_AddItemToArray(seriesArray, seriesItem);
  }
  cJSON_AddItemToObject(root, "series_array", seriesArray);

  // 转换为字符串
  char* jsonString = cJSON_Print(root);
  std::string result(jsonString);
  free(jsonString);
  cJSON_Delete(root);

  return result;
}

bool deserialize(const std::string& jsonString, SemanticMetaValue& metaValue) {
  cJSON* root = cJSON_Parse(jsonString.c_str());
  if (!root) {
    return false; // 解析失败
  }

  // 解析 semantic_meta
  cJSON* semanticMeta = cJSON_GetObjectItem(root, "semantic_meta");
  if (semanticMeta && cJSON_IsString(semanticMeta)) {
    metaValue.semantic_meta = semanticMeta->valuestring;
  }

  // 解析 series_array
  cJSON* seriesArray = cJSON_GetObjectItem(root, "series_array");
  if (seriesArray && cJSON_IsArray(seriesArray)) {
    int seriesCount = cJSON_GetArraySize(seriesArray);
    for (int i = 0; i < seriesCount; i++) {
      cJSON* seriesItem = cJSON_GetArrayItem(seriesArray, i);
      if (!seriesItem) continue;

      SemanticSeriesValue* series = new SemanticSeriesValue();
      cJSON* seriesSegment = cJSON_GetObjectItem(seriesItem, "series_segment");
      if (seriesSegment && cJSON_IsString(seriesSegment)) {
        series->series_segment = seriesSegment->valuestring;
      }

      // 解析 values
      cJSON* valuesArray = cJSON_GetObjectItem(seriesItem, "values");
      if (valuesArray && cJSON_IsArray(valuesArray)) {
        int valuesCount = cJSON_GetArraySize(valuesArray);
        for (int j = 0; j < valuesCount; j++) {
          cJSON* sampleItem = cJSON_GetArrayItem(valuesArray, j);
          if (!sampleItem) continue;

          Sample* sample = new Sample();
          cJSON* timestamp = cJSON_GetObjectItem(sampleItem, "timestamp");
          if (timestamp && cJSON_IsNumber(timestamp)) {
            sample->timestamp = timestamp->valueint;
          }

          cJSON* valueArray = cJSON_GetObjectItem(sampleItem, "value");
          if (valueArray && cJSON_IsArray(valueArray)) {
            int valueCount = cJSON_GetArraySize(valueArray);
            for (int k = 0; k < valueCount; k++) {
              cJSON* valueItem = cJSON_GetArrayItem(valueArray, k);
              if (valueItem && cJSON_IsNumber(valueItem)) {
                sample->value.push_back(valueItem->valuedouble);
              }
            }
          }
          series->values.push_back(sample);
        }
      }
      metaValue.series_array.push_back(series);
    }
  }

  cJSON_Delete(root);
  return true;
}
}