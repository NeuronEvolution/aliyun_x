# 各类结果
- 不带初始状态简单评估排序暴力插入：time= 17.443135,cost=6406.757646,machineCount=5627
- 带初始状态不解除违规约束暴力插入：time= 26.160460,cost=6203.362824,machineCount=5611
- 带初始状态遍历解除违规约束暴力插入：time= 30.019777，cost=6191.333103，machineCount=5610
- 带初始状态解除违规约束解除CPU高负载优先高配插入：time=455.700592，cost=5963.164566，machineCount=5607
- 带初始状态解除违约解除CPU高负载解除低配机器：time=504.441860，cost=6361.214919,machineCount=5655
- 带初始状态解除违约解除CPU高负载解除低配机器排序初始化：time=537.913560，cost=6404.517490,machineCount=5649
- BFD限制CPU不惩罚再BFD放开CPU：time=909.591672，cost=6189.086039,，machineCount=5602
- 磁盘BFD预分配机器 time=75.609463,cost=6132.121646,machineCount=6000
- 均值BFD预分配机器Exp展开 time=49.752896,cost=6101.959981,machineCount=6000

# 1000实例不带初始状态,
- FFD time=0.625260,cost=66.570151,machineCount=64
- BFD time=0.727273,cost=65.265338,machineCount=64
- BFD预分配机器 time=0.839353，cost=64.000000,machineCount=64

# 10000实例不带初始状态
- BFD预分配机器 time=26.318945,cost=684.000000,machineCount=684

# 可调参数
- 每个实例的资源消耗评估公式
- 每台机器的资源消耗评估公式

#正式分值版本纪录(时间倒序)
---
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2621
    cpuCost=2.315893,machineId=1581
    cpuCost=2.315893,machineId=1589
    cpuCost=2.221054,machineId=2466
    cpuCost=2.170818,machineId=844
    cpuCost=2.170818,machineId=843
    cpuCost=2.170818,machineId=842
    cpuCost=2.047093,machineId=2468
    cpuCost=2.047093,machineId=2469
    cpuCost=2.016797,machineId=2470
    cpuCost=2.016797,machineId=2465
    cpuCost=2.016797,machineId=2467
total high cpu(1.010000) count=431
MachineDeployPool.DebugPrint machineCount=5621,instanceCount=68219
cpuCost=5659.310861,totalCommands=68219
time=52.352874