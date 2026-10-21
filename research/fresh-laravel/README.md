# 未改変Laravel 12の検証

[Issue #1](https://github.com/ue555/gophp/issues/1)の再現・回帰用アプリ。
`app/` は `laravel/laravel v12.12.2`（commit
`945f4e5a9fd3695dc0ee512f497c650fb82cfbb8`）の未改変スケルトン。
Framework **v12.70.0**、開発用35パッケージを含む**112パッケージ**を
独立した `composer.lock` に固定する。
元IssueのアプリとlockはこのMacに存在しなかったため、同一lockの再現とは扱わない。
初回のPHP 8.4による解決ではSymfony 8.1系が入り、PHP 8.3で動かないため、
PHP 8.3で解決し直した。Frameworkは変更せず、Symfonyの該当6パッケージは7.4系となった。

Laravel・vendor・標準ルート・middleware・welcomeビューは変更しない。
`.env`、APP_KEY、SQLite、storage、bootstrap/cacheだけを通常の環境設定として用意する。
既存 `laravel-api` / `task-api` とvendor・DB・キャッシュ・ポートを共有しない。
検証の結果・生ログ・残る一般機能の制限は
[今回の記録](../../docs/validation/fresh-laravel-issue-1-results.md)を参照。

## 新しいチェックアウトからの準備

リポジトリルートで実行する。Go/Cコンパイラ・SQLite/PCRE2開発ライブラリ・
Python 3・稼働中のPodmanが必要。Macのネイティブビルド条件はルートREADME参照。

```sh
make build
for version in 8.3.35 8.4.26 8.5.11; do
  podman build --build-arg PHP_VERSION="$version" \
    -f containers/dependency-research/Containerfile \
    -t "localhost/gophp-dependency-research:$version" .
done
python3 research/fresh-laravel/record.py --label rebuild -- \
  python3 research/fresh-laravel/rebuild.py
```

`rebuild.py` は `.cache/fresh-laravel/rebuild-*/app` に新規配置し、固定lockから
**公式PHP 8.3.35のComposer**で `install --no-scripts` を実行する。
開発依存を外さず、公式PHP3版で拡張要件を検証し、lockとインストール済み一覧を照合する。
この依存準備をGoでのアプリ実行成功には数えない。既存vendorは再利用しない。
必要な拡張の実測一覧と `check-platform-reqs --lock` の結果を版別に保存する。
イメージはintl・mbstring・pdo_sqliteに加え公式CLI標準のopenssl、tokenizer、DOM/XML等を含む。

## 実TCP・CLI・SQLite・Bladeの比較

再構築したアプリの場所を使用する。`--app` はリポジトリ内の配置先を指定する。

```sh
fresh_app=$(cat .cache/fresh-laravel/rebuilt-app.txt)
python3 research/fresh-laravel/record.py --label http-all -- \
  python3 research/fresh-laravel/verify.py --app "$fresh_app" --binary bin/gophp
```

版を分けて再実行する場合：

```sh
python3 research/fresh-laravel/verify.py --app "$fresh_app" --version 8.3 --binary bin/gophp
python3 research/fresh-laravel/verify.py --app "$fresh_app" --version 8.4 --binary bin/gophp
python3 research/fresh-laravel/verify.py --app "$fresh_app" --version 8.5 --binary bin/gophp
```

公式PHPは順に8.3.35 / 8.4.26 / 8.5.11を使用する。`--runtime go` / `official` は
片側の調査用で、両者比較成功には数えない。Dockerの場合は両スクリプトに
`--engine docker` を指定し、イメージ作成もDockerで行う。DockerにはホストUID/GIDを渡す。

各runtimeに空DB・独立storage・bootstrap/cacheを割り当てる。Composerスクリプトの
公式実行結果や既存manifestには依存せず、Go/公式それぞれのbootstrapがpackage discoveryを行う。
標準アプリにはAPI Platformがないため、既存APIアプリで必要だったmigration後の
API Platformメタデータキャッシュ破棄は不要。

検証内容：

- 開発用autoload、`.env`のAPP_NAME・SQLiteパス・変数展開の実反映。
- `artisan list` / `about --json` / `route:list --json`。開発コマンドとweb middlewareも確認。
- 空SQLiteへの標準3 migration、9テーブル、再migration時の状態維持。
- 実サーバーの `/up` と `/`、Bladeキャッシュ空からの初回と生成済みキャッシュでの後続要求。
- Cookieを再送して同じセッションとCSRFトークンを維持し、8並列要求では独立したセッションを生成。
- 100連続GETと8並列GETを含む各runtime **111要求**。DBの他テーブルに副作用がないこと。
- HTTPステータス、Content-Type、本文、Cookie属性、CLI stderr、サーバー／アプリログ。
- 標準2ビューと生成済み2ビューの `token_get_all()` 出力を公式と比較。

比較で除外する動的値は、配置先パス、healthビューの整数ミリ秒、HTTP Date、
Cookieの暗号文・Expires・Max-Age、session ID・CSRFトークン・最終活動日時。
値の形式・期限の整合性・独立性・維持を別途検証してから比較する。
CookieのMax-AgeはSymfonyが `expires - time()` で算出するため7190〜7200秒を許容し、
Expires/HTTP Dateとの誤差も2秒以内とする。本文全体のtrim、警告除去は行わない。
HTTPのtransportヘッダー（Server/Connection/Content-Length/Transfer-Encoding）は比較対象外。
Content-Type/Cache-Control/Vary/CORSヘッダーは比較する。
Composerは準備専用で、アプリ実行時PATHを両者 `/usr/bin:/bin` に揃える。
公式PHPは絶対パスで起動し、`about`のComposer表示も両者 `-` とする。

終了時に検証スクリプトが起動したサーバーと一時DB/storageだけを片付ける。
再構築先vendorと証拠は残す。手動DBと既存サービスは変更しない。

## 手動起動

`app/` に直接依存を入れる場合も、公式PHPを使うのはComposerのみ。

```sh
podman run --rm -v "$PWD/research/fresh-laravel/app:/app" \
  localhost/gophp-dependency-research:8.3.35 \
  composer install --no-scripts --no-interaction --prefer-dist
cd research/fresh-laravel/app
test -e .env || cp .env.example .env
../../../bin/gophp --php-version=8.4 artisan key:generate
test -e database/database.sqlite || touch database/database.sqlite
../../../bin/gophp --php-version=8.4 artisan package:discover
../../../bin/gophp --php-version=8.4 artisan migrate --force
../../../bin/gophp --php-version=8.4 --serve=127.0.0.1:8093 public/index.php
```

既存の有効なAPP_KEYがある場合、`key:generate`を再実行しない。
別ターミナルで `curl -i http://127.0.0.1:8093/up` と
`curl -i http://127.0.0.1:8093/` を利用できる。終了は起動したターミナルでCtrl-C。

## 処理系の回帰検証

```sh
make test
make vet
make build
git diff --check
python3 research/fresh-laravel/compare-fixtures.py
python3 research/task-api/verify.py --interpreter bin/gophp
python3 research/laravel-api/verify-gophp.py --interpreter bin/gophp --stress-requests 100
python3 research/fresh-laravel/audit.py
```

全fixture比較は既存 `compare_php.py` と公式PHP wrapperを利用する。
wrapperの準備は `containers/php-version/README.md` を参照。
`docs/validation/fresh-laravel-*` に実行ごとに別名でコマンド・終了コード・
HEAD/差分/ソースmanifest・バイナリSHA-256・lock・HTTP/DB/Blade成果物を保存する。
失敗ログを上書きしない。`observe.py` は新規fixtureの公式観測を採取する開発用ツールで、
既存期待値の自動更新には使わない。
