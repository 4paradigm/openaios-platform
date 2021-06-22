import { IQuestionItem, IQuestionsGroup, IFileTree } from '@/interfaces';
import { APPLICATION_TEMPLATE_QUESTION_ITEM_REQURE_MESSAGE } from '@/constant/application';
import JsYmal from 'js-yaml';
import React from 'react';
import folderIcon from '@/assets/images/application/folder.svg';
import openFolderIcon from '@/assets/images/application/open-folder.svg';
import { Icon } from 'cess-ui';
/*
 * @Author: liyuying
 * @Date: 2021-06-07 11:43:08
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-16 20:02:26
 * @Description: 应用的工具方法
 */

export const questionYamlToQuestionConfig = (
  questionYaml: string,
): { questionsGroupList: IQuestionsGroup[]; answers: any } => {
  const questionsMap: any = {};
  const answers: any = {};
  const questionsGroupList: IQuestionsGroup[] = [];
  if (questionYaml) {
    const doc: any = JsYmal.load(questionYaml);
    if (doc.questions && doc.questions.length > 0) {
      for (let i = 0; i < doc.questions.length; i++) {
        const item: IQuestionItem = doc.questions[i];
        if (item.default) {
          answers[item.variable] = item.default;
        }
        setQuestionItemRule(item);
        /* 处理子问题 */
        (item.subquestions || []).forEach((subItem) => {
          setQuestionItemRule(subItem);
        });
        /* 设置分组 */
        if (questionsMap[item.group || '']) {
          questionsMap[item.group || ''].push(item);
        } else {
          questionsMap[item.group || ''] = [item];
        }
      }
    }
    for (const gropName in questionsMap) {
      questionsGroupList.push({
        groupName: gropName,
        questionList: questionsMap[gropName],
      });
    }
  }
  return { questionsGroupList, answers };
};
const setQuestionItemRule = (item: IQuestionItem) => {
  item.rules = [];
  // 必填
  if (item.required) {
    item.rules.push({
      required: true,
      message: `${APPLICATION_TEMPLATE_QUESTION_ITEM_REQURE_MESSAGE[item.type]} ${item.label}`,
    });
  }
  // valid_chars 判读值合法的正则
  if (item.valid_chars) {
    item.rules.push({
      validator: async (_: any, value: any) => {
        if (new RegExp(item.valid_chars || '').test(value)) {
          return Promise.resolve();
        }
        return Promise.reject(
          new Error(`${item.label}的值，不符合规则，合法的正则表达式为：${item.valid_chars}`),
        );
      },
    });
  }
  // invalid_chars 判读值不合法的正则
  if (item.invalid_chars) {
    item.rules.push({
      validator: async (_: any, value: any) => {
        if (new RegExp(item.invalid_chars || '').test(value)) {
          return Promise.reject(
            new Error(`${item.label}的值，不符合规则，不合法的正则表达式为：${item.invalid_chars}`),
          );
        }
        return Promise.resolve();
      },
    });
  }
  // 处理是否展示
  if (item.show_if) {
    const conditionArray = item.show_if.split('&&');
    item.showConditions = [];
    conditionArray.forEach((conditionItem) => {
      const splitIndex = conditionItem.indexOf('=');
      const confitionKey = conditionItem.substring(0, splitIndex);
      const value = conditionItem.substring(splitIndex + 1, conditionItem.length);

      (item.showConditions || []).push({ key: confitionKey, value });
    });
  }
};
/**
 * 预览yaml为树形
 * @param files key为文件的path，value为文件内容的base64加码
 */
export const filesToFileTree = (files: { [key: string]: string }): IFileTree[] => {
  const fileTree: IFileTree[] = [];
  const filePathMap = new Map<string, IFileTree>();

  if (files) {
    for (const name in files) {
      const paths = name.split('/');
      let parrantPath = '';
      for (let i = 0; i < paths.length; i++) {
        const currentPath = parrantPath + `${paths[i]}-${i}`;
        const parrentFile: IFileTree = {
          title: paths[i],
          label: paths[i],
          key: currentPath,
          parrentKey: parrantPath,
          isLeaf: false,
        };
        parrantPath = currentPath;
        if (i === paths.length - 1) {
          parrentFile.isLeaf = true;
          parrentFile.key = name; // 最后一个的key需要file的全路径
        } else {
          parrentFile.children = [];
          parrentFile.icon = ({ expanded }) => {
            return expanded ? <Icon type="open-folder" /> : <Icon type="folder" />;
          };
        }
        filePathMap.set(parrentFile.key + '', parrentFile);
      }
    }
  }
  filePathMap.forEach((value: IFileTree, key: string) => {
    if (value.parrentKey) {
      if (filePathMap.has(value.parrentKey)) {
        ((filePathMap.get(value.parrentKey) as any) || {}).children.push(value);
      }
    } else {
      fileTree.push(value);
    }
  });
  return fileTree;
};
/**
 * 上传的tar包解析为{'filePath':filecontent}
 * @param tarObject
 */
export const transformTarObjectToFiles = (tarObject: any, callBack: any) => {
  const files: any = {};
  let existChartYaml: boolean = false;
  let existReadMeMd: boolean = false;
  let chartTemplateName: string = ''; // tar的第一层名称
  let waitSetFileCount = 0;
  try {
    const iterateFile = (fileDirObj: any, parrentPath: string) => {
      for (const fileName in fileDirObj) {
        const childFileOrDir = fileDirObj[fileName];
        if (childFileOrDir instanceof File) {
          waitSetFileCount += 1;
          let reader = new FileReader();
          reader.readAsText(childFileOrDir, 'utf-8');
          // eslint-disable-next-line no-loop-func
          reader.onload = () => {
            files[`${parrentPath}${fileName}`] = reader.result;
            waitSetFileCount = waitSetFileCount - 1;
            if (waitSetFileCount === 0 && existChartYaml && existReadMeMd) {
              callBack(chartTemplateName, files);
            }
          };
          if (fileName === `Chart.yaml`) {
            existChartYaml = true;
          } else if (fileName === `README.md`) {
            existReadMeMd = true;
          }
        } else {
          iterateFile(childFileOrDir, `${parrentPath}${fileName}/`);
        }
      }
    };
    if (tarObject) {
      for (const chartName in tarObject) {
        chartTemplateName = chartName;
        const rootFilesObj = tarObject[chartName];
        iterateFile(rootFilesObj, '');
      }
    }
  } catch (error) {
    throw '上传的tar包不符合规范！';
  }

  if (!existChartYaml) {
    throw '上传的tar包根目录中不存在Chart.yaml文件，请您确认后重新上传！';
  } else if (!existReadMeMd) {
    throw '上传的tar包根目录中不存在README.md文件，请您确认后重新上传！';
  }
};
