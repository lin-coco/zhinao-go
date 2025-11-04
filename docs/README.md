# 360智脑 Go SDK 文档

本目录包含 360智脑 Go SDK 的详细文档。

## 📚 文档列表

### 1. [完整指南 (GUIDE.md)](./GUIDE.md) ⭐ 推荐

**最全面的文档**，整合了以下内容：
- 快速开始和安装
- 架构设计详解
- 核心功能使用
- 测试策略和最佳实践
- 扩展开发指南

**适合**：
- 新用户快速上手
- 了解 SDK 架构
- 学习最佳实践
- 进行扩展开发

### 2. [SDK 对比分析 (COMPARISON.md)](./COMPARISON.md)

对比 360智脑 Go SDK 与其他流行 SDK 的设计：
- **go-openai** - OpenAI 官方 SDK
- **go-moonshot** - Moonshot AI SDK
- **deepseek-go** - DeepSeek SDK

**内容包括**：
- 架构对比
- 功能对比矩阵
- 环境变量支持对比
- 设计优势分析
- 易用性评分

**适合**：
- 技术选型参考
- 了解 SDK 优势
- 学习设计模式
- 架构设计参考

## 📖 快速导航

### 我想了解...

- **如何快速开始？** → [GUIDE.md - 快速开始](./GUIDE.md#快速开始)
- **架构是怎样的？** → [GUIDE.md - 架构设计](./GUIDE.md#架构设计)
- **如何配置客户端？** → [GUIDE.md - 客户端配置](./GUIDE.md#1-客户端配置)
- **如何使用 Builder？** → [GUIDE.md - Builder 模式](./GUIDE.md#2-builder-模式)
- **如何处理错误？** → [GUIDE.md - 错误处理](./GUIDE.md#4-错误处理)
- **如何编写测试？** → [GUIDE.md - 测试指南](./GUIDE.md#测试指南)
- **最佳实践是什么？** → [GUIDE.md - 最佳实践](./GUIDE.md#最佳实践)
- **如何扩展功能？** → [GUIDE.md - 扩展开发](./GUIDE.md#扩展开发)
- **与其他 SDK 对比？** → [COMPARISON.md](./COMPARISON.md)

### 我想看...

- **代码示例** → [../examples/](../examples/)
- **API 文档** → [主 README](../README.md)
- **更新日志** → [CHANGELOG](../CHANGELOG.md)

## 🔧 维护指南

### 更新文档

当需要更新文档时：

1. **功能说明、使用指南** → 更新 `GUIDE.md`
2. **SDK 对比、设计分析** → 更新 `COMPARISON.md`
3. **示例代码** → 更新 `../examples/`
4. **API 参考** → 更新主 `README.md`

### 添加新内容

添加新内容时，考虑：

- **是否属于使用指南？** → 添加到 `GUIDE.md` 相应章节
- **是否涉及设计对比？** → 添加到 `COMPARISON.md`
- **是否需要示例代码？** → 创建新的 example
- **是否需要独立文档？** → 评估必要性，避免过度分散

## 📝 文档编写原则

1. **清晰性** - 使用简洁明了的语言
2. **完整性** - 提供足够的上下文和示例
3. **可维护性** - 避免重复，保持结构清晰
4. **实用性** - 提供可执行的代码示例
5. **及时性** - 保持文档与代码同步

## 🤝 贡献文档

欢迎改进文档！提交文档更新时：

1. 确保语言清晰准确
2. 提供完整的代码示例
3. 保持格式一致
4. 更新相关索引
5. 避免创建冗余文档

## 📮 反馈

如果发现文档问题或有改进建议：

- 提交 Issue: https://github.com/lin-coco/zhinao-go/issues
- 提交 Pull Request: https://github.com/lin-coco/zhinao-go/pulls

---

**快速链接**：[主 README](../README.md) | [完整指南](./GUIDE.md) | [SDK 对比](./COMPARISON.md) | [示例代码](../examples/)
