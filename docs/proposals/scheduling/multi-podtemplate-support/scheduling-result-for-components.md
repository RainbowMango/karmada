# Component-level Scheduling Result for Multiple Pod Templates (Design Discussion Record)

- Status: Draft (design discussion, not yet a formal proposal)
- Date: 2026-08-04
- Related: [multiple-pod-template-support.md](./multiple-pod-template-support.md), issue #5115
- Feature gate: `MultiplePodTemplatesScheduling` (Alpha, default off)

本文档整理自一次设计讨论（会话记录见文末附录），主题是：在多模板（multi-pod-template）
特性下，如何存放调度结果，以及让 `Components` 在传统工作负载场景下最终替代
`spec.replicas` / `spec.replicaRequirements` 的演进路径。

---

## 1. 背景与现状盘点

多模板特性主体已落地：

| 设计点 | 状态 |
|---|---|
| API：`ResourceBindingSpec.Components` / `Component` | ✅ |
| 特性开关 `MultiplePodTemplatesScheduling`（Alpha，默认关闭） | ✅ `pkg/features/features.go` |
| 解释器 hook `GetComponents`（`InterpretComponent`），声明式 Lua / Webhook / 第三方内置三条路径 | ✅ |
| detector 填充逻辑（优先 `GetComponents`，回退 `GetReplicas`） | ✅ `pkg/detector/detector.go` |
| 调度器多模板路径（整体调度、跳过副本分配） | ✅ `pkg/scheduler/core/common.go` |
| Estimator gRPC `MaxAvailableComponentSets` + `EstimateComponents` 插件 | ✅ 精确/通用估算器均已实现 |
| FRQ 准入按 components 计算配额 | ✅ |
| 第三方内置解释器（FlinkDeployment、Spark、Ray、PyTorch、TF、MPI、Volcano Job） | ✅ |
| E2E（`schedule_multi_template_test.go`、FRQ multi-components） | ✅ |
| Estimator 预留（AssumedWorkloads） | ✅ |

尚未完成的关键缺口：

1. 原生解释器 `GetComponents` 未实现（Deployment/StatefulSet 等内置负载不走 Components 路径）；
2. **调度结果 `spec.Clusters []TargetCluster` 仍是传统写法，多模板语境下语义不明**（本文主题）；
3. Descheduler 无 components 概念；
4. `AssignReplicas` 已知问题：单模板 workload replicas=0 时被广播到所有候选集群；
5. legacy 字段弃用路线尚未启动；PaddleJob/XGBoostJob/MXJob/TrainJob 等内置解释器缺失。

### 当前 `TargetCluster` 的语义缺陷

多组件负载走"整体调度、不做副本划分"，`TargetCluster` 只填 `Name`、`Replicas` 留空：

1. **无法表达组件级划分**：如 "jobmanager 1 副本放 A，taskmanager 10 副本 6/4 分到 A/B"。
   注意：多模板负载的一"套"即用户模板的原始声明，没有任何 API 字段表达"套数需求"，
   因此不存在 "set 级 divided 调度"（套数永远是 0 或 1）；`maxSets` 只是可行性信号
   （`maxSets >= 1` 才能装下一套），未来顶多用于打分。真正可能的划分是**拆组件内的副本**；
3. **单模板收敛受阻**：要弃用 `spec.replicas`，`TargetCluster.Replicas` "是哪个 component
   的副本数"必须有明确定义，否则 `ReviseReplica`、graceful eviction、HPA 等消费方语义悬空。

更实质的问题：当前即使调用了 `GetComponents`，Work 下发时仍是把原模板原样复制
（多组件负载没有可用的 `ReviseReplica`），集群内副本全靠模板自身声明，**调度结果并未真正
约束下发**。

---

## 2. 方案对比

### 方案 A：语义重载 —— 多组件时 `Replicas` 表示"套数（sets）"

- 零 API 改动。
- 缺点：一"套"即用户模板的原始声明，没有字段表达"套数需求"，套数实际永远是 0 或 1，
  该字段几乎不携带信息；同一字段两种语义靠 `len(Components)` 区分，消费方易踩坑；
  且永远无法表达组件级划分。
- **不推荐**。

### 方案 B：在 `TargetCluster` 内扩展组件级结果（推荐，讨论结论）

```go
// TargetCluster represents the identifier of a member cluster.
type TargetCluster struct {
	// Name of target cluster.
	Name string `json:"name"`

	// Replicas in target cluster.
	// +optional
	Replicas int32 `json:"replicas,omitempty"`

	// Components represents the per-component replica assignment in this cluster.
	// It is populated only for workloads with multiple pod templates, and only
	// when the MultiplePodTemplatesScheduling feature gate is enabled.
	// Each entry corresponds to an entry in spec.Components by Name.
	// +optional
	Components []TargetComponent `json:"components,omitempty"`
}

// TargetComponent represents the replica assignment of a component in a cluster.
type TargetComponent struct {
	// Name of the component, matching spec.components[*].name.
	// +required
	Name string `json:"name"`

	// Replicas of this component assigned to the cluster.
	// +required
	Replicas int32 `json:"replicas"`
}
```

- 与 `spec.Components` 替代 `spec.Replicas` 的模式**完全对称**：spec 侧按组件声明需求，
  结果侧按组件记录分配，弃用路线一致；
- 向后兼容：所有按 cluster name 索引的消费方不受影响；单模板场景双写
  `Replicas = Components[0].Replicas` 平滑过渡；
- 表达力完整：Duplicated（每集群一套，即每组件全量副本）与未来的组件级副本划分
  （拆某个 component 的 replicas 到多个集群）均可表示。

### 方案 C：新增顶层字段（如 `spec.ComponentAssignments`）替代 `Clusters`

- 最"干净"，但 `spec.Clusters` 是消费方最多的字段之一（binding-controller、scheduler
  patch、graceful eviction、HPA、descheduler、metrics、karmadactl…），双源并存期一致性
  维护成本极高，scheduler patch 两个字段易产生竞态。**不建议**。

设计原则一句话：**让"调度结果"与"调度需求"保持同构** —— `spec.Components` 声明什么，
`spec.Clusters[*].Components` 就记录什么；legacy `Replicas` 双写过渡，随 feature gate
毕业一并弃用。

---

## 3. 按数据流次序的改动清单（方案 B）

### 0）API 定义（源头）

`pkg/apis/work/v1alpha2/binding_types.go`：按上节扩展 `TargetCluster`，新增
`TargetComponent`；`GracefulEvictionTask` 同步增加 `Components []TargetComponent`。
配套 deepcopy / CRD / openapi 代码生成（`hack/update-codegen.sh`、`update-crdgen.sh`）。

### 1）Detector：填充需求侧（已就绪）

`detector.go` 已在 feature gate 开启且 `GetComponents` 可用时填 `spec.Components`。
"传统负载也走 Components" 在此体现为：给原生解释器实现 `GetComponents`，Deployment
产出单元素 `components: [{name: "app", replicas: N, ...}]`（过渡期与 `spec.replicas` 双写）。

### 2）Scheduler：产出结果侧（核心改动一）

改动点在 `pkg/scheduler/core`：

- **多组件（整体调度）**：`AssignReplicas` 由只填 `TargetCluster{Name}` 改为同时填
  `Components`，每组件副本 = `spec.Components[i].Replicas`（Duplicated 语义，每集群一套）：

  ```go
  targetClusters[i] = workv1alpha2.TargetCluster{
      Name:       cluster.Cluster.Name,
      Components: fullComponentSetOf(spec.Components),
  }
  ```

- **单组件 + Divided**：现有分配算法算出每集群 `Replicas` 后，双写
  `Components: [{Name: spec.Components[0].Name, Replicas: r}]`。算法内部仍按标量运算，
  仅出口处多写一份，改动面小。
- scheduler patch binding 路径无结构性改动，字段随 `Clusters` 一起序列化。
- `pkg/util/binding.go` 的聚合工具（按 TargetCluster 汇总副本、`ConvertToClusterNames` 等）
  需感知 `Components`。

### 3）Webhook：结果合法性校验

`pkg/webhook/resourcebinding/validating.go` 追加：

- `clusters[*].components[*].name` ∈ `spec.components[*].name`，且同一 TargetCluster 内不重复；
- feature gate 关闭时拒绝携带该字段（与 `spec.Components` 现行策略一致）；
- FRQ 准入计算改为优先按 `clusters[*].components` 计使用量
  （cluster × component × replicas × resourceRequest），比 "spec.Components × 集群数" 更精确。

### 4）Binding Controller → Work：同步适配（核心改动二）

位置：`pkg/controllers/binding/common.go` 的 `ensureWork`。现状：

```go
if bindingSpec.IsWorkload() {
	if resourceInterpreter.HookEnabled(gvk, InterpreterOperationReviseReplica) {
		clonedWorkload, err = resourceInterpreter.ReviseReplica(clonedWorkload, int64(targetCluster.Replicas))
	}
}
```

问题：`ReviseReplica(object, int64)` 只能改一个标量，无法把 `{jobmanager: 1, taskmanager: 4}`
写回 FlinkDeployment 的多个字段。需新增解释器操作 **`ReviseComponents`**：

```go
// ResourceInterpreter 接口新增
// ReviseComponents revises the per-component replicas of the object
// according to the scheduling result on a specific cluster.
ReviseComponents(object *unstructured.Unstructured, components []workv1alpha2.TargetComponent) (*unstructured.Unstructured, error)
```

`ensureWork` 适配逻辑（与 `GetComponents` ↔ `GetReplicas` 的回退模式对称）：

```go
if bindingSpec.IsWorkload() {
	switch {
	// 组件路径：结果带 Components 且解释器实现了 ReviseComponent
	case len(targetCluster.Components) > 0 &&
		resourceInterpreter.HookEnabled(gvk, InterpreterOperationReviseComponent):
		clonedWorkload, err = resourceInterpreter.ReviseComponents(clonedWorkload, targetCluster.Components)
	// 回退：单组件时退化为标量，复用既有 ReviseReplica
	case resourceInterpreter.HookEnabled(gvk, InterpreterOperationReviseReplica):
		clonedWorkload, err = resourceInterpreter.ReviseReplica(clonedWorkload, int64(replicasFor(targetCluster)))
	}
}
```

其中 `replicasFor` 在 `Replicas` 为空但恰有一个 component 时取 `Components[0].Replicas`，
保证旧解释器配置在新结果格式下仍工作。

解释器框架连锁改动（照 `GetComponents` 样板扩展）：

- `config/v1alpha1`：`CustomizationRules` 新增 `ComponentRevision`
  （Lua 函数 `ReviseComponents(desiredObj, components)` 返回修改后的对象）；
- 声明式 Lua、webhook 自定义、thirdparty 三条路径各加实现；thirdparty 的
  FlinkDeployment/PyTorchJob 等 yaml 同步补 `ReviseComponents` 脚本；
- 原生解释器为 Deployment/StatefulSet 等实现（等价于现有 ReviseReplica）；
- `karmadactl interpret` 与配置校验 webhook（`webhook/configuration/validating.go`
  操作白名单）注册新操作。

注意 `mergeTargetClusters(bindingSpec.Clusters, bindingSpec.RequiredBy)`：
`BindingSnapshot.Clusters` 也是 `[]TargetCluster`，需定义 Components 的合并语义
（同名 component 取 max，与现有 replicas 合并策略对齐）。

### 5）失败迁移 / 优雅驱逐

`GracefulEvictionTask` 增加 `Components` 后：

- taint-manager / 应用失败迁移建任务时，把被驱逐集群的 `TargetCluster.Components` 快照进任务；
- 重调度时 scheduler 依据 `spec.Components` 重新选择替代集群并填组件级结果（第 2 步天然覆盖）；
- eviction controller 的就绪判断暂无需按组件改动，保留组件信息便于未来做部分驱逐。

### 6）下游消费方（次序靠后，可分批）

- **FRQ enforcement controller**：使用量计算改为按 `clusters[*].components` 汇总；
- **HPA（FederatedHPA / cronFHPA）与 descheduler**：读写副本走组件路径
  （descheduler 需配套组件级 `GetUnschedulableReplicas`，属独立迭代）；
- **调度器缓存 AssumedWorkloads / estimator 请求**：assumed 结果由 `spec.Components`
  切换为直接携带 `TargetCluster.Components`，语义更精确；
- **可观测性**：`karmadactl get/describe`、metrics 展示组件级分配。

---

## 4. 建议的落地切分（PR 粒度）

1. API + codegen + webhook 校验（纯增量，无行为变化）；
2. Scheduler 双写 `Components`（feature gate 保护）；
3. `ReviseComponents` 解释器操作全链路 + `ensureWork` 适配（本轮核心）；
4. thirdparty 脚本补齐 + e2e（FlinkDeployment 端到端：模板 → binding → work → 成员集群副本正确）；
5. 下游消费方（FRQ / eviction / HPA）逐个迁移。

1–3 完成后，"多组件负载的调度结果能被正确下发"这一闭环才真正闭上。

## 5. 下一轮迭代建议（整体优先级，供参考）

1. 打通内置工作负载的 Components 路径（原生解释器 `GetComponents`，
   顺带修复 replicas=0 广播的已知问题）；
2. 本文的组件级调度结果（方案 B）；
3. Feature gate Alpha → Beta（默认开启），扩大 e2e 覆盖；
4. Descheduler 支持多模板；
5. 补齐 PaddleJob / XGBoostJob / MXJob / TrainJob 内置解释器；
6. 启动 legacy 字段（`spec.replicas` / `spec.replicaRequirements`）弃用流程；
7. （可选）利用 `maxSets` 做集群打分；（远期）组件跨集群拆分（workload affinity）。

---

## 附录：原始会话记录（摘要）

> **Q1（用户）**：我在思考多模板特性的下一步规划。请阅读设计文档并根据当前代码实现程度
> 分析，建议下一轮迭代做什么。
>
> **A1（助手）**：盘点实现现状（见第 1 节表格），指出缺口：原生解释器 `GetComponents`
> 未实现、descheduler 无组件概念、replicas=0 广播问题、feature gate 仍 Alpha、
> legacy 字段弃用未启动、部分第三方解释器缺失。建议主线为"内置负载支持 + Beta 毕业"。

> **Q2（用户）**：我也想让 Components 在传统工作负载场景下替代 replicas 和
> ReplicaRequirements。但还有个烦人的事情没解决：如何存放调度结果。当前
> `Clusters []TargetCluster` 还是传统写法，多模板语境下语义不明。
>
> **A2（助手）**：分析 `TargetCluster.Replicas` 的三个语义缺陷（无法表达套数、无法表达
> 组件级划分、单模板收敛受阻），给出三个方案（A 语义重载 / B 扩展 TargetCluster /
> C 新增顶层字段），推荐方案 B（见第 2 节），并列出配套改动：`GracefulEvictionTask`
> 对称扩展、`ReviseComponents` 解释器 hook、消费方梯次迁移、webhook 校验。

> **Q3（用户）**：我挺喜欢方案 B 这个 API。请按数据流次序依次介绍相关改动，
> 比如 ResourceBinding 同步到 Work 时如何适配。
>
> **A3（助手）**：按 "API 定义 → detector 填需求 → scheduler 填结果 → webhook 校验 →
> binding-controller ensureWork（`ReviseComponents` 新 hook + 回退 `ReviseReplica`）→
> 驱逐 → 下游消费方" 顺序展开（见第 3 节），并给出 PR 切分建议（见第 4 节）。
> 关键结论：当前调度结果并未真正约束 Work 下发，`ReviseComponents` 全链路是本设计的核心闭环。

> **Q4（用户）**：请把上面的设计保存成一个文件，最好把原始 session 会话也保留下来。
> （即本文档）

