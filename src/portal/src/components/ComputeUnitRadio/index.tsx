import React, { useEffect } from 'react';
import { Radio, OverflowToolTip, Icon } from 'cess-ui';
import { IComputeUnitRadioState, ComputeUnitRadioAction } from 'umi';
import { useSelector, useDispatch } from 'react-redux';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-05-27 18:55:02
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-16 19:46:17
 * @Description: file content
 */
const ComputeUnitRadio = () => {
  const dispatch = useDispatch();
  const { computeUnitList }: IComputeUnitRadioState = useSelector(
    (state: any) => state.computeUnitRadio,
  );
  useEffect(() => {
    dispatch({
      type: ComputeUnitRadioAction.GET_DATA,
    });
  }, []);
  return (
    <>
      {computeUnitList.map((compute) => {
        return (
          <Radio value={compute.id} key={compute.id} className="compute-unit-radio">
            <div className={`title-block`}>
              <p className="title">{compute.id}</p>
              <p className="sub-title">
                <Icon type="price" className="unit-price-icon" />
                <OverflowToolTip
                  title={`${compute.price}/min | ${compute.description}`}
                  width={220}
                ></OverflowToolTip>
              </p>
            </div>
          </Radio>
        );
      })}
    </>
  );
};
export default ComputeUnitRadio;
