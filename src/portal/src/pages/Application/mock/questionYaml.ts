/*
 * @Author: liyuying
 * @Date: 2021-05-25 14:04:40
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-28 17:31:48
 * @Description: file content
 */
export const QuestionYaml = `
categories:
- storage
namespace: openebs
labels:
  io.rancher.certified: partner
questions:
- variable: defaultImage
  default: "true"
  description: "Use default OpenEBS images"
  label: Use Default Image
  type: boolean
  show_subquestion_if: false
  group: "Container Images"
  subquestions:
  - variable: cstor.target.imageTag
    default: "2.5.0"
    description: "The image tag of cStor Storage Engine Target image"
    type: string
    invalid_chars: "[a-z]{1,}"
    label: Image Tag For OpenEBS cStor Storage Engine Target Image
  - variable: policies.monitoring.image
    default: "quay.io/openebs/m-exporter"
    description: "Default OpeneEBS Volume and pool Exporter image"
    type: multiline
    label: Monitoring Exporter Image
    show_if: "policies.monitoring.enabled=true&&defaultImage=false"
- variable: policies.monitoring.enabled
  default: true
  description: "Enable prometheus monitoring"
  type: boolean
  label: Enable Prometheus Monitoring
  group: "Monitoring Settings"
- variable: compute.unit
  default: ""
  description: ""
  type: ComputeUnit
  label: Compute Unit
  group: "Monitoring Settings"
- variable: apiserver.ports.externalPort
  default: 5656
  description: "Default External Port for OpenEBS API Server"
  type: int
  min: 0
  max: 9999
  label: OpenEBS API Server External Port
  group: "Communication Ports"
- variable: mode
  label: "Monitoring mode"
  description: "Either fullstack for full monitoring or apm for application only monitoring"
  default: "fullstack"
  type: enum
  group: "Agent Configuration (REQUIRED)"
  options:
    - "fullstack"
    - "apm"

- variable: oneagent.apiUrl
  label: "Dynatrace API URL"
  description: "Dynatrace API URL including '\api' path at the end"
  default: "https://ENVIRONMENTID.live.dynatrace.com/api"
  valid_chars: "^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+"
  type: string
  required: true
  group: "Agent Configuration (REQUIRED)"



  #################### Use custom limits settings ###################

- variable: use_custom_limits_settings
  label: "Use custom limits settings"
  description: "Use custom resource limits for the Dynatrace OneAgent"
  default: false
  type: boolean
  group: "Use custom limits settings"

- variable: oneagent.resources.requests.cpu
  label: "CPU resource request"
  description: "Defines the minimum requested CPU by the OneAgent"
  type: string
  group: "Use custom limits settings"
  show_if: "use_custom_limits_settings=true"
`;
