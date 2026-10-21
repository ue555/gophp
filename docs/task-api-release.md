# タスクAPI追加版のリリース記録

公開候補は**v0.3.0**。既存v0.1.0・v0.2.0の次のminor版として、タスク管理APIと初期manifest生成・日時JSON化のGo対応、両アプリの再構築／比較CIを追加する。既存タグ・履歴は変更しない。

## 実装と検証対象

- Laravel 12.69.3・API Platform 5.0.1を含む固定94パッケージ。両アプリの既存lockを保持する。
- Task CRUD・入力検証・状態分離・100連続／8並列GET。認証などの機能追加は行わない。
- 既存BookアプリのLaravel／API Platform／OpenAPI回帰、2000連続／8並列GET。
- 開始時のPHPソース1,173件とexpected.jsonの1,173エントリーを照合。内訳は正常745・実行時294・構文134。
- 対象は公式PHP 8.3.35・8.4.26・8.5.11とGoの8.3・8.4・8.5設定。

[実装時の検証](validation/task-api-results.md)は過去の成功証拠であり、このリリースの変更後検証とは区別する。

## 独立配置からの再構築

Go 1.27.0以上、Python 3、Cコンパイラ、SQLite・PCRE2開発用ヘッダー／ライブラリ、DockerまたはPodmanを用意する。Linuxでは`build-essential libsqlite3-dev libpcre2-dev python3`、Apple Silicon macOSでは既存のHomebrew PCRE2配置を使用する。CGOが必要。

新しいチェックアウトのルートで実行する。既存vendor・bootstrap/cacheがある場合、`--require-fresh-vendor`は処理を拒否し、削除しない。

```sh
python3 -u scripts/verify-release.py --engine podman \
  --build-images --require-fresh-vendor --output-root docs/validation
```

LinuxのDockerでは`--engine docker`へ変更する。CIは版ごとに独立したcheckoutを使い、例えば次を実行する。

```sh
python3 -u scripts/verify-release.py --engine docker --version 8.4.26 \
  --build-images --require-fresh-vendor --output-root .cache/ci
```

他版のコマンドは`--version 8.3.35`・`--version 8.5.11`。版省略は3版を順次実行する。依存は`composer install --no-dev --no-scripts --no-plugins`で固定lockから新規導入し、全パッケージの名前・版とplatform requirementsを照合する。Composer工程は公式PHPによる準備であり、Goの実行成功と区別する。

生成済みmanifestは持ち込まず、両アプリのvendorがあることを必須条件にして`make test vet build`を実行する。既存テストのvendor不足によるskipを防ぐ。fixture比較は独自の公式PHPコンテナを同じ絶対パスにmountして行い、Mac固定パスのラッパーに依存しない。

タスクAPIは既存`verify.py`へDocker／Podman選択と証拠保存先を追加したものを使用し、各runtime・版の空DB・storage・bootstrapキャッシュ分離と191要求の契約を維持する。DockerではホストUID/GIDで実行して書込権限と後片付けを確保する。既存アプリは`regression.py`経由で元のHTTP検証を呼び、bootstrapキャッシュを隔離する。

## 証拠と保護

`release-task-api-日時-識別子/`にコマンド別生ログ、`commands.json`の実終了コード、`result.json`、fixture件数、実行前・ビルド後・実行後の`source-*.json`を保存する。タスクHTTPの生本文・比較用JSON・サーバーログ・ソース識別情報は版別サブディレクトリに保存する。新しい実行は別のディレクトリを作り、失敗ログを上書きしない。`.gitattributes`の生ログ保護をネストした証拠ディレクトリにも拡張し、末尾空白・改行を加工せず保存する。

`source_sha256`はGit候補に含まれる実行ソース・設定・テストのパス、内容hash、実行権、symlink先から作る。更新される証拠`docs/`と説明Markdownは除外し、HEADと`diff_sha256`を別途記録する。コミット前の未追跡ソースも対象。実行前後でこのhashとlockが変わらず、ビルド後のバイナリが変わらないことを検査する。CIではさらにartifact名に`github.sha`を含める。

手動DB・既存lock/vendor・既存8091サーバーを保護する。検証は独立配置で行い、独自oracle・一時HTTPサーバー・一時DB/storageだけを終了時に片付ける。候補からvendor、DB、`.env`、キャッシュ、バイナリ、既知の秘密情報パターンを除外する。既存の追跡済み`public/vendor`静的アセットは保持する。ログ中のテスト用固定APP_KEYとローカルパスは検証入力・証拠として保存する。

## CIと公開手順

[GitHub Actions](../.github/workflows/ci.yml)の3版matrixすべてを必須検証とする。新規依存取得・make test/vet/build・diff検査・全fixture・タスクAPI公式／Go比較・HTTP成果物監査・既存API回帰を同じreleaseスクリプトで実行し、失敗時もログをartifactへ保存する。タグpushによる同一ソースの重複実行を避け、branch push・pull request・手動起動をトリガーとする。

公開は通常commit/push後、リモートmainとCIのheadShaを照合し、3必須ジョブの成功を確認してから、そのSHAへ新しい注釈付きv0.3.0を付ける。既存タグの上書きやforce pushは行わない。追跡文書を追加更新した場合は再commitし、新しいSHAのCI成功を確認してからタグを付ける。

## 今回の状態

独立配置での再構築と全検証を実行中。CI成功・公開参照は確認後に確定する。過去の成功結果を今回の成功へ転用しない。
