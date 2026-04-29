-- =============================================================================
-- OpenTelemetry ClickHouse 数据表初始化脚本
-- =============================================================================
-- 用途：手动创建 OTel 指标、日志、追踪数据表
-- 执行方式：
--   docker exec -i <container> clickhouse-client --multiquery < init.sql
-- 或进入容器后执行：
--   clickhouse-client --multiquery < /path/to/init.sql
-- =============================================================================
-- Schema 参考：
--   https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/clickhouseexporter
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 数据库
-- -----------------------------------------------------------------------------
CREATE DATABASE IF NOT EXISTS otel;

-- -----------------------------------------------------------------------------
-- Metrics（指标）
-- -----------------------------------------------------------------------------

-- Gauge / Sum 类型指标
CREATE TABLE IF NOT EXISTS otel.otel_metrics_gauge
(
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeDroppedAttrCount          UInt32 CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    MetricName                     String CODEC(ZSTD(1)),
    MetricDescription              String CODEC(ZSTD(1)),
    MetricUnit                     String CODEC(ZSTD(1)),
    Attributes                     Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    StartTimeUnix                  DateTime64(9) CODEC(Delta, ZSTD(1)),
    TimeUnix                       DateTime64(9) CODEC(Delta, ZSTD(1)),
    Value                          Float64 CODEC(ZSTD(1)),
    Flags                          UInt32 CODEC(ZSTD(1)),
    Exemplars                      Nested (
        FilteredAttributes         Map(LowCardinality(String), String),
        TimeUnix                   DateTime64(9),
        Value                      Float64,
        SpanId                     String,
        TraceId                    String
    ) CODEC(ZSTD(1))
)
ENGINE = MergeTree()
PARTITION BY toDate(TimeUnix)
ORDER BY (ServiceName, MetricName, Attributes, toUnixTimestamp64Nano(TimeUnix))
TTL toDateTime(TimeUnix) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS otel.otel_metrics_sum
(
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeDroppedAttrCount          UInt32 CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    MetricName                     String CODEC(ZSTD(1)),
    MetricDescription              String CODEC(ZSTD(1)),
    MetricUnit                     String CODEC(ZSTD(1)),
    Attributes                     Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    StartTimeUnix                  DateTime64(9) CODEC(Delta, ZSTD(1)),
    TimeUnix                       DateTime64(9) CODEC(Delta, ZSTD(1)),
    Value                          Float64 CODEC(ZSTD(1)),
    Flags                          UInt32 CODEC(ZSTD(1)),
    Exemplars                      Nested (
        FilteredAttributes         Map(LowCardinality(String), String),
        TimeUnix                   DateTime64(9),
        Value                      Float64,
        SpanId                     String,
        TraceId                    String
    ) CODEC(ZSTD(1)),
    IsMonotonic                    Boolean CODEC(Delta, ZSTD(1))
)
ENGINE = MergeTree()
PARTITION BY toDate(TimeUnix)
ORDER BY (ServiceName, MetricName, Attributes, toUnixTimestamp64Nano(TimeUnix))
TTL toDateTime(TimeUnix) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS otel.otel_metrics_histogram
(
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeDroppedAttrCount          UInt32 CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    MetricName                     String CODEC(ZSTD(1)),
    MetricDescription              String CODEC(ZSTD(1)),
    MetricUnit                     String CODEC(ZSTD(1)),
    Attributes                     Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    StartTimeUnix                  DateTime64(9) CODEC(Delta, ZSTD(1)),
    TimeUnix                       DateTime64(9) CODEC(Delta, ZSTD(1)),
    Count                          UInt64 CODEC(Delta, ZSTD(1)),
    Sum                            Float64 CODEC(ZSTD(1)),
    BucketCounts                   Array(UInt64) CODEC(ZSTD(1)),
    ExplicitBounds                 Array(Float64) CODEC(ZSTD(1)),
    Exemplars                      Nested (
        FilteredAttributes         Map(LowCardinality(String), String),
        TimeUnix                   DateTime64(9),
        Value                      Float64,
        SpanId                     String,
        TraceId                    String
    ) CODEC(ZSTD(1)),
    Flags                          UInt32 CODEC(ZSTD(1)),
    Min                            Float64 CODEC(ZSTD(1)),
    Max                            Float64 CODEC(ZSTD(1))
)
ENGINE = MergeTree()
PARTITION BY toDate(TimeUnix)
ORDER BY (ServiceName, MetricName, Attributes, toUnixTimestamp64Nano(TimeUnix))
TTL toDateTime(TimeUnix) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS otel.otel_metrics_exponential_histogram
(
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeDroppedAttrCount          UInt32 CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    MetricName                     String CODEC(ZSTD(1)),
    MetricDescription              String CODEC(ZSTD(1)),
    MetricUnit                     String CODEC(ZSTD(1)),
    Attributes                     Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    StartTimeUnix                  DateTime64(9) CODEC(Delta, ZSTD(1)),
    TimeUnix                       DateTime64(9) CODEC(Delta, ZSTD(1)),
    Count                          UInt64 CODEC(Delta, ZSTD(1)),
    Sum                            Float64 CODEC(ZSTD(1)),
    Scale                          Int32 CODEC(ZSTD(1)),
    ZeroCount                      UInt64 CODEC(Delta, ZSTD(1)),
    PositiveOffset                 Int32 CODEC(ZSTD(1)),
    PositiveBucketCounts           Array(UInt64) CODEC(ZSTD(1)),
    NegativeOffset                 Int32 CODEC(ZSTD(1)),
    NegativeBucketCounts           Array(UInt64) CODEC(ZSTD(1)),
    Exemplars                      Nested (
        FilteredAttributes         Map(LowCardinality(String), String),
        TimeUnix                   DateTime64(9),
        Value                      Float64,
        SpanId                     String,
        TraceId                    String
    ) CODEC(ZSTD(1)),
    Flags                          UInt32 CODEC(ZSTD(1)),
    Min                            Float64 CODEC(ZSTD(1)),
    Max                            Float64 CODEC(ZSTD(1))
)
ENGINE = MergeTree()
PARTITION BY toDate(TimeUnix)
ORDER BY (ServiceName, MetricName, Attributes, toUnixTimestamp64Nano(TimeUnix))
TTL toDateTime(TimeUnix) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS otel.otel_metrics_summary
(
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeDroppedAttrCount          UInt32 CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    MetricName                     String CODEC(ZSTD(1)),
    MetricDescription              String CODEC(ZSTD(1)),
    MetricUnit                     String CODEC(ZSTD(1)),
    Attributes                     Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    StartTimeUnix                  DateTime64(9) CODEC(Delta, ZSTD(1)),
    TimeUnix                       DateTime64(9) CODEC(Delta, ZSTD(1)),
    Count                          UInt64 CODEC(Delta, ZSTD(1)),
    Sum                            Float64 CODEC(ZSTD(1)),
    ValueAtQuantiles               Nested(
        Quantile                   Float64,
        Value                      Float64
    ) CODEC(ZSTD(1)),
    Flags                          UInt32 CODEC(ZSTD(1))
)
ENGINE = MergeTree()
PARTITION BY toDate(TimeUnix)
ORDER BY (ServiceName, MetricName, Attributes, toUnixTimestamp64Nano(TimeUnix))
TTL toDateTime(TimeUnix) + INTERVAL 90 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

-- -----------------------------------------------------------------------------
-- Logs（日志）
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS otel.otel_logs
(
    Timestamp                      DateTime64(9) CODEC(Delta, ZSTD(1)),
    TraceId                        String CODEC(ZSTD(1)),
    SpanId                         String CODEC(ZSTD(1)),
    TraceFlags                     UInt32 CODEC(ZSTD(1)),
    SeverityText                   LowCardinality(String) CODEC(ZSTD(1)),
    SeverityNumber                 Int32 CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    Body                           String CODEC(ZSTD(1)),
    ResourceSchemaUrl              String CODEC(ZSTD(1)),
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeSchemaUrl                 String CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    ScopeAttributes                Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    LogAttributes                  Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    INDEX idx_trace_id              TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_res_attr_key          mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_res_attr_value        mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_log_attr_key          mapKeys(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_log_attr_value        mapValues(LogAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_body                  Body TYPE tokenbf_v1(32768, 3, 0) GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, SeverityText, toUnixTimestamp64Nano(Timestamp))
TTL toDateTime(Timestamp) + INTERVAL 30 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

-- -----------------------------------------------------------------------------
-- Traces（追踪）
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS otel.otel_traces
(
    Timestamp                      DateTime64(9) CODEC(Delta, ZSTD(1)),
    TraceId                        String CODEC(ZSTD(1)),
    SpanId                         String CODEC(ZSTD(1)),
    ParentSpanId                   String CODEC(ZSTD(1)),
    TraceState                     String CODEC(ZSTD(1)),
    SpanName                       LowCardinality(String) CODEC(ZSTD(1)),
    SpanKind                       LowCardinality(String) CODEC(ZSTD(1)),
    ServiceName                    LowCardinality(String) CODEC(ZSTD(1)),
    ResourceAttributes             Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    ScopeName                      String CODEC(ZSTD(1)),
    ScopeVersion                   String CODEC(ZSTD(1)),
    SpanAttributes                 Map(LowCardinality(String), String) CODEC(ZSTD(1)),
    Duration                       Int64 CODEC(ZSTD(1)),
    StatusCode                     LowCardinality(String) CODEC(ZSTD(1)),
    StatusMessage                  String CODEC(ZSTD(1)),
    Events                         Nested (
        Timestamp                  DateTime64(9),
        Name                       LowCardinality(String),
        Attributes                 Map(LowCardinality(String), String)
    ) CODEC(ZSTD(1)),
    Links                          Nested (
        TraceId                    String,
        SpanId                     String,
        TraceState                 String,
        Attributes                 Map(LowCardinality(String), String)
    ) CODEC(ZSTD(1)),
    INDEX idx_trace_id              TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
    INDEX idx_res_attr_key          mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_res_attr_value        mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_span_attr_key         mapKeys(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_span_attr_value       mapValues(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
    INDEX idx_duration              Duration TYPE minmax GRANULARITY 1
)
ENGINE = MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, SpanName, toUnixTimestamp64Nano(Timestamp))
TTL toDateTime(Timestamp) + INTERVAL 30 DAY
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;
