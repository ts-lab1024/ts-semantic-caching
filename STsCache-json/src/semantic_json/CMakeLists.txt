add_library(
        semantic_json
        OBJECT
        semantic_meta_value_json.cpp
)

target_link_libraries(
        semantic_json
        cjson
)

set(ALL_OBJECT_FILES
        ${ALL_OBJECT_FILES} $<TARGET_OBJECTS:semantic_json>
        PARENT_SCOPE)