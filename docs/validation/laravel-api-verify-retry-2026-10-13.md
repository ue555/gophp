# Laravel / API Platform 再検証（2026-10-13）

前回の `/tmp/laravel-api-verify.log` では、PHP 8.3.35 の `package:discover`、`migrate`、`api-platform:install` は終了コード0だった。`api-platform:metadata:dump` も終了コード0だったが、出力は `No dump path configured` であり、ダンプは生成されていない。HTTP の API Platform `POST /api/books` は500で、応答の詳細は `Unable to generate an IRI for the item of type "App\\Models\\Book"` だった。8.4.26 と8.5.11のHTTP検証は未実施。

現在の `config/api-platform.php` の探索先は存在する `app_path('Models')` のみ。`verify.py` は `--path=/probe/metadata.dump` を指定し、ダンプファイルの存在と非空を確認する。これはHTTP成功の確認とは別である。HTTPの各応答、Laravelのstderrログ、PHP組み込みサーバーのログを出すため、500が続く場合はIRI生成時の元例外を調べる。

この環境での `python3 -B research/laravel-api/verify.py` は終了コード1。最初のPHP 8.3.35用 `podman run` が終了コード125で、Podmanソケットへの接続が `operation not permitted` により拒否された。PHP 8.3.35・8.4.26・8.5.11のダンプ生成とHTTP再検証はすべて未実施。実出力は [laravel-api-verify-retry-local.log](laravel-api-verify-retry-local.log) に保存した。Pythonテスト12件と `verify.py` の構文解析、`git diff --check` は成功。

Macのターミナルで次を実行し、`verify.py` の終了コードをパイプライン直後に保存する。ログは前回の証拠を上書きしない。

```zsh
cd /Users/kouji/dev/php
set -o pipefail
python3 research/laravel-api/verify.py 2>&1 | tee /tmp/laravel-api-verify-retry.log
verify_exit=$?
echo "verify.py 終了コード: $verify_exit"
```

終了コード0かつ3版の `PASS` が揃って初めてHTTP検証成功とする。500が続けば、同ログのHTTP応答とサーバーログにある元例外を確認する。
