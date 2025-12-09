# PAY.JP CLI

PAY.JP APIをコマンドラインから操作するためのCLIツールです。

## 機能

- 支払い (Charges) の作成・取得・更新・キャプチャ・返金
- 顧客 (Customers) の作成・取得・更新・削除
- カード (Cards) の追加・取得・更新・削除
- プラン (Plans) の作成・取得・更新・削除
- 定期課金 (Subscriptions) の作成・取得・更新・停止・再開・キャンセル
- トークン (Tokens) の取得
- 入金 (Transfers) の取得・リスト
- イベント (Events) の取得・リスト
- 取引明細 (Statements) の取得・リスト
- 集計区間 (Terms) の取得・リスト
- 残高 (Balances) の取得・リスト
- アカウント (Accounts) 情報の取得

## インストール

### Go Install

```bash
go install github.com/payjp/payjp-cli@latest
```

### ソースからビルド

```bash
git clone https://github.com/payjp/payjp-cli.git
cd payjp-cli
make build
```

### バイナリダウンロード

[Releases](https://github.com/payjp/payjp-cli/releases)ページから、お使いのプラットフォームに合ったバイナリをダウンロードしてください。

## 初期設定

### APIキーの設定

```bash
# 設定ファイルにAPIキーを保存
payjp config set api-key sk_test_xxxxxxxxxxxxx

# または環境変数で設定
export PAYJP_API_KEY=sk_test_xxxxxxxxxxxxx
```

### プロファイルの設定

複数の環境（テスト/本番）を切り替えて使用できます。

```bash
# プロファイルの作成
payjp config set-profile development --api-key sk_test_xxxxxxxxxxxxx
payjp config set-profile production --api-key sk_live_xxxxxxxxxxxxx

# プロファイルの切り替え
payjp config use-profile production

# プロファイル一覧
payjp config list-profiles
```

## 使用例

### 支払い

```bash
# 支払いの作成
payjp charges create --amount 1000 --currency jpy --card tok_xxxxx

# 支払い情報の取得
payjp charges get ch_xxxxx

# 支払いリストの取得
payjp charges list --limit 10

# 支払いの返金
payjp charges refund ch_xxxxx
```

### 顧客

```bash
# 顧客の作成
payjp customers create --email user@example.com --card tok_xxxxx

# 顧客情報の取得
payjp customers get cus_xxxxx

# 顧客リストの取得
payjp customers list --limit 10
```

### 定期課金

```bash
# プランの作成
payjp plans create --amount 1000 --currency jpy --interval month --name "Basic Plan"

# 定期課金の作成
payjp subscriptions create --customer cus_xxxxx --plan pln_xxxxx

# 定期課金の停止
payjp subscriptions pause sub_xxxxx

# 定期課金の再開
payjp subscriptions resume sub_xxxxx
```

## グローバルオプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--api-key` | `-k` | APIキー（環境変数より優先） | - |
| `--output` | `-o` | 出力形式 (json/table/yaml) | table |
| `--live` | - | 本番モード | false |
| `--verbose` | `-v` | 詳細出力 | false |
| `--quiet` | `-q` | 最小出力（IDのみ） | false |
| `--config` | `-c` | 設定ファイルパス | ~/.payjp/config.yaml |

## 出力形式

### Table形式（デフォルト）

```bash
payjp charges list -o table
```

### JSON形式

```bash
payjp charges get ch_xxxxx -o json
```

### YAML形式

```bash
payjp charges get ch_xxxxx -o yaml
```

### Quiet形式（IDのみ）

```bash
payjp charges create --amount 1000 --currency jpy --card tok_xxxxx -q
# 出力: ch_xxxxxxxxxxxxx
```

## 設定ファイル

設定ファイルは `~/.payjp/config.yaml` に保存されます。

```yaml
default_profile: development

output:
  format: table
  color: true

retry:
  max_count: 3
  initial_delay: 2
  max_delay: 32

profiles:
  development:
    api_key: sk_test_xxxxxxxxxxxxx
    mode: test
  production:
    api_key: sk_live_xxxxxxxxxxxxx
    mode: live

aliases:
  ch: charges
  cu: customers
  sub: subscriptions
```

## 環境変数

| 環境変数 | 説明 |
|----------|------|
| `PAYJP_API_KEY` | APIキー |
| `PAYJP_CONFIG` | 設定ファイルパス |
| `PAYJP_OUTPUT` | 出力形式 |
| `PAYJP_LIVE` | 本番モード (true/false) |
| `PAYJP_PROFILE` | 使用するプロファイル名 |

## 終了コード

| コード | 意味 |
|--------|------|
| 0 | 成功 |
| 1 | 一般的なエラー |
| 2 | コマンドライン引数エラー |
| 3 | 設定エラー |
| 4 | 認証エラー (401) |
| 5 | リクエストエラー (400) |
| 6 | 支払いエラー (402) |
| 7 | リソース未発見 (404) |
| 8 | レートリミット (429) |
| 9 | サーバーエラー (500) |

## コマンド一覧

```
payjp [command]

Available Commands:
  accounts      Manage account
  balances      Manage balances
  cards         Manage customer cards
  charges       Manage charges
  config        Manage CLI configuration
  customers     Manage customers
  events        Manage events
  help          Help about any command
  plans         Manage subscription plans
  statements    Manage statements
  subscriptions Manage subscriptions
  terms         Manage terms
  tokens        Manage tokens
  transfers     Manage transfers
```

詳細なヘルプは `payjp [command] --help` で確認できます。

## 開発

### 依存関係のインストール

```bash
make deps
```

### ビルド

```bash
make build
```

### テスト

```bash
make test
```

### 全プラットフォーム向けビルド

```bash
make build-all
```

## ライセンス

MIT License

## 関連リンク

- [PAY.JP](https://pay.jp/)
- [PAY.JP API ドキュメント](https://pay.jp/docs/api/)
- [payjp-go SDK](https://github.com/payjp/payjp-go)
