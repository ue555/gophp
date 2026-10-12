# ユーザー提供の第1版実比較ログと追補修正

この文書は全件成功前の履歴です。最新の結果は [全件実比較と第1版の完成確認](v1-user-comparison-complete.md) を参照してください。

出典は会話に貼られた `make compare-all` の実行ログです。以下はその要約であり、エージェントが実行したログではありません。
作業ディレクトリは `/Users/kouji/dev/php`、PHPは各版の `containers/php-version/bin/php83`・`php84`・`php85`、Goは対応する8.3・8.4・8.5設定です。

## 提供ログの結果（修正前）

GoテストとPythonテスト4件、ビルドは成功しています。

| PHP | fixture数 | MATCH | DIFF | ERROR |
| --- | ---: | ---: | ---: | ---: |
| 8.3.35 | 466 | 465 | 0 | 1 |
| 8.4.26 | 466 | 465 | 0 | 1 |
| 8.5.11 | 466 | 179 | 8 | 279 |

8.3・8.4でERRORになったのは `success/isset_short_circuit_order.php` です。
補助関数 `key()` がPHPの既存標準関数に衝突し、実行前に終了255となりました。
8.3では `Cannot redeclare key()`、8.4では `Cannot redeclare function key()`、行番号は2です。
Goはstdout `:10\n:341\n:560\n`、stderr空、終了0でした。
正常系が両方で終了0であるという比較器の検査により、検証エラーになっています。

8.5では `runtime/function_argument_count.php` まで179件がMATCHした後、タイムアウトと接続障害が発生しました。
以下のDIFF 8件はすべてPHP側の終了125で、stderrは `Cannot connect to Podman.` と
`Error: unable to connect to Podman socket: failed to connect: ssh: handshake failed: EOF` です。

- `runtime/object_scalar_write.php`
- `runtime/strict_evaluates_operands.php`
- `runtime/type_error.php`
- `runtime/type_error_before_warning.php`
- `runtime/undefined.php`
- `runtime/unset_undefined_array.php`
- `runtime/warning_before_type_error.php`
- `runtime/while_condition_warning.php`

これらはPHP言語動作を比較できた結果ではありません。ERROR 279件もタイムアウトまたは接続失敗を伴う正常系の非ゼロ終了です。
`isset_short_circuit_order.php` も8.5ではタイムアウトし、短絡評価を確認できていません。
上表は元ログの集計をそのまま記載し、修正後の比較器による件数へ読み替えていません。

## 修正

1. fixtureの補助関数と全呼び出しを `gophp_isset_order_key()` に改名しました。PHPソースの検証内容、期待出力、終了0の条件、fixture数は維持しています。Go処理系と保存済み期待値は変更していません。
2. 比較器は終了125とPodman固有の接続診断を組み合わせて検出し、DIFFではなくERRORで報告します。stdout・stderr・終了コードは報告に残します。終了125だけなら通常どおり比較します。
3. PHP側のタイムアウト等では版確認を試します。版確認に成功すれば続行し、失敗または明確な接続障害の場合はその版を停止します。残りを `not run` と明記し、終了コード2を返します。
4. Pythonテストを5件追加し、接続断の分類、残りの未実行、タイムアウト後の継続／停止、類似出力・終了125の厳密比較を検査します。疑似プロセス結果を使う比較器のテストであり、PHPの実比較ではありません。

既存337件のfixtureと期待値、追加後の全466件という件数は維持しています。
今回改名したのは追加129件のうち1件です。PHPの対象外標準関数名の予約はGo側の追加機能に含めていません。

## 修正後

[ローカル検証ログ](v1-followup-local-validation.txt) と [公式PHP比較の試行ログ](v1-followup-php-comparison.txt) を分けて保存しています。
エージェント環境はPodman接続が `operation not permitted` で拒否されるため、修正後の公式実比較は実施できていません。
新しいスクリプトがユーザー環境のPodman自体を復旧したという意味ではありません。

この時点では「実装・ローカル検証完了、修正後の実比較待ち」で、完成条件は未達でした。
その後、ユーザー環境で停止していたコンテナを起動し、466件×3版すべての成功を確認しました。
第1版の完成確認は [後続の記録](v1-user-comparison-complete.md) に記載しています。
