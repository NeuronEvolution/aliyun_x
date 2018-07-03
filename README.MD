# 各类结果
- 不带初始状态简单评估排序暴力插入：time= 17.443135,cost=6406.757646,machineCount=5627
- 带初始状态不解除违规约束暴力插入：time= 26.160460,cost=6203.362824,machineCount=5611
- 带初始状态遍历解除违规约束暴力插入：time= 30.019777，cost=6191.333103，machineCount=5610
- 带初始状态解除违规约束解除CPU高负载优先高配插入：time=455.700592，cost=5963.164566，totalCommands=46408，machineCount=5607
- 带初始状态解除违约解除CPU高负载解除低配机器：time=504.441860，cost=6361.214919,totalCommands=50749，machineCount=5655
- 带初始状态解除违约解除CPU高负载解除低配机器排序初始化：time=537.913560，cost=6404.517490,totalCommands=0，machineCount=5649
- BFD限制CPU不惩罚再BFD放开CPU：time=909.591672，cost=6189.086039,totalCommands=0，machineCount=5602
- 磁盘BFD预分配机器 time=75.609463,cost=6132.121646,totalCommands=0,machineCount=6000
- 均值BFD预分配机器Exp展开 time=49.752896,cost=6101.959981,totalCommands=40548,machineCount=6000

# 1000实例不带初始状态,
- FFD time=0.625260,cost=66.570151,totalCommands=0,machineCount=64
- BFD time=0.727273,cost=65.265338,totalCommands=0,machineCount=64
- BFD预分配机器 time=0.839353，cost=64.000000,totalCommands=0，machineCount=64

# 10000实例不带初始状态
- BFD预分配机器 time=26.318945,cost=684.000000,totalCommands=0,machineCount=684

# 可调参数
- 每个实例的资源消耗评估公式
- 每台机器的资源消耗评估公式