-- SQLite Auto-generated schema
-- This will be executed automatically on startup if tables don't exist

CREATE TABLE IF NOT EXISTS forward (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  user_name VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  tunnel_id BIGINT NOT NULL,
  remote_addr TEXT NOT NULL,
  strategy VARCHAR(100) NOT NULL DEFAULT 'fifo',
  in_flow BIGINT NOT NULL DEFAULT 0,
  out_flow BIGINT NOT NULL DEFAULT 0,
  created_time BIGINT NOT NULL,
  updated_time BIGINT NOT NULL,
  status BIGINT NOT NULL,
  inx BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS forward_port (
  id BIGSERIAL PRIMARY KEY,
  forward_id BIGINT NOT NULL,
  node_id BIGINT NOT NULL,
  port BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS node (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  secret VARCHAR(100) NOT NULL,
  server_ip VARCHAR(100) NOT NULL,
  port TEXT NOT NULL,
  interface_name VARCHAR(200),
  version VARCHAR(100),
  http BIGINT NOT NULL DEFAULT 0,
  tls BIGINT NOT NULL DEFAULT 0,
  socks BIGINT NOT NULL DEFAULT 0,
  created_time BIGINT NOT NULL,
  updated_time BIGINT,
  status BIGINT NOT NULL,
  tcp_listen_addr VARCHAR(100) NOT NULL DEFAULT '[::]',
  udp_listen_addr VARCHAR(100) NOT NULL DEFAULT '[::]'
);

CREATE TABLE IF NOT EXISTS speed_limit (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  speed BIGINT NOT NULL,
  tunnel_id BIGINT NOT NULL,
  tunnel_name VARCHAR(100) NOT NULL,
  created_time BIGINT NOT NULL,
  updated_time BIGINT,
  status BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS statistics_flow (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  flow BIGINT NOT NULL,
  total_flow BIGINT NOT NULL,
  time VARCHAR(100) NOT NULL,
  created_time BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS tunnel (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  traffic_ratio DOUBLE PRECISION NOT NULL DEFAULT 1.0,
  type BIGINT NOT NULL,
  protocol VARCHAR(10) NOT NULL DEFAULT 'tls',
  flow BIGINT NOT NULL,
  created_time BIGINT NOT NULL,
  updated_time BIGINT NOT NULL,
  status BIGINT NOT NULL,
  in_ip TEXT
);

CREATE TABLE IF NOT EXISTS chain_tunnel (
    id BIGSERIAL PRIMARY KEY,
    tunnel_id BIGINT NOT NULL ,
    chain_type VARCHAR(10),
    node_id BIGINT NOT NULL ,
    port BIGINT,
    strategy VARCHAR(10),
    inx  BIGINT,
    protocol  VARCHAR(10)
);


CREATE TABLE IF NOT EXISTS "user" (
  id BIGSERIAL PRIMARY KEY,
  "user" VARCHAR(100) NOT NULL,
  pwd VARCHAR(100) NOT NULL,
  role_id BIGINT NOT NULL,
  exp_time BIGINT NOT NULL,
  flow BIGINT NOT NULL,
  in_flow BIGINT NOT NULL DEFAULT 0,
  out_flow BIGINT NOT NULL DEFAULT 0,
  flow_reset_time BIGINT NOT NULL,
  num BIGINT NOT NULL,
  created_time BIGINT NOT NULL,
  updated_time BIGINT,
  status BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tunnel (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  tunnel_id BIGINT NOT NULL,
  speed_id BIGINT,
  num BIGINT NOT NULL,
  flow BIGINT NOT NULL,
  in_flow BIGINT NOT NULL DEFAULT 0,
  out_flow BIGINT NOT NULL DEFAULT 0,
  flow_reset_time BIGINT NOT NULL,
  exp_time BIGINT NOT NULL,
  status BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS vite_config (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL UNIQUE,
  value VARCHAR(200) NOT NULL,
  time BIGINT NOT NULL
);

