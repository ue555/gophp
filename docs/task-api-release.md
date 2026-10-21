# タスクAPI追加版のリリース記録

本リリースの版番号は**v0.3.0**。既存v0.1.0・v0.2.0の次のminor版として、タスク管理APIと初期manifest生成・日時JSON化のGo対応、両アプリの再構築／比較CIを追加する。既存タグ・履歴は変更しない。

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

## 今回の検証結果（2026-10-21）

新しい独立配置へ固定lockから両vendorを導入し、Podmanで全工程を完走した。[実行結果](validation/release-task-api-20261021-140007-49d645lt/result.json)は終了0。[工程別コマンドと実終了コード](validation/release-task-api-20261021-140007-49d645lt/commands.json)、[make test・vet・build](validation/release-task-api-20261021-140007-49d645lt/16-local-test-vet-build.log)（Goテスト・Python 15件）、差分検査もすべて成功した。

| 公式PHP / Go設定 | 全fixture比較 | タスクAPI実TCP・DB比較 | 既存API実HTTP回帰（Go） |
| --- | --- | --- | --- |
| 8.3.35 / 8.3 | 1,173一致・差異0・エラー0 | 191要求×2 runtime、差異なし | 2,000連続・8並列成功 |
| 8.4.26 / 8.4 | 1,173一致・差異0・エラー0 | 191要求×2 runtime、差異なし | 2,000連続・8並列成功 |
| 8.5.11 / 8.5 | 1,173一致・差異0・エラー0 | 191要求×2 runtime、差異なし | 2,000連続・8並列成功 |

追加の手動集計では、2,000要求すべてをprobeと仮定した監査が[終了1](validation/task-api-release-final-local-audit-20261021-144802-ljuoieba.log)となった。元の回帰はprobe 1,750回とBook取得250回の構成であり、[要求構成を照合した再監査](validation/task-api-release-final-local-audit-corrected-20261021-144840-m_aum6yq.log)は終了0。失敗ログを保持し、APIの実装・期待値・元の検証器は変更していない。

生ログとHTTP成果物は[今回の証拠ディレクトリ](validation/release-task-api-20261021-140007-49d645lt/)に保存した。Task CRUD、入力検証、失敗時のDB維持、部分更新、404、204本文なし、状態分離、100連続・8並列GETを含む。HTTP応答の型と日時の形式・整合性を検査したうえで、比較用の日時だけを正規化し、生本文はそのまま保存する。既存fixtureのstdout・stderr・終了コード比較とは別の証拠である。

独立配置の実行ソースhashは`390aead845413bc831d43ad45d34b2e381834716a5497397fb5944045db34de5`。公開候補のhashは`a97df3f886b34ddf56b07acbe7aa30d0bc33c57145c6d9d774e78c6c74b08c3e`で、差はsnapshot後に追加した`.gitattributes`のネストした生ログ属性だけ。[照合記録](validation/task-api-release-source-review-20261021-141022-jr3p8c_c.log)でアプリ・Go・テスト・CIの同一性を確認した。実行前後のソース・lockと、ビルド後／実行後のバイナリhashは一致している。

候補コミット`f5054981730e5f9044925b7392ec541949db4512`の[Docker CI](https://github.com/ue555/gophp/actions/runs/39547883982)は3必須ジョブすべて成功。[GitHubのhead SHAとジョブ結果](validation/release-task-api-candidate-ci.json)、[ダウンロードした3版artifactのコマンド・hash監査](validation/release-task-api-candidate-ci-artifact-audit.json)を保存した。各CIジョブでも独立checkoutから全工程を実行し、工程の省略はない。

今回のリリース整備でGo処理系への追加修正はない。検証基盤にはDocker／Podman対応とソース識別を追加し、終了コード保持・未追跡ファイル／実行権／dangling symlinkの識別・fixture整合性の最小回帰テストを加えた。タスクAPI実装時のGo修正と4件のPHP fixtureは[実装記録](validation/task-api-results.md)に記載している。

既存lock・vendor・手動DBの[13,349ファイルの保護監査](validation/task-api-release-protected-files-review-20261021-143600-i9x16b3k.log)は一致。8091番の手動サーバーを維持し、独自検証サービスと一時データだけを片付ける。今回の検証に未実施項目はない。一般PHP機能の未対応範囲は[互換性文書](compatibility.md)に残す。

## 最終公開の確定情報

この結果と生ログを保存した最終コミットも、push後に3版CIを実行する。候補CIの成功を最終SHAへ転用しない。最終SHA・そのSHAのCI URL・ソースhashは、全必須ジョブ成功後に作成する注釈付きタグ`v0.3.0`のメッセージに記録する。`git show --no-patch v0.3.0`で取得でき、[公開タグのソース](https://github.com/ue555/gophp/tree/v0.3.0)から本記録を参照できる。

公開操作と最終リモート参照確認のローカルログは`.cache/release-publication/`に保存し、タグ対象の検証済みコミットを変更しない。既存v0.2.0は`cc9388e6b1d600f8cae341281718c68b56f7f0c6`を維持する。
