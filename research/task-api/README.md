# Laravelタスク管理API

Go製PHPの実用アプリ検証用。Laravel 12.69.3のルート・入力検証・EloquentとSQLiteを使う、認証なしのAPI。`research/laravel-api/`の回帰用Bookアプリとは独立している。認証、関連モデル、ページング、API Platformによる同等APIは実装範囲に含めない。

`composer.json`・`composer.lock`はこのアプリ自身が所有する。既存アプリと同じ94パッケージを固定し、API Platform 5.0.1も依存として保持するが、Taskには`ApiResource`を付けない。lockのパッケージレコードは既存lockと同一で、rootの名前に対応する`content-hash`だけ異なる。既存のlock・vendorには書き込まない。検証結果とGo変更の詳細は[検証記録](../../docs/validation/task-api-results.md)を参照。

## HTTP契約

JSONを送る際は`Content-Type: application/json`、応答形式には`Accept: application/json`を指定する。

| 要求 | 成功 | 存在しないID |
| --- | --- | --- |
| `POST /api/tasks` | 201、作成したTask | — |
| `GET /api/tasks` | 200、id昇順のTask配列。空は`[]` | — |
| `GET /api/tasks/{id}` | 200、Task | 404 |
| `PATCH /api/tasks/{id}` | 200、更新後のTask | 404 |
| `DELETE /api/tasks/{id}` | 204、本文・Content-Typeなし | 404 |

Taskは次の5項目を持つ。単体応答に`data`のラッパーはなく、一覧は配列。JSONのキー順には依存しない。

```json
{"id":1,"title":"日本語のタスク","completed":false,"created_at":"2026-10-21T03:00:00.000000Z","updated_at":"2026-10-21T03:00:00.000000Z"}
```

- `id`はSQLite自動採番。入力されたidや日時を代入しない。
- `title`は必須文字列、最大120文字（UTF-8の文字数）。Laravel標準ミドルウェアで前後の空白を除去する。空・空白のみ・null・数値・真偽値・配列・オブジェクトは422。
- `completed`は省略時false。Laravelの`boolean`規則に従い、JSONの`true`/`false`、整数`0`/`1`、文字列`"0"`/`"1"`を許容する。文字列は標準ミドルウェアでtrimされた後に検証される。null・空文字・`"true"`/`"false"`・その他の数値や配列などは422。応答では常にJSON真偽値。
- 日時はUTCの`YYYY-MM-DDTHH:mm:ss.ffffffZ`。SQLite保存は秒精度の`YYYY-MM-DD HH:mm:ss`で、JSONの小数部は6桁の0。作成日時は更新で維持し、実際に変更した場合だけ更新日時を更新する。同じ秒内の更新では日時が等しい場合がある。
- PATCHは省略項目を維持する。指定したtitleは空にできない。空JSONオブジェクト・本文なし・既存値と同じ値のみ・未知の項目のみのPATCHは200で、日時も含めて無変更。未知の項目は無視する。
- 入力検証を終えてから保存するため、422でDBは変わらない。422はLaravel標準の`{"message":"...","errors":{"title":["..."]}}`形式、404は`{"message":"..."}`形式。メッセージは英語。存在しないIDのPATCHは入力検証より先に404を返す。
- `{id}`は数字のルート制約。GETのクエリ文字列はフィルタとして扱わない。

## 新規チェックアウトからの準備

必要環境はGo（ルートの`go.mod`に従う）、Python 3、稼働中のPodman。リポジトリルートで実行する。VMの再起動や既存サービスの停止は不要。

```sh
make build
# 調査用イメージが存在しない場合のみ、指定版でビルドする。
podman build --build-arg PHP_VERSION=8.3.35 \
  -f containers/dependency-research/Containerfile \
  -t localhost/gophp-dependency-research:8.3.35 .
python3 research/task-api/record.py --label install-8.3.35 -- \
  podman run --rm -v "$PWD/research/task-api:/app" \
  localhost/gophp-dependency-research:8.3.35 \
  composer install --no-dev --no-interaction --prefer-dist --no-scripts --no-plugins
python3 research/task-api/record.py --label platform-8.3.35 -- \
  podman run --rm -v "$PWD/research/task-api:/app:ro" \
  localhost/gophp-dependency-research:8.3.35 composer check-platform-reqs --no-dev
python3 research/task-api/record.py --label prepare-go-8.4 -- \
  python3 research/task-api/prepare.py --version 8.4
```

Composerによる依存取得・platformチェックには公式PHPを使う。これはGoによるアプリ実行の証拠とは分ける。`composer update`や`--ignore-platform-reqs`は不要。vendorは新規アプリ内に作成し、既存vendorへのsymlinkは作らない。

`prepare.py`は**Go製PHP**で`package:discover`→`migrate --force`→`migrate:status`を実行する。初回の空キャッシュから実行できる。手動利用のDBは`database/database.sqlite`、storageは`storage/`、bootstrapキャッシュは`bootstrap/cache/`。再実行しても既存のDB行を消去しない。これらとvendorはGit管理対象外。

既存アプリではmigration前のAPI Platformメタデータが問題になった。本アプリは`resources: []`・`cache: array`で永続的なモデルメタデータを使わず、`api-platform:install`やmetadata dumpは不要。HTTP検証では念のため、一時storageのファイルキャッシュをmigration後に空にする。既存アプリの回帰検証では従来どおりこの削除が必要。

## 手動起動・使用例

```sh
cd research/task-api
../../bin/gophp --php-version=8.4 --serve=127.0.0.1:8091 public/index.php
```

ポート8091は既存アプリの8080と分離する。使用中の場合は別ポートを指定する。停止はこのサーバーの端末でCtrl-C。手動DBは残る。

別の端末で実行する。

```sh
curl -i -H 'Accept: application/json' http://127.0.0.1:8091/api/tasks
curl -i -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' \
  --data '{"title":"日本語のタスク"}' http://127.0.0.1:8091/api/tasks
curl -i -X PATCH -H 'Accept: application/json' -H 'Content-Type: application/json' \
  --data '{"completed":true}' http://127.0.0.1:8091/api/tasks/1
curl -i -X DELETE -H 'Accept: application/json' http://127.0.0.1:8091/api/tasks/1
```

DBやstorageは環境変数`DB_DATABASE`・`LARAVEL_STORAGE_PATH`で変更できる。設定ファイルに固定の絶対パスはない。カスタムパスは起動と準備の両方で同じ絶対パスを環境変数に設定する（例：`export DB_DATABASE="$PWD/research/task-api/database/manual.sqlite"`）。

## 再実行可能な実HTTP比較

公式比較用には追加で8.4.26と8.5.11の同名調査用イメージを上のbuildコマンドの版を替えて準備する。依存を更新する必要はない。

```sh
# 公式PHP 8.3.35→Go 8.3、8.4.26→Go 8.4、8.5.11→Go 8.5の順に比較
python3 research/task-api/verify.py

# 版別。各コマンドが公式とGoの両方を実行して比較する
python3 research/task-api/verify.py --version 8.3
python3 research/task-api/verify.py --version 8.4
python3 research/task-api/verify.py --version 8.5

# 環境障害時などの独立実行。単独実行は「公式との比較済み」としない
python3 research/task-api/verify.py --runtime go --version 8.4
python3 research/task-api/verify.py --runtime official --version 8.4
```

各実行は`docs/validation/task-api-日時-ランダム文字列/`に`run.log`、版別サーバーログ、HTTP生本文を含む`*-http.jsonl`、比較用`*-normalized.json`を作る。上書きしない。コマンド・環境変数・対象版・終了コード・GoバイナリSHA-256を記録する。公式PHPのパッチ版も実行して確認する。失敗時もログは残り、後続の版を成功扱いしない。

全シナリオは同じ関数で実行する。空一覧、日本語作成、省略時false、昇順一覧・個別取得、部分更新・空PATCH、titleの120/121文字境界と不正型、completedの許容／不正入力、422時のDB維持、不存在ID、削除後の取得・一覧、リクエスト状態分離、連続100 GET・8並列GET、全件削除後の空一覧を検証する。SQLiteの実行とDB行をPython側から別接続で確認する。

各runtime・版には独立した空の一時DB、storage、bootstrapキャッシュを割り当てる。公式はPodmanの実PHPサーバー、Goは実`--serve`サーバーにTCPで接続する。起動確認はTCP接続のみで、アプリに追加のHTTP要求を送らない。並列8要求は送信順・完了順を保証しないため、記録は要求番号順に揃える。

比較するのはHTTPステータス、Content-Type、JSONの全構造・値。日時は形式、現在時刻との距離、作成・更新の順序、更新時の変化／維持、DBとの一致を確認したうえで、比較用の**応答**の`created_at`/`updated_at`だけを固定マーカーへ置き換える。JSONキー順・JSON外側の整形用空白はJSONパーサーの意味比較に従う。警告・余計な出力・不正JSONは除去せず失敗とする。Dateや接続などの転送用HTTPヘッダーは比較対象外。204は本文なし・Content-Typeなしを比較する。

終了時には、この検証が起動したプロセス／コンテナを停止し、一時DBとstorageだけを削除する。Goサーバーの`server_exit_code: -15`は検証終了後のSIGTERM。`verify.py`自身の終了コードとは別。Podman/TCPの環境障害が起きた場合は実エラーを保存し、同じ環境で繰り返さず独立実行へ進む。未実施の比較は完了と判定しない。

保存済み6実行分の生HTTP・比較用JSONを再監査する場合（日時の現在値チェックやDB実行検証は元のrun.logを参照）：

```sh
python3 research/task-api/record.py --label http-artifact-audit -- \
  python3 research/task-api/audit-evidence.py docs/validation/task-api-20261021-124313-i5cu398g
```

## Go変更の回帰検証

```sh
python3 research/task-api/record.py --label test -- make test
python3 research/task-api/record.py --label vet -- make vet
python3 research/task-api/record.py --label build -- make build
python3 research/task-api/record.py --label fixtures-all-3versions -- \
  ./containers/php-version/compare-all.sh tests/fixtures ./bin/gophp
python3 research/task-api/record.py --label legacy-go-http-3versions -- \
  python3 research/task-api/regression.py
```

`regression.py`は既存`research/laravel-api/verify-gophp.py`をそのまま利用する。bootstrapキャッシュを追加で隔離し、3設定で既存Book CRUD・API Platform・OpenAPI・100連続／8並列GETを実行する。既存アプリのDB・vendor・lockは変更しない。fixture比較のstdout/stderr/終了コード一致と、タスクAPIおよび既存アプリの実HTTP証拠は別々に記録する。

## リリース再検証とCI

Docker環境でも同じHTTPシナリオを実行できる。`verify.py --engine docker`で切り替え、コンテナ内の実行ユーザーをホストUID/GIDへ合わせる。Podmanは既定のユーザーマッピングを使用する。`--output-root`でHTTP証拠の保存先、`--image-template 'localhost/gophp-dependency-research:{version}'`で固定版イメージ名を指定できる。HTTPログと一緒にHEAD・ソースmanifestのSHA-256・未コミット差分のSHA-256・lock・バイナリの識別情報も保存する。

両アプリの新規依存導入、ローカル回帰、公式全fixture比較、タスクHTTP比較、既存アプリの2000連続／8並列GETは`scripts/verify-release.py`へまとめた。CIも同じスクリプトをDockerで版別に実行し、失敗時も証拠一式をartifactへ保存する。[リリース手順・公開記録](../../docs/task-api-release.md)を参照。既存手動vendorのある作業場所ではなく、新規チェックアウトで実行する。
