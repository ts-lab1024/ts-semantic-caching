/*
 * fatcache - memcache on ssd.
 * Copyright (C) 2013 Twitter, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <stdio.h>

#include <fc_core.h>
#include <string>
#include <vector>

extern struct string msg_strings[];

struct msg *
rsp_get(struct conn *conn)
{
    struct msg *msg;

    msg = msg_get(conn, false);
    if (msg == NULL) {
        conn->err = errno;
    }

    return msg;
}

void
rsp_put(struct msg *msg)
{
    ASSERT(!msg->request);
    ASSERT(msg->peer == NULL);
    msg_put(msg);
}

struct msg *
rsp_send_next(struct context *ctx, struct conn *conn)
{
    rstatus_t status;
    struct msg *msg, *pmsg;

    pmsg = TAILQ_FIRST(&conn->omsg_q);
    if (pmsg == NULL || !pmsg->done) {
        if (pmsg == NULL && conn->eof) {
            conn->done = 1;
        }

        status = event_del_out(ctx->ep, conn);
        if (status != FC_OK) {
            conn->err = errno;
        }

        return NULL;
    }

    msg = conn->smsg;
    if (msg != NULL) {
        pmsg = TAILQ_NEXT(msg->peer, c_tqe);
    }

    if (pmsg == NULL || !pmsg->done) {
        conn->smsg = NULL;
        return NULL;
    }

    msg = pmsg->peer;

    conn->smsg = msg;

    return msg;
}

void
rsp_send_done(struct context *ctx, struct conn *conn, struct msg *msg)
{
    struct msg *pmsg;

    pmsg = msg->peer;

    req_dequeue_omsgq(ctx, conn, pmsg);

    req_put(pmsg);
}

void rsp_send_string(struct context *ctx, struct conn *conn, struct msg *msg,
                     struct string *str)
{
    rstatus_t status;
    struct msg *pmsg;

    if(!str) return;

    pmsg = rsp_get(conn);
    if (pmsg == NULL) {
        req_process_error(ctx, conn, msg, ENOMEM);
        return;
    }

    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
    pmsg->mlen += str->len;

    msg->done = 1;
    msg->peer = pmsg;
    pmsg->peer = msg;

    status = event_add_out(ctx->ep, conn);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
}

void
rsp_send_status(struct context *ctx, struct conn *conn, struct msg *msg,
                msg_type_t rsp_type)
{
    rstatus_t status;
    struct msg *pmsg;
    struct string *str;

    if (msg->noreply) {
        req_put(msg);
        return;
    }

    pmsg = rsp_get(conn);
    if (pmsg == NULL) {
        req_process_error(ctx, conn, msg, ENOMEM);
        return;
    }

    str = &msg_strings[rsp_type];
    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
    pmsg->mlen += str->len;

    msg->done = 1;
    msg->peer = pmsg;
    pmsg->peer = msg;

    status = event_add_out(ctx->ep, conn);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
}

void
rsp_send_error(struct context *ctx, struct conn *conn, struct msg *msg,
               msg_type_t rsp_type, int err)
{
    rstatus_t status;
    struct msg *pmsg;
    struct string *str, pstr;
    char const*errstr;

    pmsg = rsp_get(conn);
    if (pmsg == NULL) {
        req_process_error(ctx, conn, msg, ENOMEM);
        return;
    }

    str = &msg_strings[rsp_type];
    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
    pmsg->mlen += str->len;

    errstr = err != 0 ? strerror(err) : "unknown";
    string_set_raw(&pstr, errstr);
    str = &pstr;
    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
    pmsg->mlen += str->len;

    str = &msg_strings[MSG_CRLF];
    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
    pmsg->mlen += str->len;

    msg->done = 1;
    msg->peer = pmsg;
    pmsg->peer = msg;

    status = event_add_out(ctx->ep, conn);
    if (status != FC_OK) {
        req_process_error(ctx, conn, msg, errno);
        return;
    }
}


void rsp_send_value(struct context *ctx, struct conn *conn, struct msg *msg, std::vector<uint8_t *>& items,std::vector<std::string> &segments,std::vector<uint32_t > &items_len,std::vector<uint32_t>&per_segment_items_num,std::vector<TimeRange>&per_segment_time_range){
  rstatus_t status;
  struct msg *pmsg;
  struct string *str;

  semantic_json::SemanticMetaValue* semMeta=new semantic_json::SemanticMetaValue();

  pmsg = rsp_get(conn);
  if (pmsg == NULL) {
    req_process_error(ctx, conn, msg, ENOMEM);
    return;
  }

  uint64_t idx=0;
  for(uint32_t i=0;i<segments.size();i++){
    semantic_json::SemanticSeriesValue* semSeries=new semantic_json::SemanticSeriesValue();
    semSeries->series_segment=segments[i];
    int field_num = CountFieldNumInSegment(segments[i]);

    for(uint32_t j=0;j<per_segment_items_num[i];j++){

      if(items[idx+j]== nullptr||items_len[j+idx]==0){
        continue;
      }


      int row_num = items_len[idx+j] / ((field_num+1)*8);
      for (int k = 0; k < row_num; k++) {
        semantic_json::Sample* smp=new semantic_json::Sample();
        for (int m = 0; m < field_num+1; m++) {
          if (m == 0) {
            smp->timestamp=*(uint64_t*)(&items[idx+j][k*(field_num+1)*8]);
          } else {
            smp->value.push_back(*(double*)(&items[idx+j][(k*(field_num+1)+m)*8]));
          }
        }
        semSeries->values.push_back(smp);
      }

    }
    idx += per_segment_items_num[i];

    semMeta->series_array.push_back(semSeries);
  }

  semMeta->semantic_meta=segments[0];

  std::string json_str = serialize(*semMeta);
  status = mbuf_copy_from(&pmsg->mhdr, reinterpret_cast<uint8_t *>(json_str.data()), json_str.size());
  pmsg->mlen += json_str.size();
  if (status != FC_OK) {
    req_process_error(ctx, conn, msg, errno);
    return;
  }

  str = &msg_strings[MSG_CRLF];
  status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
  if (status != FC_OK) {
    req_process_error(ctx, conn, msg, errno);
    return;
  }
  pmsg->mlen += str->len;

  if (msg->frag_id == 0 || msg->last_fragment) {
    str = &msg_strings[MSG_RSP_END];
    status = mbuf_copy_from(&pmsg->mhdr, str->data, str->len);
    if (status != FC_OK) {
      req_process_error(ctx, conn, msg, errno);
      return;
    }
    pmsg->mlen += str->len;
  }

  msg->done = 1;
  msg->peer = pmsg;
  pmsg->peer = msg;


  status = event_add_out(ctx->ep, conn);
  if (status != FC_OK) {
    req_process_error(ctx, conn, msg, errno);
    return;
  }
}

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