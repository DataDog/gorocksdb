#include "rocksdb/c.h"
#include "rocksdb/options.h"
#include "rocksdb/statistics.h"
#include "rocksdb/flush_block_policy.h"

using rocksdb::Options;
using rocksdb::Tickers;

using namespace rocksdb;

class FlushBlockEveryKeyPolicy : public FlushBlockPolicy {
 public:
  bool Update(const Slice& /*key*/, const Slice& /*value*/) override {
    if (!start_) {
      start_ = true;
      return false;
    }
    return true;
  }

 private:
  bool start_ = false;
};

class FlushBlockEveryKeyPolicyFactory : public FlushBlockPolicyFactory {
 public:
  explicit FlushBlockEveryKeyPolicyFactory() {}

  const char* Name() const override {
    return "FlushBlockEveryKeyPolicyFactory";
  }

  FlushBlockPolicy* NewFlushBlockPolicy(
      const BlockBasedTableOptions& /*table_options*/,
      const BlockBuilder& /*data_block_builder*/) const override {
    return new FlushBlockEveryKeyPolicy;
  }
};


extern "C" {

struct rocksdb_options_t { Options rep; };
struct rocksdb_block_based_table_options_t { BlockBasedTableOptions rep; };

uint64_t rocksdb_get_ticker_count(rocksdb_options_t *options, uint32_t ticker) {
    return options->rep.statistics->getTickerCount(ticker);
}

uint64_t rocksdb_get_and_reset_ticker_count(rocksdb_options_t *options, uint32_t ticker) {
    return options->rep.statistics->getAndResetTickerCount(ticker);
}

void rocksdb_record_tick(rocksdb_options_t *options, uint32_t ticker, uint64_t count) {
    options->rep.statistics->recordTick(ticker, count);
}

void rocksdb_set_ticker_count(rocksdb_options_t *options, uint32_t ticker, uint64_t count) {
    options->rep.statistics->recordTick(ticker, count);
}

void rocksdb_set_stats_level(rocksdb_options_t *options, uint8_t level) {
    options->rep.statistics->set_stats_level((rocksdb::StatsLevel)level);
}

uint8_t rocksdb_get_stats_level(rocksdb_options_t *options) {
    return options->rep.statistics->get_stats_level();
}

void rocksdb_block_based_options_set_flush_every_key_policy(rocksdb_block_based_table_options_t *opts) {
    opts->rep.flush_block_policy_factory = std::make_shared<FlushBlockEveryKeyPolicyFactory>();
}


}