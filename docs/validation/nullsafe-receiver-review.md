# nullsafeメソッド受け手の検証（2026-10-12）

主基準はPHP 8.4.26、追加対象は8.3.35と8.5.11。既存518件のPHPソースと期待値を変更せず、9件を追加して527件とした。比較器はstdout・stderrをバイト単位、終了コードを整数として比較する。

## 根拠と修正

提供された `/tmp/nullsafe-receiver-comparison.log` をそのまま [nullsafe-receiver-before.log](nullsafe-receiver-before.log) に保存した。修正前、`private.php` は3版とも公式PHPが空stdout・privateプロパティの致命エラー・終了255、Goは`fallback`・空stderr・終了0だった。`undefined.php` は3版とも一致し、警告なしの`fallback`・終了0だった。この結果により、未定義変数の警告が必要という以前のレビュー指摘を訂正する。

原因は `quietNullsafeChain` が `??` 等の静かなプロパティ読み取りを、後続メソッドの受け手にも適用していたこと。メソッドの受け手となるプロパティ・添字には通常の読み取りを使い、末端の存在確認は静かな読み取りを維持した。先行するnullsafe短絡状態を伝え、受け手・添字・引数を再評価しない。

追加fixture: `runtime/nullsafe_receiver_private.php`、`runtime/nullsafe_receiver_protected.php`、`runtime/nullsafe_receiver_error_order.php`、`runtime/nullsafe_receiver_missing_warning.php`、`success/nullsafe_receiver_undefined.php`、`success/nullsafe_receiver_inside.php`、`success/nullsafe_receiver_public.php`、`success/nullsafe_receiver_offset_once.php`、`success/nullsafe_receiver_prior_stop.php`。既存の `success/nullsafe_quiet_missing.php` と `success/nullsafe_quiet_short_circuit.php` も対象比較に含めた。

## 実行結果

| コマンド | 終了コード | 結果ログ |
| --- | ---: | --- |
| `podman ps -a --format '{{.Names}} {{.State}}'` | 125 | [container-check](nullsafe-receiver-container-check.log) |
| `make test` | 0 | [local-check](nullsafe-receiver-local-check.log) |
| `make vet` | 0 | [vet](nullsafe-receiver-vet.log) |
| `make build` | 0 | [build](nullsafe-receiver-build.log) |
| `git diff --check` | 0 | [diff-check](nullsafe-receiver-diff-check.log) |
| `./containers/php-version/compare-nullsafe-receiver.sh` | 2 | [targeted-comparison](nullsafe-receiver-targeted-comparison.log) |
| `make compare-all` | 2 | [all-comparison](nullsafe-receiver-all-comparison.log) |

Go fixtureテストは全527件を既存テストに含めて4つのGo対象版指定（default、8.3、8.4、8.5）で実行した。Pythonの12テストも成功した。新規fixtureの保存済み期待値はGoのローカル検証用であり、提供ログで確認済みのprivate・undefined以外は公式PHPで未確認である。`make build` は終了0でバイナリを生成したが、サンドボックス外のGoモジュールキャッシュへのstat cache書き込み警告が出た。

| PHP | 比較対象 | 一致 | 不一致 | 環境エラー | 未実施 |
| --- | ---: | ---: | ---: | ---: | ---: |
| 8.3.35 | 527 | 0 | 0 | 1（版確認） | 527 |
| 8.4.26 | 527 | 0 | 0 | 1（版確認） | 527 |
| 8.5.11 | 527 | 0 | 0 | 1（版確認） | 527 |

Podmanソケットへの接続が `operation not permitted` で拒否され、比較器は各版の版確認で停止した。従って追加11件の対象比較も各版0件実行である。これはPHPとの不一致件数ではない。作成済みコンテナの状態も取得できず、起動・削除・初期化は行っていない。公式PHPでの関連ケースの期待値確定と、全527件の再比較が残る。

Macのターミナルでは、VMが起動済みならそのまま、各コンテナの停止状態を確認して必要なものだけ起動する。`zsh` で以下を実行し、成功・失敗を問わず新しいログを保存する。

```sh
cd /Users/kouji/dev/php
podman info >/dev/null
for name in php-version-83 php-version-84 php-version-85; do
  if [ "$(podman inspect -f '{{.State.Status}}' "$name")" != running ]; then podman start "$name"; fi
done
./containers/php-version/check.sh
make build
set -o pipefail
./containers/php-version/compare-nullsafe-receiver.sh 2>&1 | tee docs/validation/nullsafe-receiver-targeted-mac.log
make compare-all 2>&1 | tee docs/validation/nullsafe-receiver-all-mac.log
```

各ログの `MATCH` / `DIFF` / `ERROR` と終了コードを確認し、必要なら暫定期待値を公式PHPの結果に基づき訂正する。修正前ログと今回の失敗ログは上書きしない。
