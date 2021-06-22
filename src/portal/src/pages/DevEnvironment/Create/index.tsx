import React, { useState, useEffect } from 'react';
import {
  Breadcrumb,
  CessCard,
  Input,
  Icon,
  Select,
  Button,
  Checkbox,
  Radio,
  CessCreateBtnBar,
  message,
} from 'cess-ui';
import { DEV_ENVIRONMENT } from '@/router/url';
import { IImage, IPath } from '@/interfaces/bussiness';
import { useSelector, useDispatch } from 'react-redux';
import { IEnvironmentCreateState, EnvironmentCreateAction } from '../models/environment-create';
import './index.less';
import { history } from 'umi';
import ComputeUnitRadio from '@/components/ComputeUnitRadio';

const { Option, OptGroup } = Select;
const { LeftBtnBar, RightBtnBar } = CessCreateBtnBar;

const breadcrumb = (
  <Breadcrumb>
    <Breadcrumb.Item href={DEV_ENVIRONMENT}>开发环境</Breadcrumb.Item>
    <Breadcrumb.Item>创建开发环境</Breadcrumb.Item>
  </Breadcrumb>
);

const EnvironmentCreate = () => {
  const dispatch = useDispatch();
  const {
    publicMirror,
    privateMirror,
    environmentData,
    loading,
    name,
    interactList,
    errorMap,
  }: IEnvironmentCreateState = useSelector((state: any) => state.environmentCreate);
  const [interactkey, setInteractKey] = useState([]);

  useEffect(() => {
    dispatch({
      type: EnvironmentCreateAction.GET_DATA,
    });
  }, []);

  // 大保存校验 (由于后端创建环境时会返回其他不确定的校验错误信息，所以自定义前端验证规则)
  const variableKey = () => {
    const errMap = { ...errorMap };
    // 校验名称
    if (name) {
      if (!name.trim()) {
        errMap.name = '环境名称不可以只包含空格';
      } else {
        const fitNameList = new RegExp(/[a-z]([-a-z0-9]*[a-z0-9])?/).exec(name) || [];
        if (fitNameList.length > 0 && fitNameList[0] === name) {
          errMap.name = '';
        } else {
          errMap.name = '以小写字母开头，只能包含小写英文字母、数字和中划线，且不能以中划线结尾';
        }
      }
    } else {
      errMap.name = '环境名称不可以为空';
    }
    // 校验算力规格是否选择
    if (environmentData.compute_unit) {
      errMap.compute_unit = '';
    } else {
      errMap.compute_unit = '请选择算力规格';
    }
    // 校验是否选择镜像
    if (environmentData.image.repo) {
      errMap.image = '';
    } else {
      errMap.image = '请选择镜像';
    }
    // TODO引入数据
    if (environmentData.mounts.find((mount: IPath) => !mount.mountpath || !mount.subpath)) {
      errMap.mount = '引入数据的路径均不能为空';
    } else {
      errMap.mount = '';
    }
    // 交互方式
    if (environmentData.ssh.enable || environmentData.jupyter.enable) {
      errMap.interact = '';
      if (
        (environmentData.ssh.enable && environmentData.ssh['id_rsa.pub']) ||
        !environmentData.ssh.enable
      ) {
        errMap.ssh = '';
      } else {
        errMap.ssh = '请输入SSH服务的id_rsa.pub';
      }
      if (environmentData.jupyter.enable) {
        if (environmentData.jupyter.token) {
          if (/^\d+$/.test(environmentData.jupyter.token)) {
            errMap.jupyter = 'JupyterLab的token不允许为纯数字';
          } else {
            errMap.jupyter = '';
          }
        } else {
          errMap.jupyter = '请输入JupyterLab的token';
        }
      } else {
        errMap.jupyter = '';
      }
    } else {
      // 允许不选择交互方式
      // errMap.interact = '请选择交互方式';
    }
    dispatch({
      type: EnvironmentCreateAction.UPDATE_STATE,
      payload: {
        errorMap: errMap,
      },
    });
    return errMap;
  };

  const changeKey = (e: any, key: string) => {
    if (key === 'interact') {
      const data: any = { ...environmentData };
      interactList.forEach((inter: any) => {
        data[inter.key].enable = e.indexOf(inter.key) > -1 ? true : false;
      });
      dispatch({
        type: EnvironmentCreateAction.UPDATE_STATE,
        payload: {
          environmentData: data,
        },
      });
      setInteractKey(e);
    } else {
      dispatch({
        type:
          key === 'name'
            ? EnvironmentCreateAction.UPDATE_STATE
            : EnvironmentCreateAction.UPDATE_DATA,
        payload: {
          [key]: e.target.value,
        },
      });
    }
  };

  const changeInterKey = (e: any, interact: any) => {
    const curInter: any = { ...(environmentData as any)[interact.key] };
    curInter[interact.pswKey] = e.target.value;
    dispatch({
      type: EnvironmentCreateAction.UPDATE_DATA,
      payload: {
        [interact.key]: curInter,
      },
    });
  };

  const handelCreateEnviron = () => {
    const data = variableKey();
    if (Object.keys(data).find((key) => data[key])) {
      message.error('环境信息输入有误，请检查');
      return;
    }
    dispatch({
      type: EnvironmentCreateAction.CREATE_ENVIRONMENT,
      payload: {
        name,
        data: environmentData,
      },
    });
  };

  const back = () => {
    dispatch({
      type: EnvironmentCreateAction.UPDATE_STATE,
      payload: {
        name: '',
      },
    });
    history.push('/devEnvironment');
  };

  const mirrorChange = (value: any, option: any) => {
    const { tag, source, repo } = option;
    dispatch({
      type: EnvironmentCreateAction.UPDATE_DATA,
      payload: {
        image: {
          repo,
          tag,
          source,
        },
      },
    });
  };

  const handelMountChange = (
    type: 'add' | 'delete' | 'edit',
    index: number = -1,
    key: 'subpath' | 'mountpath' = 'subpath',
    e?: any,
  ) => {
    const mounts = [...environmentData.mounts];
    if (type === 'edit' && e) {
      mounts[index][key] = e.target.value;
    } else {
      type === 'add' ? mounts.push({ subpath: '', mountpath: '' }) : mounts.splice(index, 1);
    }
    dispatch({
      type: EnvironmentCreateAction.UPDATE_DATA,
      payload: {
        mounts,
      },
    });
  };

  return (
    <div className="environment-create comm-create-page">
      {breadcrumb}
      <div className="environment-create-container">
        <CessCard>
          <h3 className="card-title">
            环境名称 <span>*</span>
          </h3>
          <Input
            maxLength={30}
            value={name}
            onChange={(e) => changeKey(e, 'name')}
            className={errorMap['name'] ? 'error-input' : ''}
            placeholder="请输入环境名称，以小写字母开头，只能包含小写英文字母、数字和中划线，且不能以中划线结尾"
            autoComplete="off"
          />
          {errorMap['name'] && <p className="error-code">{errorMap['name']}</p>}
        </CessCard>
        <CessCard>
          <h3 className="card-title">
            选择算力规格 <span>*</span>
          </h3>
          <p className="desc">算力规格决定了在容器内执行所提供的资源多少</p>
          <Radio.Group
            value={environmentData.compute_unit}
            onChange={(e) => changeKey(e, 'compute_unit')}
          >
            <ComputeUnitRadio></ComputeUnitRadio>
          </Radio.Group>
          {errorMap['compute_unit'] && <p className="error-code">{errorMap['compute_unit']}</p>}
        </CessCard>
        <CessCard>
          <h3 className="card-title">
            选择镜像 <span>*</span>
          </h3>
          <Select
            onChange={(val, options) => mirrorChange(val, options)}
            showSearch
            style={{ width: '516px' }}
            placeholder="请选择镜像"
            dropdownClassName="mirror-dropdown"
          >
            {publicMirror.length > 0 && (
              <OptGroup label="公有环境镜像">
                {publicMirror.map((mirror: IImage) => {
                  return (
                    <Option
                      key={'public' + mirror.repo + mirror.tags[0]}
                      value={'public' + mirror.repo + mirror.tags.join(',')}
                      source="public"
                      repo={mirror.repo}
                      tag={mirror.tags[0]}
                    >
                      {mirror.repo}：{mirror.tags}
                    </Option>
                  );
                })}
              </OptGroup>
            )}
            {privateMirror.length > 0 && (
              <OptGroup label="私有镜像">
                {privateMirror.map((mirror: IImage) => {
                  return (
                    <Option
                      key={'private' + mirror.repo + mirror.tags[0]}
                      value={'private' + mirror.repo + mirror.tags.join(',')}
                      repo={mirror.repo}
                      source="private"
                      tag={mirror.tags[0]}
                    >
                      {mirror.repo}：{mirror.tags}
                    </Option>
                  );
                })}
              </OptGroup>
            )}
          </Select>
          {errorMap['image'] && <p className="error-code">{errorMap['image']}</p>}
        </CessCard>
        <CessCard>
          <h3 className="card-title">
            引入数据<span></span>
          </h3>
          <p className="desc">
            您可以选择mount数据文件夹到容器内，双方数据即可完成同步，请注意设置绝对路径
          </p>
          <div className="mount-group">
            {environmentData.mounts.length > 0 && (
              <>
                <div className="flex-row mount-group-header">
                  <div className="col col-line"></div>
                  <div className="col flex-auto">设置路径</div>
                  <div className="col flex-auto">目标路径</div>
                  <div className="col col-operate">操作</div>
                </div>
                <div className="mount-group-body">
                  {environmentData.mounts.map((mount: IPath, index: number) => (
                    <div key={`mount${index}`} className="flex-row mount-group-list">
                      <div className="col flex-auto">
                        <Input
                          placeholder="填写相对路径"
                          value={mount.subpath}
                          onChange={(e: any) => {
                            handelMountChange('edit', index, 'subpath', e);
                          }}
                        />
                      </div>
                      <div className="col flex-auto">
                        <Input
                          placeholder="填写绝对路径"
                          value={mount.mountpath}
                          onChange={(e: any) => {
                            handelMountChange('edit', index, 'mountpath', e);
                          }}
                        />
                      </div>
                      <div className="col col-operate">
                        <Button
                          type="link"
                          onClick={() => handelMountChange('delete', index)}
                          icon={<Icon type="delete" />}
                        >
                          删除
                        </Button>
                      </div>
                    </div>
                  ))}
                  {errorMap['mount'] && <p className="error-code">{errorMap['mount']}</p>}
                  {environmentData.mounts.length > 1 && <div className="group-line-box"></div>}
                </div>
              </>
            )}
          </div>
          <span className="add-path-span" onClick={() => handelMountChange('add')}>
            <Icon type="plus-circle" /> 添加数据
          </span>
        </CessCard>
        <CessCard>
          <h3 className="card-title">
            交互方式 <span></span>
          </h3>
          <p className="desc">您可以选择使用哪种环境交互方式，交互方式可多选</p>
          <Checkbox.Group value={interactkey} onChange={(e) => changeKey(e, 'interact')}>
            {interactList.map((interact: any) => {
              const curEnvirInteract: any = (environmentData as any)[interact.key];
              return (
                <Checkbox value={interact.key} key={interact.key}>
                  <div
                    className={`title-block interact-block ${
                      curEnvirInteract.enable ? 'unit-checked' : ''
                    }`}
                  >
                    <img src={interact.icon} alt={interact.name} />
                    <p className="title">{interact.name}</p>
                    <p className="sub-title">{interact.desc}</p>
                  </div>
                  {curEnvirInteract.enable && (
                    <div className="token-wrap">
                      <Input.TextArea
                        rows={3}
                        value={(environmentData as any)[interact.key][interact.pswKey]}
                        onChange={(e: any) => changeInterKey(e, interact)}
                        placeholder={`请输入${interact.key}的${interact.pswKey}`}
                      />
                      {errorMap[interact.key] && (
                        <p className="error-code">{errorMap[interact.key]}</p>
                      )}
                    </div>
                  )}
                </Checkbox>
              );
            })}
          </Checkbox.Group>
          {errorMap.interact && <p className="error-code">{errorMap.interact}</p>}
        </CessCard>
      </div>
      <CessCreateBtnBar>
        <LeftBtnBar>
          <Button onClick={back}>退出</Button>
        </LeftBtnBar>
        <RightBtnBar>
          <Button type="primary" loading={loading} onClick={handelCreateEnviron}>
            准备环境
          </Button>
        </RightBtnBar>
      </CessCreateBtnBar>
    </div>
  );
};

export default EnvironmentCreate;
