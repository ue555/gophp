# Laravel / API Platform 依存調査の実行記録

作業ディレクトリ `/Users/kouji/dev/php`。既存の全527 fixtureと期待値は変更していない。ユーザー報告では公式PHP 8.3.35・8.4.26・8.5.11で全527件が一致したが、このターンでは比較ログの実体を受け取っておらず再実行もしていない。今回の対象は依存調査であり、Go言語機能の実装は変更していない。

候補の直接依存は Laravel Framework 12.69.3、API Platform Laravel 5.0.1、Symfony type-info / web-link 7.4系。Laravelの公式サポート表と指定タグ、API Platformの指定タグとPackagist制約で存在と直接制約を確認した。全推移的依存の共通解はComposer未実行で未確定。比較用3コンテナの構成は変更せず、独立した `containers/dependency-research/Containerfile` を使う。

[生ログ](laravel-api-dependency-attempt.log)にはコマンド・終了コード・標準出力と標準エラーを保存した。

| コマンド | 終了コード | 判定 |
| --- | ---: | --- |
| `podman info --format '{{.Version.Version}}'` | 125 | Podmanソケットへの接続が `operation not permitted`。VM状態は確認不能 |
| `./research/laravel-api/resolve.sh` | 125 | 最初のイメージビルド前に接続失敗。Composer実行0回 |
| `python3 research/laravel-api/verify.py` | 2 | lockとvendorの前提検査で停止。PHP/HTTP/DB検証0件 |
| `make test` | 0 | 既存GoテストとPythonテストが成功 |
| `git diff --check` | 0 | 実行時点の差分に空白エラーなし |

`research/laravel-api/composer.json` のJSON解析と `verify.py` のPython構文解析は成功した。PHP CLI、Composer CLIはホストに見つからなかった。`composer.lock` と `vendor/` は存在しない。PHP、Composer、拡張の実バージョン、lockのパッケージ版、`audit-vendor.php` の実行、公式PHPでの最小アプリ起動・JSON・SQLite・API Platform・検証エラーの結果はすべて未実施。したがって項目1・2は完了扱いにしない。検証用PHPファイルの構文も公式PHPでは未確認である。

Macでの残作業は[調査アプリの再構築手順](../../research/laravel-api/README.md)を使う。`resolve.sh` の成功で生成された `composer.lock` を保存し、3版の `check-platform-reqs`、HTTP検証、vendorトークン一覧、実コードの目視監査を順に行う。失敗時は各コマンドの終了コードと阻害要因を保存し、PHPや拡張要件を無視して成功としない。
