# gophp — 機能を限定したPHP CLI処理系

Goで字句解析・構文解析・評価を行うPHPインタプリタです。基準は **PHP 8.4.26 / 64ビット環境**。
同じGo実装とPHP fixtureを、8.3.35・8.4.26・8.5.11と比較します。実行時に公式PHPを呼び出しません。

**現在のfixtureは622件**（正常318・実行時230・構文74）。ユーザーがMacで実行した[成功ログ](docs/validation/named-exception-622-compare-all-mac.log)では、PHP 8.3.35・8.4.26・8.5.11の各版で622一致・0不一致・0環境エラー、`make compare-all`終了0です。605件から追加した17件の公式比較待ちは解消しました。今回は保存済みログを確認・保存したもので、比較は再実行していません。[実装範囲・過去の比較記録・残件](docs/items6-10-progress.md)を参照してください。項目6〜10全体とGo製PHPでのLaravel起動は未完了です。Laravel/API Platform調査用アプリのlock・vendorは生成済みで、公式PHPの3版でHTTP・SQLite・API作成/取得が成功しています。[証拠と残件](docs/validation/laravel-api-evidence-2026-10-13.md)を参照してください。
[今回の対応範囲](docs/phase2-compatibility.md)、[依存調査](docs/dependency-audit.md)、[Laravel/API Platformの初期目標](docs/laravel-api-target.md)、[実行時設計](docs/phase2-design.md) を参照してください。

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
| 値・演算 | int・float・string・bool・null・配列・オブジェクト、四則演算・剰余・連結・比較・論理短絡、三項・null合体・ビット演算 |
| 代入 | 式内・連鎖・複合代入、`??=`、ビット複合代入、前置／後置の増減、配列要素・プロパティの更新 |
| 制御 | if / elseif / else if / else、switch、match、while / for / do while / foreach、レベル付きbreak / continue、return |
| 関数・クラス | 通常関数、位置引数・既定値・再帰・ローカルスコープ、`global`・`static`ローカル、文字列と二要素配列の一部の変数callable。継承・interface・trait・staticメンバー・型なしクラス定数・`self`/`parent`/`static`、動的クラス名の部分集合 |
| 型宣言 | 単純な引数・戻り値型、nullable、ファイル単位の`declare(strict_types=0/1)`、基本的な型変換とTypeErrorの部分集合 |
| 引数・例外 | 名前付き引数、整数／文字列キー配列の展開、可変長引数、throw・try/catch/finally、基本的なThrowable階層の部分集合。収録fixtureの範囲で公式PHP3版と比較済み |
| 配列 | 作成・キー変換・挿入順・添字・追加・ネスト・値コピー・削除・値渡しforeach |
| 言語構造 | isset（複数対象・短絡）、empty（通常式も可）、unset |
| キャスト | `(bool)`・`(int)`・`(float)`・`(string)` |
| 標準関数 | strlen、count、gettype、is_int、is_float、is_string、is_bool、is_null、is_array、is_object、is_numeric |
| 文字列 | 二重引用符内の限定した変数展開、heredoc、nowdoc |
| 名前・定数 | 名前空間と各種use、`__NAMESPACE__`、`Foo::class`、`const`、`define` / `defined` / `constant`、COUNT_NORMAL / COUNT_RECURSIVE |
| 複数ファイル | ローカルファイルの`include` / `require` / `*_once`、`__FILE__` / `__DIR__`の部分集合 |

foreachは配列の値渡し形式が対象です。参照形式・オブジェクト反復・分解代入は明示的に拒否します。
参照全般・高度な型宣言・型付きプロパティ／定数・上表以外の標準関数と定数は対象外です。クロージャは値捕捉・arrow・static形式の部分集合です。参照捕捉、Traversableの引数展開、一般のcallable、完全な例外組み込みAPIとエラーハンドラは未対応です。クラス機能と複数ファイル実行には未対応の形式・境界条件が残ります。
キャストの別名やその他のキャスト、文字列添字の更新にも対応していません。
通常関数とメソッドを合わせた呼び出し深度は256です。
詳細な仕様・診断・版差・制約は [互換性文書](docs/compatibility.md) に記載しています。

依存調査用の[専用アプリ](research/laravel-api/README.md)には、固定したComposer依存と最小HTTP検証コードを置いています。公式PHPの3版でアプリ検証に成功していますが、既存fixtureの一致もこの公式PHPでの成功も、Go製PHPのフレームワーク起動を示すものではありません。

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

以下は追加前の第1版に関する検証記録です：

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
