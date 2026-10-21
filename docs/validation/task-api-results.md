# タスクAPI検証記録（2026-10-21）

**今回の完成条件を達成。未実施の必須検証なし。** 新規アプリは[research/task-api](../../research/task-api/README.md)。Laravel 12.69.3のルート・標準バリデーション・Eloquent、SQLiteで、指定された認証なしのTask CRUDを実装した。認証・関連モデル・ページング・API Platformの同等APIは追加していない。実行と比較の再現手順、入力・日時・空PATCHの契約はアプリREADMEを参照。

## HTTP契約と証拠

最終比較ログは[run.log](task-api-20261021-124313-i5cu398g/run.log)。公式PHP 8.3.35・8.4.26・8.5.11とGoの対応する8.3・8.4・8.5設定を同じ要求順で実行する。各runtime・版は独立した空DB・storage・bootstrapキャッシュを使用し、公式はPodman経由、Goはホストの実TCPサーバーへ接続する。in-processテストを実HTTPの証拠としていない。

各実行191要求。201の作成、200の昇順一覧・取得・部分更新、204の削除（本文・Content-Typeなし）、404、422を検証する。日本語、120/121文字、titleの不正型、completedの省略・6種類の許容値・不正値、空／本文なしPATCH、未指定項目の維持、422でのDB不変、リクエスト状態分離、連続100 GET・8並列GET、全削除後の空配列を含む。DBはPythonのSQLite接続で別途読む。

日時はUTC、JSONは6桁小数、DBは秒精度。作成時の一致、更新での作成日時維持・更新日時前進、空PATCHの日時維持、DBとの対応を検査してから比較時だけ応答の2日時項目をマーカーにする。HTTP生本文とサーバー出力を保存する。警告削除やtrimで差異を隠さない。

公式PHP／Goの3組すべて191要求が成功し、HTTP比較は差異なし、検証スクリプト終了0。合計1,146件の生応答を[型を厳密に区別する追加監査](task-api-http-artifact-audit-3versions-20261021-124944-0ltedn8g.log)でも確認した。Pythonの等値比較による`true == 1`の見落としを避け、保存したJSON値の型・204のヘッダー／本文・正規化成果物との一致を確認して終了0。今後のverify.pyにもこの型確認を組み込んだ。

## Go実装の変更

1. 空のpackage discoveryキャッシュから起動するため、`is_writable`／`is_writeable`、スカラー・配列の`var_export`、`sys_get_temp_dir`・`tempnam`・`chmod`・引数なし`umask`・`rename`・`unlink`のローカルファイル操作を追加した。Laravel自身が生成したmanifestをPHPとして再読込する。
2. EloquentのCarbon日時JSON化に必要なDateTimeの書式トークンを接続した。`Y`などの個別トークン、マイクロ秒、エスケープ、タイムゾーンを扱い、`P`のUTC表記も`+00:00`に修正した。
3. 新規アプリの`artisan`は`handleCommand`の戻り値を`exit`へ渡す。CLI処理に失敗した場合の非ゼロ終了を維持する。既存アプリのartisanは変更していない。

最小fixtureは[is_writable_basic.php](../../tests/fixtures/success/is_writable_basic.php)、[var_export_manifest.php](../../tests/fixtures/success/var_export_manifest.php)、[manifest_atomic_replace.php](../../tests/fixtures/success/manifest_atomic_replace.php)、[datetime_json_format.php](../../tests/fixtures/success/datetime_json_format.php)。期待値は公式PHP 8.4.26の実出力で確認した（[filesystem/export初期参照](task-api-new-fixtures-reference-20261021-124002.log)、[export追加境界](task-api-export-reference-20261021-124259.log)、[日時参照](task-api-datetime-reference-20261021-124126.log)）。

これは必要経路の対応でありPHP全体の実装完了ではない。`var_export`の通常オブジェクトは明示的に未対応。`umask(mask)`によるプロセス全体の変更、filesystem stream context、一般のstream wrapperやfilesystem失敗診断の完全互換は対象外。DateTimeの全書式・全日時パーサー・拡張年も網羅しない。短時間のGET検証は長時間負荷・並列書込み性能の保証ではない。

## 回帰検証

- [make test](task-api-test-20261021-124259-9tte8m5e.log)：終了0。Go全パッケージとPython 12テスト。
- [make vet](task-api-vet-20261021-124313-rhuxh2go.log)、[make build](task-api-build-20261021-124247-rr5amplg.log)：終了0。
- [全fixtureの公式3版比較](task-api-fixtures-all-3versions-20261021-124313-mipuzcdk.log)：3版それぞれ**1,173一致・0差異・0エラー、終了0**。4件追加後の内訳は正常745・実行時294・構文134。HTTP検証とは別のstdout・stderr・終了コードの厳密比較。
- [既存固定アプリのGo実HTTP回帰](task-api-legacy-go-http-3versions-20261021-124343-o0z27mz1.log)：8.3・8.4・8.5すべて終了0。BookのLaravel／API Platform CRUD、OpenAPI、連続100 GET・8並列GET。既存検証スクリプトを使い、bootstrapキャッシュだけ追加で隔離した。

## 依存・準備・保護

[環境確認](task-api-environment-20261021-123327.log)ではPodman接続・TCP bindとも成功。VM再起動・既存サービス停止を行っていない。

[依存取得](task-api-dependencies-20261021-123407.log)には公式PHP 8.3.35のComposerを使用した。Goアプリ実行とは区別する。新規アプリ内に独立したvendorを作成し、root名の変更に必要な[lockメタデータ更新](task-api-lock-refresh-20261021-123741.log)後、94パッケージのレコードが既存lockと同一であることを確認した。[Composer validate](task-api-composer-lock-8.3.35-20261021-124431-yfzfqn55.log)は終了0。license未指定と意図的な完全固定バージョンへの一般警告はログに残す。各版のplatform requirements検査はHTTP検証run.logにも含めた。

本アプリはAPI Platformリソース探索を空配列、metadata cacheをarrayにしており、metadata dump／API Platform installは不要。全HTTP実行で空bootstrapキャッシュからGoまたは公式PHP自身の`package:discover`とmigrationを実行した。migration後には検証用storageのファイルキャッシュも空にした。既存Bookアプリでは従来のmigration後キャッシュ削除を維持している。

[手動DBの準備](task-api-prepare-manual-go-8.4-20261021-124600-2xgwitm6.log)はGo 8.4で終了0。`research/task-api/database/database.sqlite`は手動利用用として保持する。設定に絶対パスは固定していない。検証用一時DB・storageは終了時に削除し、自分が起動したサーバーだけを停止した。Goの`server_exit_code: -15`は停止操作の結果で、検証スクリプトは別の終了コードを記録する。

[保護監査](task-api-protected-files-and-dependencies-20261021-124944-ojqe_br3.log)で、既存composer.json・lock・vendorの6,675ファイルのSHA-256が作業前と同一、新規vendorの94版がlockと一致、既存アプリに追跡ファイル変更なし、検証一時ディレクトリの削除、手動DBの保持を確認した。既存`research/laravel-api/`、v0.2.0・既存履歴を維持した。commit・push・タグ付け・公開は行っていない。

## 失敗の履歴

- [最初の8.4比較](task-api-20261021-123616-36fba8gn/run.log)：公式190要求は成功。Goは`is_writable`未実装でmigrationできなかった。アプリCLIが戻り値を無視して終了0になる問題もこの時点で判明し、新規artisanを修正した。検証自体はDB表不在を検出して終了1。
- [次のGo試行](task-api-20261021-123742-v_qy66_t/run.log)：`var_export`未実装を検出、終了1。
- [manifest補完後](task-api-20261021-124015-ppwkg_sb/run.log)：日時書式未対応でPOSTが500、終了1。
- [日時修正後のGo 8.4](task-api-20261021-124135-qu8_3ub9/run.log)：当時の190要求が成功、終了0。その後本文なしPATCHを1件追加し、最終比較は191要求になった。
- [Composer PHAR直接読込の試行](task-api-lock-metadata-20261021-123633.log)：内部パスを読めずPHP終了255。新規lockのみをComposerの`update --lock --no-install`で更新し、パッケージ不変を確認した。既存lockは操作していない。

これらの失敗ログは成功ログで上書きしていない。

最終の[証拠件数・成果物・リンク点検](task-api-final-evidence-review-20261021-125142-vax5mxwm.log)と[差分検査](task-api-final-diff-check-20261021-125142-36iwsge7.log)も終了0。必須検証の未実施や今回のHTTP契約に残る不具合はない。
