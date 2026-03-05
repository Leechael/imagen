# Nano Banana Image Skill

基于 [nano-banana-2-skill](https://github.com/kingbootoshi/nano-banana-2-skill/) 重写的 Go 版 AI 图片生成 CLI，使用 Gemini image models。

相比原版的主要变化：
- 重构代码结构，补充完整的 `--help` 输出
- 支持从 stdin 管道读取 prompt（`echo "prompt" | nano-banana`）
- 补全 human 模式日志（`[nano-banana]` 前缀）
- `--costs --json` 输出包含 per-model 明细
- 添加单元测试
- GitHub Actions CI / Release 工作流

## 安装

要求：Go 1.25+（透明模式另需 FFmpeg + ImageMagick）

```bash
go mod tidy
go build -o ~/.local/bin/nano-banana ./cmd/nano-banana
```

确保 `~/.local/bin` 在 PATH 中。

## API Key

按以下优先级读取 Gemini API Key：

1. `--api-key` 命令行参数
2. `GEMINI_API_KEY` 环境变量
3. 当前目录 `.env`
4. 可执行文件上级目录 `.env`
5. `~/.nano-banana/.env`

```bash
# 推荐：环境变量
export GEMINI_API_KEY=your_key_here

# 或写入文件
mkdir -p ~/.nano-banana
echo "GEMINI_API_KEY=your_key_here" > ~/.nano-banana/.env
```

获取 API Key：https://aistudio.google.com/apikey

## Usage

```bash
nano-banana "minimal dashboard UI with dark theme"
nano-banana "luxury product mockup" -o product -s 2K
nano-banana "cinematic scene" -a 16:9 -s 4K
nano-banana "change to white background" -r input.png -o output
nano-banana "robot mascot" -t -o mascot

# 从 stdin 读取 prompt
echo "a cat in a spacesuit" | nano-banana
pbpaste | nano-banana -s 2K -o result
```

## Options

| Option | Default | Description |
|--------|---------|-------------|
| `-o, --output` | `nano-gen-{ts}` | 输出文件名（不含后缀） |
| `-s, --size` | `1K` | 图片尺寸：`512`, `1K`, `2K`, `4K` |
| `-a, --aspect` | model default | 宽高比：`1:1`, `16:9`, `9:16`, `4:3`, `3:4` 等 |
| `-m, --model` | `flash` | 模型：`flash`/`nb2`, `pro`/`nb-pro`, 或任意 model ID |
| `-d, --dir` | 当前目录 | 输出目录 |
| `-r, --ref` | - | 参考图（可多次指定） |
| `-t, --transparent` | - | 绿幕抠图（需要 FFmpeg + ImageMagick） |
| `--api-key` | - | Gemini API Key（优先级最高） |
| `--costs` | - | 查看成本汇总 |
| `--json` | - | JSON 输出（stdout，脚本友好） |
| `--plain` | - | 纯文本输出（stdout） |
| `--jq EXPR` | - | 过滤 JSON 输出（需配合 `--json`） |

## Models

| Alias | Model | 用途 |
|-------|-------|------|
| `flash`, `nb2` | Gemini 3.1 Flash | 默认，快速便宜 |
| `pro`, `nb-pro` | Gemini 3 Pro | 最高质量 |

## 透明模式

`-t` 会自动在 prompt 中加入绿幕指令，生成后用 ffmpeg colorkey + despill 去背景，最后 magick trim 裁边。

```bash
brew install ffmpeg imagemagick
nano-banana "robot mascot" -t -o mascot
```

## 成本追踪

每次生成记录到 `~/.nano-banana/costs.json`：

```bash
nano-banana --costs
nano-banana --costs --json
```

## 开发

```bash
make test     # 运行测试
make build    # 构建二进制
make lint     # vet + 格式检查
make ci       # 完整 CI 流程
```
