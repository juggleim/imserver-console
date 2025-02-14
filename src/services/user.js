import { request } from './request';
import SERVER_PATH from './api';
import utils from '../common/utils';
function login(data){
  return request(SERVER_PATH.USER_LOGIN, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

function updatePwd(data){
  return request(SERVER_PATH.USER_UPDATE_PWD, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

function getUsers(params = {}){
  let { limit = 50, offset = 1 } = params;
  let url = `${SERVER_PATH.USER_LIST}?limit=${limit}&offset=${offset}`;
  return request(url, {
    method: 'GET'
  });
}

function add(data){
  return request(SERVER_PATH.USER_ADD, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}
function remove(data){
  return request(SERVER_PATH.USER_DELETE, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}
function disable(data){
  return request(SERVER_PATH.USER_DISABLE, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

export default {
  login,
  updatePwd,
  getUsers,
  add,
  remove,
  disable,
}