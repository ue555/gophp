# 比較用PHPコンテナ

PHP 8.3・8.4・8.5のCLIをpodman-composeで常駐させ、ホスト側のGo実装と比較します。Docker Official ImageのパッチバージョンとDebian bookwormを指定しています。

| サービス | PHP | コンテナ | 実行用ラッパー |
| --- | --- | --- | --- |
| php83 | 8.3.35 | php-version-83 | bin/php83 |
| php84 | 8.4.26 | php-version-84 | bin/php84 |
| php85 | 8.5.11 | php-version-85 | bin/php85 |

タグの確認元: [Docker Official ImagesのPHP定義](https://github.com/docker-library/official-images/blob/master/library/php)。パッチタグを更新する場合はcompose.yaml、check.sh、compare-all.shの期待バージョンも合わせて変更します。

## 配置と起動

構成は `/Users/kouji/dev/php/containers/php-version` にあります。このディレクトリで起動します。

```sh
cd /Users/kouji/dev/php/containers/php-version
podman-compose -f compose.yaml up -d
./check.sh
```

Podmanの仮想マシンが停止している場合は、先に `podman machine start` を実行します。

VMが起動していても既存コンテナが停止していると、`can only create exec sessions on running containers` で失敗します。
作成済みの3コンテナを再開して比較する場合は、プロジェクトルートで次を実行します。

```sh
cd /Users/kouji/dev/php &&
podman start php-version-83 php-version-84 php-version-85 &&
./containers/php-version/check.sh &&
make compare-all
```

ホストの `/Users/kouji/dev/php` を、コンテナ内の同じ絶対パスへ読み取り専用でマウントしています。既存の構成には `/Users/kouji/docker/php-version` のマウントも残しています。Podman machineがこれらのディレクトリを共有できる必要があります。Webポートは使用しません。各実行用ラッパーはTTYを使わず、PHPの出力と終了コードを呼び出し元へ渡します。

## PHPを実行する

```sh
cd /Users/kouji/dev/php
./containers/php-version/bin/php83 -n -r 'echo PHP_VERSION, PHP_EOL;'
./containers/php-version/bin/php84 -n -r 'echo PHP_VERSION, PHP_EOL;'
./containers/php-version/bin/php85 -n -r 'echo PHP_VERSION, PHP_EOL;'
```

PHP CLIの実行や比較を呼ぶ際の作業ディレクトリは、マウントした2つのディレクトリのいずれか、またはその配下にしてください。

## Go実装と比較する

主基準はPHP 8.4.26、追加確認は8.3.35・8.5.11です。
同じGo実行ファイルと同じ全466fixture（項目9〜12追加前の337件と追加129件）を各バージョンへ渡します。
比較スクリプトはGoに `--php-version=8.3`、`8.4`、`8.5` をそれぞれ指定します。
これにより8.5ではNaNのif条件・論理演算の追加警告を出し、8.3・8.4では従来どおり警告を出しません。
Go側はこのオプションと、最後の引数でPHPファイルのパスを受け取るCLIを想定しています。
Goのテストとビルドを含めて実行する場合:

```sh
cd /Users/kouji/dev/php
make compare-all
```

ビルド済みなら、次のコマンドでも実行できます。

```sh
cd /Users/kouji/dev/php/containers/php-version
./compare-all.sh tests/fixtures ./bin/gophp
```

このコマンドは3バージョンすべてについて、標準出力・標準エラー・終了コードを比較します。終了コードは0が全件一致、1が不一致、2が環境や検証のエラーです。相対パスは `/Users/kouji/dev/php` を基準に解釈します。

PHPバージョン間の差は、主基準へのGo実装の不適合と区別します。コンテナに接続できない場合は検証エラーとし、該当バージョンの実比較を未実施として記録します。各バージョンの結果と差異の分類方針は[互換性文書](../../docs/compatibility.md)にあります。

比較中のタイムアウトではPHPの版確認を再試行します。版確認も失敗した場合、または公式PHP側で終了125とPodman固有の接続エラー・コンテナ停止の診断（stderr先頭の `Error: can only create exec sessions on running containers`）が返った場合は、その版をERRORで停止し、残りを `not run` と表示します。未実行を一致扱いにはしません。終了125だけならPHPプログラム自身の終了コードとして通常どおり比較し、stdout・stderr・終了コードの厳密比較を維持します。

`ssh: handshake failed: EOF` や連続タイムアウトが出た場合は、接続状態を確認してから再比較してください。

```sh
podman system connection list
podman machine list
podman info
./containers/php-version/check.sh
make compare-all
```

上記はプロジェクトルートで実行します。VMが停止していれば `podman machine start` で起動します。
VMが稼働中でも `podman info` が失敗する場合はPodman側の接続を復旧する必要があります。
このログだけでVM停止や障害の原因までは断定できません。タイムアウトの延長では接続障害は解消しません。

1つのバージョンだけ比較する例:

```sh
cd /Users/kouji/dev/php
python3 skills/php-go-interpreter/scripts/compare_php.py \
  --php ./containers/php-version/bin/php83 \
  --php-version 8.3.35 \
  --fixtures tests/fixtures \
  --timeout 15 \
  --interpreter ./bin/gophp --php-version=8.3 '{file}'
```

## 停止

```sh
cd /Users/kouji/dev/php/containers/php-version
podman-compose -f compose.yaml down
```
