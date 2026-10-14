# Laravel / API Platform 依存調査用アプリ

比較用 PHP コンテナとは別の環境です。直接指定は `laravel/framework` 12.69.3 と `api-platform/laravel` 5.0.1、Symfony `type-info`・`web-link` は `^7.4` です。Composerのlockは通常パッケージ94件・開発用0件を固定し、vendorと一致しています。公式PHP 8.3.35・8.4.26・8.5.11で、この共通構成の検証アプリが成功しました。[検証記録](../../docs/validation/laravel-api-evidence-2026-10-13.md)に根拠と制限があります。

初回の解決に使った `resolve.sh` はPHP 8.3.35でlockを生成し、同じlockを8.4.26・8.5.11で検査しました。**確定したlockの再構築には下記の `composer install` を使い、`resolve.sh` の `composer update` は再実行しません。** 比較用コンテナの設定は変更しません。

`composer.lock` は保存対象、`vendor/`、`.cache/`、一時SQLiteと生成キャッシュはGit管理対象外です。Composerの自動スクリプトを実行しないため、依存解決とアプリ起動の検証は別です。

2026-10-13 時点で `composer.lock` と `vendor/` は作成済みです。提供されたMacでの検証ログは3版でHTTP成功を示します。このエージェント環境からのPodmanソケット接続は `operation not permitted` で拒否されるため、今回ここでの再実行はできません。[依存コード調査](../../docs/dependency-audit.md)は実vendorの代表コードまで更新しています。

このディレクトリには最小Laravelアプリも置いた。`bootstrap/app.php` は Laravel 12 形式で Composer autoload から起動し、`routes/api.php` は JSON 応答・SQLite保存/取得・入力検証を行う。`app/Models/Book.php` の `#[ApiResource]` は Laravel向けAPI Platformの公式統合方式に従う。`database/migrations/` は `books` の最小テーブルを作る。既存の比較用PHPコンテナは変更しない。

Macで既存の調査用イメージとvendorがある場合、HTTPの再検証は次のとおり。Podman VMが起動済みなら再起動せず、作成済みの比較コンテナを削除・初期化しない。ログを `tee` に渡すときは終了コードを直後に保存する。

```zsh
cd /Users/kouji/dev/php
podman info
set -o pipefail
python3 research/laravel-api/verify.py 2>&1 | tee /tmp/laravel-api-verify-retry.log
verify_exit=$?
echo "verify.py 終了コード: $verify_exit"
```

vendorをlockから再取得する場合は、調査用の独立したイメージを使う。`--ignore-platform-reqs` は使わない。各版で `composer install` と `check-platform-reqs` を行い、版ごとの終了コードを直後に保存する。調査用イメージがなければ、先に `podman build --build-arg PHP_VERSION=8.3.35 -f containers/dependency-research/Containerfile -t localhost/gophp-dependency-research:8.3.35 .` のように版ごとにビルドする。

```zsh
cd /Users/kouji/dev/php
set -o pipefail
for version in 8.3.35 8.4.26 8.5.11; do
  image="localhost/gophp-dependency-research:$version"
  podman run --rm -v "$PWD/research/laravel-api:/app" "$image" composer install --no-dev --no-interaction --prefer-dist --no-scripts --no-plugins 2>&1 | tee "/tmp/laravel-api-install-$version.log"
  install_exit=$?
  echo "$version install 終了コード: $install_exit"
  podman run --rm -v "$PWD/research/laravel-api:/app" "$image" composer check-platform-reqs --no-dev 2>&1 | tee "/tmp/laravel-api-platform-$version.log"
  platform_exit=$?
  echo "$version platform 終了コード: $platform_exit"
done
```

`verify.py` は版ごとに `.cache` 下の独立した一時SQLiteとファイルキャッシュを作り、migrate後にキャッシュを消去する。`package:discover`、`migrate`、`api-platform:install`、`metadata:dump --path=/probe/metadata.dump` を実行し、ダンプの実ファイルを確認してからHTTPとDB状態を検査する。一時DBとダンプは終了時に削除される。ダンプ生成とHTTP成功は別の判定である。

vendor監査では、`composer.lock` の全パッケージ名・版を `vendor/composer/installed.json` と照合し、autoloadとLaravel/API Platformの代表実コードを目視確認した。`audit-vendor.php` の構文候補とHTTP成功だけで個別の内部経路を実行確認済みとはしない。結果は[依存コード調査](../../docs/dependency-audit.md)に記録した。再監査は次のコマンドで、終了コードを直後に記録する。

```zsh
cd /Users/kouji/dev/php
podman run --rm -v "$PWD/research/laravel-api:/app" localhost/gophp-dependency-research:8.3.35 php audit-vendor.php > /tmp/laravel-api-vendor-audit.json 2> /tmp/laravel-api-vendor-audit.stderr.log
audit_exit=$?
echo "audit-vendor.php 終了コード: $audit_exit"
```

`audit-vendor.php` の出力はlockにあるパッケージと構文候補の件数・代表位置であり、起動時の到達やGo対応状況を自動判定しない。結果を見て実ファイルを確認し、3版のHTTPトレースと突き合わせる。

根拠: [Laravel 12 のサポート表](https://laravel.com/framework/docs/12.x/releases)、[Laravel 12.69.3 パッケージ](https://packagist.org/packages/laravel/framework)、[API Platform Laravel 5.0.1 パッケージ](https://packagist.org/packages/api-platform/laravel)、[Symfony 7.4 の PHP 8.2 以上という要件](https://symfony.com/releases/7.4)、[API Platform Laravel のインストール案内](https://api-platform.com/docs/laravel/)。
