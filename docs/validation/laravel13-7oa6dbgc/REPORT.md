# 固定Laravel 13のGo PHP 8.5起動検証（2026-10-22）

対象は `/Users/kouji/dev/gophp-app`。Laravel 13の今回指定された11条件は修正版で成功した。開発依存を維持し、標準アプリ・vendor・lockを書き換えず、公式PHPへのフォールバックも使っていない。全PHP機能やLaravel全機能への対応を意味しない。

## 基準・識別情報

- 作業開始HEAD: `e1e1ba0c76401ae84908e3f852499c3d1d094d54`。既存の未追跡 `.gophp-version` を保持。今回commit・push・タグ作成なし。
- 開始時fixture/期待値は実ファイルで1193件。現在1201件（正常771・実行時296・構文134）、追加8件。以前の1193件の期待値は一切変更していない。
- スケルトンv13.11.0 / commit `06d016a364a37430eec9cbc52209adce3d7667ce` はアプリの既存provenanceに基づく。今回上流との全ファイル照合を新たに行ったという意味ではない。framework **v13.38.0**、110パッケージ（開発用33）は実lockとinstalled.jsonから確認した。
- composer.lock SHA-256: `802c12e5dd5733a25123300cf461bd279b831f398f9a9b211b81fa20bedf4b25`
- package-lock.json SHA-256: `33ffe5ba021c97fe011c1196f8979bbdb7e4af21a88de245e87872017e4c02e3`
- 修正版: `/Users/kouji/dev/php/.cache/laravel13/gophp`
- バイナリSHA-256: `2fdf4f3f01f59a3df1fb5e0607014bb664724aa9b021762c257be59b45c36cf2`
- 検証ソースmanifest SHA-256: `2bbb66d86cc1e6d15c83a3b62799036240fe292324ffe6e56aec9d0f557bed3a`。HEADのみで今回の検証対象を再現できない。未コミット差分、新規ファイル各ハッシュ、アプリの検証スクリプトのハッシュは [tested-source.json](tested-source.json) に記録した。文書・ログは実行ソースのハッシュから除外する。
- ビルド: Go1.27.0 / darwin arm64 / CGO有効 / Apple clang21.0.0。標準のCGO flags。Goのcacheはこのcheckout内を使用。[初期識別](initial.json)。

## 到達した障害と修正

1. 元の `docs/repro/getmypid.php` は旧ローカルバイナリで終了255、公式PHP8.5.11で `integer|positive`、終了0。名前空間中の既存 `strlen()` は正常、明示したグローバル `\getmypid()` も未定義だったため、名前空間フォールバックではなく組み込み本体の欠落と切り分けた。`os.Getpid()`の実PIDを返し、名前付き引数・個数・副作用・診断スタックを既存経路に接続した。現在の元再現も終了0。
2. 次はSymfony Requestの `public ... $property { set { ... } }` で構文停止。PHP8.4/8.5のbacked block setフックを実行し、自分のbacking propertyへの書き込みは再帰させない。ReflectionProperty `setRawValue()` / `getRawValue()`を追加し、raw書き込みでも型検査を維持する。継承・フック副作用・型エラー後の状態・間接変更の拒否を公式観測から検証した。
3. Consoleの `getcwd()` を実際の作業ディレクトリ取得に接続した。
4. Viteの `array_unique(..., SORT_REGULAR)` が配列を文字列化していたため、通常比較による重複判定へ修正した。標準welcomeのpreloadを処理できる。
5. 404のBlade component処理に必要な `get_defined_vars()` を現在のスコープ・変数・参照セルに接続した。
6. `--document-root`によるpublicアセット配信を明示指定で追加。Laravelのfront controllerと併用する。PHPソース、隠しパス、トラバーサル、外向きsymlinkや拡張子を偽ったsymlinkを拒否する。新しい実TCP GoテストでCSS/JS/font、HEAD、拒否経路を確認した。

PHP側の根拠は [getmypid](https://www.php.net/manual/en/function.getmypid.php)、[property hooks](https://www.php.net/manual/en/language.oop5.property-hooks.php)、[setRawValue](https://www.php.net/manual/en/reflectionproperty.setrawvalue.php) と、このcheckoutの公式PHP3版の実観測。新規期待値は `research/fresh-laravel/observe.py` により版別stdout/stderr/終了コードから採取した。`fresh-laravel-official-observations-*` の今回作成した6ディレクトリを保持している。今回の新規fixtureは `protection-and-fixture-audit.json` に列挙する。

## Laravel 13の実CLI・TCP・DB

各ランタイムで空の独立SQLiteと、services/packages/config/routes/events/viewsの隔離キャッシュを用意。同じ要求順で実行し、自分のサーバーと一時データだけを最後に片付けた。

- Go: [47検査の集計](/Users/kouji/dev/gophp-app/docs/validation/20261022-122544-776c16d7-gophp/results.json)
- 公式PHP8.5.11: [同じ47検査](/Users/kouji/dev/gophp-app/docs/validation/20261022-121925-5367ae79-official/results.json)
- 両方終了0、各46 HTTP要求。autoload、package discovery、list/about/route:listの内容とstderr、空DBの標準3 migration・10テーブル、再migrationのdump不変、`/up`と標準welcomeのBlade冷/温キャッシュ、9アセットのバイト列とMIME、404、9種類のprivate/traversal経路、20連続・4並列要求を確認した。
- 同じCookieを送る要求はセッションを維持。独立4並列と合わせて5セッション、異なるID/CSRFトークン、JSON payload、直近のactivityとアプリテーブルの保持を検査。Laravel13の固定configはJSON session serializationなのでPHP形式を前提にしない。
- Goの全CLI stderrとサーバーログは空。公式サーバーログは通常の接続情報のみで、PHP警告・例外等はない。生成済みviewキャッシュのハッシュが温要求で変わらないことも確認。
- 日常の `scripts/artisan` / `scripts/serve` も別途local環境・debug有効で検証し、about/migrationと4 HTTP要求が成功。手動DBは使っていない。

### 比較の規則と実際の差

`scripts/compare-validation.py` は両方の検証成功、lockとverifierの同一ハッシュ、46要求のpath/status/body hashを確認する。アプリHTML・health・未登録ルート・アセットの計37本文は、実行URLと `/up` の `Response rendered in Nms.` だけを正規化して一致した。migrationのDB dumpはそのまま一致、session payloadは形式検査済みCSRF tokenと実行URLだけを正規化して一致。警告や空白は削除しない。

9種類の非公開パスは漏えいせず拒否されることを比較する。そのうち8応答はGoの配信入口でplain text 404、公式はLaravelのHTML 404という差がある。JSのMIMEは公式 `application/javascript`、Go `text/javascript; charset=utf-8`。これらは隠さず比較結果に列挙する。CSSのcharset表記の大小文字のみをMIME比較で同一視する。Date、暗号化Cookie、Server、Last-Modified等のヘッダー全体のバイト一致は主張しない。Cookieの動作は実再送とDBの状態で検証した。

## 共通実装の回帰結果

版別の最終結果と生ログへの索引は [verification-index.json](verification-index.json) に記録する。過去のIssue #1の1193件成功を今回の結果として流用していない。

| 検証 | 実測結果 | 生ログ |
|---|---|---|
| full-fixtures-8.3.35 | 1201一致・0差異・0エラー、終了0 | [stdout](20261022-122544-full-fixtures-8.3.35-4fd96ee7.stdout)・[stderr](20261022-122544-full-fixtures-8.3.35-4fd96ee7.stderr)・[コマンド/終了コード](20261022-122544-full-fixtures-8.3.35-4fd96ee7.json) |
| full-fixtures-8.4.26 | 1201一致・0差異・0エラー、終了0 | [stdout](20261022-122831-full-fixtures-8.4.26-f5a3630f.stdout)・[stderr](20261022-122831-full-fixtures-8.4.26-f5a3630f.stderr)・[コマンド/終了コード](20261022-122831-full-fixtures-8.4.26-f5a3630f.json) |
| full-fixtures-8.5.11 | 1201一致・0差異・0エラー、終了0 | [stdout](20261022-123101-full-fixtures-8.5.11-a7c9324a.stdout)・[stderr](20261022-123101-full-fixtures-8.5.11-a7c9324a.stderr)・[コマンド/終了コード](20261022-123101-full-fixtures-8.5.11-a7c9324a.json) |
| final-make-test | 終了0 | [stdout](20261022-122544-final-make-test-ca85defe.stdout)・[stderr](20261022-122544-final-make-test-ca85defe.stderr)・[コマンド/終了コード](20261022-122544-final-make-test-ca85defe.json) |
| final-make-vet | 終了0 | [stdout](20261022-122728-final-make-vet-af9c5d7c.stdout)・[stderr](20261022-122728-final-make-vet-af9c5d7c.stderr)・[コマンド/終了コード](20261022-122728-final-make-vet-af9c5d7c.json) |
| final-make-build | 終了0 | [stdout](20261022-122729-final-make-build-40621b68.stdout)・[stderr](20261022-122729-final-make-build-40621b68.stderr)・[コマンド/終了コード](20261022-122729-final-make-build-40621b68.json) |
| final-diff-check | 終了0 | [stdout](20261022-122730-final-diff-check-b9e75e12.stdout)・[stderr](20261022-122730-final-diff-check-b9e75e12.stderr)・[コマンド/終了コード](20261022-122730-final-diff-check-b9e75e12.json) |
| fixed-laravel-api-platform-three | 終了0 | [stdout](20261022-122544-fixed-laravel-api-platform-three-116525d7.stdout)・[stderr](20261022-122544-fixed-laravel-api-platform-three-116525d7.stderr)・[コマンド/終了コード](20261022-122544-fixed-laravel-api-platform-three-116525d7.json) |
| task-api-three | 終了0 | [stdout](20261022-122748-task-api-three-a99b2ccc.stdout)・[stderr](20261022-122748-task-api-three-a99b2ccc.stderr)・[コマンド/終了コード](20261022-122748-task-api-three-a99b2ccc.json) |
| fresh-laravel12-85 | 終了0 | [stdout](20261022-123249-fresh-laravel12-85-9e45306f.stdout)・[stderr](20261022-123249-fresh-laravel12-85-9e45306f.stderr)・[コマンド/終了コード](20261022-123249-fresh-laravel12-85-9e45306f.json) |
| go-final-http | 終了0 | [stdout](20261022-122544-go-final-http-73a07aad.stdout)・[stderr](20261022-122544-go-final-http-73a07aad.stderr)・[コマンド/終了コード](20261022-122544-go-final-http-73a07aad.json) |
| official-final-http | 終了0 | [stdout](20261022-121924-official-final-http-007eacd9.stdout)・[stderr](20261022-121924-official-final-http-007eacd9.stderr)・[コマンド/終了コード](20261022-121924-official-final-http-007eacd9.json) |
| paired-http-and-db | 終了0 | [stdout](20261022-122757-paired-http-and-db-2823db79.stdout)・[stderr](20261022-122757-paired-http-and-db-2823db79.stderr)・[コマンド/終了コード](20261022-122757-paired-http-and-db-2823db79.json) |
| daily-launchers | 終了0 | [stdout](20261022-122946-daily-launchers-4bdeca84.stdout)・[stderr](20261022-122946-daily-launchers-4bdeca84.stderr)・[コマンド/終了コード](20261022-122946-daily-launchers-4bdeca84.json) |
| original-repro-fixed | 終了0 | [stdout](20261022-123248-original-repro-fixed-f176ff3f.stdout)・[stderr](20261022-123248-original-repro-fixed-f176ff3f.stderr)・[コマンド/終了コード](20261022-123248-original-repro-fixed-f176ff3f.json) |

追加8 fixtureの対象比較も3版すべて一致（計24比較、全て終了0）。既存Laravel/API PlatformはGo3設定で96連続/8並列、タスクAPIは公式/Go3組で各191要求、Laravel12は8.5の公式/Go各111要求を確認した。

## 対応の境界

- property hooksはbacked block `set { ... }`の部分対応。getter、virtual、arrow、明示set引数等の未実装構文は拒否する。trait・unset・ArrayAccess等すべての組み合わせまで網羅したものではない。raw static accessも未対応。
- `get_defined_vars()`は値・参照・スコープを扱うが、変数名順に返す。公式PHPの定義順は再現しない（fixtureはキー集合/値/参照に焦点を当ててksortしている）。
- `array_unique`の対応フラグはSORT_STRINGとSORT_REGULAR。異種値の非推移的比較、循環配列、他のフラグや参照を含む全境界は未網羅。
- 静的配信は既知アセット拡張子に限定した開発用途。公開ディレクトリを非信頼主体が同時書換する競合まで防御する本番サーバーではない。public外symlinkを含む一般的なLaravel storage公開は今回の対象外。
- Laravel13アプリ自体は依存要件に適合する8.5で検証。8.3/8.4は共有処理系のfixtureと既存の対応アプリで回帰検証する。
- npm audit未完了、fontaine任意依存のビルド警告は初期構築の残件としてアプリREADMEに維持する。今回依存導入・npm再ビルドはしていない。ブラウザー描画や外部CDNの可用性の確認は含まない。

## 失敗ログ・保護

このディレクトリのtimestamp付き `.json` に実コマンドと実際の終了コード、同名 `.stdout` / `.stderr` に生出力を保存した。終了255はPHPの期待診断であるケースもあり、全てを検証失敗と数えていない。最初のLaravel停止、回帰fixtureで発見した内部スタック欠落/フック迂回、セッション形式を誤認した検証失敗も保持した。テスト名の選択誤りで0件だった初回実行を成功根拠にせず、対象選択を訂正して全件も実行した。

初回のGo buildは既定のmodule cacheに権限警告を出したが終了0。以降はcheckout内のcacheを明示して警告なしでビルドした。`ps`はsandboxに拒否され終了1（`process-list-environment-error.json`）、同じ試行を繰り返していない。PodmanとTCPは利用でき、必要な比較を環境エラーで省略していない。システム全体のプロセス列挙はできないが、自分が起動したサーバーは所有するprocess handle/container名で停止し、実際の終了コードを保存した。

開始時の8,999ファイル（vendor・標準アプリコード・lock・.env・手動DB・既存モード設定）のハッシュ一致は [保護監査](protection-and-fixture-audit.json) に記録する。アプリ変更はscripts/serve、scripts/verify、新規比較スクリプト、READMEとsetup-status、および新規検証ログ。標準アプリの変更なし。インストール済みgophpやgophpenv設定・タグは変更していない。

## 起動と再検証

```sh
cd /Users/kouji/dev/php
GOCACHE="$PWD/.cache/go-build" GOMODCACHE="$PWD/.cache/go-mod" \
  go build -o .cache/laravel13/gophp ./cmd/gophp

cd /Users/kouji/dev/gophp-app
export GOPHP_BIN=/Users/kouji/dev/php/.cache/laravel13/gophp
scripts/artisan about
scripts/serve
# http://127.0.0.1:8092（Ctrl-Cで自分のサーバーを停止）
```

検証はサーバーを別の専用ポートへ自動起動する:

```sh
cd /Users/kouji/dev/gophp-app
GOPHP_BIN=/Users/kouji/dev/php/.cache/laravel13/gophp scripts/verify
scripts/verify --official
python3 scripts/compare-validation.py \
  docs/validation/20261022-121925-5367ae79-official \
  docs/validation/20261022-122544-776c16d7-gophp
```

比較引数は再実行時に表示された新しい成功ディレクトリへ置き換える。通常ランチャーの既定は引き続き既存gophpenvなので、修正版の `GOPHP_BIN` 指定が必要。

本体の版別比較（作業ディレクトリ `/Users/kouji/dev/php`）:

```sh
python3 skills/php-go-interpreter/scripts/compare_php.py --php containers/php-version/bin/php83 --php-version 8.3.35 --fixtures tests/fixtures --timeout 30 --interpreter .cache/laravel13/gophp --php-version=8.3 '{file}'
python3 skills/php-go-interpreter/scripts/compare_php.py --php containers/php-version/bin/php84 --php-version 8.4.26 --fixtures tests/fixtures --timeout 30 --interpreter .cache/laravel13/gophp --php-version=8.4 '{file}'
python3 skills/php-go-interpreter/scripts/compare_php.py --php containers/php-version/bin/php85 --php-version 8.5.11 --fixtures tests/fixtures --timeout 30 --interpreter .cache/laravel13/gophp --php-version=8.5 '{file}'
python3 research/laravel-api/verify-gophp.py --interpreter .cache/laravel13/gophp
python3 research/task-api/verify.py --interpreter .cache/laravel13/gophp
python3 research/fresh-laravel/verify.py --binary .cache/laravel13/gophp --version 8.5
GOCACHE="$PWD/.cache/go-build" GOMODCACHE="$PWD/.cache/go-mod" make test vet build
git diff --check
```
