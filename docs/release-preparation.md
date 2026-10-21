# 固定アプリ到達版のリリース準備

この文書はv0.2.0準備時の履歴。タスクAPI追加後の現行手順・CI・公開状況は[タスクAPI追加版のリリース記録](task-api-release.md)を参照する。

対象は差分整理・別配置からの再構築・GitHub保存・継続検証の4項目。固定アプリの6段階の完成とは別の作業である。公開候補タグは既存v0.1.0に続く`v0.2.0`。

## 再構築

Go 1.27.0、Cコンパイラ、SQLiteとPCRE2のヘッダー／ライブラリ、Python 3、DockerまたはPodmanを用意する。Linuxでは`build-essential libsqlite3-dev libpcre2-dev`が必要。Apple SiliconのmacOSではPCRE2をHomebrewの`/opt/homebrew`配下へ導入する。SQLite接続と一部の正規表現処理はCGOでネイティブライブラリへ接続する。

新しいチェックアウトのルートから以下を実行する。Dockerを使う場合は`podman`を`docker`へ置き換える。固定vendorを必要とするfixtureがあるので、依存導入をテストより先に行う。

```sh
mkdir -p research/laravel-api/bootstrap/cache research/laravel-api/storage/framework/cache/data research/laravel-api/storage/framework/views research/laravel-api/storage/framework/sessions research/laravel-api/storage/logs
podman build --build-arg PHP_VERSION=8.4.26 -f containers/dependency-research/Containerfile -t gophp-dependencies .
podman run --rm -v "$PWD/research/laravel-api:/app" gophp-dependencies composer install --no-dev --no-interaction --prefer-dist --no-scripts --no-plugins
podman run --rm -v "$PWD/research/laravel-api:/app" gophp-dependencies composer check-platform-reqs --no-dev
podman run --rm -e DB_DATABASE=:memory: -v "$PWD/research/laravel-api:/app" gophp-dependencies php artisan package:discover --no-interaction
git diff --exit-code -- research/laravel-api/composer.lock
make test vet build
python3 -u research/laravel-api/verify-gophp.py --version 8.4 --stress-requests 2000
```

固定PHP比較のMac用既存ラッパーは`/Users/kouji/dev/php`マウントを前提とする。別配置ではCIと同じくチェックアウトを同一絶対パスで公式PHPコンテナへマウントする。CIはこの方式で3パッチ版を個別に比較し、パスを含む診断も厳密に比較する。

## 継続検証

[GitHub Actions](../.github/workflows/ci.yml)はpush・pull request・手動実行で、版別に固定vendor導入、Go／Pythonテスト、vet、ビルド、公式PHP全fixture比較、Go実HTTP2000回・8並列検証を行う。失敗時も取得済みログをartifactへ保存する。依存を最新版へ更新しない。CIの成功はGitHub上で実行後に確定する。

## 公開

前回のセッションでは`.git/index.lock`作成が権限エラー、GitHub接続が名前解決エラーとなり、commit・push・タグ作成は未実施だった。2026-10-21の再開時はサンドボックス外でGit書き込みが可能で、GitHub接続とpush権限を確認した。既存履歴を書き換えるforce pushは使わない。

```sh
git add -A
git diff --cached --check
git diff --cached --stat
git commit -m "feat: 固定Laravel・API Platformアプリの動作検証とCIを整備"
git push origin HEAD
```

push後、対象コミットの3版CIが成功したことを確認してからタグを付ける。タグが既にある場合は上書きしない。

```sh
git tag -a v0.2.0 -m "固定Laravel・API Platformアプリの検証完了版"
git push origin v0.2.0
```

今回の候補ファイルでは既知の秘密鍵・GitHubトークン・AWSアクセスキー形式は未検出。これは全形式の秘密情報が存在しない保証ではない。vendor、環境設定、ビルドキャッシュ、生成バイナリは既存ignoreで除外される。公開用API Platformアセットは既に追跡されているため維持する。歴史的ログ中のローカルパスとテスト用固定APP_KEYは証拠・検証入力として保持する。

## 前回の検証状況（再開前の記録）

元の作業ツリーとは別の一時ディレクトリに、Git管理候補だけをコピーして検証した。vendorなしのテストは依存不足で失敗した。Macの`/var`と`/private/var`の別名による診断パス差を避け、実パスから再実行する。

期待値1件のスタック引数に元の配置先が埋め込まれていたため、`{trace_file}`を実行対象パスの先頭15バイト＋省略記号で展開する形にした。診断の削除・出力trimは行わず、既存Mac配置での期待文字列は変わらない。

Podman接続が拒否されるため、新規composer installとコンテナビルドの再現確認は未実施。既存の固定vendorをコピーした別配置テストは、依存取得の再現性とは区別する。GitHub CIも未実行。過去の1169件公式比較・実HTTP成功の証拠は維持する。

別配置の確認で、固定アプリは事前生成されたpackage manifestにも依存すると判明した。再構築手順とCIに公式PHPでの`artisan package:discover`を明記した。これは準備工程であり、Goによるmanifest生成成功とは扱わない。また`include_path_explicit.php`の親相対パスはチェックアウト名`php`を固定していたため、現在のルート名から算出するように変更した。期待出力は維持し、このfixtureの変更後の公式比較はCIで再確認する。

別配置のテストで独自CLIテストが未設定の`DB_DATABASE`に依存することも確認した。テストに`:memory:`を明示し、既存のローカル`.env`なしで検証するようにした。アプリの実装は変更していない。

CIの構成は[GitHub Actionsのコンテナ実行](https://docs.github.com/en/actions/how-tos/write-workflows/choose-where-workflows-run/run-jobs-in-a-container)と[setup-go](https://github.com/actions/setup-go)の公開仕様を参照した。コンテナはジョブ内のDocker CLIで明示的に起動し、比較終了後に削除する。

最終の[別配置検証](validation/release-relocated-test.log)は`make test vet build`終了0（Go全パッケージ、Python12件、vet、ビルド）。[監査記録](validation/release-preparation-review.log)に環境制限と確認範囲をまとめた。再構築のソース配置・テストは確認できたが、依存の新規取得・コンテナビルド・GitHub CI実行・公開は未完了である。

## 2026-10-21の再開後の検証

Git書き込みはサンドボックス外の`git add --dry-run`で確認した。GitHubの`ue555/gophp`への接続・push権限を確認し、開始時点の`main`はローカル・リモートとも`62ba3f4`、`v0.2.0`は双方に存在しなかった。

Git管理候補だけを別配置へコピーし、vendor・`.env`・既存manifestを持ち込まずに再構築した。[依存用コンテナビルド](validation/release-resume-dependency-build.log)、[固定lockからの新規composer install・platform要件確認・package:discover](validation/release-resume-fresh-install.log)が成功し、`composer.lock`は不変だった。イメージのビルドには既存のレイヤーキャッシュを利用している。

再構築で見つかった以下の問題を修正した。

- PCRE2開発ヘッダー／ランタイムが依存一覧から漏れていたため、CI・Goコンテナ・導入手順へ追加した。
- `stream_select`の`syscall.Select`戻り値と`FdSet`要素幅がOSによって異なるため、Darwin/Linux別の呼び出しとABIに応じたビット操作へ修正した。ビット境界と実パイプの待機を回帰テストに追加した。
- 独自CLIテストは公式PHPの準備工程が生成したmetadata cacheを読み込んで失敗したため、HTTP統合テストと同じく専用の一時ストレージで検証するようにした。公式PHPとGo製PHPの永続アプリキャッシュ共有は検証範囲外であり、キャッシュの相互互換性が修正されたという意味ではない。
- コンテナ検証のCLIヘルプ期待値を、現在のスクリプト引数対応に合わせた。

[元配置](validation/release-resume-local-test.log)と[新規依存を導入した別配置](validation/release-resume-relocated-test.log)で`make test vet build`が終了0。Go全パッケージ、Python12件、vet、ビルドを通過した。元配置のビルドにはサンドボックスによるGoモジュールキャッシュ書き込みの警告が出たが、コマンドは終了0だった。[Linux/arm64コンテナビルド](validation/release-resume-container-build.log)、[3設定の21fixtureとCLI 2項目](validation/release-resume-container-check.log)、[Linuxのselect回帰テスト](validation/release-resume-linux-select.log)も終了0。

リリース時はpushしたコミットSHAと[GitHub Actions](https://github.com/ue555/gophp/actions/workflows/ci.yml)の`headSha`を照合し、3版の全ジョブ成功後にそのSHAへ`v0.2.0`を付ける。各ジョブが固定vendorの新規導入、全fixtureの公式比較、実HTTP 2000回・8並列を検証する。公開状態はリモートタグと対象コミットのCI実行記録を正とし、ローカル成功だけをCI成功とは扱わない。

[再開後の全件公式比較](validation/release-resume-compare-all.log)はPHP 8.3.35・8.4.26・8.5.11の各版で1169一致・0差異・0エラー、全体終了0。ステージ時に見つかった成功fixture 6件の末尾空行を除去し、その6件も[3版で再比較](validation/release-resume-whitespace-fixtures.log)して全件一致した。生ログの出力バイト列を変更しないため、`.gitattributes`で`docs/validation/*.log`だけ空白検査・改行変換の対象外にしている。その他のソース・文書は通常の差分検査を維持する。

最初のコミット`1e47241`を`main`へ通常pushした。[初回CI](https://github.com/ue555/gophp/actions/runs/39533730928)は3版ともGoテストを通過したが、Pythonテストで`compare-all.sh`のMac固定パスが原因となって失敗した。[失敗ログ](validation/release-ci-first-failures.log)を保存し、比較・確認スクリプトは自身の位置からプロジェクトルートを求める形に修正した。回帰テストは空白を含む別配置へスクリプトをコピーし、別ディレクトリから起動して作業ディレクトリ・版指定・終了コードを確認する。CIではテスト工程のログもartifact対象へ追加した。失敗コミットへのタグ付けは行っていない。

[修正後のPython12件・公式PHP3版の起動確認](validation/release-ci-path-fix-test.log)は終了0。別配置での[実HTTP検証](validation/release-resume-relocated-http.log)も、固定vendorの新規導入後、PHP 8.4設定でCLI・SQLite・入力検証・API Platform CRUD・OpenAPI・2000回連続／8並列要求を通過し、検証コマンドは終了0だった。ログの`server_exit_code: -15`は検証後のサーバー停止によるもので、検証コマンドの終了コードではない。
