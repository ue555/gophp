# 新規Laravel 13アプリをgophpenvのgophp v0.4.0（8.5）で起動（2026-10-22、WSL2 Arch Linux）

`composer create-project laravel/laravel:^13.0`で作成した未改変の新規アプリ（framework **v13.38.0**、通常77＋開発33の110パッケージ、`--no-scripts`）を`~/dev/laravel13-trial`に置き、Composerの実行だけ`localhost/gophp-dependency-research:8.5.11`イメージ内の公式PHPで行った。以降はすべて[gophpenv](https://github.com/ue555/gophpenv)が`v0.4.0`タグ（`602d6110106d`）からビルドしたエンジンを、グローバルモード8.5のシム`gophp`経由で実行した。上流の[固定Laravel 13の検証](laravel13-7oa6dbgc/REPORT.md)と同じframework版だが、別環境・別チェックアウトでの再現であり、上流の証拠と置き換えるものではない。

## CLI

| コマンド | 結果 |
| --- | --- |
| `artisan key:generate`（開発依存込みautoload、`.env`書き込み） | 成功、終了0 |
| `artisan package:discover` | pail・pao・tinker・carbon・collision を検出、成功 |
| `artisan about` | Laravel 13.38.0、PHP 8.5.11、Session/Cache/Queue=database、終了0 |
| `artisan migrate --force` | migrations・users・password_reset_tokens・sessions・cache・cache_locks・jobs・job_batches・failed_jobs を作成 |
| `artisan migrate --force`（再実行） | `Nothing to migrate.` |

## 実TCP（`gophp --serve=127.0.0.1:8097 --document-root=public public/index.php`）

| 要求 | gophp | 公式PHP 8.5.11（`php -S`、`scripts/php-host`経由） |
| --- | --- | --- |
| `GET /up` | 200 text/html 1,830B | 200 text/html 1,827B（応答時間表示の桁差） |
| `GET /`（Blade冷キャッシュ） | 200 text/html 70,403B、XSRF／セッションCookie | 200 text/html 70,403B |
| `GET /`（温キャッシュ） | 200、同サイズ | — |
| `GET /nonexistent-page` | 404 text/html 6,586B `<title>Not Found</title>` | 404 6,586B |
| `GET /favicon.ico` | 200 image/vnd.microsoft.icon 0B（静的配信） | 200 0B |
| `GET /robots.txt` | 200 text/plain; charset=utf-8 24B | 200 text/plain; charset=UTF-8 24B |
| `GET /.env`、`/../composer.json`、`/storage/../.env` | 404（front controllerへ委譲、ファイルは露出しない） | — |
| `GET /index.php` | 200（front controller） | — |
| 20連続 `GET /` | 全20件 200 | — |
| 4並列 `GET /` | 全4件 200 | — |

サーバーの標準エラー出力は空。コンパイル済みBladeは4ファイル、sessionsテーブルには要求数に応じた行が保存された。本文のバイト一致比較はCSRFトークン・セッションIDが毎回異なるため行っていない。CSS/JSはnpmビルドをしていないため、標準welcomeのVite manifest不在時の分岐で描画している。

## 制限

上流レポートの[対応の境界](laravel13-7oa6dbgc/REPORT.md)がそのまま当てはまる（property hooksはbacked `set`ブロックの部分対応、`get_defined_vars()`の順序、`array_unique`のフラグ、静的配信は開発用途）。`var_dump`・`file`・`mb_substr`・`mb_strlen`・`intdiv`・`array_is_list`などの内蔵関数は未実装のままで、未改変アプリの標準経路では呼ばれなかったが、アプリ固有のコードを追加すると当たる可能性がある。
