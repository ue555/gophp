# Issue #1: 未改変Laravel 12の実装・検証記録

対象は [Issue #1](https://github.com/ue555/gophp/issues/1) のLaravel起動条件。
既存固定Laravel/API Platformアプリの成功とは分けて検証した。
commit・push・タグ変更・Issue投稿/closeは実施していない。

## 開始時の事実と再現環境

開始HEADは `03325d794dd501f31d85d75eeea75286227aa4de`、ローカルのfixtureと期待値は
**1173件**。Issueが報告した未コミットのinterfaceContracts修正、追加2 fixture、
`~/dev/laravel-trial`、以下の3資料はこのMacには存在しなかった。

- `fresh-laravel-trial-2026-10-21.md`
- `fresh-laravel-12.70-missing-builtins-2026-10-21.txt`
- `host-php85-arch-1175-compare-after-interface-fix-wsl.log`

[Issue本文の取得結果](fresh-laravel-issue-1-report.json)はコメント0件、OPEN。
Issueの「WSL/Archで1175件の8.5比較成功」は外部報告であり、今回の成功証拠には使っていない。
ローカルの旧interfaceContractsで3段継承の逆向き検査を再現してから修正した。
既存差分や履歴を巻き戻す操作はしていない。

`research/fresh-laravel/app` に公式スケルトン **laravel/laravel v12.12.2** を用意し、
Framework **v12.70.0**、**112パッケージ（開発依存35）**を固定した。
Issueの元lockがないため同一依存構成とは断定しない。最初のPHP 8.4での依存解決は
Symfony 8.1系6パッケージを含み8.3非対応だったため、PHP 8.3で同じFramework指定を
解決し直し、該当パッケージをSymfony 7.4系とした。既存2アプリの依存には触れていない。

- スケルトンcommit: `945f4e5a9fd3695dc0ee512f497c650fb82cfbb8`
- lock SHA-256: `352299ca4ff6eb707455e612b499bb3606cdf0fafbf1b8eff7650e19e4c1e545`
- 最終Goバイナリ SHA-256: `04444cd35f905a4bb9d10558af4010c5513fc1c86a0ea570e24b1664d3623fa8`
- 現行fixture: **1193件**（正常764・実行時295・構文134）。新規20件、既存1173期待値は不変。

公式PHPのComposerで `install --no-scripts` を実施した独立再構築は
[版別コマンド・拡張要件・パッケージ数](fresh-laravel-rebuild-qom8n2tg/)に保存。
既存vendorと生成キャッシュのコピーによる再構築ではない。
[公式スケルトン54ファイルと保護対象13,349ファイルの照合](fresh-laravel-audit-1oz5423u/)
および[再導入vendorの8,561ファイル照合](fresh-laravel-vendor-integrity-20261021.json)は差異0。
この依存準備をGo実行の成功には数えていない。

## 実際に到達して修正した停止経路

| 経路 | 修正と回帰の焦点 |
|---|---|
| EventDispatcherのinterface継承 | 子孫の互換な再宣言で祖先契約を置き換え、独立した兄弟契約を保持。正常・不適合ケースと既存diamondを比較 |
| 開発用autoload / phpdotenv | `substr_count`の境界・版別診断、`rsort`の参照更新 |
| Console | `krsort`のキー比較、各フラグ、安定順序、参照/名前付き引数、一度だけの添字評価、失敗時の入力維持 |
| 正規表現 | preg系を実PCRE2結果へ接続。コンパイル/UTFエラーと成功後リセット、`preg_last_error[_msg]`、ライブラリから取得したPCREバージョン/JIT情報 |
| Blade | raw lexemeを保持する専用tokenizer、版別トークンID、空白・コメント・タグ・行番号・文字列補間・heredoc・TOKEN_PARSE。`substr_replace`、代替if構文、`filemtime` |
| 例外の継承 | `ErrorException`派生クラスの6引数コンストラクタ |
| セッション | 暗号学的`random_int`、DateIntervalのISO期間、CarbonからのDateTime add/sub |
| Cookie | HMAC、定時間文字列比較、AES-CBC/GCM encrypt/decrypt、tag参照更新、Cookie URL復号、gmdateのGMT書式 |
| JSON | Unicodeエスケープ、ネストしたオブジェクト内のHTML文字、encode成功/失敗後のエラー状態 |

関数名だけのスタブ、PHPへのアプリ処理転送、Laravel/vendor書き換え、開発依存除去、
middleware/Blade無効化は行っていない。
新規期待値は公式3版の観測から採取し、`fresh-laravel-official-observations-*` に
stdout・stderr・実際の終了コードを残した。既存期待値は変更していない。

## 新規Laravelの実TCP比較

[最終3版HTTPログ](fresh-laravel-final-three-http-20261021-213359-wk5fmxvm.log)、
[全HTTP/CLI/DB/Blade成果物](fresh-laravel-tcp-ljc9yvoj/)の終了コードは0。
公式PHPとGoのそれぞれに空の一時SQLite・storage・bootstrap/cacheを割り当て、
同じ要求順で実行した。各runtime **111要求**、6runtime合計666要求。

| 条件 | 8.3.35 / Go 8.3 | 8.4.26 / Go 8.4 | 8.5.11 / Go 8.5 |
|---|---|---|---|
| 開発用vendor/autoload | 成功 | 成功 | 成功 |
| .envの設定反映・変数展開 | 一致 | 一致 | 一致 |
| artisan list / about / route:list | 終了0 | 終了0 | 終了0 |
| 空SQLite標準migration / 再実行 | 終了0・状態維持 | 終了0・状態維持 | 終了0・状態維持 |
| 実TCP `/up` / `/` と標準welcome HTML | 一致 | 一致 | 一致 |
| Blade冷／温キャッシュとトークン列 | 一致 | 一致 | 一致 |
| 100連続・8並列、Cookie再送・DBセッション | 成功 | 成功 | 成功 |
| CLI stderr・サーバー／アプリ診断 | なし | なし | なし |

JSON/HTMLをtrimせず、警告も除去しない。動的値の範囲と個別検査は
[検証README](../../research/fresh-laravel/README.md#実tcpclisqlitebladeの比較)に記載。
healthのミリ秒・配置先パスだけを本文/トークン比較で正規化する。
Cookie/日時/セッションは形式・属性・整合性・独立性・再送時維持を検査する。
Composerの有無による差を隠さず、実行時PATHを揃えた結果、`about`の表示も一致した。

PCREエラー状態のリクエスト分離はLaravelとは別の実TCP回帰テストでも確認した。
[100連続・8並列×3設定](fresh-laravel-preg-tcp-regression-20261021-212327-er9iteak.log)。
手動利用手順の `key:generate` と `package:discover` も
[破棄用.envで実行](fresh-laravel-manual-commands-20261021-213850-bdun7hil.log)して終了0。

## 全回帰の集計

| 対象 | 結果・証拠 | 実際の終了コード |
|---|---|---|
| 公式8.3.35 / Go 8.3 全fixture | [1193一致・差異0・エラー0](fresh-laravel-final-fixtures-8.3.35-20261021-213358-m73ezime.log) | 0 |
| 公式8.4.26 / Go 8.4 全fixture | [1193一致・差異0・エラー0](fresh-laravel-final-fixtures-8.4.26-20261021-213653-w3quop4q.log) | 0 |
| 公式8.5.11 / Go 8.5 全fixture | [1193一致・差異0・エラー0](fresh-laravel-final-fixtures-8.5.11-20261021-213942-u0wivxuq.log) | 0 |
| タスクAPI 公式/Go 3組 | [各191要求・DB比較一致](task-api-20261021-213736-trqw46y4/run.log) / [親コマンド](fresh-laravel-final-task-api-three-20261021-213736-zsco9wyp.log) | 0 |
| 固定Laravel/API Platform Go 3設定 | [CRUD・OpenAPI・100連続・8並列](fresh-laravel-final-fixed-laravel-three-20261021-213737-xfse7cts.log) | 0 |

全fixture比較の[版別実コマンド・終了コード・ソースmanifest](fresh-laravel-fixture-matrix-dlxuo0h8/)
をHTTP証拠とは分けて保存した。最終回帰は同一の上記Goバイナリで実行した。
[最終証拠監査](fresh-laravel-results-audit-jlqegqfp/)でも、全fixtureの個別MATCH、
期待値件数、HTTP/既存アプリの成功、現在のGoソースとバイナリhashを照合して終了0。
監査ツール初回はsymlinkを通常ファイルとしてhash比較して失敗したが、
ソースmanifestと同じlink文字列のhashへ修正した。失敗ログも保持した。
サーバー終了-15は検証成功後に自分のGoサーバーをSIGTERMで片付けた実測値で、
検証コマンドの終了0とは区別している。

- [make test](fresh-laravel-final-make-test-20261021-213400-4s_b32nl.log): 終了0（Go全体・Python 15テスト）。
- [make vet](fresh-laravel-final-make-vet-20261021-213402-gmtksr_c.log): 終了0。
- [make build](fresh-laravel-final-make-build-20261021-213338-bniiyb1g.log): 終了0。
- [git diff --check](fresh-laravel-final-diff-check-20261021-213851-cyowumxy.log): 終了0。

## 失敗記録と検証上の修正

初回のautoload/Console/HTTP停止、JSON・tokenizer・参照診断の不一致は
`fresh-laravel-*` の別名ログとして保持している。
HTTPの初回500ではSymfonyエラー表示側の未対応が元の例外を隠したため、
アプリのlaravel.logも採取して本来の停止箇所を追った。

- 8.4固定Composer依存で8.3非対応となった準備結果も保存し、lockの変更理由を明示した。
- 新規期待値採取ツールの既定値を8.3からリポジトリ標準の8.4へ修正。版別公式値自体は不変。
- `json_encode`の循環配列回帰案は、既存の「循環配列参照自体が未対応」の制限に到達した。
  今回のエラー状態回帰は対応済み入力型・INF・resourceで検証し、循環参照を対応済みとはしていない。
- CookieのMax-Ageは期限設定から応答出力までに1秒進む場合があるため、Symfonyの算出式と
  HTTP Date/Expiresの整合性に基づく検査に修正。初回の固定7200秒判定失敗を保存した。
- Podmanの転送ポートがPHP起動より先に応答して接続リセットとなった1回を保存。
  PHP起動ログ確認を追加し、VM再起動・権限変更・既存サービス停止は行っていない。
- 参照診断の修正前に開始した全比較は自分のプロセスだけを中断した。
  [中断ログ](fresh-laravel-fixture-all-three-20261021-213036-kxxfl2b5.log)は子プロセス終了-15、
  recorder終了143で、成功として扱わない。修正後は版別に全件を再実行した。
- 最終診断修正前後にまたがったタスクAPI検証は最終証拠に採用せず、バイナリ固定後に取り直した。
- 最終HTTP検証後のスクリプト変更はDockerのUID/GID引数と版チェック追加、監査・手動準備
  ツールと文書のみ。Go/fixtureは変更しておらず、保存したaboutのパッチ版も最終監査で確認した。

## 残る一般機能の制限

本件の完成条件は新規Laravelの上記経路であり、PHP全仕様・Laravel全機能の完成ではない。
開発依存のautoloadは検証したが、PHPUnit等の全コマンド実行を網羅してはいない。

- tokenizerは標準Blade実入力と回帰範囲を照合。`TOKEN_PARSE`はGoの対応済みPHP構文に従う。
  全PHP字句状態、`__halt_compiler`以降、全予約語文脈や全構文エラー診断は未網羅。
- PCRE情報は実際にリンクしたPCRE2由来。PHP各イメージとライブラリのpatch版・JIT設定が
  同一とは限らない。PCREの全オプション・INI上限・警告境界は未網羅。
- 暗号関数はAES 128/192/256 CBC/GCMと今回のHMACアルゴリズムが対象。
  OpenSSL全暗号、全エラーキュー、短いGCM tag等の全仕様は未対応。
- DateIntervalは今回のISO期間とDateTime add/sub経路。日時拡張全体は未完成。
- 一般の循環配列参照、Serializableを含むSymfonyの全例外表示経路は未対応・未検証。
  正常HTTPの例外を握りつぶす変更はしていない。
- Docker用UID/GID設定を用意したが、今回の実測環境はMac/Podman。
  GitHub CIや公開操作は今回実施していない。

再実行・手動起動コマンドは[アプリREADME](../../research/fresh-laravel/README.md)にまとめた。

指定された9項目の新規Laravel完成条件と必須回帰は達成したため、Issue #1は今回の
対象範囲で解決と判定する。未実施の必須検証はない。上記の一般機能制限は残る。
公開操作は許可されていないため、GitHub上のIssueはOPENのままである。

[最終ソース識別情報](fresh-laravel-final-source-complete-20261021.json)、
[検証後の保護対象照合](fresh-laravel-audit-n6igouml/)、
[後片付け確認](fresh-laravel-cleanup-20261021.json)を保存した。
自分の検証サーバー・一時SQLite・runディレクトリは残っていない。
独立再構築した依存配置は再現用として保持し、以前からの8091番手動サーバーとDBは維持した。
作業ツリーには実装・新規アプリ・fixture・文書・ログの未コミット変更が残る。
