import { request, downloadFile } from './request';
import SERVER_PATH from './api';
import utils from '../common/utils';

function getList(params){
  let { limit = 15, offset = 1, app_key, user_id, start, end, description } = params;
  let url = `${SERVER_PATH.LOG_GET_LIST}?app_key=${app_key}&user_id=${user_id}&limit=${limit}&offset=${offset}&start=${start}&end=${end}&description=${description}`;
  return request(url, { method: 'GET' });
}

function download(params){
  let { app_key, id } = params;
  let url = `${SERVER_PATH.LOG_DOWNLOAD}?app_key=${app_key}&id=${id}`;
  return downloadFile(url, { method: 'GET' });
}

function create(data){
  return request(SERVER_PATH.LOG_CREATE_PULL, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

export default {
  getList,
  create,
  download
}