#!/bin/bash
#
# Copyright © 2021 peizhaoyou <peizhaoyou@4paradigm.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

BASE_IMAGE=nvidia/tensorflow:20.12-tf1-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf1-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf1-py3

BASE_IMAGE=nvidia/tensorflow:20.12-tf2-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf2-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/tensorflow:20.12-tf2-py3

BASE_IMAGE=nvidia/pytorch:20.12-py3
docker build --build-arg BASE_IMAGE=${BASE_IMAGE} -t m7-model-inf01:30003/fuhao/pineapple/env/pytorch:20.12-py3 .
docker push m7-model-inf01:30003/fuhao/pineapple/env/pytorch:20.12-py3
