# nullsafe演算子のレビュー修正と検証（2026-10-12）

主基準はPHP 8.4.26、追加対象はPHP 8.3.35と8.5.11。今回の11 fixtureを加え、総数は518件。修正前の507 PHPソースと507期待値を作業前スナップショットと比較し、変更は0だった。追加11件の保存済み期待値は公式PHPで採取した値ではなく、仕様・PHPソースとGoローカル実行に基づく暫定値である。

| 指摘 | 原因 | 修正 |
| --- | --- | --- |
| `?->` 後の添字が実行される | `quietAccess` はnullsafeチェーンの停止状態を受け取らず、添字式を評価した | `quietNullsafeChain` が停止状態を外側のプロパティ・添字・メソッドへ伝える。null受け手ではキーと引数を評価しない |
| 静かなアクセスで不要な警告 | `quietAccess` がnullsafe式を通常の`eval`へ渡し、`readProperty`が未定義警告を出した | 非null受け手のプロパティを静かな読み取り経路で扱う。通常の読み取りは従来の警告を維持する |
| 未選択分岐の`unset`を実行前に拒否できない | `validateUnset` が対象のnullsafeチェーンを調べなかった | `hasNullsafe`で直接・後続プロパティ・添字を含む対象を事前に拒否する |

追加fixtureは以下のとおり。

- 正常系：`success/nullsafe_quiet_short_circuit.php`、`success/nullsafe_quiet_missing.php`、`success/nullsafe_quiet_nested.php`、`success/nullsafe_quiet_side_effects.php`、`success/nullsafe_legal_unset.php`
- 警告と副作用：`runtime/nullsafe_normal_read_warning.php`、`runtime/nullsafe_key_warning.php`、`runtime/nullsafe_argument_warning.php`
- 実行前の致命エラー：`parser/unset_nullsafe_direct.php`、`parser/unset_nullsafe_property_chain.php`、`parser/unset_nullsafe_offset_chain.php`

`make test vet build` は終了コード0。Go・Python（12件）テスト、`go vet`、ビルドが通過し、518件の保存済み期待値テストはGoの既定・8.3・8.4・8.5設定で通過した。これは公式PHPとの差分比較とは区別する。

このCodex実行環境では`podman ps -a`と`check.sh`が終了コード125となった。Podmanソケットへの接続が`operation not permitted`で拒否されるため、既存コンテナの状態は取得できず、起動操作もしていない。追加11件の専用比較と全518件の`make compare-all`は、いずれも各版のPHP版確認で同じ環境エラーとなり、比較器の終了コードは2。fixtureの実行は各版0件だった。

| PHP | 専用11件：一致／不一致／版確認エラー／未実施 | 全518件：一致／不一致／版確認エラー／未実施 |
| --- | --- | --- |
| 8.3.35 | 0／0／1／11 | 0／0／1／518 |
| 8.4.26 | 0／0／1／11 | 0／0／1／518 |
| 8.5.11 | 0／0／1／11 | 0／0／1／518 |

版確認エラーはfixtureの不一致ではない。比較前のユーザー環境で各版507件が一致した結果は、今回の11件追加後の比較結果ではない。今回の正常系・警告・診断の公式PHPとの厳密一致は未確認である。

実行ログは[ローカル検査](nullsafe-local-check.log)、[コンテナ確認](nullsafe-container-check.log)、[追加11件の比較試行](nullsafe-targeted-comparison.log)、[全518件の比較試行](nullsafe-all-comparison.log)に保存した。各ログに実行コマンドと終了コードを記録している。

Macのターミナルでコンテナの稼働を確認し、停止している場合のみ既存コンテナを起動してから、次を実行する。VMが起動済みなら`podman machine start`は不要。

```sh
cd /Users/kouji/dev/php
podman ps -a --format '{{.Names}} {{.State}}'
podman start php-version-83 php-version-84 php-version-85  # 停止している場合だけ
./containers/php-version/check.sh
make build
./containers/php-version/compare-nullsafe-review.sh > docs/validation/nullsafe-targeted-comparison.log 2>&1
nullsafe_target_status=$?
printf '\nexit code: %s\n' "$nullsafe_target_status" >> docs/validation/nullsafe-targeted-comparison.log
make compare-all > docs/validation/nullsafe-all-comparison.log 2>&1
nullsafe_all_status=$?
printf '\nexit code: %s\n' "$nullsafe_all_status" >> docs/validation/nullsafe-all-comparison.log
printf 'targeted=%s all=%s\n' "$nullsafe_target_status" "$nullsafe_all_status"
```

比較スクリプトはstdout・stderr・終了コードをバイト単位で照合し、trimや警告の削除は行わない。
