# Laravel / API Platform 初期目標と判定条件

[現行の完成判定ゲート](full-stack-development.md)で、固定アプリに定義した全6段階の条件を達成した。GoのCLI・プロセス内HTTP/DB/APIに加え、公式PHP3版との各1169件一致と、Go製PHPの[版別実TCP検証](full-stack-development.md)で各2000回連続・8並列要求を含むHTTP・SQLite・API Platform契約の終了0を確認した。以前のTCPポート拒否は履歴として維持する。

現在のGo側ローカル基準はfixture・期待値各1169件。版別の[CLI検証](validation/full-stack-current-cli-matrix-local.log)では3設定ともmigration・独自コマンド・OpenAPI・SQLiteテーブル確認がPASSした。現行全件は公式PHP3版との[版別比較](full-stack-development.md)で各1169件一致した。実サーバー経由の証拠は別に[8.3](validation/full-stack-real-http-83-mac.log)・[8.4](validation/full-stack-real-http-84-mac.log)・[8.5](validation/full-stack-real-http-85-mac.log)の版別ログへ記録する。

過去のローカル確認（2026-10-21）：固定Laravel Translatorのネストした配列文言をGoの3設定で置換できた。[当時の配列走査の回帰・比較待ち](full-stack-development.md)。当時は1168件の公式PHP比較とアプリ完成判定が未達だった。現在の結果は冒頭の1169件比較と実TCP検証を参照。

DB状態の追加確認（2026-10-21）：PDO接続とStatementのエラー状態を別に保持し、固定vendorが使う`ATTR_ERRMODE`の取得・設定、silent／warning／exceptionの基本分岐をGoでローカル実行した。[対象fixtureと未比較範囲](full-stack-development.md)。固定アプリの完成条件と公式PHP3版の当時1167件比較は引き続き未達。

DB失敗経路の追加確認（2026-10-21）：SQLiteの一意制約違反が固定Laravelの`UniqueConstraintViolationException`になり、前の`PDOException`とSQLSTATEを保持する基本形をGoの3設定で実行した。[fixture・ログ・比較待ち](full-stack-development.md)。固定アプリの完成判定は引き続き未達。

追加のローカル確認（2026-10-21）：固定vendorのEloquentカーソルによるSQLite行走査と、SQL失敗時のLaravel `QueryException`連鎖をGoの3設定で実行した。[fixtureと検証状況](full-stack-development.md)。CLI・HTTP・DB・APIの完成判定と、公式PHP3版の当時1167件比較は引き続き未完了。

主基準は PHP 8.4.26、追加対象は 8.3.35 と 8.5.11。Laravel統合方式を引き継ぎ、`composer.json` の直接指定は `laravel/framework` 12.69.3、`api-platform/laravel` 5.0.1、Symfony `type-info`・`web-link` の `^7.4` である。`composer.lock` は通常パッケージ94件・開発用0件を固定し、`vendor/composer/installed.json` と全件一致した。同じlockとvendorで公式PHPの3版の検証アプリが成功したため、この**調査アプリの共通構成**は成立した。`resolve.sh` 全体の終了コードは未確認である。[Laravel 12 の公式PHP対応表](https://laravel.com/framework/docs/12.x/releases)、[API PlatformのLaravel手順](https://api-platform.com/docs/laravel/)、[依存監査](dependency-audit.md)を参照。Symfony 8系はPHP 8.3と共通にできないため7.4系を選んだ。

検証アプリは [research/laravel-api](../research/laravel-api/README.md) に置く。公式PHPではCLI版`php -S`の組み込みHTTPサーバーと版別の一時SQLiteを使用して初期契約を検証した。Go製PHPでは[HTTPブリッジの統合テスト](../cmd/gophp/main_test.go)が一時SQLiteと新しい一時storageを使い、クエリ付きGETの12回連続実行、下表の基本契約、API Platformの一覧・PATCH・DELETE・削除後404・OpenAPI JSON、追加8回のPOST/PATCH/DELETE、既定96回、[追加512回](validation/full-stack-http-512-after-callback.log)・[1024回のGET](validation/full-stack-http-1024-caching.log)をローカル実行した。[ページ分割の統合テスト](validation/full-stack-pagination-jsonld.log)では一時SQLiteの31行が30・1・0件の3ページへ分かれ、JSON-LDのCollection形式とtotalItemsを確認した。Goの3設定で固定アプリのprobe起動を確認し、8.3・8.5設定でも独立した一時DBでAPI Platformの作成・取得・更新・削除・削除後404を[版別統合テスト](validation/full-stack-api-version-matrix-local.log)で確認した。8.4設定の同経路は前記HTTPブリッジ統合テストで確認済み。リクエスト間のPDOカーソル寿命とキャッシュに保存するenum caseの同一性を修正し、現行1169件の公式PHP3版比較と実TCPでの2000回連続・8並列要求は冒頭のログで確認済み。定義済み契約を超える運用条件での長時間負荷試験は未実施である。

| 段階 | 入力・URL | 期待HTTP結果 | 期待DB状態 |
| --- | --- | --- | --- |
| autoload・起動 | `vendor/autoload.php` → `bootstrap/app.php`、`GET /api/probe` | 起動例外なし。200、`application/json`、`{"ok":true}` | `books` は空 |
| Laravel JSON | `GET /api/probe` | 200、`{"ok":true}` | 変化なし |
| SQLite保存 | `POST /api/books-manual`、`{"title":"Alpha"}` | 201、`{"id":N,"title":"Alpha"}`、Nは正の整数 | `books` に `(N, Alpha)` 1行 |
| SQLite取得 | `GET /api/books-manual/N` | 200、同じID・titleのJSON | 1行のまま |
| 入力検証 | `POST /api/books-manual`、`{}` | 422、JSONの `errors.title` に1件以上 | 1行のまま |
| API Platform作成 | `POST /api/books`、`Accept` / `Content-Type: application/ld+json`、`{"title":"Beta"}` | 201、JSON-LDの `title=Beta` と `@id` | `Beta` が追加され計2行 |
| API Platform取得 | 直前の `@id` へGET、`Accept: application/ld+json` | 200、JSON-LDの `title=Beta` | 計2行のまま |

上表の初期契約は公式PHP 8.3.35・8.4.26・8.5.11で、各版とも `verify.py` の全HTTP段階とSQLite行チェックを通過した。ログには `/api/probe` 200、手動の作成201・取得200・無効入力422、API Platformの作成201・返却IRIの取得200がある。`metadata:dump` は各版で1モデル・1123バイトのファイル生成を別に確認した。[証拠と制限](validation/laravel-api-evidence-2026-10-13.md)に記録する。Go製PHPの上記ローカル統合経路は成功したが、この記述時点（1167件）では版別公式比較や長時間継続動作まで確認したものではなかった。現行1169件の公式比較結果は冒頭に記録する。

その後、Go製PHPの同一プロセス・一時SQLiteで[2000リクエストのローカル継続試験](validation/full-stack-http-2000-stress-local.log)が終了0となった。GET応答とGC後の保持ヒープを確認したが、運用条件の長時間負荷、全HTTP契約の反復、公式PHP3版との現行実装の厳密比較は残る。GETの配列形式クエリとPOSTフォーム本文は[共通解析経路](full-stack-development.md)へ接続した。

Goの[HTTP統合テスト](../cmd/gophp/main_test.go)では、API PlatformのJSON-LD応答メディア型とOpenAPIのBookコレクションGET・POST定義も確認するようにした。[ローカル検証ログ](validation/full-stack-http-contract-strengthened-test.log)は終了0。公式PHPの出力や実サーバー経路との比較結果ではない。

Go製PHPの実CLI・実HTTPサーバーを通す[版別検証スクリプト](../research/laravel-api/verify-gophp.py)も用意した。[3設定のCLI単独実行](validation/full-stack-go-cli-openapi-subprocess-3versions-local.log)は別プロセスでmigration・独自コマンド・OpenAPI・DBテーブル作成が終了0。[8.4設定のサーバー試行](validation/full-stack-go-server-84-sandbox-attempt.log)ではHTTP待ち受けがこの環境の権限エラーで止まった。既存のin-process HTTPテストとは証拠を分ける。[Mac実行コマンド](full-stack-development.md)を参照する。

初期目標に含めない機能は、認証・認可、キュー、メール、ファイルアップロード、検索・ページネーションの詳細、GraphQL、管理画面、複数DB、運用用Webサーバー、完全なLaravel/API Platform互換である。後続の実装順序は[依存調査](dependency-audit.md)に記録する。
