// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: message.proto

#ifndef PROTOBUF_INCLUDED_message_2eproto
#define PROTOBUF_INCLUDED_message_2eproto

#include <string>

#include <google/protobuf/stubs/common.h>

#if GOOGLE_PROTOBUF_VERSION < 3006001
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please update
#error your headers.
#endif
#if 3006001 < GOOGLE_PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_table_driven.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/inlined_string_field.h>
#include <google/protobuf/metadata.h>
#include <google/protobuf/message.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
#include <google/protobuf/unknown_field_set.h>
// @@protoc_insertion_point(includes)
#define PROTOBUF_INTERNAL_EXPORT_protobuf_message_2eproto 

namespace protobuf_message_2eproto {
// Internal implementation detail -- do not use these members.
struct TableStruct {
  static const ::google::protobuf::internal::ParseTableField entries[];
  static const ::google::protobuf::internal::AuxillaryParseTableField aux[];
  static const ::google::protobuf::internal::ParseTable schema[3];
  static const ::google::protobuf::internal::FieldMetadata field_metadata[];
  static const ::google::protobuf::internal::SerializationTable serialization_table[];
  static const ::google::protobuf::uint32 offsets[];
};
void AddDescriptors();
}  // namespace protobuf_message_2eproto
class Sample;
class SampleDefaultTypeInternal;
extern SampleDefaultTypeInternal _Sample_default_instance_;
class SemanticMetaValue;
class SemanticMetaValueDefaultTypeInternal;
extern SemanticMetaValueDefaultTypeInternal _SemanticMetaValue_default_instance_;
class SemanticSeriesValue;
class SemanticSeriesValueDefaultTypeInternal;
extern SemanticSeriesValueDefaultTypeInternal _SemanticSeriesValue_default_instance_;
namespace google {
namespace protobuf {
template<> ::Sample* Arena::CreateMaybeMessage<::Sample>(Arena*);
template<> ::SemanticMetaValue* Arena::CreateMaybeMessage<::SemanticMetaValue>(Arena*);
template<> ::SemanticSeriesValue* Arena::CreateMaybeMessage<::SemanticSeriesValue>(Arena*);
}  // namespace protobuf
}  // namespace google

// ===================================================================

class SemanticMetaValue : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:SemanticMetaValue) */ {
 public:
  SemanticMetaValue();
  virtual ~SemanticMetaValue();

  SemanticMetaValue(const SemanticMetaValue& from);

  inline SemanticMetaValue& operator=(const SemanticMetaValue& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  SemanticMetaValue(SemanticMetaValue&& from) noexcept
    : SemanticMetaValue() {
    *this = ::std::move(from);
  }

  inline SemanticMetaValue& operator=(SemanticMetaValue&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const SemanticMetaValue& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const SemanticMetaValue* internal_default_instance() {
    return reinterpret_cast<const SemanticMetaValue*>(
               &_SemanticMetaValue_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  void Swap(SemanticMetaValue* other);
  friend void swap(SemanticMetaValue& a, SemanticMetaValue& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline SemanticMetaValue* New() const final {
    return CreateMaybeMessage<SemanticMetaValue>(NULL);
  }

  SemanticMetaValue* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<SemanticMetaValue>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const SemanticMetaValue& from);
  void MergeFrom(const SemanticMetaValue& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(SemanticMetaValue* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // repeated .SemanticSeriesValue series_array = 2;
  int series_array_size() const;
  void clear_series_array();
  static const int kSeriesArrayFieldNumber = 2;
  ::SemanticSeriesValue* mutable_series_array(int index);
  ::google::protobuf::RepeatedPtrField< ::SemanticSeriesValue >*
      mutable_series_array();
  const ::SemanticSeriesValue& series_array(int index) const;
  ::SemanticSeriesValue* add_series_array();
  const ::google::protobuf::RepeatedPtrField< ::SemanticSeriesValue >&
      series_array() const;

  // string semantic_meta = 1;
  void clear_semantic_meta();
  static const int kSemanticMetaFieldNumber = 1;
  const ::std::string& semantic_meta() const;
  void set_semantic_meta(const ::std::string& value);
  #if LANG_CXX11
  void set_semantic_meta(::std::string&& value);
  #endif
  void set_semantic_meta(const char* value);
  void set_semantic_meta(const char* value, size_t size);
  ::std::string* mutable_semantic_meta();
  ::std::string* release_semantic_meta();
  void set_allocated_semantic_meta(::std::string* semantic_meta);

  // @@protoc_insertion_point(class_scope:SemanticMetaValue)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::RepeatedPtrField< ::SemanticSeriesValue > series_array_;
  ::google::protobuf::internal::ArenaStringPtr semantic_meta_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_message_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class SemanticSeriesValue : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:SemanticSeriesValue) */ {
 public:
  SemanticSeriesValue();
  virtual ~SemanticSeriesValue();

  SemanticSeriesValue(const SemanticSeriesValue& from);

  inline SemanticSeriesValue& operator=(const SemanticSeriesValue& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  SemanticSeriesValue(SemanticSeriesValue&& from) noexcept
    : SemanticSeriesValue() {
    *this = ::std::move(from);
  }

  inline SemanticSeriesValue& operator=(SemanticSeriesValue&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const SemanticSeriesValue& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const SemanticSeriesValue* internal_default_instance() {
    return reinterpret_cast<const SemanticSeriesValue*>(
               &_SemanticSeriesValue_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    1;

  void Swap(SemanticSeriesValue* other);
  friend void swap(SemanticSeriesValue& a, SemanticSeriesValue& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline SemanticSeriesValue* New() const final {
    return CreateMaybeMessage<SemanticSeriesValue>(NULL);
  }

  SemanticSeriesValue* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<SemanticSeriesValue>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const SemanticSeriesValue& from);
  void MergeFrom(const SemanticSeriesValue& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(SemanticSeriesValue* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // repeated .Sample values = 2;
  int values_size() const;
  void clear_values();
  static const int kValuesFieldNumber = 2;
  ::Sample* mutable_values(int index);
  ::google::protobuf::RepeatedPtrField< ::Sample >*
      mutable_values();
  const ::Sample& values(int index) const;
  ::Sample* add_values();
  const ::google::protobuf::RepeatedPtrField< ::Sample >&
      values() const;

  // string series_segment = 1;
  void clear_series_segment();
  static const int kSeriesSegmentFieldNumber = 1;
  const ::std::string& series_segment() const;
  void set_series_segment(const ::std::string& value);
  #if LANG_CXX11
  void set_series_segment(::std::string&& value);
  #endif
  void set_series_segment(const char* value);
  void set_series_segment(const char* value, size_t size);
  ::std::string* mutable_series_segment();
  ::std::string* release_series_segment();
  void set_allocated_series_segment(::std::string* series_segment);

  // @@protoc_insertion_point(class_scope:SemanticSeriesValue)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::RepeatedPtrField< ::Sample > values_;
  ::google::protobuf::internal::ArenaStringPtr series_segment_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_message_2eproto::TableStruct;
};
// -------------------------------------------------------------------

class Sample : public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:Sample) */ {
 public:
  Sample();
  virtual ~Sample();

  Sample(const Sample& from);

  inline Sample& operator=(const Sample& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  Sample(Sample&& from) noexcept
    : Sample() {
    *this = ::std::move(from);
  }

  inline Sample& operator=(Sample&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const ::google::protobuf::Descriptor* descriptor();
  static const Sample& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const Sample* internal_default_instance() {
    return reinterpret_cast<const Sample*>(
               &_Sample_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    2;

  void Swap(Sample* other);
  friend void swap(Sample& a, Sample& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline Sample* New() const final {
    return CreateMaybeMessage<Sample>(NULL);
  }

  Sample* New(::google::protobuf::Arena* arena) const final {
    return CreateMaybeMessage<Sample>(arena);
  }
  void CopyFrom(const ::google::protobuf::Message& from) final;
  void MergeFrom(const ::google::protobuf::Message& from) final;
  void CopyFrom(const Sample& from);
  void MergeFrom(const Sample& from);
  void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) final;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const final;
  ::google::protobuf::uint8* InternalSerializeWithCachedSizesToArray(
      bool deterministic, ::google::protobuf::uint8* target) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(Sample* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // repeated double value = 2;
  int value_size() const;
  void clear_value();
  static const int kValueFieldNumber = 2;
  double value(int index) const;
  void set_value(int index, double value);
  void add_value(double value);
  const ::google::protobuf::RepeatedField< double >&
      value() const;
  ::google::protobuf::RepeatedField< double >*
      mutable_value();

  // int64 timestamp = 1;
  void clear_timestamp();
  static const int kTimestampFieldNumber = 1;
  ::google::protobuf::int64 timestamp() const;
  void set_timestamp(::google::protobuf::int64 value);

  // @@protoc_insertion_point(class_scope:Sample)
 private:

  ::google::protobuf::internal::InternalMetadataWithArena _internal_metadata_;
  ::google::protobuf::RepeatedField< double > value_;
  mutable int _value_cached_byte_size_;
  ::google::protobuf::int64 timestamp_;
  mutable ::google::protobuf::internal::CachedSize _cached_size_;
  friend struct ::protobuf_message_2eproto::TableStruct;
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// SemanticMetaValue

// string semantic_meta = 1;
inline void SemanticMetaValue::clear_semantic_meta() {
  semantic_meta_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& SemanticMetaValue::semantic_meta() const {
  // @@protoc_insertion_point(field_get:SemanticMetaValue.semantic_meta)
  return semantic_meta_.GetNoArena();
}
inline void SemanticMetaValue::set_semantic_meta(const ::std::string& value) {
  
  semantic_meta_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:SemanticMetaValue.semantic_meta)
}
#if LANG_CXX11
inline void SemanticMetaValue::set_semantic_meta(::std::string&& value) {
  
  semantic_meta_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:SemanticMetaValue.semantic_meta)
}
#endif
inline void SemanticMetaValue::set_semantic_meta(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  semantic_meta_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:SemanticMetaValue.semantic_meta)
}
inline void SemanticMetaValue::set_semantic_meta(const char* value, size_t size) {
  
  semantic_meta_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:SemanticMetaValue.semantic_meta)
}
inline ::std::string* SemanticMetaValue::mutable_semantic_meta() {
  
  // @@protoc_insertion_point(field_mutable:SemanticMetaValue.semantic_meta)
  return semantic_meta_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* SemanticMetaValue::release_semantic_meta() {
  // @@protoc_insertion_point(field_release:SemanticMetaValue.semantic_meta)
  
  return semantic_meta_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void SemanticMetaValue::set_allocated_semantic_meta(::std::string* semantic_meta) {
  if (semantic_meta != NULL) {
    
  } else {
    
  }
  semantic_meta_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), semantic_meta);
  // @@protoc_insertion_point(field_set_allocated:SemanticMetaValue.semantic_meta)
}

// repeated .SemanticSeriesValue series_array = 2;
inline int SemanticMetaValue::series_array_size() const {
  return series_array_.size();
}
inline void SemanticMetaValue::clear_series_array() {
  series_array_.Clear();
}
inline ::SemanticSeriesValue* SemanticMetaValue::mutable_series_array(int index) {
  // @@protoc_insertion_point(field_mutable:SemanticMetaValue.series_array)
  return series_array_.Mutable(index);
}
inline ::google::protobuf::RepeatedPtrField< ::SemanticSeriesValue >*
SemanticMetaValue::mutable_series_array() {
  // @@protoc_insertion_point(field_mutable_list:SemanticMetaValue.series_array)
  return &series_array_;
}
inline const ::SemanticSeriesValue& SemanticMetaValue::series_array(int index) const {
  // @@protoc_insertion_point(field_get:SemanticMetaValue.series_array)
  return series_array_.Get(index);
}
inline ::SemanticSeriesValue* SemanticMetaValue::add_series_array() {
  // @@protoc_insertion_point(field_add:SemanticMetaValue.series_array)
  return series_array_.Add();
}
inline const ::google::protobuf::RepeatedPtrField< ::SemanticSeriesValue >&
SemanticMetaValue::series_array() const {
  // @@protoc_insertion_point(field_list:SemanticMetaValue.series_array)
  return series_array_;
}

// -------------------------------------------------------------------

// SemanticSeriesValue

// string series_segment = 1;
inline void SemanticSeriesValue::clear_series_segment() {
  series_segment_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& SemanticSeriesValue::series_segment() const {
  // @@protoc_insertion_point(field_get:SemanticSeriesValue.series_segment)
  return series_segment_.GetNoArena();
}
inline void SemanticSeriesValue::set_series_segment(const ::std::string& value) {
  
  series_segment_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:SemanticSeriesValue.series_segment)
}
#if LANG_CXX11
inline void SemanticSeriesValue::set_series_segment(::std::string&& value) {
  
  series_segment_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:SemanticSeriesValue.series_segment)
}
#endif
inline void SemanticSeriesValue::set_series_segment(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  series_segment_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:SemanticSeriesValue.series_segment)
}
inline void SemanticSeriesValue::set_series_segment(const char* value, size_t size) {
  
  series_segment_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:SemanticSeriesValue.series_segment)
}
inline ::std::string* SemanticSeriesValue::mutable_series_segment() {
  
  // @@protoc_insertion_point(field_mutable:SemanticSeriesValue.series_segment)
  return series_segment_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* SemanticSeriesValue::release_series_segment() {
  // @@protoc_insertion_point(field_release:SemanticSeriesValue.series_segment)
  
  return series_segment_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void SemanticSeriesValue::set_allocated_series_segment(::std::string* series_segment) {
  if (series_segment != NULL) {
    
  } else {
    
  }
  series_segment_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), series_segment);
  // @@protoc_insertion_point(field_set_allocated:SemanticSeriesValue.series_segment)
}

// repeated .Sample values = 2;
inline int SemanticSeriesValue::values_size() const {
  return values_.size();
}
inline void SemanticSeriesValue::clear_values() {
  values_.Clear();
}
inline ::Sample* SemanticSeriesValue::mutable_values(int index) {
  // @@protoc_insertion_point(field_mutable:SemanticSeriesValue.values)
  return values_.Mutable(index);
}
inline ::google::protobuf::RepeatedPtrField< ::Sample >*
SemanticSeriesValue::mutable_values() {
  // @@protoc_insertion_point(field_mutable_list:SemanticSeriesValue.values)
  return &values_;
}
inline const ::Sample& SemanticSeriesValue::values(int index) const {
  // @@protoc_insertion_point(field_get:SemanticSeriesValue.values)
  return values_.Get(index);
}
inline ::Sample* SemanticSeriesValue::add_values() {
  // @@protoc_insertion_point(field_add:SemanticSeriesValue.values)
  return values_.Add();
}
inline const ::google::protobuf::RepeatedPtrField< ::Sample >&
SemanticSeriesValue::values() const {
  // @@protoc_insertion_point(field_list:SemanticSeriesValue.values)
  return values_;
}

// -------------------------------------------------------------------

// Sample

// int64 timestamp = 1;
inline void Sample::clear_timestamp() {
  timestamp_ = GOOGLE_LONGLONG(0);
}
inline ::google::protobuf::int64 Sample::timestamp() const {
  // @@protoc_insertion_point(field_get:Sample.timestamp)
  return timestamp_;
}
inline void Sample::set_timestamp(::google::protobuf::int64 value) {
  
  timestamp_ = value;
  // @@protoc_insertion_point(field_set:Sample.timestamp)
}

// repeated double value = 2;
inline int Sample::value_size() const {
  return value_.size();
}
inline void Sample::clear_value() {
  value_.Clear();
}
inline double Sample::value(int index) const {
  // @@protoc_insertion_point(field_get:Sample.value)
  return value_.Get(index);
}
inline void Sample::set_value(int index, double value) {
  value_.Set(index, value);
  // @@protoc_insertion_point(field_set:Sample.value)
}
inline void Sample::add_value(double value) {
  value_.Add(value);
  // @@protoc_insertion_point(field_add:Sample.value)
}
inline const ::google::protobuf::RepeatedField< double >&
Sample::value() const {
  // @@protoc_insertion_point(field_list:Sample.value)
  return value_;
}
inline ::google::protobuf::RepeatedField< double >*
Sample::mutable_value() {
  // @@protoc_insertion_point(field_mutable_list:Sample.value)
  return &value_;
}

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__
// -------------------------------------------------------------------

// -------------------------------------------------------------------


// @@protoc_insertion_point(namespace_scope)


// @@protoc_insertion_point(global_scope)

#endif  // PROTOBUF_INCLUDED_message_2eproto
