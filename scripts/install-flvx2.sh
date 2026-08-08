#!/bin/bash
# ============================================================
# FLVX2 一键脚本(Go 单二进制,免编译)
#
# 用法:
#   curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install
#   curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- install-agent
#   curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- update
#   curl -fsSL https://raw.githubusercontent.com/Fyzzp0o0/FLVX2/main/scripts/install-flvx2.sh | bash -s -- uninstall
#
# 环境变量:
#   FLVX2_VERSION          指定版本(默认 latest release)
#   FLVX2_AGENT_URL        agent 二进制下载地址(默认本仓库 Releases)
#   FRONTEND_PORT          前端端口(默认 6635)
#   BACKEND_PORT           后端端口,agent 对接(默认 6636)
#   DB_NAME/DB_USER/DB_PASSWORD  PostgreSQL 配置(默认随机,已存在则复用 .env)
#   JWT_SECRET             JWT 密钥(默认随机,已存在则复用 .env)
#   NODE_SECRET            节点密钥(install-agent 必填:面板创建节点后获得)
#   NODE_PANEL_ADDR        节点连接的面板地址(默认 本机IP:6636)
# ============================================================
set -euo pipefail

ACTION="${1:-install}"
REPO_SLUG="${FLVX2_GITHUB_REPO:-Fyzzp0o0/FLVX2}"
APP_DIR="${FLVX2_PANEL_DIR:-/opt/FLVX2}"
AGENT_DIR="${FLVX2_AGENT_DIR:-/opt/FLVX2-Agent}"
SERVICE_NAME="flvx2-backend"
AGENT_SERVICE="flvx2-agent"
ENV_FILE="$APP_DIR/.env"

FRONTEND_PORT="${FRONTEND_PORT:-6635}"
BACKEND_PORT="${BACKEND_PORT:-6636}"
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-flvx2}"
DB_USER="${DB_USER:-flvx2}"
DB_PASSWORD="${DB_PASSWORD:-$(openssl rand -hex 12)}"
JWT_SECRET="${JWT_SECRET:-$(openssl rand -hex 16)}"
NODE_SECRET="${NODE_SECRET:-}"
NODE_PANEL_ADDR="${NODE_PANEL_ADDR:-}"
# install-agent 支持位置参数: -a <面板IP:后端端口> -s <节点密钥>
if [ "$ACTION" = "install-agent" ] || [ "$ACTION" = "agent" ]; then
  shift
  while [ $# -gt 0 ]; do
    case "$1" in
      -a) NODE_PANEL_ADDR="$2"; shift 2 ;;
      -s) NODE_SECRET="$2"; shift 2 ;;
      *) shift ;;
    esac
  done
fi

GREEN='\033[0;32m'; RED='\033[0;31m'; NC='\033[0m'
ok()  { echo -e "${GREEN}[OK] $1${NC}"; }
err() { echo -e "${RED}[FAIL] $1${NC}"; exit 1; }

require_root() {
  [ "$(id -u)" = "0" ] || err "请以 root 运行"
}

# ---------- 版本与下载 ----------
latest_release_version() {
  curl -fsSL --retry 3 --connect-timeout 10 "https://api.github.com/repos/${REPO_SLUG}/releases?per_page=1" \
    | sed -nE 's/.*"tag_name"[[:space:]]*:[[:space:]]*"v?([^"]+)".*/\1/p' | head -1 \
    || err "无法获取最新版本号"
}

resolve_version() {
  if [ -n "${FLVX2_VERSION:-}" ]; then echo "${FLVX2_VERSION#v}"; else latest_release_version; fi
}

panel_url() {
  local arch
  case "$(uname -m)" in
    x86_64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) err "不支持的架构: $(uname -m)" ;;
  esac
  echo "https://github.com/${REPO_SLUG}/releases/download/v${1}/FLVX2-linux-${arch}"
}

agent_url() {
  if [ -n "${FLVX2_AGENT_URL:-}" ]; then echo "$FLVX2_AGENT_URL"; return; fi
  local arch
  case "$(uname -m)" in
    x86_64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) err "不支持的架构: $(uname -m)" ;;
  esac
  echo "https://github.com/${REPO_SLUG}/releases/download/v${1}/FLVX2-Agent-linux-${arch}"
}

# ---------- 依赖 ----------
install_base_deps() {
  local need=""
  command -v curl >/dev/null 2>&1 || need="$need curl"
  command -v tar >/dev/null 2>&1 || need="$need tar"
  command -v openssl >/dev/null 2>&1 || need="$need openssl"
  command -v psql >/dev/null 2>&1 || need="$need postgresql"
  if [ -n "$need" ]; then
    echo "[INFO] 安装系统依赖: $need"
    export DEBIAN_FRONTEND=noninteractive
    apt-get update -qq
    apt-get install -y -qq curl tar openssl postgresql >/dev/null
  fi
}

# ---------- PostgreSQL(幂等) ----------
ensure_postgres() {
  if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" | grep -q 1; then
    sudo -u postgres psql -c "CREATE USER \"${DB_USER}\" WITH PASSWORD '${DB_PASSWORD}'" >/dev/null
  fi
  if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1; then
    sudo -u postgres createdb -O "${DB_USER}" "${DB_NAME}"
  fi
  ok "PostgreSQL 就绪: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
}

# 修复 /etc/hosts 主机名映射(消除 sudo "unable to resolve host" 警告)
fix_hostname_resolution() {
  local hn; hn="$(hostname)"
  if [ -n "$hn" ] && ! grep -q "$hn" /etc/hosts; then
    echo "127.0.1.1 $hn" >> /etc/hosts
    echo "[INFO] 已修复 /etc/hosts 主机名解析"
  fi
}

# ---------- 面板 .env(复用已有,避免重跑改密导致 token 失效) ----------
write_env() {
  if [ -f "$ENV_FILE" ]; then
    echo "[INFO] 复用已有配置 $ENV_FILE"
    return
  fi
  mkdir -p "$APP_DIR"
  cat > "$ENV_FILE" <<EOF
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
JWT_SECRET=$JWT_SECRET
LOG_DIR=$APP_DIR/logs
FRONTEND_PORT=$FRONTEND_PORT
BACKEND_PORT=$BACKEND_PORT
EOF
  chmod 600 "$ENV_FILE"
  ok "已生成配置 $ENV_FILE(密钥权限 600)"
}

# ---------- systemd ----------
write_systemd() {
  cat > /etc/systemd/system/${SERVICE_NAME}.service <<EOF
[Unit]
Description=FLVX2 Panel
After=network.target postgresql.service

[Service]
Type=simple
WorkingDirectory=$APP_DIR
EnvironmentFile=$ENV_FILE
ExecStart=$APP_DIR/FLVX2
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable --now ${SERVICE_NAME} >/dev/null 2>&1
  ok "服务已启动: systemctl status ${SERVICE_NAME}"
}

write_agent_systemd() {
  cat > /etc/systemd/system/${AGENT_SERVICE}.service <<EOF
[Unit]
Description=FLVX2 Agent
After=network.target

[Service]
Type=simple
WorkingDirectory=$AGENT_DIR
ExecStart=$AGENT_DIR/FLVX2-Agent
Restart=always
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable --now ${AGENT_SERVICE} >/dev/null 2>&1
  ok "Agent 服务已启动: systemctl status ${AGENT_SERVICE}"
}

# ---------- 动作 ----------
install_panel() {
  require_root
  fix_hostname_resolution
  install_base_deps
  ensure_postgres
  write_env
  local version; version="$(resolve_version)"
  echo "[INFO] 下载面板 v${version} ..."
  local tmp; tmp="$(mktemp)"
  curl -fsSL --retry 3 -o "$tmp" "$(panel_url "$version")" || err "面板下载失败: $(panel_url "$version")"
  install -m755 "$tmp" "$APP_DIR/FLVX2"
  rm -f "$tmp"
  write_systemd
  echo
  echo "================ FLVX2 面板安装完成 ================"
  echo "  前端: http://<服务器IP>:${FRONTEND_PORT}"
  echo "  后端(agent 对接): :${BACKEND_PORT}"
  echo "  默认账号: admin_user / admin_user (请立即修改)"
  echo "  维护: systemctl status/restart ${SERVICE_NAME}"
  echo "  更新: $0 update   卸载: $0 uninstall"
  echo "===================================================="
}

install_agent() {
  require_root
  fix_hostname_resolution
  [ -n "$NODE_SECRET" ] || err "请设置 NODE_SECRET(面板 → 节点管理 → 创建节点后获得)"
  if [ -z "$NODE_PANEL_ADDR" ]; then
    NODE_PANEL_ADDR="$(hostname -I 2>/dev/null | awk '{print $1}'):${BACKEND_PORT}"
  fi
  local version; version="$(resolve_version)"
  echo "[INFO] 下载 Agent v${version} ..."
  mkdir -p "$AGENT_DIR"
  local tmp; tmp="$(mktemp)"
  curl -fsSL --retry 3 -o "$tmp" "$(agent_url "$version")" || err "Agent 下载失败: $(agent_url "$version")"
  install -m755 "$tmp" "$AGENT_DIR/FLVX2-Agent"
  rm -f "$tmp"
  cat > "$AGENT_DIR/config.json" <<EOF
{"addr": "${NODE_PANEL_ADDR}", "secret": "${NODE_SECRET}", "http": 1, "tls": 0, "socks": 1}
EOF
  chmod 600 "$AGENT_DIR/config.json"
  write_agent_systemd
  ok "Agent 已连接面板: ${NODE_PANEL_ADDR}"
}

update_panel() {
  require_root
  [ -f "$APP_DIR/FLVX2" ] || err "面板未安装($APP_DIR/FLVX2 不存在)"
  local version; version="$(resolve_version)"
  echo "[INFO] 更新面板至 v${version} ..."
  local tmp; tmp="$(mktemp)"
  curl -fsSL --retry 3 -o "$tmp" "$(panel_url "$version")" || err "面板下载失败"
  install -m755 "$tmp" "$APP_DIR/FLVX2"
  rm -f "$tmp"
  systemctl restart ${SERVICE_NAME}
  ok "面板已更新至 v${version}(配置与数据保留)"
}

uninstall_panel() {
  require_root
  echo "[WARN] 将完全卸载: 服务/二进制/配置,数据库保留(见下方提示)"
  read -r -p "确认卸载? [y/N] " ans
  [ "$ans" = "y" ] || [ "$ans" = "Y" ] || { echo "已取消"; exit 0; }
  systemctl disable --now ${SERVICE_NAME} >/dev/null 2>&1 || true
  systemctl disable --now ${AGENT_SERVICE} >/dev/null 2>&1 || true
  rm -f /etc/systemd/system/${SERVICE_NAME}.service /etc/systemd/system/${AGENT_SERVICE}.service
  systemctl daemon-reload
  rm -rf "$APP_DIR" "$AGENT_DIR"
  echo
  echo "================ 卸载完成 ================"
  echo "  已删除: 服务、$APP_DIR、$AGENT_DIR"
  echo "  数据库 ${DB_NAME} 保留;如需删除:"
  echo "    sudo -u postgres dropdb ${DB_NAME}"
  echo "    sudo -u postgres psql -c \"DROP USER IF EXISTS \\\"${DB_USER}\\\"\""
  echo "=========================================="
}

case "$ACTION" in
  install) install_panel ;;
  install-agent|agent) install_agent ;;
  update|upgrade) update_panel ;;
  uninstall|remove) uninstall_panel ;;
  *)
    echo "用法: $0 install|install-agent|update|uninstall"
    exit 1
    ;;
esac
