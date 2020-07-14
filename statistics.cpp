#include "rocksdb/c.h"
#include "rocksdb/options.h"
#include "rocksdb/statistics.h"

using rocksdb::Options;
using rocksdb::Tickers;

extern "C" {
struct rocksdb_options_t { Options rep; };

uint64_t gorocksdb_get_ticker_count(rocksdb_options_t *options, uint32_t ticker) {
    return options->rep.statistics->getTickerCount(ticker);
}

uint64_t gorocksdb_get_and_reset_ticker_count(rocksdb_options_t *options, uint32_t ticker) {
    return options->rep.statistics->getAndResetTickerCount(ticker);
}

void gorocksdb_set_stats_level(rocksdb_options_t *options, uint8_t level) {
    options->rep.statistics->set_stats_level((rocksdb::StatsLevel)level);
}

uint8_t gorocksdb_get_stats_level(rocksdb_options_t *options) {
    return options->rep.statistics->get_stats_level();
}

}