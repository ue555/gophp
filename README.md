# gophp — 機能を限定したPHP CLI処理系

Goで字句解析・構文解析・評価を行うPHPインタプリタです。基準は **PHP 8.4.26 / 64ビット環境**。
同じGo実装とPHP fixtureを、8.3.35・8.4.26・8.5.11と比較します。実行時に公式PHPを呼び出しません。

**現在の状態：限定した対応範囲を持つCLI処理系の第1版が完成。**
ユーザー環境での公式PHP実比較ログにより、8.3.35・8.4.26・8.5.11の各466件、
合計1,398件の一致（差分0・検証エラー0）を確認しました。
実装・fixture・CLI・文書と検証を含む項目9〜12は完了です。

## 実行

Go 1.22以上、Python 3を使用します。開発時の確認環境はGo 1.27.0・Python 3.9.6・darwin/arm64です。

```sh
cd /Users/kouji/dev/php
make build
./bin/gophp examples/report.php
./bin/gophp --php-version=8.5 examples/report.php
./bin/gophp --help
```

[集計スクリプト](examples/report.php)は、配列を関数に渡してforeachで選別し、結果の配列を返します。
isset・empty・unset・キャスト・標準関数を使い、元配列の保持とStatsオブジェクトの共有も確認します。
出力は以下です。

```text
a:tea:12:3
c:水:3:3
2:2:1:0:integer
```

CLIは `gophp [--php-version=8.3|8.4|8.5] <file.php>`。既定は8.4です。
標準出力はプログラムの出力、標準エラーは診断とヘルプです。
終了コードは成功0、ファイル読み込み等の失敗1、引数の誤り2、PHP構文・致命エラー255です。
ファイル名の前に `--` を指定できます。PHPの `-r`・`-l`、stdin実行、スクリプト引数は未対応です。

## Go製PHPをコンテナで実行

```sh
make container-build
make container-check
make container-run
make container-run GOPHP_TARGET=8.5 GOPHP_FILE=tests/fixtures/success/hello.php
```

PodmanでGo製インタプリタ自体をビルドし、PHPファイルをマウントして実行します。
podman-composeでの実行と全fixtureのチェックは [gophpコンテナの案内](containers/gophp/README.md) にあります。
このCodexセッションではPodman接続が拒否されたため、コンテナの実ビルド・実行は未確認です。

## 対応範囲

| 項目 | 内容 |
| --- | --- |
| 値・演算 | int・float・string・bool・null・配列・オブジェクト、四則演算・剰余・連結・比較・論理短絡 |
| 代入 | 式内・連鎖・複合代入、前置／後置の増減、配列要素・プロパティの更新 |
| 制御 | if / elseif / else if / else、while / for / do while / foreach、レベル付きbreak / continue、return |
| 関数・クラス | 通常関数、位置引数・既定値・再帰・ローカルスコープ、継承なしのクラスとメソッド・アクセス制御 |
| 配列 | 作成・キー変換・挿入順・添字・追加・ネスト・値コピー・削除・値渡しforeach |
| 言語構造 | isset（複数対象・短絡）、empty（通常式も可）、unset |
| キャスト | `(bool)`・`(int)`・`(float)`・`(string)` |
| 標準関数 | strlen、count、gettype、is_int、is_float、is_string、is_bool、is_null、is_array、is_object、is_numeric |
| 定数 | COUNT_NORMAL、COUNT_RECURSIVE |

foreachは配列の値渡し形式が対象です。参照形式・オブジェクト反復・分解代入は明示的に拒否します。
参照・可変長／名前付き引数・型宣言・クロージャ・継承・例外処理・上表以外の標準関数と定数は対象外です。
キャストの別名やその他のキャスト、文字列添字の更新にも対応していません。
通常関数とメソッドを合わせた呼び出し深度は256です。
詳細な仕様・診断・版差・制約は [互換性文書](docs/compatibility.md) に記載しています。

## 検証

```sh
make test vet build
make compare       # 主基準 PHP 8.4.26、全fixture
make compare-all   # 8.3.35 / 8.4.26 / 8.5.11、全fixture
make compare-parser
```

比較環境の起動方法は [コンテナの案内](containers/php-version/README.md) を参照してください。
[php-go-interpreterスキル](skills/php-go-interpreter/SKILL.md) の比較スクリプトを使用し、
stdout・stderrをバイト単位、終了コードを厳密に比較します。trim・警告の削除・fixtureの除外は行いません。
正常系がどちらか一方でも非ゼロ終了なら検証エラーとします。

2026-10-12の今回の結果：

- Goテスト、Pythonテスト11件、go vet、ビルドは成功。Pythonにはコンテナ停止の検出と誤検出防止の回帰テスト2件を追加しています。
- 既存337件のPHPソースと保存済み期待値を維持し、129件追加。全466件（正常208・実行時198・構文／コンパイル60）。
- ビルド済みGo CLIと保存済み期待値の照合は3設定とも466一致・差分0・検証エラー0。正常系208件は全設定で終了コード0。
- ユーザー環境での公式PHP実比較は8.3.35・8.4.26・8.5.11の各466一致・差分0・検証エラー0。正常系208件は各版とも両処理系で終了コード0。
- fixtureの補助関数名の修正後、停止していた3コンテナを起動し、`check.sh` と `make compare-all` が成功しました。

[修正後のローカル検証ログ](docs/validation/v1-followup-local-validation.txt)、
[ユーザー提供の全件実比較と完成確認](docs/validation/v1-user-comparison-complete.md) を区別して記録しています。
新規129件の期待値は公式資料・PHPソースを根拠に作成し、全件を3版の公式PHP実比較で確認しました。
比較スクリプトはPodman接続障害・コンテナ停止を検証エラーとして扱い、その版の残りを未実行と明示します。
[全件成功前の修正記録](docs/validation/v1-user-comparison-followup.md) と
[エージェント環境で接続を拒否された試行ログ](docs/validation/v1-followup-php-comparison.txt) は履歴として残しています。

追加前の337件は、ユーザー提供の公式比較ログで各版337一致・差分0・検証エラー0を確認済みです。
これは今回の変更前の履歴です。[以前の仕様・検証記録](docs/validation/pre-v1-compatibility.md) を参照してください。
PHP全機能との互換性を示すものではありません。
