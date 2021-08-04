#!/bin/bash
BASE_IMAGE=nvidia/tensorflow:20.12-tf1-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf1-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf1-py3

BASE_IMAGE=nvidia/tensorflow:20.12-tf2-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf2-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf2-py3

BASE_IMAGE=nvidia/pytorch:20.12-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/pytorch:20.12-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/pytorch:20.12-py3
