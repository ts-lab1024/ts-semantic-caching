#ifndef STSCACHE_JSON_H
#define STSCACHE_JSON_H

#include <string>
#include <vector>
#include <cjson/cJSON.h>

namespace semantic_json {

struct Sample {
  int64_t timestamp;
  std::vector<double> value;
};

struct SemanticSeriesValue {
  std::string series_segment;
  std::vector<Sample*> values;
};

struct SemanticMetaValue {
  std::string semantic_meta;
  std::vector<SemanticSeriesValue*> series_array;
};

std::string serialize(const SemanticMetaValue& metaValue);

bool deserialize(const std::string& jsonString, SemanticMetaValue& metaValue);

}

#endif  // STSCACHE_JSON_H
