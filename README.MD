# TODO
对初始状态及结果进行定量分析
针对应用冲突优化
预分配数量调度完成后，分配额外的机器处理CPU高负载的机器
调整CPU使用过低的机器
资源消耗公式调整
进行数据分析后分别处理各资源的组合
初始状态迁移前调整机器ID减少迁移
平方根展开

# 已完成
初始状态迁移至最终态的操作步骤

# 各类初始结果
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

#正式分值版本纪录(时间倒序)
---交叉匹配
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2571
total high cpu(1.010000) count=331
MachineDeployPool.DebugPrint machineCount=5571,instanceCount=68219
cpuCost=5583.059498,totalCommands=68219
time=167.991460

---fix
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2571
total high cpu(1.010000) count=728
MachineDeployPool.DebugPrint machineCount=5571,instanceCount=68219
cpuCost=5604.356598,totalCommands=68219
time=172.307525

---加入实例再计算资源消耗
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2571
    cpuCost=1.483749,machineId=1657
    cpuCost=1.405415,machineId=1714
    cpuCost=1.403423,machineId=1655
    cpuCost=1.403423,machineId=1678
    cpuCost=1.403423,machineId=1693
    cpuCost=1.403423,machineId=1672
    cpuCost=1.403423,machineId=1650
    cpuCost=1.401222,machineId=1646
    cpuCost=1.401222,machineId=1660
    cpuCost=1.398347,machineId=1665
    cpuCost=1.367922,machineId=2395
    cpuCost=1.361750,machineId=1383
    cpuCost=1.361735,machineId=1784
    cpuCost=1.333378,machineId=1389
    cpuCost=1.333378,machineId=2122
    cpuCost=1.333378,machineId=1746
    cpuCost=1.333378,machineId=1756
    cpuCost=1.333378,machineId=1794
    cpuCost=1.333378,machineId=1296
    cpuCost=1.328816,machineId=1820
    cpuCost=1.303553,machineId=1682
    cpuCost=1.303553,machineId=2037
    cpuCost=1.303553,machineId=2129
    cpuCost=1.282055,machineId=1780
    cpuCost=1.278889,machineId=1688
    cpuCost=1.259645,machineId=1677
    cpuCost=1.251318,machineId=1783
    cpuCost=1.240034,machineId=1765
    cpuCost=1.233940,machineId=1715
    cpuCost=1.230250,machineId=1385
    cpuCost=1.227271,machineId=1823
    cpuCost=1.219499,machineId=1653
    cpuCost=1.219499,machineId=1425
    cpuCost=1.213787,machineId=1900
    cpuCost=1.211728,machineId=1741
    cpuCost=1.209093,machineId=1497
total high cpu(1.010000) count=726
MachineDeployPool.DebugPrint machineCount=5571,instanceCount=68219
cpuCost=5611.009504,totalCommands=68219
time=168.180176

---初始状态merge到最终状态
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2590
    cpuCost=2.546094,machineId=2066
    cpuCost=2.546094,machineId=2065
    cpuCost=2.509665,machineId=1647
    cpuCost=2.509642,machineId=1646
    cpuCost=2.509642,machineId=1645
    cpuCost=2.190585,machineId=2442
    cpuCost=2.190437,machineId=2441
    cpuCost=2.190437,machineId=2308
    cpuCost=2.190437,machineId=2304
    cpuCost=2.190437,machineId=2309
    cpuCost=2.190437,machineId=2306
    cpuCost=2.190437,machineId=2305
    cpuCost=2.190437,machineId=2307
    cpuCost=2.190437,machineId=2440
total high cpu(1.010000) count=738
MachineDeployPool.DebugPrint machineCount=5590,instanceCount=68219
cpuCost=5647.659305,totalCommands=68219
time=49.348370

---
Strategy:BestFit
MachineDeployPool.DebugPrint
    &{92 288 1024 7 7 9} machineCount=3000
    &{32 64 600 7 3 7} machineCount=2590
    cpuCost=2.546094,machineId=2065
    cpuCost=2.546094,machineId=2066
    cpuCost=2.509665,machineId=1647
    cpuCost=2.509642,machineId=1646
    cpuCost=2.509642,machineId=1645
    cpuCost=2.190585,machineId=2442
    cpuCost=2.190437,machineId=2304
    cpuCost=2.190437,machineId=2308
    cpuCost=2.190437,machineId=2305
    cpuCost=2.190437,machineId=2309
    cpuCost=2.190437,machineId=2441
    cpuCost=2.190437,machineId=2440
    cpuCost=2.190437,machineId=2306
    cpuCost=2.190437,machineId=2307
total high cpu(1.010000) count=738
MachineDeployPool.DebugPrint machineCount=5590,instanceCount=68219
cpuCost=5647.659305,totalCommands=68219
time=49.443332
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