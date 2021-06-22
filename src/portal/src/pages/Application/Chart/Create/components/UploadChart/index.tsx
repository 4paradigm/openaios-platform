/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:41
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 12:06:40
 * @Description: file content
 */
import React from 'react';
import { encode } from 'js-base64';
import JsYmal from 'js-yaml';
import { Upload, Button, Alert, CessCard, Grid, Info, notification } from 'cess-ui';
import { transformTarObjectToFiles } from '../../../../utils';
import { CloudUploadOutlined } from '@ant-design/icons';
import { Archive } from 'libarchive.js/main.js';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationChartCreateAction,
} from 'umi';
import './index.less';

Archive.init({
  workerUrl: '/worker-bundle.js',
});

/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-04 11:17:12
 * @Description: 创建我的应用--chart 上传
 */

const ApplicationChartUpload = () => {
  const dispatch = useDispatch();
  const { loading, chartData }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  /**
   * 上传tar文件
   */
  const handlelUpload = (file: any) => {
    console.log(file);
    if (file.size > 1024 * 1024 * 4) {
      notification.error({ message: '不允许上传大于4MB的文件' });
      return false;
    }

    if (typeof FileReader === 'undefined') {
      notification.error({ message: '您的浏览器不支持FileReader接口' });
      return false;
    }
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: {
        loading: true,
        chartData: {
          metadata: {
            name: '',
            description: '',
            version: '',
            url: '',
            icon_link: '',
          },
          files: '',
        },
      },
    });
    extractTarFiles(file);
    return true;
  };
  /**
   * 设置files的回调并解析chart.yaml
   */
  const setChartFile = (chartName: string, files: any) => {
    let chartYaml = '';
    if (files) {
      for (const name in files) {
        if (name === 'Chart.yaml') {
          chartYaml = files[name] || '';
        }
        files[name] = encode(files[name]);
      }
    }
    const chartYamlObj: any = JsYmal.load(chartYaml);
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: {
        chartData: {
          metadata: {
            name: chartName,
            description: chartYamlObj.description || '',
            version: chartYamlObj.version || '',
            url: '',
            icon_link: chartYamlObj.icon || '',
          },
          files: files,
        },
      },
    });
  };
  /**
   * 解析Tar包
   * @param file
   */
  const extractTarFiles = async (file: any) => {
    const archive = await Archive.open(file);
    const chartTar = await archive.extractFiles();
    try {
      transformTarObjectToFiles(chartTar, setChartFile);
    } catch (error) {
      notification.error({ message: error });
    }
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { loading: false },
    });
  };
  return (
    <div className="application-chart-upload">
      <Alert
        type="info"
        message="创建指南"
        description={
          <div>
            创建应用需要提前准备，把您的应用按照所需格式准备好并打包完毕，详情请点击
            <Button type="link">查看</Button>
          </div>
        }
      ></Alert>
      <CessCard title="应用上传">
        <Upload
          action=""
          accept=".tar,.zip,.gz,.bz"
          beforeUpload={handlelUpload}
          showUploadList={false}
        >
          <Button loading={loading} icon={<CloudUploadOutlined />} type="primary" size="large">
            请上传你的应用
          </Button>
        </Upload>
        <div>
          <Info>
            文件格式支持 TAR, TAR.GZ, TAR.BZ 和
            ZIP)，最大支持4MB。若配置包符合开发规范，则提示上传成功
          </Info>
        </div>
      </CessCard>
      {chartData.metadata?.name ? (
        <CessCard smallGap={true} title="应用信息">
          <Grid
            list={[
              {
                label: '应用名称',
                value: chartData.metadata?.name,
              },
              {
                label: '应用版本',
                value: chartData.metadata?.version,
              },
              { label: '应用描述', value: chartData.metadata?.description },
            ]}
          />
        </CessCard>
      ) : (
        <></>
      )}
    </div>
  );
};
export default ApplicationChartUpload;
