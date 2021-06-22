/*
 * @Author: liyuying
 * @Date: 2021-05-26 11:16:18
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 13:43:13
 * @Description: file content
 */
import React, { useState, useEffect } from 'react';
import { decode } from 'js-base64';
import CodeMirror from '@uiw/react-codemirror';
import { Tree, CessSearchInput } from 'cess-ui';
import { IFileTree } from '@/interfaces';
import { filesToFileTree } from '@/pages/Application/utils';
import 'codemirror/keymap/sublime';
import 'codemirror/theme/darcula.css';
import './index.less';
import { useSelector, IApplicationInstanceCreateState } from 'umi';

interface IProp {
  files: { [key: string]: any };
}
const PreviewYaml = () => {
  const { applicationInstance, applicationTemplate }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  const [code, setCode] = useState('');
  const [fileTree, setFileTree] = useState<IFileTree[]>([]);
  const [selectFile, setSelectFile] = useState('');
  /**
   * 切换展示的file
   */
  const onSelect = (keys: React.Key[], info: any) => {
    setSelectFile((keys[0] as any) || '');
    if (applicationTemplate.files) {
      setCode(
        keys[0] === 'answers.yaml'
          ? applicationInstance.answersYaml || ''
          : decode(applicationTemplate.files[keys[0]] || ''),
      );
    }
  };
  /**
   * 搜索筛选
   * @param e
   */
  const handleTreeFilter = (e: any) => {
    const searchValue = e.target.value;
    const setShowTitle = (file: IFileTree[]) => {
      file.forEach((item: IFileTree) => {
        const index = item.label.toLocaleLowerCase().indexOf(searchValue.toLocaleLowerCase());
        const beforeStr = item.label.substring(0, index);
        const fitStr = item.label.substring(index, index + searchValue.length);
        const afterStr = item.label.substring(index + searchValue.length, item.label.length);
        item.title =
          index > -1 ? (
            <span>
              {beforeStr}
              <span className="site-tree-search-value">{fitStr}</span>
              {afterStr}
            </span>
          ) : (
            <span>{item.label}</span>
          );
        if (item.children) {
          setShowTitle(item.children as IFileTree[]);
        }
      });
    };
    setShowTitle(fileTree);
    setFileTree([...fileTree]);
  };
  useEffect(() => {
    if (selectFile === 'answers.yaml') {
      setCode(applicationInstance.answersYaml || '');
    }
  }, [applicationInstance.answersYaml]);
  useEffect(() => {
    if (applicationTemplate.files) {
      const fileTreeData = filesToFileTree(applicationTemplate.files);
      setFileTree(fileTreeData);
    }
  }, []);
  return (
    <div className="preview-yaml">
      <div className="preview-yaml-tree">
        <CessSearchInput onPressEnter={handleTreeFilter} onChange={() => {}}></CessSearchInput>
        {fileTree && fileTree.length > 0 ? (
          <Tree.DirectoryTree
            defaultExpandAll
            selectedKeys={[selectFile]}
            onSelect={onSelect}
            treeData={fileTree}
          />
        ) : (
          <></>
        )}
      </div>
      <CodeMirror
        value={code}
        options={{
          theme: 'darcula',
          keyMap: 'sublime',
          mode: 'YAML',
          readOnly: true,
        }}
      />
    </div>
  );
};
export default PreviewYaml;
