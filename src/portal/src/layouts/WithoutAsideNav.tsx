import React from 'react';
import styles from './index.less';

const BasicLayout: React.FC = (props) => {
  return (
    <div className={styles.layoutWithoutAsideNav}>
      <div>{props.children}</div>
    </div>
  );
};

export default BasicLayout;
