# 第1版の全件実比較と完成確認（2026-10-12）

限定した対応範囲を持つCLI処理系の第1版の完成条件を満たしました。
出典はユーザーが会話に貼った公式PHP実比較ログです。この文書はその要約と抜粋であり、
エージェントが公式PHPを実行して生成したログではありません。同じログの重複貼付は1回の実行として集計しています。

## 実行と結果

fixtureの補助関数名を修正した後、最初の再試行はPHPコンテナが停止していたため失敗しました。
そのログの `podman info` は `running: 0`・`stopped: 3` を示し、PHP実行は
`can only create exec sessions on running containers: container state improper` で停止しています。
続いてユーザーが以下を実行し、3コンテナの起動・版確認・全件比較に成功しました。

```sh
cd /Users/kouji/dev/php &&
podman start php-version-83 php-version-84 php-version-85 &&
./containers/php-version/check.sh &&
make compare-all
```

版確認の抜粋：

```text
PASS: PHP 8.3.35, shared fixture executed
PASS: PHP 8.4.26, shared fixture executed
PASS: PHP 8.5.11, shared fixture executed
```

各版の集計行はすべて次のとおりです。

```text
466 matched, 0 different, 0 errors; 466 fixtures
```

| 公式PHP | Go対象設定 | fixture数 | 一致 | 不一致 | 検証エラー |
| --- | --- | ---: | ---: | ---: | ---: |
| 8.3.35 | 8.3 | 466 | 466 | 0 | 0 |
| 8.4.26 | 8.4 | 466 | 466 | 0 | 0 |
| 8.5.11 | 8.5 | 466 | 466 | 0 | 0 |
| 合計 | | 1,398 | 1,398 | 0 | 0 |

同じGo実装・共通のPHP fixtureを使い、stdout・stderrをバイト単位、終了コードを厳密比較した結果です。
全件のMATCH行があり、未実行・除外はありません。正常系208件は各版とも両処理系で終了0です。
修正した `success/isset_short_circuit_order.php` も3版すべてで `MATCH ... (exit 0)` を確認しています。

同じ `make compare-all` のログでGoテスト、Pythonテスト9件、ビルドの成功も確認しました。
このコマンドはgo vetを実行しません。go vetの成功は先の
[ローカル検証ログ](v1-followup-local-validation.txt) に記録しています。

## 完成条件

| 項目 | 確認結果 |
| --- | --- |
| 9：配列の値渡しforeach | 実装・単体テスト・公式3版の実比較が完了 |
| 10：isset / empty / unset、4種のキャスト | 実装・単体テスト・公式3版の実比較が完了 |
| 11：指定11標準関数、COUNT_NORMAL / COUNT_RECURSIVE | 実装・単体テスト・公式3版の実比較が完了 |
| 12：統合例、CLI、文書、全件検証 | 統合fixture・CLI検証・文書更新と公式3版の全件実比較が完了 |
| 既存fixtureの維持 | 既存337件のPHPソースと期待値を維持。ローカル検証でSHA-256と期待値を照合済み |
| 追加fixture | 129件を追加し、計466件（正常208・実行時198・構文／コンパイル60）を全件実比較 |

新規129件の保存済み期待値は公式資料・対象版PHPソースを根拠に手動で作成し、
Go CLIとの照合後、今回のユーザー環境での公式実比較で全件を確認しました。
完成条件に関する未実施項目はありません。対応範囲と未対応機能は [互換性文書](../compatibility.md) を参照してください。
全件一致は収録したfixtureについての結果であり、PHP全機能との互換性を示すものではありません。

## 履歴との区別

- [全件成功前の実比較と追補修正](v1-user-comparison-followup.md)：補助関数名の衝突、接続障害、比較器の修正の記録。
- [修正後のローカル検証](v1-followup-local-validation.txt)：Go・Python・vet・ビルドと保存済み期待値の検査。
- [エージェント環境での公式比較試行](v1-followup-php-comparison.txt)：接続権限の制約で版確認時に停止し、各版のfixture実比較は0件。

過去の失敗や未実行を一致へ読み替えず、後続のユーザー環境での全件実行を完成の根拠としています。
