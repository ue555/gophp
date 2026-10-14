# Laravel / API Platform 初期目標と判定条件

主基準は PHP 8.4.26、追加対象は 8.3.35 と 8.5.11。Laravel統合方式を引き継ぎ、`composer.json` の直接指定は `laravel/framework` 12.69.3、`api-platform/laravel` 5.0.1、Symfony `type-info`・`web-link` の `^7.4` である。`composer.lock` は通常パッケージ94件・開発用0件を固定し、`vendor/composer/installed.json` と全件一致した。同じlockとvendorで公式PHPの3版の検証アプリが成功したため、この**調査アプリの共通構成**は成立した。`resolve.sh` 全体の終了コードは未確認である。[Laravel 12 の公式PHP対応表](https://laravel.com/framework/docs/12.x/releases)、[API PlatformのLaravel手順](https://api-platform.com/docs/laravel/)、[依存監査](dependency-audit.md)を参照。Symfony 8系はPHP 8.3と共通にできないため7.4系を選んだ。

検証アプリは [research/laravel-api](../research/laravel-api/README.md) に置く。公式PHPではこのアプリをCLI版 `php -S` の組み込みHTTPサーバーで実行し、各版ごとに独立した一時SQLiteファイルを使う。Go製PHPの現行CLIで確認できるのは、収録済みの部分構文を単一ファイルとして処理することまでである。Composer autoload、Laravel起動、HTTPリクエスト、PDO SQLite、API Platform実行はいずれも未対応であり、今後のHTTP実行環境も必要となる。

| 段階 | 入力・URL | 期待HTTP結果 | 期待DB状態 |
| --- | --- | --- | --- |
| autoload・起動 | `vendor/autoload.php` → `bootstrap/app.php`、`GET /api/probe` | 起動例外なし。200、`application/json`、`{"ok":true}` | `books` は空 |
| Laravel JSON | `GET /api/probe` | 200、`{"ok":true}` | 変化なし |
| SQLite保存 | `POST /api/books-manual`、`{"title":"Alpha"}` | 201、`{"id":N,"title":"Alpha"}`、Nは正の整数 | `books` に `(N, Alpha)` 1行 |
| SQLite取得 | `GET /api/books-manual/N` | 200、同じID・titleのJSON | 1行のまま |
| 入力検証 | `POST /api/books-manual`、`{}` | 422、JSONの `errors.title` に1件以上 | 1行のまま |
| API Platform作成 | `POST /api/books`、`Accept` / `Content-Type: application/ld+json`、`{"title":"Beta"}` | 201、JSON-LDの `title=Beta` と `@id` | `Beta` が追加され計2行 |
| API Platform取得 | 直前の `@id` へGET、`Accept: application/ld+json` | 200、JSON-LDの `title=Beta` | 計2行のまま |

上表の初期契約は公式PHP 8.3.35・8.4.26・8.5.11で、各版とも `verify.py` の全HTTP段階とSQLite行チェックを通過した。ログには `/api/probe` 200、手動の作成201・取得200・無効入力422、API Platformの作成201・返却IRIの取得200がある。`metadata:dump` は各版で1モデル・1123バイトのファイル生成を別に確認した。[証拠と制限](validation/laravel-api-evidence-2026-10-13.md)に記録する。Go製PHPではautoload、Laravel起動、HTTP、PDO SQLite、API Platformは未対応であり、公式PHPでの成功をGo側の成功としない。

初期目標に含めない機能は、認証・認可、キュー、メール、ファイルアップロード、検索・ページネーションの詳細、GraphQL、管理画面、複数DB、運用用Webサーバー、完全なLaravel/API Platform互換である。後続の実装順序は[依存調査](dependency-audit.md)に記録する。
