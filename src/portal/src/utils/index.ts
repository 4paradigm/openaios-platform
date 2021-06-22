/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 14:20:25
 * @Description: file content
 */

import { FILE_ICON } from '@/constant/environment';

/**
 * 将params与URL结合转换成带参数的url
 * @param {*} url： URL中的配置的带参数的路径， params:路径所需参数
 */
export function paramsToPath(url: string, params: any) {
  let path: string = url;
  for (let key in params) {
    path = path.replace(`:${key}`, params[key]);
  }
  return path;
}

/**
 * 获取formate时间
 * @param {*} time 时间
 * @param {*} fmt yyyy-MM-dd hh:mm:ss
 */
export function unixTimestampFormate(time: any, fmt = 'yyyy-MM-dd hh:mm:ss') {
  // 处理空
  if (!time) {
    return String(time);
  }
  time = new Date(Number(time));
  const o = {
    'M+': time.getMonth() + 1, // 月份
    'd+': time.getDate(), // 日
    'h+': time.getHours(), // 小时
    'm+': time.getMinutes(), // 分
    's+': time.getSeconds(), // 秒
    'q+': Math.floor((time.getMonth() + 3) / 3), // 季度
  };
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(RegExp.$1, (time.getFullYear() + '').substr(4 - RegExp.$1.length));
  }
  for (const k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length === 1
          ? (o as any)[k]
          : ('00' + (o as any)[k]).substr(('' + (o as any)[k]).length),
      );
    }
  }

  return fmt;
}

// 数组分割，最后5个和前面剩余的为一组
export function arrSplit(arr: any) {
  const results = [];
  const pre = arr.splice(0, arr.length - 1);
  results.push(pre, arr);
  return results;
}
// model中替代setTimeout
export const delay = (ms: number) =>
  new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
// 千分位
export function toThousands(num: number, fractionDigits: number = 2) {
  return num.toFixed(fractionDigits).replace(/\d{1,3}(?=(\d{3})+(\.\d*)?$)/g, '$&,');
}

// 获取文档图标
export function getFileIcon(filePath: string) {
  let fileType = '';
  var startIndex = filePath.lastIndexOf('.');
  if (startIndex !== -1) {
    fileType = filePath.substring(startIndex + 1, filePath.length).toLowerCase();
  }
  switch (fileType) {
    case 'mp4':
    case 'm2v':
    case 'mkv':
    case 'rmvb':
    case 'flv':
    case 'mov':
    case 'avi':
      return FILE_ICON.VIDEO;
    case 'bmp':
    case 'png':
    case 'jpg':
    case 'jpeg':
    case 'gif':
      return FILE_ICON.IMAGE;
    case 'mp3':
    case 'wav':
    case 'wmv':
      return FILE_ICON.VIDEO;
    case 'xls':
    case 'xlsx':
      return FILE_ICON.EXCEL;
    case 'flash':
      return FILE_ICON.FLASH;
    case 'pdf':
      return FILE_ICON.PDF;
    case 'ppt':
    case 'pptx':
      return FILE_ICON.PPT;
    case 'doc':
    case 'docx':
      return FILE_ICON.WORD;
    case 'txt':
      return FILE_ICON.DOC;
    default:
      return FILE_ICON.OTHERS;
  }
}
/**
 *  将对象由{a:{b:{c:2}}} 转化为{a.b.c:2}
 */
export function transferObjectFromLevlesToFlat(levelsObject: any) {
  const flatObject: any = {};
  const iterate = (tempObj: any, parentKeys: string[]) => {
    const keys = Object.keys(tempObj);
    if (keys && keys.length > 0 && keys[0] !== '0') {
      keys.forEach((key) => {
        iterate(tempObj[key], [...parentKeys, key]);
      });
    } else {
      flatObject[parentKeys.join('.')] = tempObj;
    }
  };
  iterate(levelsObject, []);
  return flatObject;
}
/**
 *  将对象由{a.b.c:2}转化为{a:{b:{c:2}}}
 */
export function transferObjectFromFlatToLevles(flatObject: any) {
  const levelsObject: any = {};
  const keys = Object.keys(flatObject);
  const iterate = (tempObj: any, keyList: string[], value: any) => {
    if (keyList && keyList.length > 1) {
      const currentKey = keyList[0];
      if (!tempObj[currentKey]) {
        tempObj[currentKey] = {};
      }
      keyList.shift();
      iterate(tempObj[currentKey], keyList, value);
    } else {
      tempObj[keyList[0]] = value;
    }
  };
  if (keys && keys.length > 0 && keys[0] !== '0') {
    keys.forEach((flatKey) => {
      // 处理转译的【.】
      const decodeFlateKey = flatKey.replaceAll('\\.', '\\&AIOS&\\');
      const originKeyList = decodeFlateKey.split('.');
      const keyList: string[] = [];
      for (let i = 0; i < originKeyList.length; i++) {
        keyList.push(originKeyList[i].replaceAll('\\&AIOS&\\', '.'));
      }
      iterate(levelsObject, keyList, flatObject[flatKey]);
    });
  }
  return levelsObject;
}
