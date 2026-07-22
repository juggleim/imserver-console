import { request } from './request';
import SERVER_PATH from './api';

// the im node only keeps its log files for a bounded window and grep-scans
// them per query, so every request carries an explicit start and count
const MAX_LOOKBACK = 24 * 60 * 60 * 1000;

function toStart(start) {
  let earliest = Date.now() - MAX_LOOKBACK;
  start = Number(start) || 0;
  return start < earliest ? earliest : start;
}

function buildUrl(path, params) {
  return `${path}?${new URLSearchParams(params)}`;
}

function getConns(params) {
  let { count = 50, start, app_key, user_id } = params;
  return request(buildUrl(SERVER_PATH.CONN_GET_LIST, { app_key, user_id, count, start: toStart(start) }), { method: 'GET' });
}

// the node filters by session but routes the query by user id, so user_id is
// required here even though the session already identifies the connection
function getConn(params) {
  let { count = 50, start, app_key, session, user_id } = params;
  return request(buildUrl(SERVER_PATH.CONN_GET_ONE, { app_key, session, user_id, count, start: toStart(start) }), { method: 'GET' });
}

function getSignalLogs(params) {
  let { count = 50, start, app_key, session, user_id, index = 0 } = params;
  return request(buildUrl(SERVER_PATH.CONN_GET_SIGNAL_LOGS, { app_key, session, user_id, index, count, start: toStart(start) }), { method: 'GET' });
}

export default {
  getConns,
  getConn,
  getSignalLogs,
  MAX_LOOKBACK,
}
