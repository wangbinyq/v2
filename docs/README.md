# Web框架分析文档

查看 [web-frameworks-analysis.md](./web-frameworks-analysis.md) 获取完整的Golang和Rust Web项目分析报告。

## 快速链接

### Golang 推荐项目
- [Fiber](https://github.com/gofiber/fiber) - Express风格的高性能框架
- [Echo](https://github.com/labstack/echo) - 简洁的微服务框架
- [GORM](https://github.com/go-gorm/gorm) - 强大的ORM库

### Rust 推荐项目
- [Actix-Web](https://github.com/actix/actix-web) - 极速Web框架
- [Rocket](https://github.com/rwf2/Rocket) - 类型安全的Web框架
- [Yew](https://github.com/yewstack/yew) - Rust前端框架
- [Dioxus](https://github.com/DioxusLabs/dioxus) - 跨平台全栈框架

## 快速决策

| 需求 | 推荐方案 |
|------|---------|
| 快速开发 | Golang (Fiber/Echo) |
| 极致性能 | Rust (Actix-Web) |
| 全栈开发 | Rust (Leptos) 或 Golang (Buffalo) |
| WebAssembly | Rust (Yew/Dioxus) |
| 微服务 | Golang (Echo/Fiber) |
