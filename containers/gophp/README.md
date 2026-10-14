# Go製PHPインタプリタのコンテナ

この構成はGoで作成したgophp自体をコンテナで実行します。対象PHP設定は既定8.4で、
同じイメージに `--php-version=8.3` / `8.4` / `8.5` を渡せます。
対応する言語機能は [互換性文書](../../docs/compatibility.md) の第1版の範囲です。

## ビルドと実行

macOSのターミナルで、プロジェクトルートから実行します。

```sh
cd /Users/kouji/dev/php
make container-build
make container-check
make container-run
```

Podman machineが停止している場合は、先に `podman machine start` を実行してください。
初回ビルドではGoのベースイメージを取得します。ホスト側にGoは不要です。
`container-check` はホストのPython 3を使用します。
ビルドはコンテナ内でGoテストを実行し、成功した場合に実行ファイルを作ります。
Goのソースを変更したら `make container-build` でイメージを更新してください。

イメージ名は `localhost/gophp:v1`。実行ごとにコンテナを作成し、終了後に削除します。
`make container-run` は `examples/report.php` を実行し、次を出力します。

```text
a:tea:12:3
c:水:3:3
2:2:1:0:integer
```

ファイルや対象PHP設定を選ぶ場合:

```sh
make container-run GOPHP_FILE=tests/fixtures/success/hello.php GOPHP_TARGET=8.5
./containers/gophp/run.sh --php-version=8.3 examples/report.php
./containers/gophp/run.sh --help
```

`run.sh` の相対ファイルパスは、呼び出したディレクトリによらずプロジェクトルートを基準にします。
絶対パスはプロジェクト内のものを渡してください。ホストとコンテナに同じ絶対パスをマウントし、
診断のファイル名を維持します。TTYを使わず、stdout・stderr・終了コードを呼び出し元へ渡します。
イメージが未作成の場合は自動取得せず失敗するため、先にビルドしてください。

## podman-composeで実行

同じイメージをcompose経由でもビルド・実行できます。

```sh
cd /Users/kouji/dev/php/containers/gophp
podman-compose -f compose.yaml build
podman-compose -f compose.yaml run --rm -T gophp
podman-compose -f compose.yaml run --rm -T gophp --php-version=8.5 examples/report.php
```

composeではプロジェクトを `/workspace` にマウントします。診断の絶対パスは `/workspace/...` です。
CLIのため実行終了とともに停止します。常駐用の `up -d` は不要です。

## 動作確認と公式PHP比較

```sh
cd /Users/kouji/dev/php
make container-check
python3 -B containers/gophp/check.py --all
```

通常のチェックは7fixtureを3設定で実行し、さらにヘルプ・引数不足の終了コードを検査します。
日本語の出力、警告、構文／致命エラー、版によって異なる診断・構文を確認します。
`--all` は`tests/fixtures/expected.json`に登録された全fixture（現在588件）を3設定で実行します。これは保存済み期待値との照合であり、[ホストの公式PHP全件比較](../../docs/validation/item7-588-user-comparison-after-fix.log)とは別です。
エラーを含めてstdout・stderrをバイト単位で比較し、終了コードも検査します。
終了コードはチェック成功0、不一致1、コンテナ接続失敗等の検証エラー2です。

公式PHPとコンテナ版gophpを実比較する場合は、比較用の3コンテナを起動してから次を実行します。
このコマンドはfixtureごとにgophpコンテナを起動するため、ホストのバイナリとの比較より時間がかかります。

```sh
cd /Users/kouji/dev/php
podman start php-version-83 php-version-84 php-version-85
./containers/php-version/check.sh
./containers/php-version/compare-all.sh tests/fixtures ./containers/gophp/run.sh
```

## 構成と検証状況

[Containerfile](Containerfile) は、Goの公式イメージ `golang:1.27.2-alpine3.24` で
`CGO_ENABLED=0` のLinuxバイナリをビルドし、scratchイメージへコピーします。
実行イメージにはgophpの実行ファイルだけを含め、UID/GID 65532で実行します。
プロジェクトは読み取り専用、ネットワークなしで実行します。Webポートは使用しません。
Apple SiliconではLinux/arm64、Intelの64ビット環境ではLinux/amd64のビルドを想定しています。
ベースイメージのタグ・対応アーキテクチャは [Go公式イメージ定義](https://github.com/docker-library/official-images/blob/master/library/golang)、
実行オプションは [Podman公式資料](https://docs.podman.io/en/latest/markdown/podman-run.1.html) を参照しています。

必要に応じてビルド・実行の両方に同じ `GOPHP_IMAGE` 環境変数を渡せます。
別のGoイメージを使う場合は `./containers/gophp/build.sh --build-arg GO_IMAGE=...` を指定します。

2026-10-12、このCodexセッションからのPodman接続は `operation not permitted` で拒否されました。
コンテナの実ビルド・実行・公式PHPとのコンテナ版実比較は未実施です。
過去のホスト版の全1,398件一致とは区別します。ターミナルで上記のビルド・チェックを実行してください。

同セッションで確認できた内容は、Goテスト、Pythonテスト12件、`go vet`、シェル構文、
`podman-compose config` の読み込み、Linux/arm64・Linux/amd64への静的バイナリのビルドです。
新しいチェック処理自体はホストでビルドしたgophpを使い、21fixture実行とCLI確認2件が成功しました。
実行ラッパーの引数・stdout・stderr・終了コードの保持はPodmanを模擬するテストで確認しています。
これらはコンテナ内での動作確認を示すものではありません。
