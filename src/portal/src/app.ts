import 'antd/dist/antd.less';
import 'cess-ui/lib/less/index.less';
export const dva = {
  config: {
    onError(err: ErrorEvent) {
      err.preventDefault();
      console.error(err.message);
    },
  },
};
