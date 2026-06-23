<div align="center">

# 🛰️ tg-proxy-list — 免费 Telegram 代理列表（已验证的 MTProto 代理）

**新鲜、自动验证的 Telegram MTProto 代理 —— 经过连接测试、按速度排序、标注国家，每 6 小时更新一次。**
一键连接。无需 App、无需注册、无广告。

### 👉 **[打开列表并连接 →](https://zakky8.github.io/tg-proxy-list/)** 👈

[English](README.md) · [Русский](README_RU.md) · [فارسی](README_FA.md) · 中文

</div>

---

## 这是什么

**tg-proxy-list 是一个免费、持续更新的公共 Telegram MTProto 代理列表。** 它能帮助你在 Telegram 被**屏蔽或限速**的地方使用它，无需安装任何东西。

与普通的代理列表不同，这里的每个代理在发布前都会**真正验证**：

- ✅ DNS 解析 + TCP 连接，并**测量延迟**
- ✅ 对 `ee` 代理进行 **FakeTLS 握手**测试
- 🌍 标注**国家**
- 📈 跟踪**在线率（uptime）**
- 🔁 **每 6 小时**重新检测 —— 失效代理自动剔除

> **诚实标注。** 我们从不声称代理「100% 可用」。状态为 **`reachable`**（TCP 有响应）或 **`handshake_ok`**（还通过了 FakeTLS 握手）。公共代理随时可能下线 —— 一个不行就换下一个。

## 快速开始

- **网站：** [zakky8.github.io/tg-proxy-list](https://zakky8.github.io/tg-proxy-list/) —— 按国家筛选、按速度排序、点击 **Connect**。
- **列表：** 打开 [`all_proxies.txt`](all_proxies.txt)，在手机上点击任意 `https://t.me/proxy?...` 链接。
- **按国家：** `by_country/CN.txt`、`by_country/IR.txt` 等。

## 免责声明

这些是来自公开来源的**公共**代理，按「现状」提供，不作任何保证。本项目不运营这些代理。请合法使用。国家数据 © [DB-IP](https://db-ip.com)（CC-BY-4.0）。许可证：[MIT](LICENSE)。
