import utils from '../common/utils';
import { EVENT_NAME, STORAGE } from '../common/enum';
import Storage from '@/common/storage';
import emitter from '../common/emmit';

function request(url, options = {}) {
  let user = Storage.get(STORAGE.USER_TOKEN);
  let _headers = options.headers || {},
    headers = {};
  let authorization = user.authorization || '';
  if (user.authorization) {
    headers['Authorization'] = authorization;
  }
  headers['appkey'] = Storage.get(STORAGE.APP_KEY);

  options.headers = utils.extend(_headers, headers);
  return fetch(url, options).then((res) => {
    if (res.status === 401) {
      Storage.remove(STORAGE.USER_TOKEN);
      emitter.$emit(EVENT_NAME.ON_LOGOUT);
      return Promise.reject(new Error('Unauthorized'));
    }
    return res.json();
  });
}

function get(url, params = {}, options = {}) {
  options.method = 'GET';
  return request(url + '?' + new URLSearchParams(params), options);
}

function downloadFile(url, options = {}) {
  let user = Storage.get(STORAGE.USER_TOKEN);
  let _headers = options.headers || {}, headers = {};
  let authorization = user.authorization || '';
  if(user.authorization){
    headers['Authorization'] = authorization
  }
  options.headers = utils.extend(_headers, headers);
  return fetch(url, options).then((res) => {return res.blob()});
}

export { request, get, downloadFile };
