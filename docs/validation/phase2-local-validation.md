# 基本文法・名前空間追加後の検証（2026-10-12）

ユーザー環境で `./containers/php-version/check.sh` は PHP 8.3.35・8.4.26・8.5.11 の3版とも PASS。`make compare-all` の `/tmp/php-compare-all.log` を確認した結果、各版に異なる507件の `MATCH` 行があり、集計はそれぞれ **507 matched・0 different・0 errors**（合計1,521一致）だった。正常系242件は各版とも終了コード0。GoテストとPythonの12テストも比較前に通過している。比較器は同一fixtureの stdout・stderr・終了コードを厳密に照合する。ログ中に `DIFF`、`ERROR`、`FAIL`、`make:` の失敗行はない。ユーザーが `make compare-all` の正常終了も報告した。`/tmp/php-compare-setup.log` の `podman machine start` エラーはVMが起動済みであることを示すもので、後続の比較には影響しなかった。

`make test vet build` は終了コード0。Go の `cmd/gophp` と `internal/php` のテスト、Python の12テスト、`go vet`、Go CLI のビルドが通過した。`tests/fixtures/expected.json` に収録した507件（着手前466件、新規41件）を Go の既定設定・8.3・8.4・8.5 の4設定で実行し、保存済み期待値との不一致は0。新規41件は正常34、実行時1、構文／コンパイル6である。着手前の466 PHPソースの SHA-256 と466期待値を作業開始時のスナップショットと比較し、変更は0だった。

このCodex実行環境からの先行試行では、`make compare-all` が Go / Python のローカル検査とビルドを通した後、3版とも公式 PHP の版確認で終了コード125となり、環境エラーとして停止した。stderr は `Cannot connect to Podman ... dial tcp 127.0.0.1:60028: connect: operation not permitted`。比較器は `ERROR` を報告し、DIFF に変換していない。この先行試行のfixture実比較は各0件だったが、上記のユーザー環境での実比較で各507件を確認できた。

`make container-build` は同じ Podman 接続拒否で失敗し、Go 製 PHP コンテナの再ビルド・実行は未実施。`research/laravel-api/resolve.sh` も最初のイメージビルドで同じエラー（125）により停止した。ホストに `php` と `composer` のコマンドはなく、Packagist の名前解決も失敗した。したがって Composer の依存解決、`composer.lock`、vendor 監査、3版のプラットフォーム要件確認も未実施。

`research/laravel-api/composer.json` の JSON 構文、`resolve.sh` の `sh -n`、`git diff --check` は通過した。Go のビルドは終了コード0だが、外部のモジュールキャッシュへ stat cache を書けないという警告を出す。生成物のビルドとテストは成功している。

追加41件の保存済み期待値は作成時点で公式 PHP の出力を採取したものではなかったが、今回すべて公式PHPの3版とバイト単位で一致した。依存調査は別途 `./research/laravel-api/resolve.sh` により Composer lock を生成し、vendorコードを調べる必要がある。変更前466件の期待値は変更していない。
